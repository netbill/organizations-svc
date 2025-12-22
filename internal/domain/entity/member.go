package entity

import "github.com/google/uuid"

type Member struct {
	ID              uuid.UUID `json:"id"`
	AccountID       uuid.UUID `json:"account_id"`
	AgglomerationID uuid.UUID `json:"agglomeration_id"`
	Position        string    `json:"position"`
	Label           string    `json:"label"`

	Roles []MemberRole `json:"roles"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

type MemberRole struct {
	RoleID uuid.UUID `json:"role_id"`
	Name   string    `json:"name"`
}
