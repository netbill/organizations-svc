package city

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
	"github.com/umisto/cities-svc/internal/domain/errx"
	"github.com/umisto/cities-svc/internal/domain/models"
)

type UpdateParams struct {
	Name   *string
	Icon   *string
	Banner *string
	Point  *orb.Point
}

func (s Service) UpdateCity(ctx context.Context, id uuid.UUID, params UpdateParams) (city models.City, err error) {
	if err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		city, err = s.repo.UpdateCity(ctx, id, params)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to update city: %w", err),
			)
		}

		err = s.messanger.WriteUpdateCity(ctx, city)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to send update city message: %w", err),
			)
		}

		return nil
	}); err != nil {
		return models.City{}, err
	}

	return city, nil
}

func (s Service) UpdateCityByUser(
	ctx context.Context,
	accountID, cityID uuid.UUID,
	params UpdateParams,
) (models.City, error) {
	city, err := s.GetCityByID(ctx, cityID)
	if err != nil {
		return models.City{}, err
	}

	if city.AgglomerationID == nil {
		return models.City{}, errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("city %s has no agglomeration", city.ID),
		)
	}

	_, err = s.checkAgglomerationIsActiveAndExists(ctx, *city.AgglomerationID)
	if err != nil {
		return models.City{}, err
	}

	err = s.checkPermissionForManageCity(ctx, accountID, *city.AgglomerationID)
	if err != nil {
		return models.City{}, err
	}

	return s.UpdateCity(ctx, cityID, params)
}

func (s Service) UpdateCitySlug(
	ctx context.Context,
	id uuid.UUID,
	newSlug *string,
) (city models.City, err error) {
	city, err = s.GetCityByID(ctx, id)
	if err != nil {
		return models.City{}, err
	}

	if newSlug != nil {
		err = s.checkSlugIsAvailable(ctx, *newSlug)
		if err != nil {
			return models.City{}, err
		}
	}

	oldSlug := city.Slug

	if err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		city, err = s.repo.UpdateCitySlug(ctx, id, newSlug)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to update city slug: %w", err),
			)
		}

		err = s.messanger.WriteUpdateCitySlug(ctx, city.ID, oldSlug, newSlug)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to send update city slug message: %w", err),
			)
		}

		return nil
	}); err != nil {
		return models.City{}, err
	}

	return city, nil
}

func (s Service) UpdateCitySlugByUser(
	ctx context.Context,
	accountID, cityID uuid.UUID,
	newSlug *string,
) (models.City, error) {
	city, err := s.GetCityByID(ctx, cityID)
	if err != nil {
		return models.City{}, err
	}

	if city.AgglomerationID == nil {
		return models.City{}, errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("city %s has no agglomeration", city.ID),
		)
	}

	_, err = s.checkAgglomerationIsActiveAndExists(ctx, *city.AgglomerationID)
	if err != nil {
		return models.City{}, err
	}

	err = s.checkPermissionForManageCity(ctx, accountID, *city.AgglomerationID)
	if err != nil {
		return models.City{}, err
	}

	return s.UpdateCitySlug(ctx, cityID, newSlug)
}

func (s Service) UpdateCityAgglomeration(
	ctx context.Context,
	id uuid.UUID,
	newAggloID *uuid.UUID,
) (city models.City, err error) {
	city, err = s.GetCityByID(ctx, id)
	if err != nil {
		return models.City{}, err
	}

	oldAggloID := city.AgglomerationID

	if err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		city, err = s.repo.UpdateCityAgglomeration(ctx, id, newAggloID)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to update city agglomeration: %w", err),
			)
		}

		err = s.messanger.WriteUpdateCityAgglomeration(
			ctx,
			city.ID,
			oldAggloID,
			newAggloID,
		)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to send update city agglomeration message: %w", err),
			)
		}

		return nil
	}); err != nil {
		return models.City{}, err
	}

	return city, nil
}

func (s Service) UpdateCityAgglomerationByUser(
	ctx context.Context,
	accountID, cityID uuid.UUID,
	newAggloID *uuid.UUID,
) (models.City, error) {
	city, err := s.GetCityByID(ctx, cityID)
	if err != nil {
		return models.City{}, err
	}

	if city.AgglomerationID == nil {
		return models.City{}, errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("city %s has no agglomeration", city.ID),
		)
	}

	_, err = s.checkAgglomerationIsActiveAndExists(ctx, *city.AgglomerationID)
	if err != nil {
		return models.City{}, err
	}

	err = s.checkPermissionForManageCity(ctx, accountID, *city.AgglomerationID)
	if err != nil {
		return models.City{}, err
	}

	return s.UpdateCityAgglomeration(ctx, cityID, newAggloID)
}
