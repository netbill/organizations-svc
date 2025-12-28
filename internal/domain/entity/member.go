package entity

import (
	"time"

	"github.com/google/uuid"
)

type Member struct {
	ID              uuid.UUID `json:"id"`
	AccountID       uuid.UUID `json:"account_id"`
	AgglomerationID uuid.UUID `json:"agglomeration_id"`
	Position        *string   `json:"position,omitempty"`
	Label           *string   `json:"label,omitempty"`

	Username  string  `json:"username"`
	Pseudonym *string `json:"pseudonym,omitempty"`
	Official  bool    `json:"official"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (m Member) IsNil() bool {
	return m.ID == uuid.Nil
}
