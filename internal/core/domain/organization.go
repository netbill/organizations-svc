package domain

import (
	"time"

	"github.com/google/uuid"
)

const (
	OrganizationStatusActive   = "active"
	OrganizationStatusInactive = "inactive"
)

type Organization struct {
	ID        uuid.UUID `json:"id"`
	Status    string    `json:"status"`
	Name      string    `json:"name"`
	IconKey   *string   `json:"icon_key,omitempty"`
	BannerKey *string   `json:"banner_key,omitempty"`
	MaxRoles  uint      `json:"max_roles"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UploadOrgMediaLinks struct {
	Icon   UploadMediaLink `json:"icon"`
	Banner UploadMediaLink `json:"banner"`
}

type UploadMediaLink struct {
	Key        string `json:"key"`
	UploadURL  string `json:"upload_url"`
	PreloadUrl string `json:"preload_url"`
}
