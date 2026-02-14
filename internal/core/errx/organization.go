package errx

import "github.com/netbill/ape"

var (
	ErrorOrganizationNotFound = ape.DeclareError("ORGANIZATION_NOT_FOUND")

	ErrorOrganizationIsNotActive = ape.DeclareError("AGLOMERATION_IS_NOT_ACTIVE")

	ErrorOrganizationIconInvalid   = ape.DeclareError("AGLOMERATION_ICON_INVALID")
	ErrorOrganizationBannerInvalid = ape.DeclareError("AGLOMERATION_BANNER_INVALID")

	ErrorNoContentUploaded = ape.DeclareError("NO_CONTENT_UPLOADED")
)
