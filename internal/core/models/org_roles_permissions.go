package models

import (
	"fmt"

	"github.com/netbill/organizations-svc/internal/core/errx"
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

type Permission struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

func (p Permission) IsNil() bool {
	return p.Code == ""
}

func ValidateOrganizationPermission(s string) error {
	for _, e := range allRolePermissions {
		if e == s {
			return nil
		}
	}

	return errx.ErrorOrganizationPermissionIsInvalid.Raise(
		fmt.Errorf("organization permission '%s' is invalid", s),
	)
}
