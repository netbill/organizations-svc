package repository

import (
	"context"
)

type Repository struct {
	OrganizationsSql OrganizationsQ
	OrgMembersSql    OrgMembersQ
	OrgInvitesSql    OrgInvitesQ
	ProfilesSql      ProfilesQ
	PlacesSql        PlacesQ
	TombstonesSql
	TransactionSql
}

type TransactionSql interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}
