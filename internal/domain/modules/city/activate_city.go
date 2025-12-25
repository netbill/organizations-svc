package city

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/entity"
	"github.com/umisto/cities-svc/internal/domain/errx"
)

func (s Service) ActivateCity(ctx context.Context, ID uuid.UUID) (city entity.City, err error) {
	if err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		err = s.repo.UpdateCityStatus(ctx, ID, entity.CityStatusActive)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to activate city: %w", err),
			)
		}

		if err = s.messanger.ActivateCity(ctx, city); err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to send activate city message: %w", err),
			)
		}

		return nil
	}); err != nil {
		return entity.City{}, err
	}

	return s.GetCity(ctx, ID)
}

func (s Service) ActivateCityByUser(
	ctx context.Context,
	accountID, cityID uuid.UUID,
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

	err = s.checkPermissionByCode(ctx, accountID, *city.AgglomerationID, entity.RolePermissionManageCities)
	if err != nil {
		return entity.City{}, err
	}

	return s.ActivateCity(ctx, cityID)
}
