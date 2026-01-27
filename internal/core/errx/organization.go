package errx

import "github.com/netbill/ape"

var (
	ErrorOrganizationNotFound = ape.DeclareError("ORGANIZATION_NOT_FOUND")
	//
	//ErrorOrganizationIsNotActive = ape.DeclareError("AGLOMERATION_IS_NOT_ACTIVE")
	//ErrorOrganizationInactive    = ape.DeclareError("ORGANIZATION_INACTIVE")
	//
	//ErrorOrganizationPermissionIsInvalid = ape.DeclareError("ORGANIZATION_PERMISSION_IS_INVALID")
)
