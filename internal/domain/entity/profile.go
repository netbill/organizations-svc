package entity

import (
	"github.com/google/uuid"
)

type Profile struct {
	AccountID uuid.UUID `json:"account_id"`
	Username  string    `json:"username"`
	Official  bool      `json:"official"`
	Pseudonym *string   `json:"pseudonym"`
}

func (p Profile) IsNil() bool {
	return p.AccountID == uuid.Nil
}
