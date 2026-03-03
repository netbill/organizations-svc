package handler

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/netbill/eventbox"
	"github.com/netbill/evtypes"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/modules/place"
)

const operationPlaceCreated = "place_created"

func (h *Handler) PlaceCreated(
	ctx context.Context,
	event eventbox.InboxEvent,
) error {
	var payload evtypes.PlaceCreatedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return err
	}

	log := h.log.With("operation", operationPlaceCreated).
		With("place_id", payload.PlaceID)

	err := h.modules.Place.Create(ctx, place.CreateParams{
		ID:             payload.PlaceID,
		OrganizationID: payload.OrganizationID,
		CreatedAt:      payload.CreatedAt,
	})
	switch {
	case errors.Is(err, errx.ErrorPlaceAlreadyExists):
		log.Debug("received place created event for already existing place")
		return nil
	case errors.Is(err, errx.ErrorPlaceDeleted):
		log.Debug("received place created event for deleted place")
		return nil
	case err != nil:
		log.WithError(err).Error("failed to create place")
		return err
	default:
		log.Debug("place created successfully")
		return nil
	}
}

const operationPlaceDeleted = "place_deleted"

func (h *Handler) PlaceDeleted(
	ctx context.Context,
	event eventbox.InboxEvent,
) error {
	var payload evtypes.PlaceDeletedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return err
	}

	log := h.log.With("operation", operationPlaceDeleted).
		With("place_id", payload.PlaceID)

	err := h.modules.Place.Delete(ctx, payload.PlaceID)
	switch {
	case errors.Is(err, errx.ErrorPlaceDeleted):
		log.Debug("received place deleted event for already deleted place")
		return nil
	case err != nil:
		log.WithError(err).Error("failed to delete place")
		return err
	default:
		log.Debug("place deleted successfully")
		return nil
	}
}
