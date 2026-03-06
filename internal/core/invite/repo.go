package invite

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/models"
	"github.com/netbill/restkit/pagi"
)

type repo interface {
	Create(ctx context.Context, params CreateParams) (models.Invite, error)

	Get(ctx context.Context, inviteID uuid.UUID) (models.Invite, error)
	GetList(
		ctx context.Context,
		params FilterParams,
		limit, offset uint,
	) (pagi.Page[[]models.Invite], error)
	ExistActiveForAccountInOrg(ctx context.Context, accountID, organizationID uuid.UUID) (bool, error)

	UpdateStatus(ctx context.Context, inviteID uuid.UUID, status string) (models.Invite, error)

	Delete(ctx context.Context, inviteID uuid.UUID) error
}

type member interface {
	Create(
		ctx context.Context,
		accountID, organizationID uuid.UUID,
		head bool,
	) (models.Member, error)

	ExistsForAccountAndOrg(ctx context.Context, accountID, organizationID uuid.UUID) (bool, error)
}

type profile interface {
	ExistsByAccountID(ctx context.Context, accountID uuid.UUID) (bool, error)
}

type transactor interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
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
