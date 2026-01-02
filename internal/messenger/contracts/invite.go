package contracts

import "github.com/netbill/organizations-svc/internal/core/models"

const InviteCreatedEvent = "invite.created"

type InviteCreatedPayload struct {
	Invite models.Invite `json:"invite"`
}

const InviteAcceptedEvent = "invite.accepted"

type InviteAcceptedPayload struct {
	Invite models.Invite `json:"invite"`
}

const InviteDeclinedEvent = "invite.declined"

type InviteDeclinedPayload struct {
	Invite models.Invite `json:"invite"`
}

const InviteDeletedEvent = "invite.deleted"

type InviteDeletedPayload struct {
	Invite models.Invite `json:"invite"`
}
