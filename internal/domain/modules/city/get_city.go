package city

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/errx"
	"github.com/umisto/cities-svc/internal/domain/models"
)

func (s Service) GetCityByID(ctx context.Context, id uuid.UUID) (models.City, error) {
	res, err := s.repo.GetCityByID(ctx, id)
	if err != nil {
		return models.City{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get city by id: %w", err),
		)
	}
	if res.IsNil() {
		return models.City{}, errx.ErrorCityNotFound.Raise(
			fmt.Errorf("city with id %s not found", id),
		)
	}

	return res, nil
}

func (s Service) GetCityBySlug(ctx context.Context, slug string) (models.City, error) {
	res, err := s.repo.GetCityBySlug(ctx, slug)
	if err != nil {
		return models.City{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get city by slug: %w", err),
		)
	}
	if res.IsNil() {
		return models.City{}, errx.ErrorCityNotFound.Raise(
			fmt.Errorf("city with slug %s not found", slug),
		)
	}

	return res, nil
}
