package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/netbill/organizations-svc/internal/core/models"
)

type OrganizationRolePermissionRow struct {
	Code        string `db:"code"`
	Description string `db:"description"`
}

func (r OrganizationRolePermissionRow) IsNil() bool {
	return r.Code == ""
}

type OrgRolePermissionsQ interface {
	New() OrgRolePermissionsQ

	Insert(ctx context.Context, input OrganizationRolePermissionRow) (OrganizationRolePermissionRow, error)

	FilterByRoleID(roleID uuid.UUID) OrgRolePermissionsQ
	FilterByCode(codes ...string) OrgRolePermissionsQ
	FilterLikeDescription(description string) OrgRolePermissionsQ

	UpdateOne(ctx context.Context) (OrganizationRolePermissionRow, error)
	UpdateMany(ctx context.Context) (int64, error)

	Select(ctx context.Context) ([]OrganizationRolePermissionRow, error)
	Get(ctx context.Context) (OrganizationRolePermissionRow, error)

	Delete(ctx context.Context) error
	Count(ctx context.Context) (uint, error)
	Page(limit, offset uint) OrgRolePermissionsQ
}

type OrganizationRolePermissionLinkRow struct {
	RoleID         uuid.UUID `db:"role_id"`
	PermissionCode string    `db:"permission_code"`
}

func (r OrganizationRolePermissionLinkRow) IsNil() bool {
	return r.RoleID == uuid.Nil
}

type OrgRolePermissionLinksQ interface {
	New() OrgRolePermissionLinksQ

	Insert(ctx context.Context, input ...OrganizationRolePermissionLinkRow) error

	Select(ctx context.Context) ([]OrganizationRolePermissionLinkRow, error)
	Get(ctx context.Context) (OrganizationRolePermissionLinkRow, error)

	FilterByAccountID(accountID uuid.UUID) OrgRolePermissionLinksQ
	FilterByOrganizationID(organizationID uuid.UUID) OrgRolePermissionLinksQ
	FilterByMemberID(memberID uuid.UUID) OrgRolePermissionLinksQ
	FilterByRoleID(roleID uuid.UUID) OrgRolePermissionLinksQ
	FilterByPermissionCode(codes ...string) OrgRolePermissionLinksQ

	Delete(ctx context.Context) error

	Count(ctx context.Context) (uint, error)
	Page(limit, offset uint) OrgRolePermissionLinksQ
	Exists(ctx context.Context) (bool, error)
}

func (r Repository) GetRolePermissions(ctx context.Context, roleID uuid.UUID) (map[models.Permission]bool, error) {
	rolePerm, err := r.orgRolePermissionsQ().FilterByRoleID(roleID).Select(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get role permissions, cause: %w", err)
	}

	perm, err := r.orgRolePermissionsQ().Select(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all permissions, cause: %w", err)
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
	permissions, err := r.orgRolePermissionsQ().Select(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all permissions, cause: %w", err)
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
		err := r.orgRolePermissionLinksQ().
			FilterByRoleID(roleID).
			FilterByPermissionCode(deletePermissions...).
			Delete(ctx)
		if err != nil {
			return fmt.Errorf("failed to delete role permissions, cause: %w", err)
		}
	}

	if len(addPermissions) > 0 {
		p, err := r.orgRolePermissionsQ().FilterByCode(addPermissions...).Select(ctx)
		if err != nil {
			return fmt.Errorf("filed to getting existing permissions, cause: %w", err)
		}

		existingPermissionsMap := make([]OrganizationRolePermissionLinkRow, len(p))
		for i, perm := range p {
			existingPermissionsMap[i] = OrganizationRolePermissionLinkRow{
				RoleID:         roleID,
				PermissionCode: perm.Code,
			}
		}
		if err = r.orgRolePermissionLinksQ().Insert(ctx, existingPermissionsMap...); err != nil {
			return fmt.Errorf("failed to adding role permissions, cause: %w", err)
		}
	}

	return nil
}

func (r Repository) CheckMemberHavePermission(
	ctx context.Context,
	memberID uuid.UUID,
	permissionCode string,
) (bool, error) {
	have, err := r.orgMembersQ().
		FilterByID(memberID).
		FilterByPermissionCode(permissionCode).
		Exists(ctx)
	if err != nil {
		return false, fmt.Errorf("checking member permission: %w", err)
	}

	return have, nil
}
