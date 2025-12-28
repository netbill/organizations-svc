package entity

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID              uuid.UUID `json:"id"`
	AgglomerationID uuid.UUID `json:"agglomeration_id"`
	Head            bool      `json:"head"`
	Editable        bool      `json:"editable"`
	Rank            uint      `json:"rank"`
	Name            string    `json:"name"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
