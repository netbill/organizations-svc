package callbacker

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/netbill/kafkakit/box"
	"github.com/netbill/organizations-svc/internal/domain/errx"
	"github.com/netbill/organizations-svc/internal/messenger/contracts"
)

func (c Callbacker) AccountUsernameChanged(
	ctx context.Context,
	event box.InboxEvent,
) string {
	var p contracts.AccountUsernameChangePayload
	if err := json.Unmarshal(event.Payload, &p); err != nil {
		c.log.Errorf("bad payload for %s, key %s, id: %s, error: %v", event.Type, event.Key, event.ID, err)
		return box.InboxStatusFailed
	}

	if _, err := c.domain.UpdateUsername(ctx, p.Account.ID, p.Account.Username); err != nil {
		switch {
		case errors.Is(err, errx.ErrorInternal):
			c.log.Errorf(
				"failed to update username due to internal error, key %s, id: %s, error: %v",
				event.Key, event.ID, err,
			)
			return box.InboxStatusPending
		default:
			c.log.Errorf("failed to update username, key %s, id: %s, error: %v", event.Key, event.ID, err)
			return box.InboxStatusFailed
		}
	}

	return box.InboxStatusProcessed
}
