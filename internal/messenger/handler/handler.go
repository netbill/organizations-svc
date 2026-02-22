package handler

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/domain"
	"github.com/netbill/organizations-svc/internal/core/modules/profile"
)

type Handler struct {
	modules *Modules
}

type Modules struct {
	Profile profileMod
}

func New(modules Modules) *Handler {
	return &Handler{
		modules: &modules,
	}
}

type profileMod interface {
	Create(ctx context.Context, profile domain.Profile) (domain.Profile, error)
	Update(ctx context.Context, accountID uuid.UUID, params profile.UpdateParams) (domain.Profile, error)

	Delete(ctx context.Context, accountID uuid.UUID) error
}
