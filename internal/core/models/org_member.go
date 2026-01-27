package models

import (
	"time"

	"github.com/google/uuid"
)

type Member struct {
	ID             uuid.UUID `json:"id"`
	AccountID      uuid.UUID `json:"account_id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Position       *string   `json:"position,omitempty"`
	Label          *string   `json:"label,omitempty"`

	Username  string  `json:"username"`
	Official  bool    `json:"official"`
	Pseudonym *string `json:"pseudonym,omitempty"`
	Icon      *string `json:"icon,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (m Member) IsNil() bool {
	return m.ID == uuid.Nil
}
