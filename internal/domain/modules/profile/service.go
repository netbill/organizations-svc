package profile

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/domain/models"
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

	DeleteMembersByAccountID(ctx context.Context, accountID uuid.UUID) error

	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

func New(repo repo) Service {
	return Service{repo: repo}
}
