package models

import (
	"github.com/google/uuid"
	"github.com/netbill/orgperm"
)

func GetOrgRolePermissionLength() int {
	return len(orgperm.GetAllPermissions())
}

type OrgRolePermission struct {
	ID          uuid.UUID `json:"id"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
}

type OrgRolePermissionsWithDetailsForRole map[uuid.UUID]OrgRolePermissionDetails

type OrgRolePermissionDetails struct {
	Code        string `json:"code"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}
