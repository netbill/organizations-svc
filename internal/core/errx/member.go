package errx

import "github.com/netbill/ape"

var (
	ErrorMemberNotFound = ape.DeclareError("MEMBER_NOT_FOUND")

	ErrorCannotDeleteOrganizationHeadMember = ape.DeclareError("CANNOT_DELETE_ORGANIZATION_HEAD_MEMBER")
)
