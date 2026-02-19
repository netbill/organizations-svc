package handlers

import (
	"context"
	"encoding/json"

	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/internal/core/modules/profile"
	"github.com/netbill/organizations-svc/pkg/evtypes"
	"github.com/segmentio/kafka-go"
)

func (h *Handlers) ProfileCreated(
	ctx context.Context,
	message kafka.Message,
) error {
	var payload evtypes.ProfileCreatedPayload
	if err := json.Unmarshal(message.Value, &payload); err != nil {
		return err
	}

	if _, err := h.modules.Profile.Create(ctx, models.Profile{
		AccountID: payload.AccountID,
		Username:  payload.Username,
		CreatedAt: payload.CreatedAt,
	}); err != nil {
		return err
	}

	return nil

}

func (h *Handlers) ProfileUpdated(
	ctx context.Context,
	message kafka.Message,
) error {
	var payload evtypes.ProfileUpdatedPayload
	if err := json.Unmarshal(message.Value, &payload); err != nil {
		return err
	}

	if _, err := h.modules.Profile.Update(ctx, payload.AccountID, profile.UpdateParams{
		Username:  payload.Username,
		Pseudonym: payload.Pseudonym,
		AvatarKey: payload.AvatarKey,
		Official:  payload.Official,
		UpdatedAt: payload.UpdatedAt,
	}); err != nil {
		return err
	}

	return nil
}

func (h *Handlers) ProfileDeleted(
	ctx context.Context,
	message kafka.Message,
) error {
	var payload evtypes.ProfileDeletedPayload
	if err := json.Unmarshal(message.Value, &payload); err != nil {
		return err
	}

	if err := h.modules.Profile.Delete(ctx, payload.AccountID); err != nil {
		return err
	}

	return nil
}
