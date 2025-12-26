package city

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
	"github.com/umisto/cities-svc/internal/domain/entity"
	"github.com/umisto/cities-svc/internal/domain/errx"
	"github.com/umisto/pagi"
)

type FilterParams struct {
	AgglomerationID *uuid.UUID
	Status          *string
	Name            *string
}

func (s Service) FilterCities(
	ctx context.Context,
	params FilterParams,
	pagination pagi.Params,
) (pagi.Page[[]entity.City], error) {
	res, err := s.repo.FilterCities(ctx, params, pagination)
	if err != nil {
		return pagi.Page[[]entity.City]{}, errx.ErrorInternal.Raise(
			fmt.Errorf("filter cities: %w", err),
		)
	}

	return res, nil
}

func (s Service) FilterCitiesNearest(
	ctx context.Context,
	point orb.Point,
	filter FilterParams,
	pagination pagi.Params,
) (pagi.Page[map[int64]entity.City], error) {
	res, err := s.repo.FilterCitiesNearest(ctx, point, filter, pagination)
	if err != nil {
		return pagi.Page[map[int64]entity.City]{}, errx.ErrorInternal.Raise(
			fmt.Errorf("filter cities by distance: %w", err),
		)
	}

	return res, nil
}
