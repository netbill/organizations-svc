package profile

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/domain"
)

type Module struct {
	repo repo
}

type repo interface {
	CreateProfile(ctx context.Context, profile domain.Profile) (domain.Profile, error)
	UpdateProfile(ctx context.Context, accountID uuid.UUID, params UpdateParams) (domain.Profile, error)

	DeleteProfileByAccountID(ctx context.Context, accountID uuid.UUID) error
	DeleteMembersByAccountID(ctx context.Context, accountID uuid.UUID) error

	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

func New(repo repo) *Module {
	return &Module{repo: repo}
}
