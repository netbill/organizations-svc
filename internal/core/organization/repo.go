package organization

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/models"
	"github.com/netbill/restkit/pagi"
)

type repo interface {
	Create(
		ctx context.Context,
		params CreateParams,
	) (models.Organization, error)

	Get(
		ctx context.Context,
		ID uuid.UUID,
	) (models.Organization, error)
	GetListByIds(
		ctx context.Context,
		IDs []uuid.UUID,
	) ([]models.Organization, error)
	GetList(
		ctx context.Context,
		filter FilterParams,
		limit, offset uint,
	) (pagi.Page[[]models.Organization], error)
	GetForAccountAndOrg(
		ctx context.Context,
		accountID uuid.UUID,
		limit, offset uint,
	) (pagi.Page[[]models.Organization], error)

	Update(
		ctx context.Context,
		organizationID uuid.UUID,
		params UpdateParams,
	) (models.Organization, error)
	UpdateStatus(
		ctx context.Context,
		ID uuid.UUID,
		status string,
	) (models.Organization, error)

	Delete(ctx context.Context, ID uuid.UUID) error
}

type tombstone interface {
	BuryOrganization(ctx context.Context, organizationID uuid.UUID) error
	OrganizationIsBuried(ctx context.Context, organizationID uuid.UUID) (bool, error)
}

type memberRepo interface {
	Create(
		ctx context.Context,
		accountID, organizationID uuid.UUID,
		head bool,
	) (models.Member, error)

	GetForAccountAndOrg(
		ctx context.Context,
		accountID models.AccountActor,
		organizationID uuid.UUID,
	) (models.Member, error)
}

type placeRepo interface {
	ExistsForOrg(ctx context.Context, organizationID uuid.UUID) (bool, error)
}

type transactor interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}
