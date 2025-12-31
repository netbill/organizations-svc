package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/umisto/cities-svc/internal/domain/models"
	"github.com/umisto/cities-svc/internal/domain/modules/agglomeration"
	"github.com/umisto/cities-svc/internal/repository/pgdb"
	"github.com/umisto/pagi"
)

func (s Service) CreateAgglomeration(ctx context.Context, params agglomeration.CreateParams) (models.Agglomeration, error) {
	row, err := s.agglomerationsQ().Insert(ctx, pgdb.AgglomerationsQInsertInput{
		Name: params.Name,
		Icon: params.Icon,
	})
	if err != nil {
		return models.Agglomeration{}, err
	}

	return Agglomeration(row), nil
}

func (s Service) UpdateAgglomeration(
	ctx context.Context,
	ID uuid.UUID,
	params agglomeration.UpdateParams,
) (models.Agglomeration, error) {
	q := s.agglomerationsQ().FilterByID(ID)
	if params.Name != nil {
		q = q.UpdateName(*params.Name)
	}
	if params.Icon != nil {
		q = q.UpdateIcon(*params.Icon)
	}

	row, err := q.UpdateOne(ctx)
	if err != nil {
		return models.Agglomeration{}, err
	}

	return Agglomeration(row), nil
}

func (s Service) UpdateAgglomerationStatus(
	ctx context.Context,
	ID uuid.UUID,
	status string,
) (models.Agglomeration, error) {
	row, err := s.agglomerationsQ().FilterByID(ID).UpdateStatus(status).UpdateOne(ctx)
	if err != nil {
		return models.Agglomeration{}, err
	}

	return Agglomeration(row), nil
}

func (s Service) GetAgglomerationByID(ctx context.Context, ID uuid.UUID) (models.Agglomeration, error) {
	row, err := s.agglomerationsQ().FilterByID(ID).Get(ctx)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return models.Agglomeration{}, nil
	case err != nil:
		return models.Agglomeration{}, err
	}

	return Agglomeration(row), nil
}

func (s Service) DeleteAgglomeration(ctx context.Context, ID uuid.UUID) error {
	return s.agglomerationsQ().FilterByID(ID).Delete(ctx)
}

func (s Service) FilterAgglomerations(
	ctx context.Context,
	filter agglomeration.FilterParams,
	offset, limit uint,
) (pagi.Page[[]models.Agglomeration], error) {
	q := s.agglomerationsQ()
	if filter.Name != nil {
		q = q.FilterNameLike(*filter.Name)
	}
	if filter.Status != nil {
		q = q.FilterByStatus(*filter.Status)
	}

	rows, err := q.Page(limit, offset).Select(ctx)
	if err != nil {
		return pagi.Page[[]models.Agglomeration]{}, err
	}

	total, err := q.Count(ctx)
	if err != nil {
		return pagi.Page[[]models.Agglomeration]{}, err
	}

	agglomerations := make([]models.Agglomeration, len(rows))
	for i, row := range rows {
		agglomerations[i] = Agglomeration(row)
	}

	return pagi.Page[[]models.Agglomeration]{
		Data:  agglomerations,
		Page:  uint(offset/limit) + 1,
		Size:  uint(len(agglomerations)),
		Total: total,
	}, nil

}

func Agglomeration(db pgdb.Agglomeration) models.Agglomeration {
	return models.Agglomeration{
		ID:        db.ID,
		Status:    db.Status,
		Name:      db.Name,
		Icon:      db.Icon,
		CreatedAt: db.CreatedAt,
		UpdatedAt: db.UpdatedAt,
	}
}
