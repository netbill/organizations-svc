package city

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/entity"
	"github.com/umisto/cities-svc/internal/domain/errx"
)

func (s Service) ArchiveCity(ctx context.Context, ID uuid.UUID) (city entity.City, err error) {
	if err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		city, err = s.repo.UpdateCityStatus(ctx, ID, entity.CityStatusArchived)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to archived city: %w", err),
			)
		}

		if err = s.messanger.ArchivedCity(ctx, city); err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to send archived city message: %w", err),
			)
		}

		return nil
	}); err != nil {
		return entity.City{}, err
	}

	return city, nil
}
