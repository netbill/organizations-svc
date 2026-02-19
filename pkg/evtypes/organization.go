package evtypes

import (
	"time"

	"github.com/google/uuid"
)

const OrganizationsTopicV1 = "organizations.v1"

const OrganizationCreatedEvent = "organization.created"

type OrganizationCreatedPayload struct {
	OrganizationID uuid.UUID `json:"organization_id"`
	Status         string    `json:"status"`
	Name           string    `json:"name"`
	IconKey        *string   `json:"icon_key,omitempty"`
	BannerKey      *string   `json:"banner_key,omitempty"`
	MaxRoles       uint      `json:"max_roles"`
	CreatedAt      time.Time `json:"created_at"`
}

const OrganizationUpdatedEvent = "organization.updated"

type OrganizationUpdatedPayload struct {
	OrganizationID uuid.UUID `json:"organization_id"`
	Status         string    `json:"status"`
	Name           string    `json:"name"`
	IconKey        *string   `json:"icon_key,omitempty"`
	BannerKey      *string   `json:"banner_key,omitempty"`
	MaxRoles       uint      `json:"max_roles"`
	UpdatedAt      time.Time `json:"updated_at"`
}

const OrganizationDeletedEvent = "organization.deleted"

type OrganizationDeletedPayload struct {
	OrganizationID uuid.UUID `json:"organization_id"`
	DeletedAt      time.Time `json:"deleted_at"`
}

const OrganizationActivatedEvent = "organization.status.activated"

type OrganizationActivatedPayload struct {
	OrganizationID uuid.UUID `json:"organization_id"`
	ActivatedAt    time.Time `json:"activated_at"`
}

const OrganizationDeactivatedEvent = "organization.status.deactivated"

type OrganizationDeactivatedPayload struct {
	OrganizationID uuid.UUID `json:"organization_id"`
	DeactivatedAt  time.Time `json:"deactivated_at"`
}
