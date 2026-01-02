package callbacker

import (
	"context"

	"github.com/google/uuid"
	"github.com/umisto/agglomerations-svc/internal/domain/models"
	"github.com/umisto/logium"
)

type Callbacker struct {
	log    logium.Logger
	domain domain
}

func New(log logium.Logger, domain domain) Callbacker {
	return Callbacker{
		log:    log,
		domain: domain,
	}
}

type domain interface {
	UpsertProfile(ctx context.Context, profile models.Profile) (models.Profile, error)
	UpdateUsername(ctx context.Context, accountID uuid.UUID, username string) (models.Profile, error)
	DeleteProfile(
		ctx context.Context,
		accountID uuid.UUID,
	) error
}
