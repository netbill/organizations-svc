package errx

import "github.com/netbill/ape"

var (
	ErrorOrganizationNotExists   = ape.DeclareError("ORGANIZATION_NOT_EXISTS")
	ErrorOrganizationDeleted     = ape.DeclareError("ORGANIZATION_DELETED")
	ErrorOrganizationIsSuspended = ape.DeclareError("ORGANIZATION_IS_SUSPENDED")

	ErrorOrganizationUploadedIconInvalid   = ape.DeclareError("ORGANIZATION_UPLOADED_ICON_INVALID")
	ErrorOrganizationUploadedBannerInvalid = ape.DeclareError("ORGANIZATION_UPLOADED_BANNER_INVALID")
)
