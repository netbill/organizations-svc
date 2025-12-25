package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
	"github.com/umisto/cities-svc/internal/domain/entity"
	"github.com/umisto/cities-svc/internal/domain/modules/city"
	"github.com/umisto/cities-svc/internal/repository/models"
	"github.com/umisto/cities-svc/internal/repository/pgdb"
	"github.com/umisto/nilx"
	"github.com/umisto/pagi"
)

func (s Service) CreateCity(ctx context.Context, params city.CreateParams) error {
	err := s.sql(ctx).CreateCity(ctx, pgdb.CreateCityParams{
		AgglomerationID: nilx.UUID(params.AgglomerationID),
		Slug:            nilx.String(params.Slug),
		Name:            params.Name,
		Icon:            nilx.String(params.Icon),
		Banner:          nilx.String(params.Banner),
		PointLat:        params.Point[1],
		PointLng:        params.Point[0],
	})
	if err != nil {
		return err
	}

	return nil
}

func (s Service) GetCityByID(ctx context.Context, ID uuid.UUID) (entity.City, error) {
	row, err := s.sql(ctx).GetCityByID(ctx, ID)
	if err != nil {
		return entity.City{}, err
	}

	return models.GetCityByID(row), nil
}

func (s Service) GetCityBySlug(ctx context.Context, slug string) (entity.City, error) {
	row, err := s.sql(ctx).GetCityBySlug(ctx, slug)
	if err != nil {
		return entity.City{}, err
	}

	return models.GetCityBySlug(row), nil
}

func (s Service) UpdateCity(ctx context.Context, ID uuid.UUID, params city.UpdateParams) error {
	stmt := pgdb.UpdateCityParams{
		ID:              ID,
		AgglomerationID: nilx.UUID(params.AgglomerationID),
		Name:            nilx.String(params.Name),
		Slug:            nilx.String(params.Slug),
		Icon:            nilx.String(params.Icon),
		Banner:          nilx.String(params.Banner),
	}
	if params.Point != nil {
		stmt.PointLat = nilx.Float64(&params.Point[1])
		stmt.PointLng = nilx.Float64(&params.Point[0])
	}

	return s.sql(ctx).UpdateCity(ctx, stmt)
}

func (s Service) UpdateCityStatus(ctx context.Context, ID uuid.UUID, status string) error {
	return s.sql(ctx).UpdateCityStatus(ctx, pgdb.UpdateCityStatusParams{
		ID:     ID,
		Status: pgdb.CitiesStatus(status),
	})
}

func (s Service) FilterCities(
	ctx context.Context,
	filter city.FilterParams,
	pagination pagi.Params,
) (pagi.Page[[]entity.City], error) {
	params := pgdb.FilterCitiesParams{
		AgglomerationID: nilx.UUID(filter.AgglomerationID),
		NameLike:        nilx.String(filter.NameLike),
	}

	if filter.Status != nil {
		params.Status = pgdb.NullCitiesStatus{
			Valid:        true,
			CitiesStatus: pgdb.CitiesStatus(*filter.Status),
		}
	}

	if pagination.Cursor != nil {
		createdAtStr, ok := pagination.Cursor["created_at"]
		if !ok || createdAtStr == "" {
			return pagi.Page[[]entity.City]{}, fmt.Errorf("cursor missing created_at")
		}

		idStr, ok := pagination.Cursor["id"]
		if !ok || idStr == "" {
			return pagi.Page[[]entity.City]{}, fmt.Errorf("cursor missing id")
		}

		afterT, err := time.Parse(time.RFC3339Nano, createdAtStr)
		if err != nil {
			return pagi.Page[[]entity.City]{}, err
		}

		afterID, err := uuid.Parse(idStr)
		if err != nil {
			return pagi.Page[[]entity.City]{}, err
		}

		params.AfterCreatedAt = sql.NullTime{
			Time:  afterT,
			Valid: true,
		}
		params.AfterID = uuid.NullUUID{
			UUID:  afterID,
			Valid: true,
		}
	}

	limit := pagi.CalculateLimit(pagination.Limit, 20, 100)
	params.Limit = int32(limit)

	rows, err := s.sql(ctx).FilterCities(ctx, params)
	if err != nil {
		return pagi.Page[[]entity.City]{}, err
	}

	count, err := s.sql(ctx).CountCities(ctx, pgdb.CountCitiesParams{
		AgglomerationID: params.AgglomerationID,
		Status:          params.Status,
		NameLike:        params.NameLike,
	})
	if err != nil {
		return pagi.Page[[]entity.City]{}, err
	}

	items := models.FilterCitiesRow(rows)

	var nextCursor map[string]string
	if len(items) == limit {
		lastItem := items[len(items)-1]
		nextCursor = map[string]string{
			"created_at": lastItem.CreatedAt.Format(time.RFC3339Nano),
			"id":         lastItem.ID.String(),
		}
	}

	return pagi.Page[[]entity.City]{
		Data:       items,
		NextCursor: nextCursor,
		Total:      int(count),
	}, nil
}

