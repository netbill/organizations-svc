package models

const (
	RolePermissionManageOrganization = "organization.manage"
	RolePermissionManageInvites      = "invites.manage"
	RolePermissionManageMembers      = "members.manage"
	RolePermissionManageRoles        = "roles.manage"
)

var allRolePermissions = []string{
	RolePermissionManageOrganization,
	RolePermissionManageRoles,
	RolePermissionManageInvites,
	RolePermissionManageMembers,
}

type Permission struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

func (p Permission) IsNil() bool {
	return p.Code == ""
}
