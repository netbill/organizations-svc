package errx

import "github.com/netbill/ape"

var (
	ErrorPlaceNotFound         = ape.DeclareError("PLACE_NOT_FOUND")
	ErrorOrganizationHavePlace = ape.DeclareError("ORGANIZATION_HAVE_PLACE")
)
