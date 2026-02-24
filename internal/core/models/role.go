package models

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Color          string    `json:"color"`
	Version        int32     `json:"version"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Role Rank

type OrgRoleRank struct {
	OrganizationID uuid.UUID `json:"organization_id"`
	RoleID         uuid.UUID `json:"role_id"`
	Rank           int32     `json:"rank"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type OrgRoleRanksRevision struct {
	OrganizationID uuid.UUID `json:"organization_id"`
	Revision       int32     `json:"revision"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Role members links

type OrgMemberRoleLink struct {
	MemberID  uuid.UUID `json:"member_id"`
	RoleID    uuid.UUID `json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
}

type OrgMemberRoleLinkRevision struct {
	MemberID  uuid.UUID `json:"member_id"`
	Revision  int32     `json:"revision"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Role permissions

type OrgRolePermission struct {
	ID          uuid.UUID `json:"id"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
}

type OrgRolePermissionsWithDetailsForRole struct {
	Permissions map[uuid.UUID]OrgRolePermissionDetails `json:"permissions"`
	
	Revision          int32     `json:"revision"`
	RevisionCreatedAt time.Time `json:"revision_created_at"`
	RevisionUpdatedAt time.Time `json:"revision_updated_at"`
}

type OrgRolePermissionDetails struct {
	Code        string `json:"code"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

type OrgRolePermissionLink struct {
	RoleID       uuid.UUID `json:"role_id"`
	PermissionID uuid.UUID `json:"permission_id"`
	CreatedAt    time.Time `json:"created_at"`
}

type OrgRolePermissionsLinksRevision struct {
	RoleID    uuid.UUID `json:"role_id"`
	Revision  int32     `json:"revision"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
