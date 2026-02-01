package inbound

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/netbill/ape"
	"github.com/netbill/evebox/box/inbox"
	"github.com/netbill/organizations-svc/internal/core/modules/profile"
	"github.com/netbill/organizations-svc/internal/messenger/contracts"
)

func (i *Inbound) ProfileUpdated(
	ctx context.Context,
	event inbox.Event,
) inbox.EventStatus {
	var payload contracts.ProfileUpdatedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		i.log.Errorf("bad payload for %s, key: %s, id: %s, error: %v", event.Type, event.Key, event.ID, err)
		return inbox.EventStatusFailed
	}

	if _, err := i.domain.UpdateProfile(ctx, payload.AccountID, profile.UpdateParams{
		Username:  payload.Username,
		Official:  payload.Official,
		Pseudonym: payload.Pseudonym,
		UpdatedAt: payload.UpdatedAt,
	}); err != nil {
		var ae *ape.Error
		if errors.As(err, &ae) {
			i.log.Errorf(
				"failed to upsert profile due to internal error, key: %s, id: %s, error: %v",
				event.Key, event.ID, err,
			)
			return inbox.EventStatusPending
		}

		i.log.Errorf("failed to upsert profile, key: %s, id: %s, error: %v", event.Key, event.ID, err)
		return inbox.EventStatusFailed
	}

	return inbox.EventStatusProcessed
}
