package models

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID              uuid.UUID `json:"id"`
	AgglomerationID uuid.UUID `json:"agglomeration_id"`
	Head            bool      `json:"head"`
	Rank            uint      `json:"rank"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Color           string    `json:"color"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (r Role) IsNil() bool {
	return r.ID == uuid.Nil
}
