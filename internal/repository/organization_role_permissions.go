package repository

import (
	"context"
	"fmt"
	"time"

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
	CreatedAt      time.Time `db:"created_at"`
}

func (r OrganizationRolePermissionLinkRow) IsNil() bool {
	return r.RoleID == uuid.Nil
}

type OrgRolePermissionLinksQ interface {
	New() OrgRolePermissionLinksQ

	Insert(
		ctx context.Context,
		roleID uuid.UUID,
		codes ...string,
	) ([]OrganizationRolePermissionLinkRow, error)

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

func (r *Repository) GetRolePermissions(
	ctx context.Context,
	roleID uuid.UUID,
) (models.OrgRolePermissionLinks, error) {
	dict, err := r.orgRolePermissionsQ().Select(ctx)
	if err != nil {
		return models.OrgRolePermissionLinks{}, fmt.Errorf("failed to get permissions dict: %w", err)
	}

	links, err := r.orgRolePermissionLinksQ().
		FilterByRoleID(roleID).
		Select(ctx)
	if err != nil {
		return models.OrgRolePermissionLinks{}, fmt.Errorf("failed to get role permission links: %w", err)
	}

	enabled := make(map[string]struct{}, len(links))
	for i := range links {
		enabled[links[i].PermissionCode] = struct{}{}
	}

	desc := make(map[string]string, len(dict))
	for i := range dict {
		desc[dict[i].Code] = dict[i].Description
	}

	out := models.OrgRolePermissionLinks{
		ManageOrganization: models.OrgRolePermissionWithFlag{
			Code:        models.RolePermissionManageOrganization,
			Description: desc[models.RolePermissionManageOrganization],
		},
		ManageInvites: models.OrgRolePermissionWithFlag{
			Code:        models.RolePermissionManageInvites,
			Description: desc[models.RolePermissionManageInvites],
		},
		ManageMembers: models.OrgRolePermissionWithFlag{
			Code:        models.RolePermissionManageMembers,
			Description: desc[models.RolePermissionManageMembers],
		},
		ManageRoles: models.OrgRolePermissionWithFlag{
			Code:        models.RolePermissionManageRoles,
			Description: desc[models.RolePermissionManageRoles],
		},
	}

	if _, ok := enabled[models.RolePermissionManageOrganization]; ok {
		out.ManageOrganization.Enabled = true
	}
	if _, ok := enabled[models.RolePermissionManageInvites]; ok {
		out.ManageInvites.Enabled = true
	}
	if _, ok := enabled[models.RolePermissionManageMembers]; ok {
		out.ManageMembers.Enabled = true
	}
	if _, ok := enabled[models.RolePermissionManageRoles]; ok {
		out.ManageRoles.Enabled = true
	}

	return out, nil
}

func (r *Repository) GetAllPermissions(
	ctx context.Context,
) ([]models.OrgRolePermission, error) {
	permissions, err := r.orgRolePermissionsQ().Select(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all permissions, cause: %w", err)
	}

	result := make([]models.OrgRolePermission, len(permissions))
	for i, perm := range permissions {
		result[i] = models.OrgRolePermission{
			Code:        perm.Code,
			Description: perm.Description,
		}
	}

	return result, nil
}

func (r *Repository) SetRolePermissions(
	ctx context.Context,
	roleID uuid.UUID,
	permissions models.OrgRolePermissionDict,
) (models.OrgRolePermissionLinks, error) {
	codes := make([]string, 0, 4)

	if permissions.ManageOrganization {
		codes = append(codes, models.RolePermissionManageOrganization)
	}
	if permissions.ManageInvites {
		codes = append(codes, models.RolePermissionManageInvites)
	}
	if permissions.ManageMembers {
		codes = append(codes, models.RolePermissionManageMembers)
	}
	if permissions.ManageRoles {
		codes = append(codes, models.RolePermissionManageRoles)
	}

	rows, err := r.orgRolePermissionLinksQ().Insert(ctx, roleID, codes...)
	if err != nil {
		return models.OrgRolePermissionLinks{}, fmt.Errorf("set role permissions: %w", err)
	}

	enabled := make(map[string]struct{}, len(rows))
	for i := range rows {
		enabled[rows[i].PermissionCode] = struct{}{}
	}

	dict, err := r.orgRolePermissionsQ().Select(ctx)
	if err != nil {
		return models.OrgRolePermissionLinks{}, fmt.Errorf("select permissions dict: %w", err)
	}

	desc := make(map[string]string, len(dict))
	for i := range dict {
		desc[dict[i].Code] = dict[i].Description
	}

	out := models.OrgRolePermissionLinks{
		ManageOrganization: models.OrgRolePermissionWithFlag{
			Code:        models.RolePermissionManageOrganization,
			Description: desc[models.RolePermissionManageOrganization],
			Enabled:     false,
		},
		ManageInvites: models.OrgRolePermissionWithFlag{
			Code:        models.RolePermissionManageInvites,
			Description: desc[models.RolePermissionManageInvites],
			Enabled:     false,
		},
		ManageMembers: models.OrgRolePermissionWithFlag{
			Code:        models.RolePermissionManageMembers,
			Description: desc[models.RolePermissionManageMembers],
			Enabled:     false,
		},
		ManageRoles: models.OrgRolePermissionWithFlag{
			Code:        models.RolePermissionManageRoles,
			Description: desc[models.RolePermissionManageRoles],
			Enabled:     false,
		},
	}

	if _, ok := enabled[models.RolePermissionManageOrganization]; ok {
		out.ManageOrganization.Enabled = true
	}
	if _, ok := enabled[models.RolePermissionManageInvites]; ok {
		out.ManageInvites.Enabled = true
	}
	if _, ok := enabled[models.RolePermissionManageMembers]; ok {
		out.ManageMembers.Enabled = true
	}
	if _, ok := enabled[models.RolePermissionManageRoles]; ok {
		out.ManageRoles.Enabled = true
	}

	return out, nil
}

func (r *Repository) CheckMemberHavePermission(
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
