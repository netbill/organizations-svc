package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/repository/pgdb"
)

func (s Service) CreatePermissionForRole(ctx context.Context, roleID, permissionID uuid.UUID) error {
	return s.sql.CreateRolePermission(ctx, pgdb.CreateRolePermissionParams{
		RoleID:       roleID,
		PermissionID: permissionID,
	})
}

func (s Service) DeletePermissionForRole(ctx context.Context, roleID, permissionID uuid.UUID) error {
	return s.sql.DeleteRolePermission(ctx, pgdb.DeleteRolePermissionParams{
		RoleID:       roleID,
		PermissionID: permissionID,
	})
}

func (s Service) CheckMemberHavePermissionByID(ctx context.Context, memberID, permissionID uuid.UUID) (bool, error) {
	has, err := s.sql.CheckMemberHavePermissionByID(ctx, pgdb.CheckMemberHavePermissionByIDParams{
		MemberID:     memberID,
		PermissionID: permissionID,
	})
	if err != nil {
		return false, err
	}

	return has, nil
}

func (s Service) CheckMemberHavePermissionByCode(ctx context.Context, memberID uuid.UUID, permissionKey string) (bool, error) {
	has, err := s.sql.CheckMemberHavePermissionByCode(ctx, pgdb.CheckMemberHavePermissionByCodeParams{
		MemberID: memberID,
		Code:     permissionKey,
	})
	if err != nil {
		return false, err
	}

	return has, nil
}
