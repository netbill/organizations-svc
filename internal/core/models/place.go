package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

const (
	PlaceStatusActive    = "active"
	PlaceStatusInactive  = "inactive"
	PlaceStatusSuspended = "suspended"
)

type Place struct {
	ID             uuid.UUID `json:"id"`
	ClassID        uuid.UUID `json:"class_id"`
	OrganizationID uuid.UUID `json:"organization_id"`

	Status   string    `json:"status"`
	Verified bool      `json:"verified"`
	Point    orb.Point `json:"point,omitempty"`
	Address  string    `json:"address,omitempty"`
	Name     string    `json:"name"`

	Description *string `json:"description"`
	IconKey     *string `json:"icon_key"`
	BannerKey   *string `json:"banner_key"`
	Website     *string `json:"website"`
	Phone       *string `json:"phone"`

	Version   int32     `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
