package errx

import "github.com/umisto/ape"

var ErrorInternal = ape.DeclareError("INTERNAL_ERROR")

var ErrorNotEnoughRightsForAgglomeration = ape.DeclareError("NOT_ENOUGH_RIGHTS_FOR_AGGLOMERATION")