func (s Service) FilterCitiesNearest(
	ctx context.Context,
	point orb.Point,
	filter city.FilterParams,
	pagination pagi.Params,
) (pagi.Page[map[int64]entity.City], error) {
	params := pgdb.FilterCitiesNearestParams{
		AgglomerationID: nilx.UUID(filter.AgglomerationID),
		NameLike:        nilx.String(filter.NameLike),
		UserLat:         point[1],
		UserLng:         point[0],
	}

	if filter.Status != nil {
		params.Status = pgdb.NullCitiesStatus{
			Valid:        true,
			CitiesStatus: pgdb.CitiesStatus(*filter.Status),
		}
	}

	if pagination.Cursor != nil {
		distStr, ok := pagination.Cursor["distance"]
		if !ok || distStr == "" {
			return pagi.Page[map[int64]entity.City]{}, fmt.Errorf("cursor missing distance")
		}

		idStr, ok := pagination.Cursor["id"]
		if !ok || idStr == "" {
			return pagi.Page[map[int64]entity.City]{}, fmt.Errorf("cursor missing after_id")
		}

		dist, err := strconv.ParseInt(distStr, 10, 64)
		if err != nil {
			return pagi.Page[map[int64]entity.City]{}, err
		}

		afterID, err := uuid.Parse(idStr)
		if err != nil {
			return pagi.Page[map[int64]entity.City]{}, err
		}

		params.AfterDistanceM = sql.NullInt64{
			Int64: dist,
			Valid: true,
		}
		params.AfterID = uuid.NullUUID{
			UUID:  afterID,
			Valid: true,
		}
	}

	limit := pagi.CalculateLimit(pagination.Limit, 20, 100)
	params.Limit = int32(limit)

	rows, err := s.sql(ctx).FilterCitiesNearest(ctx, params)
	if err != nil {
		return pagi.Page[map[int64]entity.City]{}, err
	}

	count, err := s.sql(ctx).CountCitiesNearest(ctx, pgdb.CountCitiesNearestParams{
		AgglomerationID: params.AgglomerationID,
		Status:          params.Status,
		NameLike:        params.NameLike,
	})
	if err != nil {
		return pagi.Page[map[int64]entity.City]{}, err
	}

	items := models.FilterCitiesNearestRow(rows)

	var nextCursor map[string]string
	if len(items) == limit {
		lastItem := rows[len(rows)-1]
		nextCursor = map[string]string{
			"created_at": lastItem.CreatedAt.Format(time.RFC3339Nano),
			"id":         lastItem.ID.String(),
		}
	}

	return pagi.Page[map[int64]entity.City]{
		Data:       items,
		NextCursor: nextCursor,
		Total:      int(count),
	}, nil
}
