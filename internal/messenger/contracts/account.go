package contracts

import (
	"github.com/netbill/organizations-svc/internal/core/models"
)

const AccountCreatedEvent = "account.created"

type AccountCreatedPayload struct {
	Account models.Account
}

const AccountUsernameChangedEvent = "account.username.changed"

type AccountUsernameChangedPayload struct {
	Account models.Account
}

const AccountDeletedEvent = "account.deleted"

//TODO remove AccountID from payload at other services

type AccountDeletedPayload struct {
	Account models.Account
}

const AccountProfileUpdatedEvent = "account.profile.updated"

type AccountProfileUpdatedPayload struct {
	Profile models.Profile `json:"profile"`
}
