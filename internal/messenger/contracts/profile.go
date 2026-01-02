package contracts

import (
	"github.com/netbill/organizations-svc/internal/domain/models"
)

const ProfileUpdatedEvent = "profile.updated"

type ProfileUpdatedPayload struct {
	Profile models.Profile `json:"profile"`
}
