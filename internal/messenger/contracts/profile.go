package contracts

import (
	"time"

	"github.com/google/uuid"
)

const ProfileUpdatedEvent = "profile.updated"

type ProfileUpdatedPayload struct {
	Profile struct {
		AccountID   uuid.UUID `json:"account_id"`
		Username    string    `json:"username"`
		Official    bool      `json:"official"`
		Pseudonym   *string   `json:"pseudonym,omitempty"`
		Description *string   `json:"description,omitempty"`
		Avatar      *string   `json:"avatar,omitempty"`

		UpdatedAt time.Time `json:"updated_at"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"profile"`
}
