package inbound

import (
	"context"
	"encoding/json"

	"github.com/netbill/organizations-svc/internal/core/modules/profile"
	"github.com/netbill/organizations-svc/internal/messenger/evtypes"
	"github.com/segmentio/kafka-go"
)

func (i *Inbound) ProfileUpdated(
	ctx context.Context,
	message kafka.Message,
) error {
	var payload evtypes.ProfileUpdatedPayload
	if err := json.Unmarshal(message.Value, &payload); err != nil {
		return err
	}

	if _, err := i.core.profile.Update(ctx, payload.AccountID, profile.UpdateParams{
		Username:  payload.Username,
		Official:  payload.Official,
		Avatar:    payload.Avatar,
		Pseudonym: payload.Pseudonym,
	}); err != nil {
		return err
	}

	return nil
}
