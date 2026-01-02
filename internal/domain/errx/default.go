package errx

import "github.com/umisto/ape"

var ErrorInternal = ape.DeclareError("INTERNAL_ERROR")

var (
	ErrorNotEnoughRights     = ape.DeclareError("NOT_ENOUGH_RIGHTS")
	ErrorNotAccessToResource = ape.DeclareError("NOT_ACCESS_TO_RESOURCE")
)
