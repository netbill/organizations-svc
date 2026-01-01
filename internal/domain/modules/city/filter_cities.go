package city

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
	"github.com/umisto/cities-svc/internal/domain/errx"
	"github.com/umisto/cities-svc/internal/domain/models"
	"github.com/umisto/pagi"
)

type FilterParams struct {
	AgglomerationID *uuid.UUID
	Name            *string
}

func (s Service) FilterCities(
	ctx context.Context,
	params FilterParams,
	offset, limit uint,
) (pagi.Page[[]models.City], error) {
	res, err := s.repo.FilterCities(ctx, params, offset, limit)
	if err != nil {
		return pagi.Page[[]models.City]{}, errx.ErrorInternal.Raise(
			fmt.Errorf("filter cities: %w", err),
		)
	}

	return res, nil
}

func (s Service) FilterCitiesNearest(
	ctx context.Context,
	filter FilterParams,
	point orb.Point,
	offset, limit uint,
) (pagi.Page[map[float64]models.City], error) {
	res, err := s.repo.FilterCitiesNearest(ctx, filter, point, offset, limit)
	if err != nil {
		return pagi.Page[map[float64]models.City]{}, errx.ErrorInternal.Raise(
			fmt.Errorf("filter cities by distance: %w", err),
		)
	}

	return res, nil
}
