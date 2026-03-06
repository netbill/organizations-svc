package evcontrollers

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/netbill/eventbox"
	"github.com/netbill/evtypes"
	"github.com/netbill/organizations-svc/internal/core/place"
	"github.com/netbill/organizations-svc/internal/errx"
	"github.com/netbill/organizations-svc/pkg/log"
)

type placeCore interface {
	Create(ctx context.Context, params place.CreateParams) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type PlaceController struct {
	log  *log.Logger
	core placeCore
}

func NewPlaceController(log *log.Logger, core placeCore) *PlaceController {
	return &PlaceController{
		log:  log,
		core: core,
	}
}

const operationPlaceCreated = "place_created"

func (c *PlaceController) Created(
	ctx context.Context,
	event eventbox.InboxEvent,
) error {
	var payload evtypes.PlaceCreatedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return err
	}

	logger := c.log.With("operation", operationPlaceCreated).
		With("place_id", payload.PlaceID)

	err := c.core.Create(ctx, place.CreateParams{
		ID:             payload.PlaceID,
		OrganizationID: payload.OrganizationID,
		CreatedAt:      payload.CreatedAt,
	})
	switch {
	case errors.Is(err, errx.ErrorPlaceAlreadyExists):
		logger.Warn("received place created event for already existing place")
		return nil
	case errors.Is(err, errx.ErrorPlaceDeleted):
		logger.Warn("received place created event for deleted place")
		return nil
	case err != nil:
		logger.WithError(err).Error("failed to create place")
		return err
	default:
		logger.Info("place created successfully")
		return nil
	}
}

const operationPlaceDeleted = "place_deleted"

func (c *PlaceController) Deleted(
	ctx context.Context,
	event eventbox.InboxEvent,
) error {
	var payload evtypes.PlaceDeletedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return err
	}

	logger := c.log.With("operation", operationPlaceDeleted).
		With("place_id", payload.PlaceID)

	err := c.core.Delete(ctx, payload.PlaceID)
	switch {
	case errors.Is(err, errx.ErrorPlaceDeleted):
		logger.Warn("received place deleted event for already deleted place")
		return nil
	case err != nil:
		logger.WithError(err).Error("failed to delete place")
		return err
	default:
		logger.Info("place deleted successfully")
		return nil
	}
}
