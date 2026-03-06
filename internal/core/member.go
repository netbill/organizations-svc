package core

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/restkit/pagi"
)

type memberRepo interface {
	Create(
		ctx context.Context,
		accountID, organizationID uuid.UUID,
		head bool,
	) (models.Member, error)

	GetByID(ctx context.Context, memberID uuid.UUID) (models.Member, error)

	GetListForAccountAndOrgs(
		ctx context.Context,
		accountID uuid.UUID,
		organizationIDs []uuid.UUID,
	) ([]models.Member, error)

	GetList(
		ctx context.Context,
		filter MemberFilterParams,
		limit, offset uint,
	) (pagi.Page[[]models.Member], error)

	ExistsForAccountAndOrg(
		ctx context.Context,
		accountID, organizationID uuid.UUID,
	) (bool, error)

	Update(ctx context.Context, ID uuid.UUID, params MemberUpdateParams) (models.Member, error)

	Delete(ctx context.Context, memberID uuid.UUID) error
	DeleteForAccountAndOrg(ctx context.Context, accountID uuid.UUID) error
}

type memberTombstoneRepo interface {
	BuryMember(ctx context.Context, memberID uuid.UUID) error
	MemberIsBuried(ctx context.Context, memberID uuid.UUID) (bool, error)
}

type orgAuth interface {
	authorizeOrgHead(
		ctx context.Context,
		actor models.AccountActor,
		organizationID uuid.UUID,
	) (models.Member, error)

	authorizeOrgMember(
		ctx context.Context,
		accountID uuid.UUID,
		organizationID uuid.UUID,
	) (models.Member, error)

	validateOrg(
		ctx context.Context,
		organizationID uuid.UUID,
	) (models.Organization, error)
}

type memberMessenger interface {
	WriteOrgMemberCreated(ctx context.Context, member models.Member) error
	WriteOrgMemberUpdated(ctx context.Context, member models.Member) error
	WriteOrgMemberDeleted(ctx context.Context, memberID uuid.UUID) error
}

type MemberModule struct {
	auth      orgAuth
	repo      memberRepo
	tombstone memberTombstoneRepo
	tx        transactor
	messenger memberMessenger
}

type MemberDeps struct {
	Auth      orgAuth
	Repo      memberRepo
	Tombstone memberTombstoneRepo
	Tx        transactor
	Messenger memberMessenger
}

func NewMemberModule(deps MemberDeps) *MemberModule {
	return &MemberModule{
		auth:      deps.Auth,
		repo:      deps.Repo,
		tombstone: deps.Tombstone,
		tx:        deps.Tx,
		messenger: deps.Messenger,
	}
}

func (m *MemberModule) GetByID(
	ctx context.Context,
	memberID uuid.UUID,
) (models.Member, error) {
	return m.repo.GetByID(ctx, memberID)
}

func (m *MemberModule) GetByAccountAndOrgs(
	ctx context.Context,
	actor models.AccountActor,
	organizationIDs []uuid.UUID,
) ([]models.Member, error) {
	return m.repo.GetListForAccountAndOrgs(ctx, actor, organizationIDs)
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
	return m.repo.GetList(ctx, filter, limit, offset)
}

type MemberUpdateParams struct {
	Position *string
	Label    *string
}

func (m *MemberUpdateParams) HasChanges(model models.Member) bool {
	if m.Position != nil && !ptrEqual(m.Position, model.Position) {
		return true
	}
	if m.Label != nil && !ptrEqual(m.Label, model.Label) {
		return true
	}
	return false
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

	if !params.HasChanges(member) {
		return member, nil
	}

	_, err = m.auth.validateOrg(ctx, member.OrganizationID)
	if err != nil {
		return models.Member{}, err
	}

	_, err = m.auth.authorizeOrgHead(ctx, actor, member.OrganizationID)
	if err != nil {
		return models.Member{}, err
	}

	err = m.tx.Transaction(ctx, func(ctx context.Context) error {
		member, err = m.repo.Update(ctx, memberID, params)
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
	buried, err := m.tombstone.MemberIsBuried(ctx, memberID)
	if err != nil {
		return err
	}
	if buried {
		return errx.ErrorMemberDeleted.Raise(
			fmt.Errorf("member with id %s is already deleted", memberID),
		)
	}

	member, err := m.repo.GetByID(ctx, memberID)
	if err != nil {
		return err
	}
	if member.Head {
		return errx.ErrorCannotDeleteOrganizationHeadMember.Raise(
			fmt.Errorf("cannot delete organization head member %s", member.ID),
		)
	}

	_, err = m.auth.authorizeOrgHead(ctx, actor, member.OrganizationID)
	if err != nil {
		return err
	}

	_, err = m.auth.validateOrg(ctx, member.OrganizationID)
	if err != nil {
		return err
	}

	return m.tx.Transaction(ctx, func(ctx context.Context) error {
		if err = m.tombstone.BuryMember(ctx, memberID); err != nil {
			return err
		}

		if err = m.repo.Delete(ctx, memberID); err != nil {
			return err
		}

		return m.messenger.WriteOrgMemberDeleted(ctx, memberID)
	})
}
