package errx

import "github.com/umisto/ape"

var (
	ErrorCityNotFound = ape.DeclareError("CITY_NOT_FOUND")

	ErrorCityIsSuspended = ape.DeclareError("CITY_IS_SUSPENDED")
	ErrorCityIsNotActive = ape.DeclareError("CITY_IS_NOT_ACTIVE")

	ErrorCityWithSlugAlreadyExists = ape.DeclareError("CITY_WITH_SLUG_ALREADY_EXISTS")
)
