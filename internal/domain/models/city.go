package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

type City struct {
	ID              uuid.UUID  `json:"id"`
	AgglomerationID *uuid.UUID `json:"agglomeration_id,omitempty"`
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
