package processor

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/netbill/kafkakit/box"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/messenger/contracts"
)

func (p Processor) AccountUsernameChanged(
	ctx context.Context,
	event box.InboxEvent,
) string {
	var payload contracts.AccountUsernameChangePayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		p.log.Errorf("bad payload for %s, key %s, id: %s, error: %v", event.Type, event.Key, event.ID, err)
		return box.InboxStatusFailed
	}

	if _, err := p.domain.UpdateUsername(ctx, payload.Account.ID, payload.Account.Username); err != nil {
		switch {
		case errors.Is(err, errx.ErrorInternal):
			p.log.Errorf(
				"failed to update username due to internal error, key %s, id: %s, error: %v",
				event.Key, event.ID, err,
			)
			return box.InboxStatusPending
		default:
			p.log.Errorf("failed to update username, key %s, id: %s, error: %v", event.Key, event.ID, err)
			return box.InboxStatusFailed
		}
	}

	return box.InboxStatusProcessed
}
