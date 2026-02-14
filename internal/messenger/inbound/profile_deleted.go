package inbound

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/netbill/ape"
	"github.com/netbill/organizations-svc/internal/messenger/evtypes"
	"github.com/segmentio/kafka-go"
)

func (i *Inbound) ProfileDeleted(
	ctx context.Context,
	message kafka.Message,
) error {
	var payload evtypes.ProfileDeletedPayload
	if err := json.Unmarshal(message.Value, &payload); err != nil {
		return err
	}

	if err := i.core.profile.Delete(ctx, payload.AccountID); err != nil {
		return err
	}

	return nil
}
