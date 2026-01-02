package errx

import "github.com/netbill/ape"

var (
	ErrorOrganizationNotFound = ape.DeclareError("ORGANIZATION_NOT_FOUND")

	ErrorOrganizationIsNotActive = ape.DeclareError("AGLOMERATION_IS_NOT_ACTIVE")
	ErrorOrganizationIsSuspended = ape.DeclareError("ORGANIZATION_IS_SUSPENDED")
	ErrorOrganizationInactive    = ape.DeclareError("ORGANIZATION_INACTIVE")
)
