package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/umisto/agglomerations-svc/internal/domain/models"
	"github.com/umisto/agglomerations-svc/internal/domain/modules/role"
	"github.com/umisto/agglomerations-svc/internal/repository/pgdb"
	"github.com/umisto/pagi"
)

func (s Service) CreateRole(ctx context.Context, params role.CreateParams) (models.Role, error) {
	row, err := s.rolesQ().Insert(ctx, pgdb.InsertRoleParams{
		AgglomerationID: params.AgglomerationID,
		Rank:            params.Rank,
		Name:            params.Name,
		Description:     params.Description,
		Color:           params.Color,
	})
	if err != nil {
		return models.Role{}, err
	}

	return Role(row), nil
}

func (s Service) CreateHeadRole(ctx context.Context, agglomerationID uuid.UUID) (models.Role, error) {
	row, err := s.rolesQ().Insert(ctx, pgdb.InsertRoleParams{
		AgglomerationID: agglomerationID,
		Head:            true,
		Rank:            1,
		Name:            "Head",
		Description:     "Head role with all permissions",
		Color:           "#000000",
	})
	if err != nil {
		return models.Role{}, err
	}

	return Role(row), nil
}
func (s Service) GetRole(ctx context.Context, roleID uuid.UUID) (models.Role, error) {
	row, err := s.rolesQ().FilterByID(roleID).Get(ctx)
	if err != nil {
		return models.Role{}, err
	}

	return Role(row), nil
}

func (s Service) GetRoles(
	ctx context.Context,
	filter role.FilterParams,
	offset uint,
	limit uint,
) (pagi.Page[[]models.Role], error) {
	q := s.rolesQ()
	if filter.AgglomerationID != nil {
		q = q.FilterByAgglomerationID(*filter.AgglomerationID)
	}
	if filter.RolesID != nil && len(*filter.RolesID) > 0 {
		q = q.FilterByID(*filter.RolesID...)
	}
	if filter.Head != nil {
		q = q.FilterHead(*filter.Head)
	}
	if filter.Rank != nil {
		q = q.FilterByRank(*filter.Rank)
	}
	if filter.Name != nil {
		q = q.FilterLikeName(*filter.Name)
	}

	rows, err := q.Page(limit, offset).Select(ctx)
	if err != nil {
		return pagi.Page[[]models.Role]{}, err
	}

	total, err := q.Count(ctx)
	if err != nil {
		return pagi.Page[[]models.Role]{}, err
	}

	collection := make([]models.Role, 0, len(rows))
	for _, row := range rows {
		collection = append(collection, Role(row))
	}

	return pagi.Page[[]models.Role]{
		Data:  collection,
		Total: uint(total),
		Page:  uint(offset/limit) + 1,
		Size:  uint(len(collection)),
	}, nil
}

func (s Service) UpdateRole(ctx context.Context, roleID uuid.UUID, params role.UpdateParams) (models.Role, error) {
	q := s.rolesQ().FilterByID(roleID)
	if params.Name != nil {
		q = q.UpdateName(*params.Name)
	}
	if params.Description != nil {
		q = q.UpdateDescription(*params.Description)
	}
	if params.Color != nil {
		q = q.UpdateColor(*params.Color)
	}

	row, err := q.UpdateOne(ctx)
	if err != nil {
		return models.Role{}, err
	}

	return Role(row), nil
}

func (s Service) UpdateRoleRank(ctx context.Context, roleID uuid.UUID, newRank uint) (models.Role, error) {
	row, err := s.rolesQ().UpdateRoleRank(ctx, roleID, newRank)
	if err != nil {
		return models.Role{}, err
	}

	return Role(row), nil
}

func (s Service) UpdateRolesRanks(
	ctx context.Context,
	agglomerationID uuid.UUID,
	order map[uuid.UUID]uint,
) error {
	_, err := s.rolesQ().UpdateRolesRanks(ctx, agglomerationID, order)
	if err != nil {
		return err
	}

	return nil
}

func (s Service) DeleteRole(ctx context.Context, roleID uuid.UUID) error {
	return s.rolesQ().DeleteAndShiftRanks(ctx, roleID)
}

func (s Service) GetAccountMaxRoleInAgglomeration(
	ctx context.Context,
	accountID, agglomerationID uuid.UUID,
) (models.Role, error) {
	res, err := s.rolesQ().
		FilterByAgglomerationID(agglomerationID).
		FilterByAccountID(accountID).
		OrderByRoleRank(true).
		Get(ctx)
	if err != nil {
		return models.Role{}, err
	}
	return Role(res), nil
}

func (s Service) GetMemberMaxRole(
	ctx context.Context,
	memberID uuid.UUID,
) (models.Role, error) {
	res, err := s.rolesQ().
		FilterByAccountID(memberID).
		OrderByRoleRank(true).
		Get(ctx)
	if err != nil {
		return models.Role{}, err
	}
	return Role(res), nil
}

func (s Service) GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]models.Permission, error) {
	rows, err := s.permissionsQ().FilterByRoleID(roleID).Select(ctx)
	if err != nil {
		return nil, err
	}

	permissions := make([]models.Permission, 0, len(rows))
	for _, row := range rows {
		permissions = append(permissions, Permission(row))
	}

	return permissions, nil
}

func Role(r pgdb.Role) models.Role {
	return models.Role{
		ID:              r.ID,
		AgglomerationID: r.AgglomerationID,
		Head:            r.Head,
		Rank:            r.Rank,
		Name:            r.Name,
		Description:     r.Description,
		Color:           r.Color,
		CreatedAt:       r.CreatedAt,
		UpdatedAt:       r.UpdatedAt,
	}
}
