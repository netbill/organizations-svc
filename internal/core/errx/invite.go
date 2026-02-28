package errx

import "github.com/netbill/ape"

var (
	ErrorInviteNotFound        = ape.DeclareError("INVITE_NOT_FOUND")
	ErrorInviteDeleted         = ape.DeclareError("INVITE_DELETED")
	ErrorInviteNotForInitiator = ape.DeclareError("INVITE_NOT_FOR_INITIATOR")

	ErrorInviteAlreadyAnswered = ape.DeclareError("INVITE_ALREADY_ANSWERED")
	ErrorInviteExpired         = ape.DeclareError("INVITE_EXPIRED")

	ErrorAccountAlreadyMember = ape.DeclareError("ACCOUNT_ALREADY_MEMBER")
)
