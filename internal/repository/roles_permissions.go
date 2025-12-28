package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/entity"
	"github.com/umisto/cities-svc/internal/repository/pgdb"
)

func (s Service) CreatePermissionForRole(ctx context.Context, roleID, permissionID uuid.UUID) error {
	_, err := s.rolePermissionsQ().Insert(ctx, pgdb.RolePermission{
		RoleID:       roleID,
		PermissionID: permissionID,
	})
	return err
}

func (s Service) GetPermissionsForRole(ctx context.Context, roleID uuid.UUID) ([]entity.Permission, error) {
	roes, err := s.permissionsQ().FilterByRoleID(roleID).Select(ctx)
	if err != nil {
		return nil, err
	}

	permissions := make([]entity.Permission, 0, len(roes))
	for _, r := range roes {
		permissions = append(permissions, entity.Permission{
			ID:          r.ID,
			Code:        entity.CodeRolePermission(r.Code),
			Description: r.Description,
		})
	}

	return permissions, nil
}

func (s Service) DeletePermissionForRole(ctx context.Context, roleID, permissionID uuid.UUID) error {
	return s.rolePermissionsQ().
		FilterByRoleID(roleID).
		FilterByPermissionID(permissionID).
		Delete(ctx)
}

func (s Service) CheckMemberHavePermissionsInAgglomerationByCode(
	ctx context.Context,
	memberID, agglomerationID uuid.UUID,
	permission string,
) (bool, error) {
	return s.rolePermissionsQ().
		FilterByMemberID(memberID).
		FilterByAgglomerationID(agglomerationID).
		FilterByPermissionCode(permission).
		Exists(ctx)
}

func (s Service) CheckMemberHavePermissionsInAgglomerationByID(
	ctx context.Context,
	memberID, agglomerationID uuid.UUID,
	permissionID uuid.UUID,
) (bool, error) {
	return s.rolePermissionsQ().
		FilterByMemberID(memberID).
		FilterByAgglomerationID(agglomerationID).
		FilterByPermissionID(permissionID).
		Exists(ctx)
}

func (s Service) CheckAccountHavePermissionByID(
	ctx context.Context,
	accountID, agglomerationID, permissionID uuid.UUID,
) (bool, error) {
	return s.rolePermissionsQ().
		FilterByAccountID(accountID).
		FilterByAgglomerationID(agglomerationID).
		FilterByPermissionID(permissionID).
		Exists(ctx)
}

func (s Service) CheckAccountHavePermissionByCode(
	ctx context.Context,
	accountID, agglomerationID uuid.UUID,
	permission string,
) (bool, error) {
	return s.rolePermissionsQ().
		FilterByAccountID(accountID).
		FilterByAgglomerationID(agglomerationID).
		FilterByPermissionCode(permission).
		Exists(ctx)
}
