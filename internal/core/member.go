package core

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/restkit/pagi"
)

type memberRepository interface {
	CreateMember(ctx context.Context, accountID, organizationID uuid.UUID) (models.Member, error)
	CreateMemberHead(ctx context.Context, accountID, organizationID uuid.UUID) (models.Member, error)

	GetMember(ctx context.Context, memberID uuid.UUID) (models.Member, error)
	GetMemberByAccountAndOrganization(ctx context.Context, accountID, organizationID uuid.UUID) (models.Member, error)
	GetMembersByAccountAndOrgs(ctx context.Context, accountID uuid.UUID, organizationIDs []uuid.UUID) ([]models.Member, error)
	GetMembers(ctx context.Context, filter MemberFilterParams, limit, offset uint) (pagi.Page[[]models.Member], error)
	MemberExists(ctx context.Context, accountID, organizationID uuid.UUID) (bool, error)

	UpdateMember(ctx context.Context, ID uuid.UUID, params MemberUpdateParams) (models.Member, error)

	DeleteMember(ctx context.Context, memberID uuid.UUID) error
	DeleteMembersByAccountID(ctx context.Context, accountID uuid.UUID) error
}

type memberOrgRepository interface {
	GetOrganizationByID(ctx context.Context, organizationID uuid.UUID) (models.Organization, error)
}

type memberTombstoneRepository interface {
	BuryMember(ctx context.Context, memberID uuid.UUID) error
	MemberIsBuried(ctx context.Context, memberID uuid.UUID) (bool, error)
}

type memberMessenger interface {
	WriteOrgMemberCreated(ctx context.Context, member models.Member) error
	WriteOrgMemberUpdated(ctx context.Context, member models.Member) error
	WriteOrgMemberDeleted(ctx context.Context, memberID uuid.UUID) error
}

type MemberModule struct {
	memberRepo    memberRepository
	orgRepo       memberOrgRepository
	tombstoneRepo memberTombstoneRepository
	tx            transactor
	messenger     memberMessenger
}

func NewMemberModule(
	memberRepo memberRepository,
	orgRepo memberOrgRepository,
	tombstoneRepo memberTombstoneRepository,
	tx transactor,
	messenger memberMessenger,
) *MemberModule {
	return &MemberModule{
		memberRepo:    memberRepo,
		orgRepo:       orgRepo,
		tombstoneRepo: tombstoneRepo,
		tx:            tx,
		messenger:     messenger,
	}
}

func (m *MemberModule) authorizeOrgHead(
	ctx context.Context,
	actor models.AccountActor,
	organizationID uuid.UUID,
) (models.Organization, models.Member, error) {
	org, err := m.orgRepo.GetOrganizationByID(ctx, organizationID)
	if err != nil {
		return models.Organization{}, models.Member{}, err
	}

	if org.Status == models.OrganizationStatusSuspended {
		return models.Organization{}, models.Member{}, errx.ErrorOrganizationIsSuspended.Raise(
			fmt.Errorf("organization with id %s is suspended", organizationID),
		)
	}

	member, err := m.memberRepo.GetMemberByAccountAndOrganization(ctx, actor, organizationID)
	if errors.Is(err, errx.ErrorMemberNotExists) {
		return models.Organization{}, models.Member{}, errx.ErrorInitiatorNotMemberOfOrganization.Raise(
			fmt.Errorf("initiator with account id %s is not a member of organization %s", actor, organizationID),
		)
	}
	if err != nil {
		return models.Organization{}, models.Member{}, err
	}

	if !member.Head {
		return models.Organization{}, models.Member{}, errx.ErrorNotOrganizationHead.Raise(
			fmt.Errorf(
				"only organization head member can manage members, but member %s is not head", member.ID,
			),
		)
	}

	return org, member, nil
}

func (m *MemberModule) GetByID(
	ctx context.Context,
	memberID uuid.UUID,
) (models.Member, error) {
	return m.memberRepo.GetMember(ctx, memberID)
}

