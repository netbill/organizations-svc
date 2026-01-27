package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/internal/repository/pgdb"
)

func (r Repository) GetRolePermissions(ctx context.Context, roleID uuid.UUID) (map[models.Permission]bool, error) {
	rolePerm, err := r.orgRolePermissionsQ(ctx).FilterByRoleID(roleID).Select(ctx)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		// Role has no permissions
	case err != nil:
		return nil, fmt.Errorf("failed to fetch role permissions, cause: %w", err)
	}

	perm, err := r.orgRolePermissionsQ(ctx).Select(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch permissions, cause: %w", err)
	}

	rolePermMap := make(map[models.Permission]bool, len(perm))
	for _, p := range perm {
		exist := false
		for _, rp := range rolePerm {
			if p.Code == rp.Code {
				exist = true
				break
			}
		}

		rolePermMap[models.Permission{
			Code:        p.Code,
			Description: p.Description,
		}] = exist
	}

	return rolePermMap, nil
}

func (r Repository) GetAllPermissions(ctx context.Context) ([]models.Permission, error) {
	permissions, err := r.orgRolePermissionsQ(ctx).Select(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch permissions, cause: %w", err)
	}

	result := make([]models.Permission, len(permissions))
	for i, perm := range permissions {
		result[i] = models.Permission{
			Code:        perm.Code,
			Description: perm.Description,
		}
	}

	return result, nil
}

func (r Repository) SetRolePermissions(
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
		err := r.orgRolePermissionLinksQ(ctx).
			FilterByRoleID(roleID).
			FilterByPermissionCode(deletePermissions...).
			Delete(ctx)
		if err != nil {
			return fmt.Errorf("failed to delete role permissions, cause: %w", err)
		}
	}

	if len(addPermissions) > 0 {
		p, err := r.orgRolePermissionsQ(ctx).FilterByCode(addPermissions...).Select(ctx)
		if err != nil {
			return fmt.Errorf("failed to select permissions to add, cause: %w", err)
		}

		existingPermissionsMap := make([]pgdb.OrganizationRolePermissionLink, len(p))
		for i, perm := range p {
			existingPermissionsMap[i] = pgdb.OrganizationRolePermissionLink{
				RoleID:         roleID,
				PermissionCode: perm.Code,
			}
		}
		if err = r.orgRolePermissionLinksQ(ctx).Insert(ctx, existingPermissionsMap...); err != nil {
			return fmt.Errorf("failed to insert role permissions, cause: %w", err)
		}
	}

	return nil
}

func (r Repository) CheckMemberHavePermission(
	ctx context.Context,
	memberID uuid.UUID,
	permissionCode string,
) (bool, error) {
	have, err := r.orgMembersQ(ctx).
		FilterByID(memberID).
		FilterByPermissionCode(permissionCode).
		Exists(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to check member permission, cause: %w", err)
	}

	return have, nil
}
