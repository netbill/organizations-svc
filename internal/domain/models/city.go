package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

const (
	CityStatusActive   = "active"
	CityStatusInactive = "inactive"
	CityStatusArchived = "archived"
)

type City struct {
	ID              uuid.UUID  `json:"id"`
	AgglomerationID *uuid.UUID `json:"agglomeration_id,omitempty"`
	Status          string     `json:"status"`
	Slug            *string    `json:"slug,omitempty"`
	Name            string     `json:"name"`
	Icon            *string    `json:"icon,omitempty"`
	Banner          *string    `json:"banner,omitempty"`

	Point orb.Point `json:"point"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (e City) IsNil() bool {
	return e.ID == uuid.Nil
}

func (e City) CanInteract() bool {
	return e.Status == CityStatusActive
}
