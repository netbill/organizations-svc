package profile

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/models"
)

type Service struct {
	repo repo
}

type repo interface {
	CreateProfile(ctx context.Context, profile models.Profile) (models.Profile, error)
	UpdateProfile(ctx context.Context, accountID uuid.UUID, params UpdateParams) (models.Profile, error)

	DeleteProfileByAccountID(ctx context.Context, accountID uuid.UUID) error
	DeleteMembersByAccountID(ctx context.Context, accountID uuid.UUID) error

	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

func New(repo repo) Service {
	return Service{repo: repo}
}
