package callbacker

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/umisto/agglomerations-svc/internal/domain/errx"
	"github.com/umisto/agglomerations-svc/internal/messenger/contracts"
	"github.com/umisto/kafkakit/box"
)

func (c Callbacker) AccountDeleted(
	ctx context.Context,
	event box.InboxEvent,
) string {
	var p contracts.AccountDeletedPayload
	if err := json.Unmarshal(event.Payload, &p); err != nil {
		c.log.Errorf("bad payload for %s, key %s, id: %s, error: %v", event.Type, event.Key, event.ID, err)
		return box.InboxStatusFailed
	}

	if err := c.domain.DeleteProfile(ctx, p.AccountID); err != nil {
		switch {
		case errors.Is(err, errx.ErrorInternal):
			c.log.Errorf(
				"failed to delete profile due to internal error, key %s, id: %s, error: %v",
				event.Key, event.ID, err,
			)
			return box.InboxStatusPending
		default:
			c.log.Errorf("failed to delete profile, key %s, id: %s, error: %v", event.Key, event.ID, err)
			return box.InboxStatusFailed
		}
	}

	return box.InboxStatusProcessed
}
