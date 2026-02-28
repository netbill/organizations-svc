package errx

import "github.com/netbill/ape"

var (
	ErrorProfileNotFound      = ape.DeclareError("PROFILE_NOT_FOUND")
	ErrorProfileDeleted       = ape.DeclareError("PROFILE_DELETED")
	ErrorProfileAlreadyExists = ape.DeclareError("PROFILE_ALREADY_EXISTS")
)
