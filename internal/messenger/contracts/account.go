package contracts

import (
	"time"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/models"
)

const AccountCreatedEvent = "account.created"

type AccountCreatedPayload struct {
	Account struct {
		ID       uuid.UUID `json:"id"`
		Username string    `json:"username"`
		Role     string    `json:"role"`
		Status   string    `json:"status"`

		CreatedAt         time.Time `json:"created_at"`
		UpdatedAt         time.Time `json:"updated_at"`
		UsernameUpdatedAt time.Time `json:"username_name_updated_at"`
	} `json:"account"`
}

const AccountUsernameChangedEvent = "account.username.changed"

type AccountUsernameChangedPayload struct {
	Account struct {
		ID       uuid.UUID `json:"id"`
		Username string    `json:"username"`
		Role     string    `json:"role"`
		Status   string    `json:"status"`

		CreatedAt         time.Time `json:"created_at"`
		UpdatedAt         time.Time `json:"updated_at"`
		UsernameUpdatedAt time.Time `json:"username_name_updated_at"`
	} `json:"account"`
}

const AccountDeletedEvent = "account.deleted"

type AccountDeletedPayload struct {
	AccountID uuid.UUID `json:"account_id"`
	Account   struct {
		ID       uuid.UUID `json:"id"`
		Username string    `json:"username"`
		Role     string    `json:"role"`
		Status   string    `json:"status"`

		CreatedAt         time.Time `json:"created_at"`
		UpdatedAt         time.Time `json:"updated_at"`
		UsernameUpdatedAt time.Time `json:"username_name_updated_at"`
	} `json:"account"`
}

const AccountProfileUpdatedEvent = "account.profile.updated"

type AccountProfileUpdatedPayload struct {
	Profile models.Profile `json:"profile"`
}
