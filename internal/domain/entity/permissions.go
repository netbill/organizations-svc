package entity

import (
	"strings"

	"github.com/google/uuid"
)

type CodeRolePermission string

func (p CodeRolePermission) IsNil() bool {
	return p == ""
}

func (p CodeRolePermission) String() string {
	return string(p)
}

func (p CodeRolePermission) split() []string {
	return strings.Split(string(p), ".")
}

const (
	RolePermissionManageAgglomeration CodeRolePermission = "agglomeration.manage"
	RolePermissionManageCities        CodeRolePermission = "cities.manage"
	RolePermissionManageRoles         CodeRolePermission = "roles.manage"
	RolePermissionManageInvites       CodeRolePermission = "invites.manage"
	RolePermissionManageMembers       CodeRolePermission = "members.manage"
)

type Permission struct {
	ID          uuid.UUID          `json:"id"`
	Code        CodeRolePermission `json:"code"`
	Description string             `json:"description"`
}
