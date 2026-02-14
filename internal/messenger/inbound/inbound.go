package inbound

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/internal/core/modules/profile"
)

type Inbound struct {
	core *core
}

type core struct {
	profile profileSvc
}

func New(profile profileSvc) *Inbound {
	return &Inbound{
		core: &core{
			profile: profile,
		},
	}
}

type profileSvc interface {
	Create(ctx context.Context, profile models.Profile) (models.Profile, error)
	Update(ctx context.Context, accountID uuid.UUID, params profile.UpdateParams) (models.Profile, error)

	Delete(ctx context.Context, accountID uuid.UUID) error
}
