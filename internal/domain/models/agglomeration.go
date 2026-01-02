package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	AgglomerationStatusActive    = "active"
	AgglomerationStatusInactive  = "inactive"
	AgglomerationStatusSuspended = "suspended"
)

type Agglomeration struct {
	ID       uuid.UUID `json:"id"`
	Status   string    `json:"status"`
	Name     string    `json:"name"`
	Icon     string    `json:"icon"`
	MaxRoles uint      `json:"max_roles"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (e Agglomeration) IsNil() bool {
	return e.ID == uuid.Nil
}
