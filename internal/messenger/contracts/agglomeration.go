package contracts

import (
	"github.com/umisto/cities-svc/internal/domain/models"
)

const AgglomerationCreatedEvent = "agglomeration.created"

type AgglomerationCreatedPayload struct {
	Agglomeration models.Agglomeration `json:"agglomeration"`
}

const AgglomerationUpdatedEvent = "agglomeration.updated"

type AgglomerationUpdatedPayload struct {
	Agglomeration models.Agglomeration `json:"agglomeration"`
}

const AgglomerationActivatedEvent = "agglomeration.activated"

type AgglomerationActivatedPayload struct {
	Agglomeration models.Agglomeration `json:"agglomeration"`
}

const AgglomerationDeactivatedEvent = "agglomeration.deactivated"

type AgglomerationDeactivatedPayload struct {
	Agglomeration models.Agglomeration `json:"agglomeration"`
}

const AgglomerationSuspendedEvent = "agglomeration.suspended"

type AgglomerationSuspendedPayload struct {
	Agglomeration models.Agglomeration `json:"agglomeration"`
}

const AgglomerationDeletedEvent = "agglomeration.deleted"

type AgglomerationDeletedPayload struct {
	Agglomeration models.Agglomeration `json:"agglomeration"`
}
