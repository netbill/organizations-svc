package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/umisto/cities-svc/internal/domain/entity"
	"github.com/umisto/cities-svc/internal/domain/modules/agglomeration"
	"github.com/umisto/cities-svc/internal/repository/pgdb"
	"github.com/umisto/pagi"
)

func (s Service) CreateAgglomeration(ctx context.Context, name string) (entity.Agglomeration, error) {
	row, err := s.sql.CreateAgglomeration(ctx, pgdb.CreateAgglomerationParams{
		Name: name,
	})
	if err != nil {
		return entity.Agglomeration{}, err
	}

	return row.ToEntity(), nil
}

func (s Service) UpdateAgglomeration(
	ctx context.Context,
	ID uuid.UUID,
	params agglomeration.UpdateParams,
) (entity.Agglomeration, error) {
	row, err := s.sql.UpdateAgglomeration(ctx, pgdb.UpdateAgglomerationParams{
		ID:   ID,
		Name: nullString(params.Name),
		Icon: nullString(params.Icon),
	})
	if err != nil {
		return entity.Agglomeration{}, err
	}

	return row.ToEntity(), nil
}

func (s Service) ActivateAgglomeration(ctx context.Context, ID uuid.UUID) (entity.Agglomeration, error) {
	res, err := s.sql.ActivateAgglomeration(ctx, ID)
	if err != nil {
		return entity.Agglomeration{}, err
	}

	return res.ToEntity(), nil
}

func (s Service) DeactivateAgglomeration(ctx context.Context, ID uuid.UUID) (entity.Agglomeration, error) {
	res, err := s.sql.DeactivateAgglomeration(ctx, ID)
	if err != nil {
		return entity.Agglomeration{}, err
	}

	return res.ToEntity(), nil
}

func (s Service) GetAgglomerationByID(ctx context.Context, ID uuid.UUID) (entity.Agglomeration, error) {
	row, err := s.sql.GetAgglomerationByID(ctx, ID)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return entity.Agglomeration{}, nil
	case err != nil:
		return entity.Agglomeration{}, err
	}

	return row.ToEntity(), nil
}

func (s Service) DeleteAgglomeration(ctx context.Context, ID uuid.UUID) error {
	return s.sql.DeleteAgglomeration(ctx, ID)
}

func (s Service) FilterAgglomerations(
	ctx context.Context,
	params agglomeration.FilterParams,
	pagination pagi.Params,
) (pagi.Page[entity.Agglomeration], error) {
	sqlParam := pgdb.FilterAgglomerationsParams{
		NameLike: nullString(params.NameLike),
	}

	if params.Status != nil {
		sqlParam.Status = pgdb.NullAdministrationStatus{
			AdministrationStatus: pgdb.AdministrationStatus(*params.Status),
			Valid:                true,
		}
	}

	if pagination.Cursor != nil {
		createdAtStr, ok := pagination.Cursor["after_created_at"]
		if !ok || createdAtStr == "" {
			return pagi.Page[entity.Agglomeration]{}, fmt.Errorf("cursor missing after_created_at")
		}

		idStr, ok := pagination.Cursor["after_id"]
		if !ok || idStr == "" {
			return pagi.Page[entity.Agglomeration]{}, fmt.Errorf("cursor missing after_id")
		}

		afterT, err := time.Parse(time.RFC3339Nano, createdAtStr)
		if err != nil {
			return pagi.Page[entity.Agglomeration]{}, err
		}

		afterID, err := uuid.Parse(idStr)
		if err != nil {
			return pagi.Page[entity.Agglomeration]{}, err
		}

		sqlParam.AfterCreatedAt = sql.NullTime{
			Time:  afterT,
			Valid: true,
		}
		sqlParam.AfterID = uuid.NullUUID{
			UUID:  afterID,
			Valid: true,
		}
	}

	limit := pagination.Limit
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	sqlParam.Limit = int32(limit)

	rows, err := s.sql.FilterAgglomerations(ctx, sqlParam)
	if err != nil {
		return pagi.Page[entity.Agglomeration]{}, err
	}

	count, err := s.sql.CountAgglomerations(ctx, pgdb.CountAgglomerationsParams{
		Status:   sqlParam.Status,
		NameLike: sqlParam.NameLike,
	})
	if err != nil {
		return pagi.Page[entity.Agglomeration]{}, err
	}

	collection := make([]entity.Agglomeration, 0, len(rows))
	for _, row := range rows {
		collection = append(collection, row.ToEntity())
	}

	var nextCursor map[string]string
	if len(rows) == limit {
		last := rows[len(rows)-1]
		nextCursor = map[string]string{
			"after_created_at": last.CreatedAt.UTC().Format(time.RFC3339Nano),
			"after_id":         last.ID.String(),
		}
	}

	return pagi.Page[entity.Agglomeration]{
		Data:       collection,
		NextCursor: nextCursor,
		Total:      int(count),
	}, nil
}
