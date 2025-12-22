package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/entity"
	"github.com/umisto/cities-svc/internal/domain/modules/city"
	"github.com/umisto/cities-svc/internal/repository/pgdb"
	"github.com/umisto/pagi"
)

func (s Service) CreateCity(ctx context.Context, params city.CreateParams) (entity.City, error) {
	row, err := s.sql.CreateCity(ctx, pgdb.CreateCityParams{
		AgglomerationID: nullUUID(params.AgglomerationID),
		Slug:            nullString(params.Slug),
		Name:            params.Name,
		Icon:            nullString(params.Icon),
		Banner:          nullString(params.Banner),
	})
	if err != nil {
		return entity.City{}, err
	}

	return row.ToEntity(), nil
}

func (s Service) GetCityByID(ctx context.Context, ID uuid.UUID) (entity.City, error) {
	row, err := s.sql.GetCityByID(ctx, ID)
	if err != nil {
		return entity.City{}, err
	}

	return row.ToEntity(), nil
}

func (s Service) GetCityBySlug(ctx context.Context, slug string) (entity.City, error) {
	row, err := s.sql.GetCityBySlug(ctx, slug)
	if err != nil {
		return entity.City{}, err
	}

	return row.ToEntity(), nil
}

func (s Service) UpdateCity(ctx context.Context, ID uuid.UUID, params city.UpdateParams) (entity.City, error) {
	row, err := s.sql.UpdateCity(ctx, pgdb.UpdateCityParams{
		ID:              ID,
		AgglomerationID: nullUUID(params.AgglomerationID),
		Name:            nullString(params.Name),
		Slug:            nullString(params.Slug),
		Icon:            nullString(params.Icon),
		Banner:          nullString(params.Banner),
	})
	if err != nil {
		return entity.City{}, err
	}

	return row.ToEntity(), nil
}

func (s Service) ActivateCity(ctx context.Context, ID uuid.UUID) (entity.City, error) {
	res, err := s.sql.ActivateCity(ctx, ID)
	if err != nil {
		return entity.City{}, err
	}

	return res.ToEntity(), nil
}

func (s Service) DeactivateCity(ctx context.Context, ID uuid.UUID) (entity.City, error) {
	res, err := s.sql.DeactivateCity(ctx, ID)
	if err != nil {
		return entity.City{}, err
	}

	return res.ToEntity(), nil
}

func (s Service) FilterCities(
	ctx context.Context,
	params city.FilterParams,
	pagination pagi.Params,
) (pagi.Page[entity.City], error) {
	sqlParams := pgdb.FilterCitiesParams{
		AgglomerationID: nullUUID(params.AgglomerationID),
		NameLike:        nullString(params.NameLike),
	}

	if params.Status != nil {
		sqlParams.Status = pgdb.NullCitiesStatus{
			Valid:        true,
			CitiesStatus: pgdb.CitiesStatus(*params.Status),
		}
	}

	if pagination.Cursor != nil {
		createdAtStr, ok := pagination.Cursor["after_created_at"]
		if !ok || createdAtStr == "" {
			return pagi.Page[entity.City]{}, fmt.Errorf("cursor missing after_created_at")
		}

		idStr, ok := pagination.Cursor["after_id"]
		if !ok || idStr == "" {
			return pagi.Page[entity.City]{}, fmt.Errorf("cursor missing after_id")
		}

		afterT, err := time.Parse(time.RFC3339Nano, createdAtStr)
		if err != nil {
			return pagi.Page[entity.City]{}, err
		}

		afterID, err := uuid.Parse(idStr)
		if err != nil {
			return pagi.Page[entity.City]{}, err
		}

		sqlParams.AfterCreatedAt = sql.NullTime{
			Time:  afterT,
			Valid: true,
		}
		sqlParams.AfterID = uuid.NullUUID{
			UUID:  afterID,
			Valid: true,
		}
	}

	limit := pagination.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	sqlParams.Limit = int32(limit)

	rows, err := s.sql.FilterCities(ctx, sqlParams)
	if err != nil {
		return pagi.Page[entity.City]{}, err
	}

	count, err := s.sql.CountCities(ctx, pgdb.CountCitiesParams{
		AgglomerationID: sqlParams.AgglomerationID,
		Status:          sqlParams.Status,
		NameLike:        sqlParams.NameLike,
	})
	if err != nil {
		return pagi.Page[entity.City]{}, err
	}

	var items []entity.City
	for _, r := range rows {
		items = append(items, r.ToEntity())
	}

	var nextCursor map[string]string
	if len(items) == limit {
		lastItem := items[len(items)-1]
		nextCursor = map[string]string{
			"after_created_at": lastItem.CreatedAt.Format(time.RFC3339Nano),
			"after_id":         lastItem.ID.String(),
		}
	}

	return pagi.Page[entity.City]{
		Data:       items,
		NextCursor: nextCursor,
		Total:      int(count),
	}, nil
}
