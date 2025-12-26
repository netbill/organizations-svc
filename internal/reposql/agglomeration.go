package repol

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/umisto/cities-svc/internal/domain/entity"
	"github.com/umisto/cities-svc/internal/domain/modules/agglomeration"
	"github.com/umisto/cities-svc/internal/reposql/pgdbsql"

	"github.com/umisto/cities-svc/internal/repository/models"
	"github.com/umisto/nilx"
	"github.com/umisto/pagi"
)

func (s Service) CreateAgglomeration(ctx context.Context, name string) (entity.Agglomeration, error) {
	row, err := s.sql(ctx).CreateAgglomeration(ctx, pgdbsql.CreateAgglomerationParams{
		Name: name,
	})
	if err != nil {
		return entity.Agglomeration{}, err
	}

	return models.AgglomerationRow(row), nil
}

func (s Service) UpdateAgglomeration(
	ctx context.Context,
	ID uuid.UUID,
	params agglomeration.UpdateParams,
) (entity.Agglomeration, error) {
	row, err := s.sql(ctx).UpdateAgglomeration(ctx, pgdbsql.UpdateAgglomerationParams{
		ID:   ID,
		Name: nilx.String(params.Name),
		Icon: nilx.String(params.Icon),
	})
	if err != nil {
		return entity.Agglomeration{}, err
	}

	return models.AgglomerationRow(row), nil
}

func (s Service) UpdateAgglomerationStatus(ctx context.Context, ID uuid.UUID, status string) (entity.Agglomeration, error) {
	res, err := s.sql(ctx).UpdateAgglomerationStatus(ctx, pgdbsql.UpdateAgglomerationStatusParams{
		ID:     ID,
		Status: pgdbsql.AdministrationStatus(status),
	})
	if err != nil {
		return entity.Agglomeration{}, err
	}

	return models.AgglomerationRow(res), nil
}

func (s Service) GetAgglomerationByID(ctx context.Context, ID uuid.UUID) (entity.Agglomeration, error) {
	row, err := s.sql(ctx).GetAgglomerationByID(ctx, ID)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return entity.Agglomeration{}, nil
	case err != nil:
		return entity.Agglomeration{}, err
	}

	return models.AgglomerationRow(row), nil
}

func (s Service) DeleteAgglomeration(ctx context.Context, ID uuid.UUID) error {
	return s.sql(ctx).DeleteAgglomeration(ctx, ID)
}

func (s Service) FilterAgglomerations(
	ctx context.Context,
	filter agglomeration.FilterParams,
	pagination pagi.Params,
) (pagi.Page[[]entity.Agglomeration], error) {
	params := pgdbsql.FilterAgglomerationsParams{
		NameLike: nilx.String(filter.Name),
	}

	if filter.Status != nil {
		params.Status = pgdbsql.NullAdministrationStatus{
			AdministrationStatus: pgdbsql.AdministrationStatus(*filter.Status),
			Valid:                true,
		}
	}

	if pagination.Cursor != nil {
		createdAtStr, ok := pagination.Cursor["created_at"]
		if !ok || createdAtStr == "" {
			return pagi.Page[[]entity.Agglomeration]{}, fmt.Errorf("cursor missing created_at")
		}

		idStr, ok := pagination.Cursor["id"]
		if !ok || idStr == "" {
			return pagi.Page[[]entity.Agglomeration]{}, fmt.Errorf("cursor missing id")
		}

		afterT, err := time.Parse(time.RFC3339Nano, createdAtStr)
		if err != nil {
			return pagi.Page[[]entity.Agglomeration]{}, err
		}

		afterID, err := uuid.Parse(idStr)
		if err != nil {
			return pagi.Page[[]entity.Agglomeration]{}, err
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

	rows, err := s.sql(ctx).FilterAgglomerations(ctx, params)
	if err != nil {
		return pagi.Page[[]entity.Agglomeration]{}, err
	}

	count, err := s.sql(ctx).CountAgglomerations(ctx, pgdbsql.CountAgglomerationsParams{
		Status:   params.Status,
		NameLike: params.NameLike,
	})
	if err != nil {
		return pagi.Page[[]entity.Agglomeration]{}, err
	}

	collection := make([]entity.Agglomeration, 0, len(rows))
	for _, row := range rows {
		collection = append(collection, models.AgglomerationRow(row))
	}

	var nextCursor map[string]string
	if len(rows) == limit {
		last := rows[len(rows)-1]
		nextCursor = map[string]string{
			"created_at": last.CreatedAt.UTC().Format(time.RFC3339Nano),
			"id":         last.ID.String(),
		}
	}

	return pagi.Page[[]entity.Agglomeration]{
		Data:       collection,
		NextCursor: nextCursor,
		Total:      int(count),
	}, nil
}
