package contracts

import (
	"github.com/netbill/places-svc/internal/core/models"
)

const OrganizationCreatedEvent = "organization.created"

type OrganizationCreatedPayload struct {
	Organization models.Organization `json:"organization"`
}

const OrganizationUpdatedEvent = "organization.updated"

type OrganizationUpdatedPayload struct {
	Organization models.Organization `json:"organization"`
}

const OrganizationActivatedEvent = "organization.activated"

type OrganizationActivatedPayload struct {
	Organization models.Organization `json:"organization"`
}

const OrganizationDeactivatedEvent = "organization.deactivated"

type OrganizationDeactivatedPayload struct {
	Organization models.Organization `json:"organization"`
}

const OrganizationSuspendedEvent = "organization.suspended"

type OrganizationSuspendedPayload struct {
	Organization models.Organization `json:"organization"`
}

const OrganizationDeletedEvent = "organization.deleted"

type OrganizationDeletedPayload struct {
	Organization models.Organization `json:"organization"`
}
