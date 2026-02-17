package models

import (
	"time"

	"github.com/google/uuid"
)

type AccountActor = uuid.UUID

type UploadScope = uuid.UUID

type Profile struct {
	AccountID uuid.UUID `json:"account_id"`
	Username  string    `json:"username"`
	Official  bool      `json:"official"`
	Pseudonym *string   `json:"pseudonym,omitempty"`
	AvatarKey *string   `json:"avatar_key,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
