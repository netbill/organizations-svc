package place

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/paulmach/orb"
)

type CreateParams struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	ClassID        uuid.UUID `json:"class_id"`

	Status   string    `json:"status"`
	Verified bool      `json:"verified"`
	Point    orb.Point `json:"point"`
	Address  string    `json:"address"`
	Name     string    `json:"name"`

	IconKey     *string `json:"icon_key,omitempty"`
	BannerKey   *string `json:"banner_key,omitempty"`
	Description *string `json:"description,omitempty"`
	Website     *string `json:"website,omitempty"`
	Phone       *string `json:"phone,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}

func (m *Module) Create(
	ctx context.Context,
	params CreateParams,
) error {
	exists, err := m.repo.PlaceExists(ctx, params.ID)
	if err != nil {
		return err
	}
	if exists {
		return errx.ErrorPlaceAlreadyExists.Raise(
			errors.New("place with given id already exists"),
		)
	}

	buried, err := m.repo.PlaceIsBuried(ctx, params.ID)
	if err != nil {
		return err
	}
	if buried {
		return errx.ErrorPlaceDeleted.Raise(
			errors.New("place with given id is deleted"),
		)
	}

	return m.repo.CreatePlace(ctx, params)
}

func (m *Module) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (place models.Place, err error) {
	return m.repo.GetPlaceByID(ctx, id)
}

func (m *Module) GetByIDs(
	ctx context.Context,
	ids []uuid.UUID,
) (places []models.Place, err error) {
	return m.repo.GetPlacesByIDs(ctx, ids)
}

type UpdateParams struct {
	ClassID  uuid.UUID `json:"class_id"`
	Name     string    `json:"name"`
	Address  string    `json:"address"`
	Status   string    `json:"status"`
	Verified bool      `json:"verified"`

	Description *string `json:"description,omitempty"`
	Website     *string `json:"website,omitempty"`
	Phone       *string `json:"phone,omitempty"`

	IconKey   *string `json:"icon_key,omitempty"`
	BannerKey *string `json:"banner_key,omitempty"`

	Version   int32     `json:"version"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (m *Module) Update(
	ctx context.Context,
	id uuid.UUID,
	params UpdateParams,
) error {
	current, err := m.repo.GetPlaceByID(ctx, id)
	if err != nil {
		return err
	}
	if current.Version >= params.Version {
		return nil
	}

	return m.repo.UpdatePlaceByID(ctx, id, params)
}

func (m *Module) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	buried, err := m.repo.PlaceIsBuried(ctx, id)
	if err != nil {
		return err
	}
	if buried {
		return errx.ErrorPlaceDeleted.Raise(
			fmt.Errorf("place with id %s is already deleted", id),
		)
	}

	return m.repo.Transaction(ctx, func(ctx context.Context) error {
		if err = m.repo.BuryPlace(ctx, id); err != nil {
			return err
		}

		return m.repo.DeletePlaceByID(ctx, id)
	})
}
