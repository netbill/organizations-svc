package models

import (
	"time"

	"github.com/google/uuid"
)

type AccountActor = uuid.UUID

type Profile struct {
	AccountID uuid.UUID `json:"account_id"`
	Username  string    `json:"username"`
	Pseudonym *string   `json:"pseudonym,omitempty"`
	AvatarKey *string   `json:"avatar_key,omitempty"`

	Version   int32     `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
