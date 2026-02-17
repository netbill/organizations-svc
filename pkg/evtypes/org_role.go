package evtypes

import (
	"time"

	"github.com/google/uuid"
)

const OrgRoleCreatedEvent = "role.created"

type OrgRoleCreatedPayload struct {
	RoleID         uuid.UUID `json:"role_id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Rank           uint      `json:"rank"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Color          string    `json:"color"`
	CreatedAt      time.Time `json:"created_at"`
}

const OrgRoleUpdatedEvent = "role.updated"

type OrgRoleUpdatedPayload struct {
	RoleID      uuid.UUID `json:"role_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Color       string    `json:"color"`
	UpdatedAt   time.Time `json:"updated_at"`
}

const OrgRoleDeletedEvent = "role.deleted"

type OrgRoleDeletedPayload struct {
	RoleID    uuid.UUID `json:"role_id"`
	DeletedAt time.Time `json:"deleted_at"`
}

const OrgRolesRanksUpdatedEvent = "roles.ranks.updated"

type OrgRolesRanksUpdatedPayload struct {
	OrganizationID uuid.UUID          `json:"organization_id"`
	Ranks          map[uuid.UUID]uint `json:"ranks"`
	UpdatedAt      time.Time          `json:"updated_at"`
}

const OrgRolePermissionsUpdatedEvent = "role.permissions.updated"

type OrgRolePermissionsUpdatedPayload struct {
	RoleID      uuid.UUID          `json:"role_id"`
	Permissions map[uuid.UUID]bool `json:"permissions"`
}

const OrgMemberRoleAddedEvent = "member_role.added"

type OrgMemberRoleAddedPayload struct {
	MemberID uuid.UUID `json:"member_id"`
	RoleID   uuid.UUID `json:"role_id"`
	AddedAt  time.Time `json:"added_at"`
}

const OrgMemberRoleRemovedEvent = "member_role.remove"

type OrgMemberRoleRemovedPayload struct {
	MemberID  uuid.UUID `json:"member_id"`
	RoleID    uuid.UUID `json:"role_id"`
	RemovedAt time.Time `json:"removed_at"`
}
