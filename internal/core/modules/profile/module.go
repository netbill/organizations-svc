package profile

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/models"
)

type Module struct {
	repo repo
}

type repo interface {
	CreateProfile(ctx context.Context, profile models.Profile) (models.Profile, error)
	UpdateProfile(ctx context.Context, accountID uuid.UUID, params UpdateParams) (models.Profile, error)
	GetProfileByAccountID(ctx context.Context, accountID uuid.UUID) (models.Profile, error)
	GetProfilesByAccountIDs(ctx context.Context, accountIDs []uuid.UUID) ([]models.Profile, error)
	GetProfileByUsername(ctx context.Context, username string) (models.Profile, error)
	ExistsProfileByUsername(ctx context.Context, username string) (bool, error)
	ExistsProfileByAccountID(ctx context.Context, accountID uuid.UUID) (bool, error)

	DeleteProfileByAccountID(ctx context.Context, accountID uuid.UUID) error

	BuryProfile(ctx context.Context, accountID uuid.UUID) error
	ProfileIsBuried(ctx context.Context, accountID uuid.UUID) (bool, error)

	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

func New(repo repo) *Module {
	return &Module{repo: repo}
}
