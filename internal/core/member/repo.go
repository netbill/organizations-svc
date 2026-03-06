package member

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/models"
	"github.com/netbill/restkit/pagi"
)

type repo interface {
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
		filter FilterParams,
		limit, offset uint,
	) (pagi.Page[[]models.Member], error)

	ExistsForAccountAndOrg(
		ctx context.Context,
		accountID, organizationID uuid.UUID,
	) (bool, error)

	Update(ctx context.Context, ID uuid.UUID, params UpdateParams) (models.Member, error)

	Delete(ctx context.Context, memberID uuid.UUID) error
	DeleteForAccountAndOrg(ctx context.Context, accountID uuid.UUID) error
}

type tombstone interface {
	BuryMember(ctx context.Context, memberID uuid.UUID) error
	MemberIsBuried(ctx context.Context, memberID uuid.UUID) (bool, error)
}

type org interface {
	AuthorizeOrgHead(
		ctx context.Context,
		actor models.AccountActor,
		organizationID uuid.UUID,
	) (models.Member, error)

	AuthorizeOrgMember(
		ctx context.Context,
		accountID uuid.UUID,
		organizationID uuid.UUID,
	) (models.Member, error)

	ValidateOrg(
		ctx context.Context,
		organizationID uuid.UUID,
	) (models.Organization, error)
}

type transactor interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}
