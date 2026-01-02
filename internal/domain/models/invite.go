package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	InviteStatusSent     = "sent"
	InviteStatusAccepted = "accepted"
	InviteStatusDeclined = "declined"
)

type Invite struct {
	ID              uuid.UUID `json:"id"`
	AgglomerationID uuid.UUID `json:"agglomeration_id"`
	AccountID       uuid.UUID `json:"account_id"`
	Status          string    `json:"status"`
	ExpiresAt       time.Time `json:"expires_at"`
	CreatedAt       time.Time `json:"created_at"`
}

func (i Invite) IsNil() bool {
	return i.ID == uuid.Nil
}
