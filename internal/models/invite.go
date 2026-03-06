package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	InviteStatusSent      = "sent"
	InviteStatusAccepted  = "accepted"
	InviteStatusDeclined  = "declined"
	InviteStatusCancelled = "cancelled"
)

type Invite struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	AccountID      uuid.UUID `json:"account_id"`
	Status         string    `json:"status"`

	UpdatedAt time.Time `json:"updated_at"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}
