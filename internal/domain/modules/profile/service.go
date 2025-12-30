package profile

import (
	"context"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/models"
)

type Service struct {
	repo repo
}

type repo interface {
	UpsertProfile(ctx context.Context, profile models.Profile) (models.Profile, error)
	UpdateUsername(ctx context.Context, accountID uuid.UUID, username string) (models.Profile, error)
	GetProfileByAccountID(ctx context.Context, accountID uuid.UUID) (models.Profile, error)
	GetProfileByUsername(ctx context.Context, username string) (models.Profile, error)
	DeleteProfileByAccountID(ctx context.Context, accountID uuid.UUID) error

	DeleteMembershipsByAccountID(ctx context.Context, accountID uuid.UUID) error

	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

func NewService(repo repo) Service {
	return Service{repo: repo}
}
