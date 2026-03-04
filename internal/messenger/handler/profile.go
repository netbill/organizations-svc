package handler

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/netbill/eventbox"
	"github.com/netbill/evtypes"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/internal/core/modules/profile"
)

const operationProfileCreated = "profile_created"

func (h *Handler) ProfileCreated(
	ctx context.Context,
	event eventbox.InboxEvent,
) error {
	var payload evtypes.ProfileCreatedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return err
	}

	log := h.log.WithOperation(operationProfileCreated).
		With("account_id", payload.AccountID)

	_, err := h.modules.Profile.Create(ctx, models.Profile{
		AccountID: payload.AccountID,
		Username:  payload.Username,
		CreatedAt: payload.CreatedAt,
	})
	switch {
	case errors.Is(err, errx.ErrorProfileAlreadyExists):
		log.Debug("received profile created event for already existing profile")
		return nil
	case errors.Is(err, errx.ErrorProfileDeleted):
		log.Debug("received profile created event for deleted profile")
		return nil
	case err != nil:
		log.WithError(err).Error("failed to create profile")
		return err
	default:
		log.Debug("profile created successfully")
		return nil
	}
}

const operationProfileDeleted = "profile_deleted"

func (h *Handler) ProfileUpdated(
	ctx context.Context,
	event eventbox.InboxEvent,
) error {
	var payload evtypes.ProfileUpdatedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return err
	}

	log := h.log.WithOperation(operationProfileDeleted).
		With("account_id", payload.AccountID)

	_, err := h.modules.Profile.Update(ctx, payload.AccountID, profile.UpdateParams{
		Username:  payload.Username,
		Pseudonym: payload.Pseudonym,
		AvatarKey: payload.AvatarKey,
		Version:   payload.Version,
		UpdatedAt: payload.UpdatedAt,
	})
	switch {
	case errors.Is(err, errx.ErrorProfileDeleted):
		log.Debug("received profile updated event for deleted profile")
		return nil
	case err != nil:
		log.WithError(err).Error("failed to update profile")
		return err
	default:
		log.Debug("profile updated successfully")
		return nil
	}
}

func (h *Handler) ProfileDeleted(
	ctx context.Context,
	event eventbox.InboxEvent,
) error {
	var payload evtypes.ProfileDeletedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return err
	}

	log := h.log.WithOperation(operationProfileDeleted).
		With("account_id", payload.AccountID)

	err := h.modules.Profile.Delete(ctx, payload.AccountID)
	switch {
	case errors.Is(err, errx.ErrorProfileDeleted):
		log.Debug("received profile deleted event for already deleted profile")
		return nil
	case err != nil:
		log.WithError(err).Error("failed to delete profile")
		return err
	default:
		log.Debug("profile deleted successfully")
		return nil
	}
}
