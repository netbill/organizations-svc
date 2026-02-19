package handlers

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/internal/core/modules/profile"
)

type Handlers struct {
	modules *Modules
}

type Modules struct {
	Profile profileMod
}

func New(modules Modules) *Handlers {
	return &Handlers{
		modules: &modules,
	}
}

type profileMod interface {
	Create(ctx context.Context, profile models.Profile) (models.Profile, error)
	Update(ctx context.Context, accountID uuid.UUID, params profile.UpdateParams) (models.Profile, error)

	Delete(ctx context.Context, accountID uuid.UUID) error
}
