package place

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/models"
)

type repo interface {
	Get(ctx context.Context, id uuid.UUID) (models.Place, error)
	GetListByIDs(ctx context.Context, ids []uuid.UUID) ([]models.Place, error)
	ExistsForOrg(ctx context.Context, organizationID uuid.UUID) (bool, error)
	Exists(ctx context.Context, id uuid.UUID) (bool, error)

	Create(ctx context.Context, params CreateParams) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type tombstone interface {
	BuryPlace(ctx context.Context, placeID uuid.UUID) error
	PlaceIsBuried(ctx context.Context, placeID uuid.UUID) (bool, error)
}

type transactor interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}
