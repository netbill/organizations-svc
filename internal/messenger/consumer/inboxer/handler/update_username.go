package handler

import (
	"context"
	"encoding/json"

	"github.com/umisto/cities-svc/internal/messenger/contracts"
	"github.com/umisto/kafkakit/box"
)

func (s Service) UpdateUsername(
	ctx context.Context,
	event box.InboxEvent,
) Status {
	var p contracts.AccountUsernameChangePayload
	if err := json.Unmarshal(event.Payload, &p); err != nil {
		s.log.Errorf("bad payload for %s, key %s, id: %s, error: %v", event.Type, event.Key, event.ID, err)
		return StatusFailed
	}

	if _, err := s.domain.UpdateUsername(ctx, p.Account.ID, p.Account.Username); err != nil {
		s.log.Errorf("failed to update username, key %s, id: %s, error: %v", event.Key, event.ID, err)
		return StatusDelayed
	}

	return StatusProcessed
}
