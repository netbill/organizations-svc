package handler

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/internal/core/modules/place"
	"github.com/netbill/organizations-svc/internal/core/modules/profile"
	"github.com/netbill/organizations-svc/pkg/log"
)

type Handler struct {
	log     *log.Logger
	modules *Modules
}

type Modules struct {
	Profile profileMod
	Place   placeMod
}

func New(log *log.Logger, modules Modules) *Handler {
	return &Handler{
		log:     log,
		modules: &modules,
	}
}

type profileMod interface {
	Create(ctx context.Context, profile models.Profile) (models.Profile, error)
	Update(ctx context.Context, accountID uuid.UUID, params profile.UpdateParams) (models.Profile, error)

	Delete(ctx context.Context, accountID uuid.UUID) error
}

type placeMod interface {
	Create(ctx context.Context, params place.CreateParams) error
	Update(ctx context.Context, id uuid.UUID, params place.UpdateParams) error
	Delete(ctx context.Context, id uuid.UUID) error
}
