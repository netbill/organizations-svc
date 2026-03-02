package place

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
)

type CreateParams struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	CreatedAt      time.Time `json:"created_at"`
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
