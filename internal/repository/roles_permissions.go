package repository

import (
	"context"

	"github.com/google/uuid"
)

func (s Service) CreatePermissionForRole(ctx context.Context, roleID, permissionID uuid.UUID) error {
	return s.rolePermissionsQ().Insert(ctx, roleID, permissionID)

}

func (s Service) DeletePermissionForRole(ctx context.Context, roleID, permissionID uuid.UUID) error {
	return s.rolePermissionsQ().Delete(ctx, roleID, permissionID)
}

func (s Service) CheckMemberHavePermissionByID(ctx context.Context, memberID, permissionID uuid.UUID) (bool, error) {
	return s.rolePermissionsQ().CheckMemberHavePermissionByID(ctx, memberID, permissionID)
}

func (s Service) CheckMemberHavePermissionByCode(ctx context.Context, memberID uuid.UUID, permissionKey string) (bool, error) {
	return s.rolePermissionsQ().CheckMemberHavePermissionByCode(ctx, memberID, permissionKey)
}

func (s Service) CheckAccountHavePermissionByCode(ctx context.Context, accountID, agglomerationID uuid.UUID, permissionKey string) (bool, error) {
	return s.rolePermissionsQ().CheckAccountHavePermissionByCode(ctx, accountID, agglomerationID, permissionKey)
}

func (s Service) CheckAccountHavePermissionByID(ctx context.Context, accountID, agglomerationID, permissionID uuid.UUID) (bool, error) {
	return s.rolePermissionsQ().CheckAccountHavePermissionByID(ctx, accountID, agglomerationID, permissionID)
}
