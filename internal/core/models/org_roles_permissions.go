package models

type OrgRolePermissionCode string

const (
	RolePermissionOrganizationUpdate = "organization.update"
	RolePermissionRolesManage        = "roles.manage"
	RolePermissionInvitesManage      = "invites.manage"

	RolePermissionMembersDelete = "members.delete"
	RolePermissionMembersUpdate = "members.update"
)

var allRolePermissions = []string{
	RolePermissionOrganizationUpdate,
	RolePermissionRolesManage,
	RolePermissionInvitesManage,
	RolePermissionMembersDelete,
	RolePermissionMembersUpdate,
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
	OrganizationUpdate bool `json:"organization.update"`
	RolesManage        bool `json:"roles.manage"`
	InvitesManage      bool `json:"invites.manage"`
	MembersDelete      bool `json:"members.delete"`
	MembersUpdate      bool `json:"members.update"`
}

func (p OrgRolePermissionAccess) ToMap() map[string]bool {
	perms := make(map[string]bool)

	perms[RolePermissionOrganizationUpdate] = p.OrganizationUpdate
	perms[RolePermissionRolesManage] = p.RolesManage
	perms[RolePermissionInvitesManage] = p.InvitesManage
	perms[RolePermissionMembersDelete] = p.MembersDelete
	perms[RolePermissionMembersUpdate] = p.MembersUpdate

	return perms
}

type OrgRolePermissionDictWithDetails struct {
	OrganizationUpdate OrgRolePermissionDetails `json:"organization.update"`
	RolesManage        OrgRolePermissionDetails `json:"roles.manage"`
	InvitesManage      OrgRolePermissionDetails `json:"invites.manage"`
	MembersDelete      OrgRolePermissionDetails `json:"members.delete"`
	MembersUpdate      OrgRolePermissionDetails `json:"members.update"`
}

func (p OrgRolePermissionDictWithDetails) ToMap() map[string]OrgRolePermissionDetails {
	perms := make(map[string]OrgRolePermissionDetails)

	perms[RolePermissionOrganizationUpdate] = p.OrganizationUpdate
	perms[RolePermissionRolesManage] = p.RolesManage
	perms[RolePermissionInvitesManage] = p.InvitesManage
	perms[RolePermissionMembersDelete] = p.MembersDelete
	perms[RolePermissionMembersUpdate] = p.MembersUpdate

	return perms
}
