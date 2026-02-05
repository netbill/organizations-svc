package models

import (
	"time"

	"github.com/google/uuid"
)

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

type OrgRolePermission struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

func (p OrgRolePermission) IsNil() bool {
	return p.Code == ""
}

type OrgRolePermissionLink struct {
	RoleID         uuid.UUID `json:"role_id"`
	PermissionCode string    `json:"permission_code"`
	CreatedAt      time.Time `json:"created_at"`
}

func (r OrgRolePermissionLink) IsNil() bool {
	return r.RoleID == uuid.Nil
}

type OrgRolePermissionWithFlag struct {
	Code        string `json:"code"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

type OrgRolePermissionDict struct {
	ManageOrganization bool `json:"manage_organization,omitempty"`
	ManageInvites      bool `json:"manage_invites,omitempty"`
	ManageMembers      bool `json:"manage_members,omitempty"`
	ManageRoles        bool `json:"manage_roles,omitempty"`
}

type OrgRolePermissionLinks struct {
	ManageOrganization OrgRolePermissionWithFlag
	ManageInvites      OrgRolePermissionWithFlag
	ManageMembers      OrgRolePermissionWithFlag
	ManageRoles        OrgRolePermissionWithFlag
}
