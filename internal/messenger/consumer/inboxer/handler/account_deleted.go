package handler

import (
	"context"
	"encoding/json"

	"github.com/umisto/cities-svc/internal/messenger/contracts"
	"github.com/umisto/kafkakit/box"
)

func (s Service) AccountDeleted(
	ctx context.Context,
	event box.InboxEvent,
) Status {
	var p contracts.AccountDeletedPayload
	if err := json.Unmarshal(event.Payload, &p); err != nil {
		s.log.Errorf("bad payload for %s, key %s, id: %s, error: %v", event.Type, event.Key, event.ID, err)
		return StatusFailed
	}

	if err := s.domain.DeleteProfile(ctx, p.AccountID); err != nil {
		s.log.Errorf("failed to delete profile, key %s, id: %s, error: %v", event.Key, event.ID, err)
		return StatusDelayed
	}

	return StatusProcessed
}
