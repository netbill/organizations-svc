package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/umisto/agglomerations-svc/internal/domain/models"
	"github.com/umisto/agglomerations-svc/internal/repository/pgdb"
	"github.com/umisto/pagi"
)

func (s Service) GetPermission(ctx context.Context, ID uuid.UUID) (models.Permission, error) {
	res, err := s.permissionsQ().FilterByID(ID).Get(ctx)
	if err != nil {
		return models.Permission{}, err
	}

	return Permission(res), nil
}

func (s Service) GetPermissionByCode(ctx context.Context, code models.CodeRolePermission) (models.Permission, error) {
	res, err := s.permissionsQ().FilterByCode(string(code)).Get(ctx)
	if err != nil {
		return models.Permission{}, err
	}

	return Permission(res), nil
}

type FilterPermissionsParams struct {
	Description *string
	Code        *models.CodeRolePermission
}

func (s Service) GetPermissions(
	ctx context.Context,
	filter FilterPermissionsParams,
	offset uint,
	limit uint,
) (pagi.Page[[]models.Permission], error) {
	q := s.permissionsQ()
	if filter.Description != nil {
		q = q.FilterLikeDescription(*filter.Description)
	}
	if filter.Code != nil {
		q = q.FilterByCode(string(*filter.Code))
	}

	rows, err := q.Page(limit, offset).Select(ctx)
	if err != nil {
		return pagi.Page[[]models.Permission]{}, err
	}

	total, err := q.Count(ctx)
	if err != nil {
		return pagi.Page[[]models.Permission]{}, err
	}

	collection := make([]models.Permission, 0, len(rows))
	for i, row := range rows {
		collection[i] = Permission(row)
	}

	return pagi.Page[[]models.Permission]{
		Data:  collection,
		Page:  uint(offset/limit) + 1,
		Size:  uint(len(collection)),
		Total: uint(total),
	}, nil
}

func Permission(p pgdb.Permission) models.Permission {
	return models.Permission{
		ID:          p.ID,
		Code:        models.CodeRolePermission(p.Code),
		Description: p.Description,
	}
}
