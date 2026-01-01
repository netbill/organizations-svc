package handler

import (
	"context"
	"encoding/json"

	"github.com/umisto/agglomerations-svc/internal/domain/models"
	"github.com/umisto/agglomerations-svc/internal/messenger/contracts"
	"github.com/umisto/kafkakit/box"
)

func (s Service) AccountCreated(
	ctx context.Context,
	event box.InboxEvent,
) Status {
	var p contracts.AccountCreatedPayload
	if err := json.Unmarshal(event.Payload, &p); err != nil {
		s.log.Errorf("bad payload for %s, key %s, id: %s, error: %v", event.Type, event.Key, event.ID, err)
		return StatusFailed
	}
	profile := models.Profile{
		AccountID: p.Account.ID,
		Username:  p.Account.Username,
	}
	if _, err := s.domain.UpsertProfile(ctx, profile); err != nil {
		s.log.Errorf("failed to upsert profile, key %s, id: %s, error: %v", event.Key, event.ID, err)
		return StatusDelayed
	}

	return StatusProcessed
}
