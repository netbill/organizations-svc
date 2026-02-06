package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/netbill/organizations-svc/internal/core/models"
)

type OrganizationRolePermissionRow struct {
	Code         string     `db:"code"`
	Description  string     `db:"description"`
	DeprecatedAt *time.Time `db:"deprecated_at"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
}

func (r OrganizationRolePermissionRow) IsNil() bool {
	return r.Code == ""
}

type OrgRolePermissionsQ interface {
	New() OrgRolePermissionsQ

	Insert(ctx context.Context, input OrganizationRolePermissionRow) (OrganizationRolePermissionRow, error)

	FilterByCode(codes ...string) OrgRolePermissionsQ
	FilterByDeprecated(deprecated bool) OrgRolePermissionsQ

	UpdateOne(ctx context.Context) (OrganizationRolePermissionRow, error)
	UpdateMany(ctx context.Context) (int64, error)

	UpdateDeprecatedAt(timestamp *time.Time) OrgRolePermissionsQ

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

	Upsert(
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
) (models.OrgRolePermissionDictWithDetails, error) {
	dict, err := r.orgRolePermissionsQ().Select(ctx)
	if err != nil {
		return models.OrgRolePermissionDictWithDetails{}, fmt.Errorf("failed to get permissions dict: %w", err)
	}

	links, err := r.orgRolePermissionLinksQ().
		FilterByRoleID(roleID).
		Select(ctx)
	if err != nil {
		return models.OrgRolePermissionDictWithDetails{}, fmt.Errorf("failed to get role permission links: %w", err)
	}

	enabled := make(map[string]struct{}, len(links))
	for i := range links {
		enabled[links[i].PermissionCode] = struct{}{}
	}

	desc := make(map[string]string, len(dict))
	for i := range dict {
		desc[dict[i].Code] = dict[i].Description
	}

	out := models.OrgRolePermissionDictWithDetails{
		OrganizationUpdate: models.OrgRolePermissionDetails{
			Description: desc[models.RolePermissionOrgUpdate],
		},
		InvitesManage: models.OrgRolePermissionDetails{
			Description: desc[models.RolePermissionInvitesManage],
		},
		RolesManage: models.OrgRolePermissionDetails{
			Description: desc[models.RolePermissionRolesManage],
		},
		MembersDelete: models.OrgRolePermissionDetails{
			Description: desc[models.RolePermissionMembersDelete],
		},
		MembersUpdate: models.OrgRolePermissionDetails{
			Description: desc[models.RolePermissionMembersUpdate],
		},
		PlaceCreate: models.OrgRolePermissionDetails{
			Description: desc[models.RolePermissionPlaceCreate],
		},
		PlaceDelete: models.OrgRolePermissionDetails{
			Description: desc[models.RolePermissionPlaceDelete],
		},
		PlaceUpdate: models.OrgRolePermissionDetails{
			Description: desc[models.RolePermissionPlaceUpdate],
		},
	}

	if _, ok := enabled[models.RolePermissionOrgUpdate]; ok {
		out.OrganizationUpdate.Enabled = true
	}
	if _, ok := enabled[models.RolePermissionInvitesManage]; ok {
		out.InvitesManage.Enabled = true
	}
	if _, ok := enabled[models.RolePermissionRolesManage]; ok {
		out.RolesManage.Enabled = true
	}
	if _, ok := enabled[models.RolePermissionMembersDelete]; ok {
		out.MembersDelete.Enabled = true
	}
	if _, ok := enabled[models.RolePermissionMembersUpdate]; ok {
		out.MembersUpdate.Enabled = true
	}
	if _, ok := enabled[models.RolePermissionPlaceCreate]; ok {
		out.PlaceCreate.Enabled = true
	}
	if _, ok := enabled[models.RolePermissionPlaceDelete]; ok {
		out.PlaceDelete.Enabled = true
	}
	if _, ok := enabled[models.RolePermissionPlaceUpdate]; ok {
		out.PlaceUpdate.Enabled = true
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
	permissions models.OrgRolePermissionAccess,
) (models.OrgRolePermissionDictWithDetails, error) {
	codes := make([]string, 0, models.GetOrgRolePermissionLength())

	if permissions.OrgUpdate {
		codes = append(codes, models.RolePermissionOrgUpdate)
	}
	if permissions.InvitesManage {
		codes = append(codes, models.RolePermissionInvitesManage)
	}
	if permissions.MembersDelete {
		codes = append(codes, models.RolePermissionMembersDelete)
	}
	if permissions.MembersUpdate {
		codes = append(codes, models.RolePermissionMembersUpdate)
	}
	if permissions.RolesManage {
		codes = append(codes, models.RolePermissionRolesManage)
	}
	if permissions.PlaceCreate {
		codes = append(codes, models.RolePermissionPlaceCreate)
	}
	if permissions.PlaceDelete {
		codes = append(codes, models.RolePermissionPlaceDelete)
	}
	if permissions.PlaceUpdate {
		codes = append(codes, models.RolePermissionPlaceUpdate)
	}

	rows, err := r.orgRolePermissionLinksQ().Upsert(ctx, roleID, codes...)
	if err != nil {
		return models.OrgRolePermissionDictWithDetails{}, fmt.Errorf("set role permissions: %w", err)
	}

	enabled := make(map[string]struct{}, len(rows))
	for i := range rows {
		enabled[rows[i].PermissionCode] = struct{}{}
	}

	dict, err := r.orgRolePermissionsQ().Select(ctx)
	if err != nil {
		return models.OrgRolePermissionDictWithDetails{}, fmt.Errorf("select permissions dict: %w", err)
	}

	desc := make(map[string]string, len(dict))
	for i := range dict {
		desc[dict[i].Code] = dict[i].Description
	}

	out := models.OrgRolePermissionDictWithDetails{
		OrganizationUpdate: models.OrgRolePermissionDetails{
			Description: desc[models.RolePermissionOrgUpdate],
		},
		InvitesManage: models.OrgRolePermissionDetails{
			Description: desc[models.RolePermissionInvitesManage],
		},
		RolesManage: models.OrgRolePermissionDetails{
			Description: desc[models.RolePermissionRolesManage],
		},
		MembersDelete: models.OrgRolePermissionDetails{
			Description: desc[models.RolePermissionMembersDelete],
		},
		MembersUpdate: models.OrgRolePermissionDetails{
			Description: desc[models.RolePermissionMembersUpdate],
		},
		PlaceCreate: models.OrgRolePermissionDetails{
			Description: desc[models.RolePermissionPlaceCreate],
		},
		PlaceDelete: models.OrgRolePermissionDetails{
			Description: desc[models.RolePermissionPlaceDelete],
		},
		PlaceUpdate: models.OrgRolePermissionDetails{
			Description: desc[models.RolePermissionPlaceUpdate],
		},
	}

	if _, ok := enabled[models.RolePermissionOrgUpdate]; ok {
		out.OrganizationUpdate.Enabled = true
	}
	if _, ok := enabled[models.RolePermissionInvitesManage]; ok {
		out.InvitesManage.Enabled = true
	}
	if _, ok := enabled[models.RolePermissionRolesManage]; ok {
		out.RolesManage.Enabled = true
	}
	if _, ok := enabled[models.RolePermissionMembersDelete]; ok {
		out.MembersDelete.Enabled = true
	}
	if _, ok := enabled[models.RolePermissionMembersUpdate]; ok {
		out.MembersUpdate.Enabled = true
	}
	if _, ok := enabled[models.RolePermissionPlaceCreate]; ok {
		out.PlaceCreate.Enabled = true
	}
	if _, ok := enabled[models.RolePermissionPlaceDelete]; ok {
		out.PlaceDelete.Enabled = true
	}
	if _, ok := enabled[models.RolePermissionPlaceUpdate]; ok {
		out.PlaceUpdate.Enabled = true
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
