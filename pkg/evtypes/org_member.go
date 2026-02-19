package evtypes

import (
	"time"

	"github.com/google/uuid"
)

const OrgMembersTopicV1 = "organization.members.v1"

const OrgMemberCreatedEvent = "member.created"

type OrgMemberCreatedPayload struct {
	MemberID       uuid.UUID `json:"member_id"`
	AccountID      uuid.UUID `json:"account_id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Head           bool      `json:"head"`
	Position       *string   `json:"position,omitempty"`
	Label          *string   `json:"label,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

const OrgMemberUpdatedEvent = "member.updated"

type OrgMemberUpdatedPayload struct {
	MemberID  uuid.UUID `json:"member_id"`
	Position  *string   `json:"position,omitempty"`
	Label     *string   `json:"label,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
}

const OrgMemberDeletedEvent = "member.deleted"

type OrgMemberDeletedPayload struct {
	MemberID  uuid.UUID `json:"member_id"`
	DeletedAt time.Time `json:"deleted_at"`
}
