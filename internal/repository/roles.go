package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/entity"
	"github.com/umisto/cities-svc/internal/repository/models"
	"github.com/umisto/cities-svc/internal/repository/pgdb"
	"github.com/umisto/pagi"
)

type InsertRoleParams struct {
	AgglomerationID uuid.UUID `json:"agglomeration_id"`
	Head            bool      `json:"head"`
	Editable        bool      `json:"editable"`
	Rank            int       `json:"rank"`
	Name            string    `json:"name"`
}

func (s Service) CreateRole(ctx context.Context, params InsertRoleParams) (entity.Role, error) {
	row, err := s.rolesQ().Insert(ctx, pgdb.InsertRoleParams{
		AgglomerationID: params.AgglomerationID,
		Head:            params.Head,
		Editable:        params.Editable,
		Rank:            params.Rank,
		Name:            params.Name,
	})
	if err != nil {
		return entity.Role{}, err
	}

	return models.Role(row), nil
}

func (s Service) GetRole(ctx context.Context, roleID uuid.UUID) (entity.Role, error) {
	row, err := s.rolesQ().FilterByID(roleID).Get(ctx)
	if err != nil {
		return entity.Role{}, err
	}

	return models.Role(row), nil
}

type FilterParams struct {
	AgglomerationID *uuid.UUID
	Head            *bool
	Editable        *bool
	Rank            *int
	Name            *string
}

func (s Service) FilterRoles(
	ctx context.Context,
	filter FilterParams,
	offset uint,
	limit uint,
) (pagi.Page[[]entity.Role], error) {
	q := s.rolesQ()
	if filter.AgglomerationID != nil {
		q = q.FilterByAgglomerationID(*filter.AgglomerationID)
	}
	if filter.Head != nil {
		q = q.FilterHead(*filter.Head)
	}
	if filter.Editable != nil {
		q = q.FilterEditable(*filter.Editable)
	}
	if filter.Rank != nil {
		q = q.FilterByRank(*filter.Rank)
	}
	if filter.Name != nil {
		q = q.FilterLikeName(*filter.Name)
	}

	limit = pagi.CalculateLimit(limit, 20, 100)

	rows, err := q.Page(limit, offset).Select(ctx)
	if err != nil {
		return pagi.Page[[]entity.Role]{}, err
	}

	total, err := q.Count(ctx)
	if err != nil {
		return pagi.Page[[]entity.Role]{}, err
	}

	collection := make([]entity.Role, 0, len(rows))
	for _, row := range rows {
		collection = append(collection, models.Role(row))
	}

	return pagi.Page[[]entity.Role]{
		Data:  collection,
		Total: uint(total),
		Page:  uint(offset/limit) + 1,
		Size:  uint(len(collection)),
	}, nil
}

type UpdateRoleParams struct {
	Name *string `json:"name"`
}

func (s Service) UpdateRole(ctx context.Context, roleID uuid.UUID, params UpdateRoleParams) (entity.Role, error) {
	q := s.rolesQ().FilterByID(roleID)
	if params.Name != nil {
		q = q.UpdateName(*params.Name)
	}

	row, err := q.UpdateOne(ctx)
	if err != nil {
		return entity.Role{}, err
	}

	return models.Role(row), nil
}

func (s Service) UpdateRoleRank(ctx context.Context, roleID uuid.UUID, newRank uint) (entity.Role, error) {
	row, err := s.rolesQ().UpdateRoleRank(ctx, roleID, newRank)
	if err != nil {
		return entity.Role{}, err
	}

	return models.Role(row), nil
}

func (s Service) DeleteRole(ctx context.Context, roleID uuid.UUID) error {
	return s.rolesQ().DeleteAndShiftRanks(ctx, roleID)
}