func (m *MemberModule) GetByAccountAndOrgs(
	ctx context.Context,
	actor models.AccountActor,
	organizationIDs []uuid.UUID,
) ([]models.Member, error) {
	return m.memberRepo.GetMembersByAccountAndOrgs(ctx, actor, organizationIDs)
}

func (m *MemberModule) GetByAccountAndOrganization(
	ctx context.Context,
	actor models.AccountActor,
	organizationID uuid.UUID,
) (models.Member, error) {
	return m.memberRepo.GetMemberByAccountAndOrganization(ctx, actor, organizationID)
}

type MemberFilterParams struct {
	OrganizationID *uuid.UUID
	AccountID      *uuid.UUID
	Head           *bool
	Username       *string
	BestMatch      *string
	Label          *string
	Position       *string
}

func (m *MemberModule) GetList(
	ctx context.Context,
	filter MemberFilterParams,
	limit, offset uint,
) (pagi.Page[[]models.Member], error) {
	return m.memberRepo.GetMembers(ctx, filter, limit, offset)
}

type MemberUpdateParams struct {
	Position *string
	Label    *string
}

func (m *MemberModule) Update(
	ctx context.Context,
	actor models.AccountActor,
	memberID uuid.UUID,
	params MemberUpdateParams,
) (models.Member, error) {
	member, err := m.GetByID(ctx, memberID)
	if err != nil {
		return models.Member{}, err
	}

	_, _, err = m.authorizeOrgHead(ctx, actor, member.OrganizationID)
	if err != nil {
		return models.Member{}, err
	}

	err = m.tx.Transaction(ctx, func(ctx context.Context) error {
		member, err = m.memberRepo.UpdateMember(ctx, memberID, params)
		if err != nil {
			return err
		}

		return m.messenger.WriteOrgMemberUpdated(ctx, member)
	})

	return member, err
}

func (m *MemberModule) Delete(
	ctx context.Context,
	actor models.AccountActor,
	memberID uuid.UUID,
) error {
	member, err := m.GetByID(ctx, memberID)
	if errors.Is(err, errx.ErrorMemberNotExists) {
		buried, err := m.tombstoneRepo.MemberIsBuried(ctx, memberID)
		if err != nil {
			return err
		}
		if buried {
			return errx.ErrorMemberDeleted.Raise(
				fmt.Errorf("member with id %s is already deleted", memberID),
			)
		}
	}
	if err != nil {
		return err
	}

	if member.Head {
		return errx.ErrorCannotDeleteOrganizationHeadMember.Raise(
			fmt.Errorf("cannot delete organization head member %s", member.ID),
		)
	}
	if member.AccountID == actor {
		return errx.ErrorCannotDeleteSelf.Raise(
			fmt.Errorf("account cannot delete itself as member %s", member.ID),
		)
	}

	_, _, err = m.authorizeOrgHead(ctx, actor, member.OrganizationID)
	if err != nil {
		return err
	}

	return m.tx.Transaction(ctx, func(ctx context.Context) error {
		if err = m.tombstoneRepo.BuryMember(ctx, memberID); err != nil {
			return err
		}

		if err = m.memberRepo.DeleteMember(ctx, memberID); err != nil {
			return err
		}

		return m.messenger.WriteOrgMemberDeleted(ctx, memberID)
	})
}

func (m *MemberModule) DeleteSelf(
	ctx context.Context,
	actor models.AccountActor,
	orgID uuid.UUID,
) error {
	member, err := m.memberRepo.GetMemberByAccountAndOrganization(ctx, actor, orgID)
	if err != nil {
		return err
	}
	if member.Head {
		return errx.ErrorCannotDeleteOrganizationHeadMember.Raise(
			fmt.Errorf("cannot delete organization head member %s", member.ID),
		)
	}

	return m.tx.Transaction(ctx, func(ctx context.Context) error {
		if err = m.tombstoneRepo.BuryMember(ctx, member.ID); err != nil {
			return err
		}

		if err = m.memberRepo.DeleteMember(ctx, member.ID); err != nil {
			return err
		}

		return m.messenger.WriteOrgMemberDeleted(ctx, member.ID)
	})
}
