package errx

import "github.com/umisto/ape"

var (
	ErrorAgglomerationNotFound    = ape.DeclareError("AGGLOMERATION_NOT_FOUND")
	ErrorAgglomerationIsSuspended = ape.DeclareError("AGGLOMERATION_IS_SUSPENDED")
	ErrorAgglomerationInactive    = ape.DeclareError("AGGLOMERATION_INACTIVE")
)
