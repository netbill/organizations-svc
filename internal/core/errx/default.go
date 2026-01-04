package errx

import "github.com/netbill/ape"

var ErrorInternal = ape.DeclareError("INTERNAL_ERROR")

var (
	ErrorNotEnoughRights     = ape.DeclareError("NOT_ENOUGH_RIGHTS")
	ErrorNotAccessToResource = ape.DeclareError("NOT_ACCESS_TO_RESOURCE")
	ErrorInvalidInput        = ape.DeclareError("INVALID_INPUT")
)
