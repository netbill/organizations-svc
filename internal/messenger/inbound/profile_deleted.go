package inbound

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/netbill/ape"
	"github.com/netbill/evebox/box/inbox"
	"github.com/netbill/organizations-svc/internal/messenger/contracts"
)

func (i *Inbound) ProfileDeleted(
	ctx context.Context,
	event inbox.Event,
) inbox.EventStatus {
	var payload contracts.ProfileDeletedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		i.log.Errorf("bad payload for %s, key %s, id: %s, error: %v", event.Type, event.Key, event.ID, err)
		return inbox.EventStatusFailed
	}

	if err := i.domain.DeleteProfile(ctx, payload.AccountID); err != nil {
		var ae *ape.Error
		if errors.As(err, &ae) {
			i.log.Errorf(
				"failed to delete profile due to internal error, key %s, id: %s, error: %v",
				event.Key, event.ID, err,
			)
			return inbox.EventStatusPending
		}

		i.log.Errorf("failed to delete profile, key %s, id: %s, error: %v", event.Key, event.ID, err)
		return inbox.EventStatusFailed
	}

	return inbox.EventStatusProcessed
}
