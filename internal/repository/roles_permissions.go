package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/models"
	"github.com/umisto/cities-svc/internal/repository/pgdb"
)

func (s Service) GetAllPermissions(ctx context.Context) ([]models.Permission, error) {
	permissions, err := s.permissionsQ().Select(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]models.Permission, len(permissions))
	for i, perm := range permissions {
		result[i] = models.Permission{
			ID:          perm.ID,
			Code:        models.CodeRolePermission(perm.Code),
			Description: perm.Description,
		}
	}

	return result, nil
}

func (s Service) SetRolePermissions(
	ctx context.Context,
	roleID uuid.UUID,
	permissions map[models.CodeRolePermission]bool,
) error {
	deletePermissions := make([]string, 0)
	addPermissions := make([]string, 0)

	for perm, toSet := range permissions {
		if toSet {
			addPermissions = append(addPermissions, perm.String())
		} else {
			deletePermissions = append(deletePermissions, perm.String())
		}
	}

	if len(deletePermissions) > 0 {
		if err := s.rolePermissionsQ().
			FilterByRoleID(roleID).
			FilterByPermissionCode(deletePermissions...).
			Delete(ctx); err != nil {
			return err
		}
	}

	if len(addPermissions) > 0 {
		p, err := s.permissionsQ().FilterByCode(addPermissions...).Select(ctx)
		if err != nil {
			return err
		}
		existingPermissionsMap := make([]pgdb.RolePermission, len(p))
		for i, perm := range p {
			existingPermissionsMap[i] = pgdb.RolePermission{
				RoleID:       roleID,
				PermissionID: perm.ID,
			}
		}
		if err := s.rolePermissionsQ().Insert(ctx, existingPermissionsMap...); err != nil {
			return err
		}
	}

	return nil
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
