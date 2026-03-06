package evcontroller

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/netbill/eventbox"
	"github.com/netbill/evtypes"
	"github.com/netbill/organizations-svc/internal/core"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/pkg/log"
)

type profileCore interface {
	Create(ctx context.Context, profile models.Profile) (models.Profile, error)
	Update(ctx context.Context, accountID uuid.UUID, params core.ProfileUpdateParams) (models.Profile, error)
	Delete(ctx context.Context, accountID uuid.UUID) error
}

type ProfileController struct {
	log  *log.Logger
	core profileCore
}

func NewProfileController(log *log.Logger, core profileCore) *ProfileController {
	return &ProfileController{
		log:  log,
		core: core,
	}
}

const operationProfileCreated = "profile_created"

func (c *ProfileController) Created(
	ctx context.Context,
	event eventbox.InboxEvent,
) error {
	var payload evtypes.ProfileCreatedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return err
	}

	logger := c.log.WithOperation(operationProfileCreated).
		With("account_id", payload.AccountID)

	_, err := c.core.Create(ctx, models.Profile{
		AccountID: payload.AccountID,
		Username:  payload.Username,
		CreatedAt: payload.CreatedAt,
	})
	switch {
	case errors.Is(err, errx.ErrorProfileAlreadyExists):
		logger.Debug("received profile created event for already existing profile")
		return nil
	case errors.Is(err, errx.ErrorProfileDeleted):
		logger.Debug("received profile created event for deleted profile")
		return nil
	case err != nil:
		logger.WithError(err).Error("failed to create profile")
		return err
	default:
		logger.Debug("profile created successfully")
		return nil
	}
}

const operationProfileDeleted = "profile_deleted"

func (c *ProfileController) Updated(
	ctx context.Context,
	event eventbox.InboxEvent,
) error {
	var payload evtypes.ProfileUpdatedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return err
	}

	logger := c.log.WithOperation(operationProfileDeleted).With("account_id", payload.AccountID)

	_, err := c.core.Update(ctx, payload.AccountID, core.ProfileUpdateParams{
		Username:  payload.Username,
		Pseudonym: payload.Pseudonym,
		AvatarKey: payload.AvatarKey,
		Version:   payload.Version,
		UpdatedAt: payload.UpdatedAt,
	})
	switch {
	case errors.Is(err, errx.ErrorProfileDeleted):
		logger.Debug("received profile updated event for deleted profile")
		return nil
	case err != nil:
		logger.WithError(err).Error("failed to update profile")
		return err
	default:
		logger.Debug("profile updated successfully")
		return nil
	}
}

func (c *ProfileController) Deleted(
	ctx context.Context,
	event eventbox.InboxEvent,
) error {
	var payload evtypes.ProfileDeletedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return err
	}

	logger := c.log.WithOperation(operationProfileDeleted).With("account_id", payload.AccountID)

	err := c.core.Delete(ctx, payload.AccountID)
	switch {
	case errors.Is(err, errx.ErrorProfileDeleted):
		logger.Debug("received profile deleted event for already deleted profile")
		return nil
	case err != nil:
		logger.WithError(err).Error("failed to delete profile")
		return err
	default:
		logger.Debug("profile deleted successfully")
		return nil
	}
}
