package city

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/errx"
	"github.com/umisto/cities-svc/internal/domain/models"
)

func (s Service) ArchiveCity(ctx context.Context, ID uuid.UUID) (city models.City, err error) {
	if err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		city, err = s.repo.UpdateCityStatus(ctx, ID, models.CityStatusArchived)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to archived city: %w", err),
			)
		}

		if err = s.messanger.WriteArchivedCity(ctx, city); err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to send archived city message: %w", err),
			)
		}

		return nil
	}); err != nil {
		return models.City{}, err
	}

	return city, nil
}
