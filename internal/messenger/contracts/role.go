package contracts

import (
	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/models"
)

const RoleCreatedEvent = "role.created"

type RoleCreatedPayload struct {
	Role models.Role `json:"role"`
}

const RoleUpdatedEvent = "role.updated"

type RoleUpdatedPayload struct {
	Role models.Role `json:"role"`
}

const RoleDeletedEvent = "role.deleted"

type RoleDeletedPayload struct {
	Role models.Role `json:"role"`
}

const RolesRanksUpdatedEvent = "roles.ranks.updated"

type RolesRanksUpdatedPayload struct {
	OrganizationID uuid.UUID          `json:"organization_id"`
	Ranks          map[uuid.UUID]uint `json:"ranks"`
}

const RolePermissionsUpdatedEvent = "role.permissions.updated"

type RolePermissionsUpdatedPayload struct {
	RoleID      uuid.UUID                  `json:"role_id"`
	Permissions map[models.Permission]bool `json:"permissions"`
}
