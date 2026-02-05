package models

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Rank           uint      `json:"rank"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Color          string    `json:"color"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (r Role) IsNil() bool {
	return r.ID == uuid.Nil
}

type OrgMemberRolesLink struct {
	MemberID  uuid.UUID `json:"member_id"`
	RoleID    uuid.UUID `json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (r OrgMemberRolesLink) IsNil() bool {
	return r.MemberID == uuid.Nil
}
