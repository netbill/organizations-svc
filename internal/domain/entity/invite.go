package entity

import (
	"time"

	"github.com/google/uuid"
)

type Invite struct {
	ID              uuid.UUID `json:"id"`
	AgglomerationID uuid.UUID `json:"agglomeration_id"`
	Status          string    `json:"status"`
	ExpiresAt       time.Time `json:"expires_at"`
	CreatedAt       time.Time `json:"created_at"`
}

func (i Invite) IsNil() bool {
	return i.ID == uuid.Nil
}
