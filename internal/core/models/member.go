package models

import (
	"time"

	"github.com/google/uuid"
)

type Member struct {
	ID             uuid.UUID `json:"id"`
	AccountID      uuid.UUID `json:"account_id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Head           bool      `json:"head"`
	Position       *string   `json:"position,omitempty"`
	Label          *string   `json:"label,omitempty"`

	Username  string  `json:"username"`
	Official  bool    `json:"official"`
	Pseudonym *string `json:"pseudonym,omitempty"`
	AvatarKey *string `json:"avatar_key,omitempty"`

	Version   int32     `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
