package contracts

import (
	"github.com/umisto/agglomerations-svc/internal/domain/models"
)

const ProfileUpdatedEvent = "profile.updated"

type ProfileUpdatedPayload struct {
	Profile models.Profile `json:"profile"`
}
