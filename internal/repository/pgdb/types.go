package pgdb

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type RoleDTO struct {
	RoleID uuid.UUID `json:"role_id"`
	Head   bool      `json:"head"`
	Rank   uint      `json:"rank"`
	Name   string    `json:"name"`
}

type RoleDTOs []RoleDTO

func (r *RoleDTOs) Scan(src any) error {
	if src == nil {
		*r = RoleDTOs{}
		return nil
	}

	var b []byte
	switch v := src.(type) {
	case []byte:
		b = v
	case string:
		b = []byte(v)
	default:
		return fmt.Errorf("RoleDTOs: unsupported Scan type %T", src)
	}

	return json.Unmarshal(b, r)
}

type PermissionDTO struct {
	PermissionID uuid.UUID `json:"permission_id"`
	Code         string    `json:"code"`
	Description  string    `json:"description"`
}

type PermissionDTOs []PermissionDTO

func (p *PermissionDTOs) Scan(src any) error {
	if src == nil {
		*p = PermissionDTOs{}
		return nil
	}

	var b []byte
	switch v := src.(type) {
	case []byte:
		b = v
	case string:
		b = []byte(v)
	default:
		return fmt.Errorf("PermissionDTO: unsupported Scan type %T", src)
	}

	return json.Unmarshal(b, p)
}
