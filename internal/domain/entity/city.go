package entity

import (
	"time"

	"github.com/google/uuid"
)

const (
	CityStatusActive   = "active"
	CityStatusInactive = "inactive"
)

type City struct {
	ID              uuid.UUID `json:"id"`
	AgglomerationID uuid.UUID `json:"agglomeration_id"`
	Status          string    `json:"status"`
	Slug            string    `json:"slug"`
	Name            string    `json:"name"`
	Icon            *string   `json:"icon,omitempty"`
	Banner          *string   `json:"banner,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (e City) IsNil() bool {
	return e.ID == uuid.Nil
}

func (e City) CanInteract() bool {
	return e.Status == CityStatusActive
}
