package errx

import "github.com/netbill/ape"

var (
	ErrorRoleNotFound = ape.DeclareError("ROLE_NOT_FOUND")

	ErrorCannotUpdatePermissionsHeadRole = ape.DeclareError("CANNOT_UPDATE_PERMISSIONS_HEAD_ROLE")
	ErrorCannotDeleteHeadRole            = ape.DeclareError("CANNOT_DELETE_HEAD_ROLE")
	ErrorCannotRemoveHeadRoleFromMember  = ape.DeclareError("CANNOT_REMOVE_HEAD_ROLE_FROM_MEMBER")
	ErrorCannotAddHeadRoleToMember       = ape.DeclareError("CANNOT_ADD_HEAD_ROLE_TO_MEMBER")
	ErrorCannotUpdateHeadRoleRank        = ape.DeclareError("CANNOT_UPDATE_HEAD_ROLE_RANK")
)
