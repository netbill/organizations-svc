package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/entity"
	"github.com/umisto/cities-svc/internal/repository/pgdb"
	"github.com/umisto/pagi"
)

func (s Service) GetPermission(ctx context.Context, ID uuid.UUID) (entity.Permission, error) {
	res, err := s.permissionsQ().FilterByID(ID).Get(ctx)
	if err != nil {
		return entity.Permission{}, err
	}

	return Permission(res), nil
}

func (s Service) GetPermissionByCode(ctx context.Context, code entity.CodeRolePermission) (entity.Permission, error) {
	res, err := s.permissionsQ().FilterByCode(string(code)).Get(ctx)
	if err != nil {
		return entity.Permission{}, err
	}

	return Permission(res), nil
}

type FilterPermissionsParams struct {
	Description *string
	Code        *entity.CodeRolePermission
}

func (s Service) FilterPermissions(
	ctx context.Context,
	filter FilterPermissionsParams,
	offset uint,
	limit uint,
) (pagi.Page[[]entity.Permission], error) {
	q := s.permissionsQ()
	if filter.Description != nil {
		q = q.FilterLikeDescription(*filter.Description)
	}
	if filter.Code != nil {
		q = q.FilterByCode(string(*filter.Code))
	}

	rows, err := q.Page(limit, offset).Select(ctx)
	if err != nil {
		return pagi.Page[[]entity.Permission]{}, err
	}

	total, err := q.Count(ctx)
	if err != nil {
		return pagi.Page[[]entity.Permission]{}, err
	}

	collection := make([]entity.Permission, 0, len(rows))
	for i, row := range rows {
		collection[i] = Permission(row)
	}

	return pagi.Page[[]entity.Permission]{
		Data:  collection,
		Page:  uint(offset/limit) + 1,
		Size:  uint(len(collection)),
		Total: uint(total),
	}, nil
}

func Permission(p pgdb.Permission) entity.Permission {
	return entity.Permission{
		ID:          p.ID,
		Code:        entity.CodeRolePermission(p.Code),
		Description: p.Description,
	}
}
