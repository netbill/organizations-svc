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

func (s Service) CreateRole(ctx context.Context, params role.CreateParams) (entity.Role, error) {
	row, err := s.sql(ctx).CreateRole(ctx, pgdb.CreateRoleParams{})
	if err != nil {
		return entity.Role{}, err
	}

	return row.ToEntity(), nil
}

func (s Service) GetRole(ctx context.Context, roleID uuid.UUID) (entity.Role, error) {
	row, err := s.sql(ctx).GetRole(ctx, roleID)
	if err != nil {
		return entity.Role{}, err
	}

	return row.ToEntity(), nil
}

func (s Service) FilterRoles(
	ctx context.Context,
	filter role.FilterParams,
	pagination pagi.Params,
) (pagi.Page[entity.Role], error) {
	params := pgdb.FilterRolesParams{
		AgglomerationID: filter.AgglomerationID,
		MemberID:        nilx.UUID(filter.MemberID),
		PermissionCodes: filter.PermissionCodes,
	}

	if pagination.Cursor != nil {
		cursorRankStr, ok := pagination.Cursor["rank"]
		if !ok || cursorRankStr == "" {
			return pagi.Page[entity.Role]{}, fmt.Errorf("missing rank in pagination cursor")
		}

		var cursorRank int
		_, err := fmt.Sscanf(cursorRankStr, "%d", &cursorRank)
		if err != nil {
			return pagi.Page[entity.Role]{}, fmt.Errorf("invalid rank in pagination cursor: %w", err)
		}

		cursorIDStr, ok := pagination.Cursor["id"]
		if !ok || cursorIDStr == "" {
			return pagi.Page[entity.Role]{}, fmt.Errorf("missing id in pagination cursor")
		}

		cursorID, err := uuid.Parse(cursorIDStr)
		if err != nil {
			return pagi.Page[entity.Role]{}, fmt.Errorf("invalid id in pagination cursor: %w", err)
		}

		params.CursorRank = sql.NullInt32{Int32: int32(cursorRank), Valid: true}
		params.CursorID = uuid.NullUUID{UUID: cursorID, Valid: true}
	}

	limit := pagi.CalculateLimit(pagination.Limit, 50, 100)
	params.Limit = int32(limit)

	rows, err := s.sql(ctx).FilterRoles(ctx, params)
	if err != nil {
		return pagi.Page[entity.Role]{}, err
	}

	count, err := s.sql(ctx).CountRoles(ctx, pgdb.CountRolesParams{
		AgglomerationID: filter.AgglomerationID,
		MemberID:        nilx.UUID(filter.MemberID),
		PermissionCodes: filter.PermissionCodes,
	})
	if err != nil {
		return pagi.Page[entity.Role]{}, err
	}

	entities := make([]entity.Role, len(rows))
	for i, row := range rows {
		entities[i] = row.ToEntity()
	}

	var nextCursor map[string]string
	if len(rows) == limit {
		lastRow := rows[len(rows)-1]
		nextCursor = map[string]string{
			"rank": fmt.Sprintf("%d", lastRow.Rank),
			"id":   lastRow.ID.String(),
		}
	}

	return pagi.Page[entity.Role]{
		Data:       entities,
		Total:      int(count),
		NextCursor: nextCursor,
	}, nil
}

func (s Service) UpdateRole(ctx context.Context, roleID uuid.UUID, params role.UpdateParams) (entity.Role, error) {
	row, err := s.sql(ctx).UpdateRole(ctx, pgdb.UpdateRoleParams{
		ID:   roleID,
		Name: nilx.String(params.Name),
	})
	if err != nil {
		return entity.Role{}, err
	}

	return row.ToEntity(), nil
}

func (s Service) UpdateRoleRank(ctx context.Context, roleID uuid.UUID, newRank int32) (entity.Role, error) {
	row, err := s.sql(ctx).UpdateRoleRank(ctx, pgdb.UpdateRoleRankParams{
		ID:      roleID,
		NewRank: newRank,
	})
	if err != nil {
		return entity.Role{}, err
	}

	return row.ToEntity(), nil
}

func (s Service) DeleteRole(ctx context.Context, roleID uuid.UUID) error {
	return s.sql(ctx).DeleteRole(ctx, roleID)
}
