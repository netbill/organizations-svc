package repository

import (
	"context"
)

type Repository struct {
	Transactioner
	OrganizationsSql OrganizationsQ
	OrgMembersSql    OrgMembersQ

	OrgRolesSql     OrgRolesQ
	OrgRoleRanksSql OrgRoleRanksQ

	OrgInvitesSql OrgInvitesQ
	ProfilesSql   ProfilesQ
}

type Transactioner interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}
