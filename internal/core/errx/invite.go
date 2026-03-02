package errx

import "github.com/netbill/ape"

var (
	ErrorInviteNotExists       = ape.DeclareError("INVITE_NOT_EXISTS")
	ErrorInviteNotForInitiator = ape.DeclareError("INVITE_NOT_FOR_INITIATOR")

	ErrorActiveInviteAlreadyExists = ape.DeclareError("ACTIVE_INVITE_ALREADY_EXISTS")
	ErrorInviteAlreadyAnswered     = ape.DeclareError("INVITE_ALREADY_ANSWERED")
	ErrorInviteExpired             = ape.DeclareError("INVITE_EXPIRED")

	ErrorAccountAlreadyMember = ape.DeclareError("ACCOUNT_ALREADY_MEMBER")
)
