package inbound

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/netbill/ape"
	"github.com/netbill/evebox/box/inbox"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/internal/messenger/contracts"
)

func (i Inbound) ProfileCreated(
	ctx context.Context,
	event inbox.Event,
) inbox.EventStatus {
	var payload contracts.ProfileCreatedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		i.log.Errorf("bad payload for %s, key %s, id: %s, error: %v", event.Type, event.Key, event.ID, err)
		return inbox.EventStatusFailed
	}

	profile := models.Profile{
		AccountID: payload.AccountID,
		Username:  payload.Username,
		CreatedAt: payload.CreatedAt,
	}
	if _, err := i.domain.CreateProfile(ctx, profile); err != nil {
		var ae *ape.Error
		if errors.As(err, &ae) {
			i.log.Errorf(
				"failed to upsert profile due to internal error, key %s, id: %s, error: %v",
				event.Key, event.ID, err,
			)
			return inbox.EventStatusPending
		}

		i.log.Errorf("failed to upsert profile, key %s, id: %s, error: %v", event.Key, event.ID, err)
		return inbox.EventStatusFailed
	}

	return inbox.EventStatusProcessed
}
