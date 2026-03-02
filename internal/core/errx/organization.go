package errx

import "github.com/netbill/ape"

var (
	ErrorOrganizationNotExists   = ape.DeclareError("ORGANIZATION_NOT_EXISTS")
	ErrorOrganizationDeleted     = ape.DeclareError("ORGANIZATION_DELETED")
	ErrorOrganizationIsNotActive = ape.DeclareError("AGLOMERATION_IS_NOT_ACTIVE")
	ErrorOrganizationIsSuspended = ape.DeclareError("ORGANIZATION_IS_SUSPENDED")

	ErrorNoContentUploaded = ape.DeclareError("NO_CONTENT_UPLOADED")

	ErrorOrganizationIconKeyIsInvalid        = ape.DeclareError("ORGANIZATION_ICON_KEY_IS_INVALID")
	ErrorOrganizationIconContentIsExceedsMax = ape.DeclareError("ORGANIZATION_ICON_CONTENT_EXCEEDS_MAX")
	ErrorOrganizationIconResolutionIsInvalid = ape.DeclareError("ORGANIZATION_ICON_RESOLUTION_IS_INVALID")
	ErrorOrganizationIconFormatIsNotAllowed  = ape.DeclareError("ORGANIZATION_ICON_FORMAT_IS_NOT_ALLOWED")

	ErrorOrganizationBannerKeyIsInvalid        = ape.DeclareError("ORGANIZATION_BANNER_KEY_IS_INVALID")
	ErrorOrganizationBannerContentIsExceedsMax = ape.DeclareError("ORGANIZATION_BANNER_CONTENT_EXCEEDS_MAX")
	ErrorOrganizationBannerResolutionIsInvalid = ape.DeclareError("ORGANIZATION_BANNER_RESOLUTION_IS_INVALID")
	ErrorOrganizationBannerFormatIsNotAllowed  = ape.DeclareError("ORGANIZATION_BANNER_FORMAT_IS_NOT_ALLOWED")
)
