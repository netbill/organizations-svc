package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/internal/repository/pgdb"
)

func (s Service) GetRolePermissions(ctx context.Context, roleID uuid.UUID) (map[models.Permission]bool, error) {
	rows, err := s.permissionsQ().GetForRole(ctx, roleID)
	if err != nil {
		return nil, err
	}

	result := make(map[models.Permission]bool, len(rows))
	for el, row := range rows {
		perm := models.Permission{
			ID:          el.ID,
			Code:        el.Code,
			Description: el.Description,
		}
		result[perm] = row
	}

	return result, nil
}

func (s Service) GetAllPermissions(ctx context.Context) ([]models.Permission, error) {
	permissions, err := s.permissionsQ().Select(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]models.Permission, len(permissions))
	for i, perm := range permissions {
		result[i] = models.Permission{
			ID:          perm.ID,
			Code:        perm.Code,
			Description: perm.Description,
		}
	}

	return result, nil
}

func (s Service) SetRolePermissions(
	ctx context.Context,
	roleID uuid.UUID,
	permissions map[string]bool,
) error {
	deletePermissions := make([]string, 0)
	addPermissions := make([]string, 0)

	for perm, toSet := range permissions {
		if toSet {
			addPermissions = append(addPermissions, perm)
		} else {
			deletePermissions = append(deletePermissions, perm)
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
		if err = s.rolePermissionsQ().Insert(ctx, existingPermissionsMap...); err != nil {
			return err
		}
	}

	return nil
}

func (s Service) CheckMemberHavePermission(
	ctx context.Context,
	memberID uuid.UUID,
	permissionCode string,
) (bool, error) {
	have, err := s.membersQ().
		FilterByID(memberID).
		FilterByPermissionCode(permissionCode).
		Exists(ctx)
	if err != nil {
		return false, fmt.Errorf("checking member have permission: %w", err)
	}

	return have, nil
}
