package controller

import (
	"context"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
	"github.com/umisto/cities-svc/internal/domain/models"
	"github.com/umisto/cities-svc/internal/domain/modules/city"
	"github.com/umisto/logium"
	"github.com/umisto/pagi"
)

type City interface {
	CreateCity(ctx context.Context, params city.CreateParams) (models.City, error)

	GetCityByID(ctx context.Context, id uuid.UUID) (models.City, error)
	GetCityBySlug(ctx context.Context, slug string) (models.City, error)
	FilterCities(
		ctx context.Context,
		params city.FilterParams,
		offset, limit uint,
	) (pagi.Page[[]models.City], error)
	FilterCitiesNearest(
		ctx context.Context,
		filter city.FilterParams,
		point orb.Point,
		offset, limit uint,
	) (pagi.Page[map[float64]models.City], error)

	UpdateCity(ctx context.Context, id uuid.UUID, params city.UpdateParams) (models.City, error)
	UpdateCityByUser(
		ctx context.Context,
		accountID, cityID uuid.UUID,
		params city.UpdateParams,
	) (models.City, error)

	UpdateCitySlug(
		ctx context.Context,
		id uuid.UUID,
		newSlug *string,
	) (city models.City, err error)
	UpdateCitySlugByUser(
		ctx context.Context,
		accountID, cityID uuid.UUID,
		newSlug *string,
	) (models.City, error)

	UpdateCityAgglomeration(
		ctx context.Context,
		id uuid.UUID,
		newAggloID *uuid.UUID,
	) (city models.City, err error)
}

type CityController struct {
	domain City
	log    logium.Logger
}
