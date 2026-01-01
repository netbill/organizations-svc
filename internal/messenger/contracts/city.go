package contracts

import (
	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/models"
)

const CityCreatedEvent = "city.created"

type CityCreatedPayload struct {
	City models.City `json:"city"`
}

const CityUpdatedEvent = "city.updated"

type CityUpdatedPayload struct {
	City models.City `json:"city"`
}

const CitySlugUpdatedEvent = "city.slug.updated"

type CitySlugUpdatedPayload struct {
	City    models.City `json:"city"`
	OldSlug *string     `json:"old_slug"`
}

const CityAgglomerationUpdatedEvent = "city.agglomeration.updated"

type CityAgglomerationUpdatedPayload struct {
	City               models.City `json:"city"`
	OldAgglomerationID *uuid.UUID  `json:"old_agglomeration_id"`
}

const CityDeletedEvent = "city.deleted"

type CityDeletedPayload struct {
	City models.City `json:"city"`
}
