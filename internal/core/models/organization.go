package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	OrganizationStatusActive   = "active"
	OrganizationStatusInactive = "inactive"
)

type Organization struct {
	ID       uuid.UUID `json:"id"`
	Status   string    `json:"status"`
	Name     string    `json:"name"`
	Icon     *string   `json:"icon,omitempty"`
	Banner   *string   `json:"banner,omitempty"`
	MaxRoles uint      `json:"max_roles"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (e Organization) IsNil() bool {
	return e.ID == uuid.Nil
}

type OrganizationUploadMediaLinks struct {
	IconUploadLink   string `json:"icon_upload_link"`
	IconGetLink      string `json:"icon_get_link"`
	BannerUploadLink string `json:"banner_upload_link"`
	BannerGetLink    string `json:"banner_get_link"`
	UploadToken      string `json:"upload_token"`
}
