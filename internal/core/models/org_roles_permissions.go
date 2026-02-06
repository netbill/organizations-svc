package models

type OrgRolePermissionCode string

const (
	RolePermissionOrgUpdate     = "organization.update"
	RolePermissionRolesManage   = "roles.manage"
	RolePermissionInvitesManage = "invites.manage"

	RolePermissionMembersDelete = "members.delete"
	RolePermissionMembersUpdate = "members.update"

	RolePermissionPlaceCreate = "places.create"
	RolePermissionPlaceDelete = "places.delete"
	RolePermissionPlaceUpdate = "places.update"
)

var allRolePermissions = []string{
	RolePermissionOrgUpdate,
	RolePermissionRolesManage,
	RolePermissionInvitesManage,

	RolePermissionMembersDelete,
	RolePermissionMembersUpdate,

	RolePermissionPlaceCreate,
	RolePermissionPlaceDelete,
	RolePermissionPlaceUpdate,
}

func GetAllOrgRolePermissions() []string {
	return allRolePermissions
}

func GetOrgRolePermissionLength() int {
	return len(allRolePermissions)
}

type OrgRolePermission struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

type OrgRolePermissionDetails struct {
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

type OrgRolePermissionAccess struct {
	OrgUpdate     bool `json:"organization.update"`
	RolesManage   bool `json:"roles.manage"`
	InvitesManage bool `json:"invites.manage"`
	MembersDelete bool `json:"members.delete"`
	MembersUpdate bool `json:"members.update"`
	PlaceCreate   bool `json:"create.place"`
	PlaceDelete   bool `json:"delete.place"`
	PlaceUpdate   bool `json:"update.place"`
}

func (p OrgRolePermissionAccess) ToMap() map[string]bool {
	perms := make(map[string]bool)

	perms[RolePermissionOrgUpdate] = p.OrgUpdate
	perms[RolePermissionRolesManage] = p.RolesManage
	perms[RolePermissionInvitesManage] = p.InvitesManage

	perms[RolePermissionMembersDelete] = p.MembersDelete
	perms[RolePermissionMembersUpdate] = p.MembersUpdate

	perms[RolePermissionPlaceCreate] = p.PlaceCreate
	perms[RolePermissionPlaceDelete] = p.PlaceDelete
	perms[RolePermissionPlaceUpdate] = p.PlaceUpdate

	return perms
}

type OrgRolePermissionDictWithDetails struct {
	OrganizationUpdate OrgRolePermissionDetails `json:"organization.update"`
	RolesManage        OrgRolePermissionDetails `json:"roles.manage"`
	InvitesManage      OrgRolePermissionDetails `json:"invites.manage"`

	MembersDelete OrgRolePermissionDetails `json:"members.delete"`
	MembersUpdate OrgRolePermissionDetails `json:"members.update"`

	PlaceCreate OrgRolePermissionDetails `json:"place.create"`
	PlaceDelete OrgRolePermissionDetails `json:"place.delete"`
	PlaceUpdate OrgRolePermissionDetails `json:"place.update"`
}

func (p OrgRolePermissionDictWithDetails) ToMap() map[string]OrgRolePermissionDetails {
	perms := make(map[string]OrgRolePermissionDetails)

	perms[RolePermissionOrgUpdate] = p.OrganizationUpdate
	perms[RolePermissionRolesManage] = p.RolesManage
	perms[RolePermissionInvitesManage] = p.InvitesManage

	perms[RolePermissionMembersDelete] = p.MembersDelete
	perms[RolePermissionMembersUpdate] = p.MembersUpdate

	perms[RolePermissionPlaceCreate] = p.PlaceCreate
	perms[RolePermissionPlaceDelete] = p.PlaceDelete
	perms[RolePermissionPlaceUpdate] = p.PlaceUpdate

	return perms
}
