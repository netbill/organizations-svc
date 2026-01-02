package errx

import "github.com/umisto/ape"

var (
	ErrorInviteNotFound = ape.DeclareError("INVITE_NOT_FOUND")

	ErrorInviteAlreadyAnswered = ape.DeclareError("INVITE_ALREADY_ANSWERED")
	ErrorInviteExpired         = ape.DeclareError("INVITE_EXPIRED")
)
