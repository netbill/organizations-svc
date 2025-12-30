package city

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/errx"
	"github.com/umisto/cities-svc/internal/domain/models"
)

func (s Service) ActivateCity(ctx context.Context, ID uuid.UUID) (city models.City, err error) {
	if err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		city, err = s.repo.UpdateCityStatus(ctx, ID, models.CityStatusActive)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to activate city: %w", err),
			)
		}

		if err = s.messanger.WriteCityActivated(ctx, city); err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to send activate city message: %w", err),
			)
		}

		return nil
	}); err != nil {
		return models.City{}, err
	}

	return city, nil
}

func (s Service) ActivateCityByUser(
	ctx context.Context,
	accountID, cityID uuid.UUID,
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

	return s.ActivateCity(ctx, cityID)
}
