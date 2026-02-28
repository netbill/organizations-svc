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

func (h *Handler) PlaceCreated(
	ctx context.Context,
	event eventbox.InboxEvent,
) error {
	var payload evtypes.PlaceCreatedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return err
	}

	log := h.log.WithInboxEvent(event).With("place_id", payload.PlaceID)

	err := h.modules.Place.Create(ctx, place.CreateParams{
		ID:             payload.PlaceID,
		ClassID:        payload.ClassID,
		OrganizationID: payload.OrganizationID,

		Status:   payload.Status,
		Verified: payload.Verified,
		Point:    payload.Point,
		Address:  payload.Address,
		Name:     payload.Name,

		Description: payload.Description,
		IconKey:     payload.IconKey,
		BannerKey:   payload.BannerKey,
		Website:     payload.Website,
		Phone:       payload.Phone,

		CreatedAt: payload.CreatedAt,
	})
	switch {
	case errors.Is(err, errx.ErrorPlaceAlreadyExists):
		log.Debug("received place created event for already existing place")
		return nil
	case errors.Is(err, errx.ErrorPlaceDeleted):
		log.Debug("received place created event for deleted place")
		return nil
	case err != nil:
		return err
	default:
		log.Debug("place created successfully")
		return nil
	}
}

func (h *Handler) PlaceUpdated(
	ctx context.Context,
	event eventbox.InboxEvent,
) error {
	var payload evtypes.PlaceUpdatedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return err
	}

	log := h.log.WithInboxEvent(event).With("place_id", payload.PlaceID)

	err := h.modules.Place.Update(ctx, payload.PlaceID, place.UpdateParams{
		ClassID:  payload.ClassID,
		Name:     payload.Name,
		Address:  payload.Address,
		Status:   payload.Status,
		Verified: payload.Verified,

		Description: payload.Description,
		IconKey:     payload.IconKey,
		BannerKey:   payload.BannerKey,
		Website:     payload.Website,
		Phone:       payload.Phone,

		Version:   payload.Version,
		UpdatedAt: payload.UpdatedAt,
	})
	switch {
	case errors.Is(err, errx.ErrorPlaceDeleted):
		log.Debug("received place update event for deleted place")
		return nil
	case err != nil:
		return err
	default:
		log.Debug("place updated successfully")
		return nil
	}
}

func (h *Handler) PlaceDeleted(
	ctx context.Context,
	event eventbox.InboxEvent,
) error {
	var payload evtypes.PlaceDeletedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return err
	}

	log := h.log.WithInboxEvent(event).With("place_id", payload.PlaceID)

	err := h.modules.Place.Delete(ctx, payload.PlaceID)
	switch {
	case errors.Is(err, errx.ErrorPlaceDeleted):
		log.Debug("received place deleted event for already deleted place")
		return nil
	case err != nil:
		return err
	default:
		log.Debug("place deleted successfully")
		return nil
	}
}
