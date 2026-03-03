package member

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/restkit/pagi"
)

type Module struct {
	repo      repo
	messenger messenger
}

func New(repo repo, messenger messenger) *Module {
	return &Module{
		repo:      repo,
		messenger: messenger,
	}
}

type repo interface {
	GetOrganizationByID(ctx context.Context, organizationID uuid.UUID) (models.Organization, error)

	CreateMember(ctx context.Context, accountID, organizationID uuid.UUID) (models.Member, error)
	GetMember(ctx context.Context, memberID uuid.UUID) (models.Member, error)
	GetMemberByAccountAndOrganization(
		ctx context.Context,
		accountID, organizationID uuid.UUID,
	) (models.Member, error)
	GetMembersByAccountAndOrgs(
		ctx context.Context,
		accountID uuid.UUID,
		organizationIDs []uuid.UUID,
	) ([]models.Member, error)
	UpdateMember(ctx context.Context, memberID uuid.UUID, params UpdateParams) (models.Member, error)
	DeleteMember(ctx context.Context, memberID uuid.UUID) error

	GetMembers(
		ctx context.Context,
		filter FilterParams,
		limit uint,
		offset uint,
	) (pagi.Page[[]models.Member], error)

	BuryMember(ctx context.Context, memberID uuid.UUID) error
	MemberIsBuried(ctx context.Context, accountID uuid.UUID) (bool, error)

	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type messenger interface {
	WriteOrgMemberCreated(ctx context.Context, member models.Member) error
	WriteOrgMemberUpdated(ctx context.Context, member models.Member) error
	WriteOrgMemberDeleted(ctx context.Context, memberID uuid.UUID) error
}

func (m *Module) getInitiator(
	ctx context.Context,
	actor models.AccountActor,
	organizationID uuid.UUID,
) (models.Member, error) {
	row, err := m.repo.GetMemberByAccountAndOrganization(ctx, actor, organizationID)
	if errors.Is(err, errx.ErrorMemberNotExists) {
		return models.Member{}, errx.ErrorInitiatorNotMemberOfOrganization.Raise(
			fmt.Errorf("initiator with account id %s is not a member of organization %s", actor, organizationID),
		)
	}

	return row, err
}
