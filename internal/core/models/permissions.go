package models

import (
	"github.com/google/uuid"
)

const (
	RolePermissionManageOrganization = "organization.manage"
	RolePermissionManageInvites      = "invites.manage"
	RolePermissionManageMembers      = "members.manage"
	RolePermissionManageRoles        = "roles.manage"
)

var AllRolePermissions = []string{
	RolePermissionManageOrganization,
	RolePermissionManageRoles,
	RolePermissionManageInvites,
	RolePermissionManageMembers,
}

type Permission struct {
	ID          uuid.UUID `json:"id"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
}

func (p Permission) IsNil() bool {
	return p.ID == uuid.Nil
}
