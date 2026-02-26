package repository

import (
	"context"
)

type Repository struct {
	TransactionSql   Transaction
	OrganizationsSql OrganizationsQ
	OrgMembersSql    OrgMembersQ
	OrgInvitesSql    OrgInvitesQ
	ProfilesSql      ProfilesQ
	PlacesSql        PlacesQ
}

type Transaction interface {
	Begin(ctx context.Context, fn func(ctx context.Context) error) error
}

func (r *Repository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.TransactionSql.Begin(ctx, fn)
}
