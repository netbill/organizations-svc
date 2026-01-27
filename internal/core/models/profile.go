package models

import (
	"time"

	"github.com/google/uuid"
)

type Profile struct {
	AccountID uuid.UUID `json:"account_id"`
	Username  string    `json:"username"`
	Pseudonym *string   `json:"pseudonym,omitempty"`
	Avatar    *string   `json:"avatar,omitempty"`
	Official  bool      `json:"official"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (p Profile) IsNil() bool {
	return p.AccountID == uuid.Nil
}

func (p Profile) GetAccountID() uuid.UUID {
	return p.AccountID
}

func (p Profile) GetUsername() string {
	return p.Username
}

func (p Profile) GetOfficial() bool {
	return p.Official
}

func (p Profile) GetPseudonym() *string {
	return p.Pseudonym
}
