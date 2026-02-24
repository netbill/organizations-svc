package errx

import "github.com/netbill/ape"

var (
	ErrorRoleNotFound = ape.DeclareError("ROLE_NOT_FOUND")

	ErrorInvalidOrgRoleRanksOrder = ape.DeclareError("INVALID_ORG_ROLE_RANKS_ORDER")

	ErrorCannotUpdatePermissionsHeadRole = ape.DeclareError("CANNOT_UPDATE_PERMISSIONS_HEAD_ROLE")
	ErrorCannotDeleteHeadRole            = ape.DeclareError("CANNOT_DELETE_HEAD_ROLE")
	ErrorCannotRemoveHeadRoleFromMember  = ape.DeclareError("CANNOT_REMOVE_HEAD_ROLE_FROM_MEMBER")
	ErrorCannotAddHeadRoleToMember       = ape.DeclareError("CANNOT_ADD_HEAD_ROLE_TO_MEMBER")
)
