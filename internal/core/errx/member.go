package errx

import "github.com/netbill/ape"

var (
	ErrorMemberNotFound = ape.DeclareError("MEMBER_NOT_FOUND")
	ErrorMemberDeleted  = ape.DeclareError("MEMBER_DELETED")

	ErrorCannotDeleteSelf = ape.DeclareError("CANNOT_DELETE_SELF")

	ErrorCannotDeleteOrganizationHeadMember = ape.DeclareError("CANNOT_DELETE_ORGANIZATION_HEAD_MEMBER")
)
