package city

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
	"github.com/umisto/cities-svc/internal/domain/entity"
	"github.com/umisto/cities-svc/internal/domain/errx"
)

type UpdateParams struct {
	AgglomerationID *uuid.UUID
	Name            *string
	Slug            *string
	Icon            *string
	Banner          *string
	Point           *orb.Point
}

func (s Service) UpdateCity(ctx context.Context, id uuid.UUID, params UpdateParams) (city entity.City, err error) {
	if err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		city, err = s.repo.UpdateCity(ctx, id, params)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to update city: %w", err),
			)
		}

		err = s.messanger.UpdateCity(ctx, city)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to send update city message: %w", err),
			)
		}

		return nil
	}); err != nil {
		return entity.City{}, err
	}

	return city, nil
}

func (s Service) UpdateCityByUser(
	ctx context.Context,
	accountID, cityID uuid.UUID,
	params UpdateParams,
) (entity.City, error) {
	city, err := s.GetCity(ctx, cityID)
	if err != nil {
		return entity.City{}, err
	}

	if city.AgglomerationID == nil {
		return entity.City{}, errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("city %s has no agglomeration", city.ID),
		)
	}

	_, err = s.checkAgglomerationIsActiveAndExists(ctx, *city.AgglomerationID)
	if err != nil {
		return entity.City{}, err
	}

	err = s.checkPermissionForManageCity(ctx, accountID, *city.AgglomerationID)
	if err != nil {
		return entity.City{}, err
	}

	return s.UpdateCity(ctx, cityID, params)
}
