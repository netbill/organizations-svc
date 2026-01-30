package inbound

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/logium"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/internal/core/modules/profile"
)

type Inbound struct {
	log    *logium.Logger
	domain domain
}

func New(log *logium.Logger, domain domain) Inbound {
	return Inbound{
		log:    log,
		domain: domain,
	}
}

type domain interface {
	CreateProfile(ctx context.Context, profile models.Profile) (models.Profile, error)
	UpdateProfile(ctx context.Context, accountID uuid.UUID, params profile.UpdateParams) (models.Profile, error)

	DeleteProfile(ctx context.Context, accountID uuid.UUID) error
}
