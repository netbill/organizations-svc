package profile

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/models"
)

type repo interface {
	Create(ctx context.Context, profile CreateParams) (models.Profile, error)

	GetByID(ctx context.Context, accountID uuid.UUID) (models.Profile, error)
	GetByUsername(ctx context.Context, username string) (models.Profile, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)

	GetListByIDs(ctx context.Context, accountIDs []uuid.UUID) ([]models.Profile, error)
	ExistsByAccountID(ctx context.Context, accountID uuid.UUID) (bool, error)

	Update(ctx context.Context, accountID uuid.UUID, params UpdateParams) (models.Profile, error)
	Delete(ctx context.Context, accountID uuid.UUID) error
}

type tombstone interface {
	BuryProfile(ctx context.Context, accountID uuid.UUID) error
	ProfileIsBuried(ctx context.Context, accountID uuid.UUID) (bool, error)
}

type transactor interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}
