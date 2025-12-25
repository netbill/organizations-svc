package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/entity"
	"github.com/umisto/cities-svc/internal/repository/pgdb"
	"github.com/umisto/nilx"
	"github.com/umisto/pagi"
)

func (s Service) GetPermission(ctx context.Context, ID uuid.UUID) (entity.Permission, error) {
	res, err := s.sql(ctx).GetPermissionByID(ctx, ID)
	if err != nil {
		return entity.Permission{}, err
	}

	return res.ToEntity(), nil
}

func (s Service) GetPermissionByCode(ctx context.Context, code entity.CodeRolePermission) (entity.Permission, error) {
	res, err := s.sql(ctx).GetPermissionByCode(ctx, string(code))
	if err != nil {
		return entity.Permission{}, err
	}

	return res.ToEntity(), nil
}

func (s Service) ListPermissions(
	ctx context.Context,
	filter role.FilterPermissionsParams,
	pagination pagi.Params,
) (pagi.Page[entity.Permission], error) {
	params := pgdb.FilterPermissionsParams{
		Description: nilx.String(filter.Description),
		Code:        nilx.String(filter.Code),
	}

	if pagination.Cursor != nil {
		cursorCode, ok := pagination.Cursor["code"]
		if !ok {
			return pagi.Page[entity.Permission]{}, fmt.Errorf("invalid code in cursor")
		}

		cursorIDStr, ok := pagination.Cursor["id"]
		if !ok {
			return pagi.Page[entity.Permission]{}, fmt.Errorf("invalid id in cursor")
		}

		cursorID, err := uuid.Parse(cursorIDStr)
		if err != nil {
			return pagi.Page[entity.Permission]{}, fmt.Errorf("invalid id in cursor: %w", err)
		}

		params.AfterCode = sql.NullString{String: cursorCode, Valid: true}
		params.AfterID = uuid.NullUUID{UUID: cursorID, Valid: true}
	}

	rows, err := s.sql(ctx).FilterPermissions(ctx, params)
	if err != nil {
		return pagi.Page[entity.Permission]{}, err
	}

	count, err := s.sql(ctx).CountPermissions(ctx, pgdb.CountPermissionsParams{
		Description: nilx.String(filter.Description),
		Code:        nilx.String(filter.Code),
	})
	if err != nil {
		return pagi.Page[entity.Permission]{}, err
	}

	permissions := make([]entity.Permission, len(rows))
	for i, row := range rows {
		permissions[i] = row.ToEntity()
	}

	var nextCursor map[string]string
	if len(permissions) == pagination.Limit {
		last := permissions[len(permissions)-1]
		nextCursor = map[string]string{
			"code": string(last.Code),
			"id":   last.ID.String(),
		}
	}

	return pagi.Page[entity.Permission]{
		Data:       permissions,
		Total:      int(count),
		NextCursor: nextCursor,
	}, nil
}
