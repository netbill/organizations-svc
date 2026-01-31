package errx

import "github.com/netbill/ape"

var (
	ErrorOrganizationNotFound = ape.DeclareError("ORGANIZATION_NOT_FOUND")

	ErrorOrganizationIsNotActive = ape.DeclareError("AGLOMERATION_IS_NOT_ACTIVE")
	//ErrorOrganizationInactive    = ape.DeclareError("ORGANIZATION_INACTIVE")
	//
	//ErrorOrganizationPermissionIsInvalid = ape.DeclareError("ORGANIZATION_PERMISSION_IS_INVALID")

	ErrorOrganizationIconTooLarge   = ape.DeclareError("ORGANIZATION_ICON_TOO_LARGE")
	ErrorOrganizationBannerTooLarge = ape.DeclareError("ORGANIZATION_BANNER_TOO_LARGE")

	ErrorOrganizationIconContentTypeNotAllowed   = ape.DeclareError("ORGANIZATION_ICON_CONTENT_TYPE_NOT_ALLOWED")
	ErrorOrganizationBannerContentTypeNotAllowed = ape.DeclareError("ORGANIZATION_BANNER_CONTENT_TYPE_NOT_ALLOWED")

	ErrorOrganizationIconContentFormatNotAllowed   = ape.DeclareError("ORGANIZATION_ICON_CONTENT_FORMAT_NOT_ALLOWED")
	ErrorOrganizationBannerContentFormatNotAllowed = ape.DeclareError("ORGANIZATION_BANNER_CONTENT_FORMAT_NOT_ALLOWED")

	ErrorOrganizationIconResolutionNotAllowed   = ape.DeclareError("ORGANIZATION_ICON_RESOLUTION_NOT_ALLOWED")
	ErrorOrganizationBannerResolutionNotAllowed = ape.DeclareError("ORGANIZATION_BANNER_RESOLUTION_NOT_ALLOWED")
)
