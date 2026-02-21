package handler

import (
	"context"
	"encoding/json"

	"github.com/netbill/eventbox"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/internal/core/modules/profile"
	"github.com/netbill/organizations-svc/pkg/evtypes"
)

func (h *Handler) ProfileCreated(
	ctx context.Context,
	event eventbox.InboxEvent,
) error {
	var payload evtypes.ProfileCreatedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
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

func (h *Handler) ProfileUpdated(
	ctx context.Context,
	event eventbox.InboxEvent,
) error {
	var payload evtypes.ProfileUpdatedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
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

func (h *Handler) ProfileDeleted(
	ctx context.Context,
	event eventbox.InboxEvent,
) error {
	var payload evtypes.ProfileDeletedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return err
	}

	if err := h.modules.Profile.Delete(ctx, payload.AccountID); err != nil {
		return err
	}

	return nil
}
