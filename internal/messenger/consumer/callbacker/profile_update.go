package callbacker

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/netbill/kafkakit/box"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/internal/messenger/contracts"
)

func (c Callbacker) ProfileUpdate(
	ctx context.Context,
	event box.InboxEvent,
) string {
	var p contracts.ProfileUpdatedPayload
	if err := json.Unmarshal(event.Payload, &p); err != nil {
		c.log.Errorf("bad payload for %s, key: %s, id: %s, error: %v", event.Type, event.Key, event.ID, err)
		return box.InboxStatusFailed
	}

	if _, err := c.domain.UpsertProfile(ctx, models.Profile{
		AccountID: p.Profile.AccountID,
		Username:  p.Profile.Username,
		Official:  p.Profile.Official,
		Pseudonym: p.Profile.Pseudonym,
	}); err != nil {
		switch {
		case errors.Is(err, errx.ErrorInternal):
			c.log.Errorf(
				"failed to upsert profile due to internal error, key: %s, id: %s, error: %v",
				event.Key, event.ID, err,
			)
			return box.InboxStatusPending
		default:
			c.log.Errorf("failed to upsert profile, key: %s, id: %s, error: %v", event.Key, event.ID, err)
			return box.InboxStatusFailed
		}
	}

	return box.InboxStatusProcessed
}
