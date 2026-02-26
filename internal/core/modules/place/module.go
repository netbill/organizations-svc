package place

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/models"
)

type Module struct {
	repo repo
}

func New(
	repo repo,
) *Module {
	return &Module{
		repo: repo,
	}
}

type repo interface {
	CreatePlace(ctx context.Context, params CreateParams) error

	GetPlaceByID(ctx context.Context, id uuid.UUID) (models.Place, error)
	GetPlacesByIDs(ctx context.Context, ids []uuid.UUID) ([]models.Place, error)

	UpdatePlaceByID(ctx context.Context, id uuid.UUID, params UpdateParams) error
	DeletePlaceByID(ctx context.Context, id uuid.UUID) error
}
