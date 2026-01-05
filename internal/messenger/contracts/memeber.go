package contracts

import "github.com/netbill/places-svc/internal/core/models"

const MemberCreatedEvent = "member.created"

type MemberCreatedPayload struct {
	Member models.Member `json:"member"`
}

const MemberUpdatedEvent = "member.updated"

type MemberUpdatedPayload struct {
	Member models.Member `json:"member"`
}

const MemberDeletedEvent = "member.deleted"

type MemberDeletedPayload struct {
	Member models.Member `json:"member"`
}
