package models

import (
	"time"

	"github.com/google/uuid"
)

type Initiator interface {
	GetAccountID() uuid.UUID
}

type Profile struct {
	AccountID uuid.UUID `json:"account_id"`
	Username  string    `json:"username"`
	Official  bool      `json:"official"`
	Pseudonym *string   `json:"pseudonym,omitempty"`
	Avatar    *string   `json:"avatar,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (p Profile) IsNil() bool {
	return p.AccountID == uuid.Nil
}
