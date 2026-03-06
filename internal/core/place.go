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

type placeRepo interface {
	Get(ctx context.Context, id uuid.UUID) (models.Place, error)
	GetListByIDs(ctx context.Context, ids []uuid.UUID) ([]models.Place, error)
	ExistsForOrg(ctx context.Context, organizationID uuid.UUID) (bool, error)
	Exists(ctx context.Context, id uuid.UUID) (bool, error)

	Create(ctx context.Context, params PlaceCreateParams) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type placeTombstone interface {
	BuryPlace(ctx context.Context, placeID uuid.UUID) error
	PlaceIsBuried(ctx context.Context, placeID uuid.UUID) (bool, error)
}

type PlaceModule struct {
	repo      placeRepo
	tombstone placeTombstone
	tx        transactor
}

type PlaceDeps struct {
	Repo      placeRepo
	Tombstone placeTombstone
	Tx        transactor
}

func NewPlaceModule(deps PlaceDeps) *PlaceModule {
	return &PlaceModule{
		repo:      deps.Repo,
		tombstone: deps.Tombstone,
		tx:        deps.Tx,
	}
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
	exists, err := m.repo.Exists(ctx, params.ID)
	if err != nil {
		return err
	}
	if exists {
		return errx.ErrorPlaceAlreadyExists.Raise(
			errors.New("place with given id already exists"),
		)
	}

	buried, err := m.tombstone.PlaceIsBuried(ctx, params.ID)
	if err != nil {
		return err
	}
	if buried {
		return errx.ErrorPlaceDeleted.Raise(
			errors.New("place with given id is deleted"),
		)
	}

	return m.repo.Create(ctx, params)
}

func (m *PlaceModule) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	buried, err := m.tombstone.PlaceIsBuried(ctx, id)
	if err != nil {
		return err
	}
	if buried {
		return errx.ErrorPlaceDeleted.Raise(
			fmt.Errorf("place with id %s is already deleted", id),
		)
	}

	return m.tx.Transaction(ctx, func(ctx context.Context) error {
		if err := m.tombstone.BuryPlace(ctx, id); err != nil {
			return err
		}

		return m.repo.Delete(ctx, id)
	})
}
