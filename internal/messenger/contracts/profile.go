package contracts

import (
	"time"

	"github.com/google/uuid"
)

const ProfilesTopicV1 = "profiles.v1"

const ProfileUpdatedEvent = "profile.updated"

type ProfileUpdatedPayload struct {
	AccountID   uuid.UUID `json:"account_id"`
	Username    string    `json:"username"`
	Official    bool      `json:"official"`
	Pseudonym   *string   `json:"pseudonym,omitempty"`
	Description *string   `json:"description,omitempty"`

	UpdatedAt time.Time `json:"updated_at"`
}

const ProfileCreatedEvent = "profile.created"

type ProfileCreatedPayload struct {
	AccountID uuid.UUID `json:"account_id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

const ProfileDeletedEvent = "profile.deleted"

type ProfileDeletedPayload struct {
	AccountID uuid.UUID `json:"account_id"`
	DeletedAt time.Time `json:"deleted_at"`
}
