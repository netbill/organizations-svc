package member

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/errx"
	"github.com/netbill/organizations-svc/internal/models"
	"github.com/netbill/restkit/pagi"
)

type memberMessenger interface {
	WriteOrgMemberCreated(ctx context.Context, member models.Member) error
	WriteOrgMemberUpdated(ctx context.Context, member models.Member) error
	WriteOrgMemberDeleted(ctx context.Context, memberID uuid.UUID) error
}

type Service struct {
	org       org
	repo      repo
	tombstone tombstone
	tx        transactor
	messenger memberMessenger
}

type ServiceDeps struct {
	Auth      org
	Repo      repo
	Tombstone tombstone
	Tx        transactor
	Messenger memberMessenger
}

func NewMemberModule(deps ServiceDeps) *Service {
	return &Service{
		org:       deps.Auth,
		repo:      deps.Repo,
		tombstone: deps.Tombstone,
		tx:        deps.Tx,
		messenger: deps.Messenger,
	}
}

func (m *Service) GetByID(
	ctx context.Context,
	memberID uuid.UUID,
) (models.Member, error) {
	return m.repo.GetByID(ctx, memberID)
}

func (m *Service) GetByAccountAndOrgs(
	ctx context.Context,
	actor models.AccountActor,
	organizationIDs []uuid.UUID,
) ([]models.Member, error) {
	return m.repo.GetListForAccountAndOrgs(ctx, actor, organizationIDs)
}

type FilterParams struct {
	OrganizationID *uuid.UUID
	AccountID      *uuid.UUID
	Head           *bool
	Username       *string
	BestMatch      *string
	Label          *string
	Position       *string
}

func (m *Service) GetList(
	ctx context.Context,
	filter FilterParams,
	limit, offset uint,
) (pagi.Page[[]models.Member], error) {
	return m.repo.GetList(ctx, filter, limit, offset)
}

type UpdateParams struct {
	Position *string
	Label    *string
}

func (m *UpdateParams) HasChanges(model models.Member) bool {
	if m.Position != nil && !ptrEqual(m.Position, model.Position) {
		return true
	}
	if m.Label != nil && !ptrEqual(m.Label, model.Label) {
		return true
	}
	return false
}

func ptrEqual[T comparable](a, b *T) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

func (m *Service) Update(
	ctx context.Context,
	actor models.AccountActor,
	memberID uuid.UUID,
	params UpdateParams,
) (models.Member, error) {
	member, err := m.GetByID(ctx, memberID)
	if err != nil {
		return models.Member{}, err
	}

	if !params.HasChanges(member) {
		return member, nil
	}

	_, err = m.org.ValidateOrg(ctx, member.OrganizationID)
	if err != nil {
		return models.Member{}, err
	}

	_, err = m.org.AuthorizeOrgHead(ctx, actor, member.OrganizationID)
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

func (m *Service) Delete(
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

	_, err = m.org.AuthorizeOrgHead(ctx, actor, member.OrganizationID)
	if err != nil {
		return err
	}

	_, err = m.org.ValidateOrg(ctx, member.OrganizationID)
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
