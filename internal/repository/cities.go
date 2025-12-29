package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
	"github.com/umisto/cities-svc/internal/domain/models"
	"github.com/umisto/cities-svc/internal/domain/modules/city"
	"github.com/umisto/cities-svc/internal/repository/pgdb"
	"github.com/umisto/pagi"
)

func (s Service) CreateCity(ctx context.Context, params city.CreateParams) (models.City, error) {
	row, err := s.citiesQ().Insert(ctx, pgdb.CityInsertParams{
		AgglomerationID: params.AgglomerationID,
		Name:            params.Name,
		Slug:            params.Slug,
		Icon:            params.Icon,
		Banner:          params.Banner,
		Point:           params.Point,
	})
	if err != nil {
		return models.City{}, err
	}

	return City(row), nil
}

func (s Service) GetCityByID(ctx context.Context, ID uuid.UUID) (models.City, error) {
	row, err := s.citiesQ().FilterByID(ID).Get(ctx)
	if err != nil {
		return models.City{}, err
	}

	return City(row), nil
}

func (s Service) GetCityBySlug(ctx context.Context, slug string) (models.City, error) {
	row, err := s.citiesQ().FilterBySlug(slug).Get(ctx)
	if err != nil {
		return models.City{}, err
	}

	return City(row), nil
}

func (s Service) UpdateCity(ctx context.Context, ID uuid.UUID, params city.UpdateParams) (models.City, error) {
	q := s.citiesQ().FilterByID(ID)
	if params.Name != nil {
		q = q.UpdateName(*params.Name)
	}
	if params.Icon != nil {
		if *params.Icon == "" {
			q = q.UpdateIcon(sql.NullString{Valid: false})
		} else {
			q = q.UpdateIcon(sql.NullString{String: *params.Icon, Valid: true})
		}
	}
	if params.Banner != nil {
		if *params.Banner == "" {
			q = q.UpdateBanner(sql.NullString{Valid: false})
		} else {
			q = q.UpdateBanner(sql.NullString{String: *params.Banner, Valid: true})
		}
	}
	if params.Point != nil {
		q = q.UpdatePoint(*params.Point)
	}

	row, err := q.UpdateOne(ctx)

	return City(row), err
}

func (s Service) UpdateCityStatus(ctx context.Context, ID uuid.UUID, status string) (models.City, error) {
	row, err := s.citiesQ().FilterByID(ID).UpdateStatus(status).UpdateOne(ctx)
	if err != nil {
		return models.City{}, err
	}

	return City(row), nil
}

func (s Service) UpdateCityAgglomeration(
	ctx context.Context,
	cityID uuid.UUID,
	agglomerationID *uuid.UUID,
) (models.City, error) {
	q := s.citiesQ().FilterByID(cityID)
	if agglomerationID == nil {
		q = q.UpdateAgglomerationID(uuid.NullUUID{Valid: false})
	} else {
		q = q.UpdateAgglomerationID(uuid.NullUUID{UUID: *agglomerationID, Valid: true})
	}

	row, err := q.UpdateOne(ctx)
	if err != nil {
		return models.City{}, err
	}

	return City(row), nil
}

func (s Service) UpdateCitySlug(ctx context.Context, ID uuid.UUID, slug *string) (models.City, error) {
	q := s.citiesQ().FilterByID(ID)
	if slug == nil {
		q = q.UpdateSlug(sql.NullString{Valid: false})
	} else {
		q = q.UpdateSlug(sql.NullString{String: *slug, Valid: true})
	}

	row, err := q.UpdateOne(ctx)
	if err != nil {
		return models.City{}, err
	}

	return City(row), nil
}

func (s Service) DeleteCity(ctx context.Context, ID uuid.UUID) error {
	return s.citiesQ().FilterByID(ID).Delete(ctx)
}

func (s Service) FilterCities(
	ctx context.Context,
	filter city.FilterParams,
	offset, limit uint,
) (pagi.Page[[]models.City], error) {
	q := s.citiesQ()
	if filter.AgglomerationID != nil {
		q = q.FilterByAgglomerationID(*filter.AgglomerationID)
	}
	if filter.Status != nil {
		q = q.FilterByStatus(*filter.Status)
	}
	if filter.Name != nil {
		q = q.FilterLikeName(*filter.Name)
	}

	rows, err := q.Page(limit, offset).Select(ctx)
	if err != nil {
		return pagi.Page[[]models.City]{}, err
	}

	total, err := q.Count(ctx)
	if err != nil {
		return pagi.Page[[]models.City]{}, err
	}

	collection := make([]models.City, 0, len(rows))
	for _, row := range rows {
		collection = append(collection, City(row))
	}

	return pagi.Page[[]models.City]{
		Data:  collection,
		Page:  uint(offset/limit) + 1,
		Size:  uint(len(collection)),
		Total: uint(total),
	}, nil
}

func (s Service) FilterCitiesNearest(
	ctx context.Context,
	filter city.FilterParams,
	point orb.Point,
	offset, limit uint,
) (pagi.Page[map[float64]models.City], error) {
	q := s.citiesQ()
	if filter.AgglomerationID != nil {
		q = q.FilterByAgglomerationID(*filter.AgglomerationID)
	}
	if filter.Status != nil {
		q = q.FilterByStatus(*filter.Status)
	}
	if filter.Name != nil {
		q = q.FilterLikeName(*filter.Name)
	}

	rows, err := q.OrderNearest(limit, point[1], point[0]).Page(limit, offset).SelectNearest(ctx)
	if err != nil {
		return pagi.Page[map[float64]models.City]{}, err
	}

	total, err := q.Count(ctx)
	if err != nil {
		return pagi.Page[map[float64]models.City]{}, err
	}

	collection := make(map[float64]models.City, len(rows))
	for _, row := range rows {
		collection[row.DistanceMeters] = CityDistance(row)
	}

	return pagi.Page[map[float64]models.City]{
		Data:  collection,
		Page:  uint(offset/limit) + 1,
		Size:  uint(len(collection)),
		Total: uint(total),
	}, nil
}

func City(c pgdb.City) models.City {
	ent := models.City{
		ID:        c.ID,
		Status:    c.Status,
		Name:      c.Name,
		Point:     c.Point,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}

	if c.AgglomerationID != nil {
		ent.AgglomerationID = c.AgglomerationID
	}
	if c.Slug != nil {
		ent.Slug = c.Slug
	}
	if c.Icon != nil {
		ent.Icon = c.Icon
	}
	if c.Banner != nil {
		ent.Banner = c.Banner
	}

	return ent
}

func CityDistance(cd pgdb.CityDistance) models.City {
	ent := models.City{
		ID:        cd.ID,
		Status:    cd.Status,
		Name:      cd.Name,
		Point:     cd.Point,
		CreatedAt: cd.CreatedAt,
		UpdatedAt: cd.UpdatedAt,
	}

	if cd.AgglomerationID != nil {
		ent.AgglomerationID = cd.AgglomerationID
	}
	if cd.Slug != nil {
		ent.Slug = cd.Slug
	}
	if cd.Icon != nil {
		ent.Icon = cd.Icon
	}
	if cd.Banner != nil {
		ent.Banner = cd.Banner
	}

	return ent
}
