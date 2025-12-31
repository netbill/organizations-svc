package handler

import (
	"context"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/models"
	"github.com/umisto/cities-svc/internal/messenger/contracts"
	"github.com/umisto/kafkakit/box"
	"github.com/umisto/logium"
)

type Status uint8

const (
	StatusProcessed Status = iota
	StatusDelayed
	StatusFailed
)

type Service struct {
	log    logium.Logger
	domain domain
}

func New(log logium.Logger, domain domain) Service {
	return Service{
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

func (s Service) Handle(ctx context.Context, event box.InboxEvent) Status {
	switch event.Type {
	case contracts.AccountCreatedEvent:
		return s.AccountCreated(ctx, event)
	case contracts.AccountDeletedEvent:
		return s.AccountDeleted(ctx, event)
	case contracts.AccountUsernameChangeEvent:
		return s.UpdateUsername(ctx, event)
	case contracts.ProfileUpdatedEvent:
		return s.ProfileUpdate(ctx, event)
	default:
		s.log.Errorf("unknown inbox event type: %s, id: %s, key: %s", event.Type, event.ID, event.Key)
		return StatusFailed
	}
}
