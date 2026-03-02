package errx

import "github.com/netbill/ape"

var (
	ErrorPlaceNotExists     = ape.DeclareError("PLACE_NOT_EXISTS")
	ErrorPlaceAlreadyExists = ape.DeclareError("PLACE_ALREADY_EXISTS")
	ErrorPlaceDeleted       = ape.DeclareError("PLACE_DELETED")

	ErrorOrganizationHavePlace = ape.DeclareError("ORGANIZATION_HAVE_PLACE")
)
