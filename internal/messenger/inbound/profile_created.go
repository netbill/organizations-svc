package inbound

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/netbill/ape"
	"github.com/netbill/evebox/box/inbox"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/internal/messenger/evtypes"
	"github.com/segmentio/kafka-go"
)

func (i *Inbound) ProfileCreated(
	ctx context.Context,
	message kafka.Message,
) error {
	var payload evtypes.ProfileCreatedPayload
	if err := json.Unmarshal(message.Value, &payload); err != nil {
		return err
	}

	if _, err := i.core.profile.Create(ctx, models.Profile{
		AccountID: payload.AccountID,
		Username:  payload.Username,
		Official:  payload.Official,
		Avatar:    payload.Avatar,
		Pseudonym: payload.Pseudonym,
		CreatedAt: payload.CreatedAt,
	}); err != nil {
		return err
	}

	return nil

}
