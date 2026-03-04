package core

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
)

type placeRepository interface {
	CreatePlace(ctx context.Context, params PlaceCreateParams) error

	GetPlaceByID(ctx context.Context, id uuid.UUID) (models.Place, error)
	GetPlacesByIDs(ctx context.Context, ids []uuid.UUID) ([]models.Place, error)
	GetPlaceExistsForOrganization(ctx context.Context, organizationID uuid.UUID) (bool, error)
	PlaceExists(ctx context.Context, id uuid.UUID) (bool, error)

	DeletePlaceByID(ctx context.Context, id uuid.UUID) error
}

type placeTombstoneRepository interface {
	BuryPlace(ctx context.Context, placeID uuid.UUID) error
	PlaceIsBuried(ctx context.Context, placeID uuid.UUID) (bool, error)
}

type PlaceModule struct {
	placeRepo     placeRepository
	tombstoneRepo placeTombstoneRepository
	tx            transactor
}

type PlaceCreateParams struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	CreatedAt      time.Time `json:"created_at"`
}

func (m *PlaceModule) Create(
	ctx context.Context,
	params PlaceCreateParams,
) error {
	exists, err := m.placeRepo.PlaceExists(ctx, params.ID)
	if err != nil {
		return err
	}
	if exists {
		return errx.ErrorPlaceAlreadyExists.Raise(
			errors.New("place with given id already exists"),
		)
	}

	buried, err := m.tombstoneRepo.PlaceIsBuried(ctx, params.ID)
	if err != nil {
		return err
	}
	if buried {
		return errx.ErrorPlaceDeleted.Raise(
			errors.New("place with given id is deleted"),
		)
	}

	return m.placeRepo.CreatePlace(ctx, params)
}

func (m *PlaceModule) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	buried, err := m.tombstoneRepo.PlaceIsBuried(ctx, id)
	if err != nil {
		return err
	}
	if buried {
		return errx.ErrorPlaceDeleted.Raise(
			fmt.Errorf("place with id %s is already deleted", id),
		)
	}

	return m.tx.Transaction(ctx, func(ctx context.Context) error {
		if err = m.tombstoneRepo.BuryPlace(ctx, id); err != nil {
			return err
		}

		return m.placeRepo.DeletePlaceByID(ctx, id)
	})
}
