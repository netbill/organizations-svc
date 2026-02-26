package handler

import (
	"context"
	"encoding/json"

	"github.com/netbill/eventbox"
	"github.com/netbill/evtypes"
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

	return h.modules.Place.Create(ctx, place.CreateParams{
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
}

func (h *Handler) PlaceUpdated(
	ctx context.Context,
	event eventbox.InboxEvent,
) error {
	var payload evtypes.PlaceUpdatedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return err
	}

	return h.modules.Place.Update(ctx, payload.PlaceID, place.UpdateParams{
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
}

func (h *Handler) PlaceDeleted(
	ctx context.Context,
	event eventbox.InboxEvent,
) error {
	var payload evtypes.PlaceDeletedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return err
	}

	return h.modules.Place.Delete(ctx, payload.PlaceID)
}
