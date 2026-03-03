package errx

import "github.com/netbill/ape"

var (
	ErrorProfileNotExists     = ape.DeclareError("PROFILE_NOT_EXISTS")
	ErrorProfileDeleted       = ape.DeclareError("PROFILE_DELETED")
	ErrorProfileAlreadyExists = ape.DeclareError("PROFILE_ALREADY_EXISTS")
)
