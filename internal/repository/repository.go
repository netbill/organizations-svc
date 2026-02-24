package repository

import (
	"context"
)

type Repository struct {
	Transactioner
	OrganizationsSql OrganizationsQ
	OrgMembersSql    OrgMembersQ

	OrgInvitesSql OrgInvitesQ
	ProfilesSql   ProfilesQ
}

type Transactioner interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}
