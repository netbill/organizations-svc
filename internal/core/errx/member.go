package errx

import "github.com/netbill/ape"

var (
	ErrorMemberNotExists = ape.DeclareError("MEMBER_NOT_EXISTS")
	ErrorMemberDeleted   = ape.DeclareError("MEMBER_DELETED")

	ErrorCannotDeleteSelf                   = ape.DeclareError("CANNOT_DELETE_SELF")
	ErrorNotOrganizationHead                = ape.DeclareError("NOT_ORGANIZATION_HEAD")
	ErrorInitiatorNotMemberOfOrganization   = ape.DeclareError("INITIATOR_NOT_MEMBER_OF_ORGANIZATION")
	ErrorCannotDeleteOrganizationHeadMember = ape.DeclareError("CANNOT_DELETE_ORGANIZATION_HEAD_MEMBER")
)
