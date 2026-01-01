package models

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
	RolePermissionManageAgglomeration  CodeRolePermission = "agglomeration.manage"
	RolePermissionManageagglomerations CodeRolePermission = "agglomerations.manage"
	RolePermissionManageRoles          CodeRolePermission = "roles.manage"
	RolePermissionManageInvites        CodeRolePermission = "invites.manage"
	RolePermissionManageMembers        CodeRolePermission = "members.manage"
)

var AllRolePermissions = []CodeRolePermission{
	RolePermissionManageAgglomeration,
	RolePermissionManageagglomerations,
	RolePermissionManageRoles,
	RolePermissionManageInvites,
	RolePermissionManageMembers,
}

type Permission struct {
	ID          uuid.UUID          `json:"id"`
	Code        CodeRolePermission `json:"code"`
	Description string             `json:"description"`
}

func (p Permission) IsNil() bool {
	return p.ID == uuid.Nil
}
