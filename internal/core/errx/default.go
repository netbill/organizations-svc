package errx

import "github.com/netbill/ape"

var (
	ErrorNotEnoughRights = ape.DeclareError("NOT_ENOUGH_RIGHTS")

	ErrorInitiatorNotMemberOfOrganization = ape.DeclareError("INITIATOR_NOT_MEMBER_OF_ORGANIZATION")
)
