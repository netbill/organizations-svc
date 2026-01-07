package processor

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/messenger/contracts"
)

func (p Processor) AccountDeleted(
	ctx context.Context,
	event box.InboxEvent,
) string {
	var payload contracts.AccountDeletedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		p.log.Errorf("bad payload for %s, key %s, id: %s, error: %v", event.Type, event.Key, event.ID, err)
		return box.InboxStatusFailed
	}

	if err := p.domain.DeleteProfile(ctx, payload.AccountID); err != nil {
		switch {
		case errors.Is(err, errx.ErrorInternal):
			p.log.Errorf(
				"failed to delete profile due to internal error, key %s, id: %s, error: %v",
				event.Key, event.ID, err,
			)
			return box.InboxStatusPending
		default:
			p.log.Errorf("failed to delete profile, key %s, id: %s, error: %v", event.Key, event.ID, err)
			return box.InboxStatusFailed
		}
	}

	return box.InboxStatusProcessed
}
