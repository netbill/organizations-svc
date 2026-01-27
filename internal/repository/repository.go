package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/netbill/organizations-svc/internal/repository/pgdb"
	"github.com/netbill/pgxtx"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) Repository {
	return Repository{pool: pool}
}

func (r Repository) organizationsQ(ctx context.Context) pgdb.OrganizationsQ {
	return pgdb.NewOrganizationsQ(pgxtx.Exec(r.pool, ctx))
}

func (r Repository) orgMembersQ(ctx context.Context) pgdb.OrgMembersQ {
	return pgdb.NewOrgMembersQ(pgxtx.Exec(r.pool, ctx))
}

func (r Repository) orgMemberRolesQ(ctx context.Context) pgdb.OrgMemberRolesQ {
	return pgdb.NewOrgMemberRolesQ(pgxtx.Exec(r.pool, ctx))
}

func (r Repository) orgRolesQ(ctx context.Context) pgdb.OrgRolesQ {
	return pgdb.NewOrgRolesQ(pgxtx.Exec(r.pool, ctx))
}

func (r Repository) orgRolePermissionLinksQ(ctx context.Context) pgdb.OrgRolePermissionLinksQ {
	return pgdb.NewOrgRolePermissionLinksQ(pgxtx.Exec(r.pool, ctx))
}

func (r Repository) orgRolePermissionsQ(ctx context.Context) pgdb.OrgRolePermissionsQ {
	return pgdb.NewOrgRolePermissionsQ(pgxtx.Exec(r.pool, ctx))
}

func (r Repository) orgInvitesQ(ctx context.Context) pgdb.OrgInvitesQ {
	return pgdb.NewOrgInvitesQ(pgxtx.Exec(r.pool, ctx))
}

func (r Repository) profilesQ(ctx context.Context) pgdb.ProfilesQ {
	return pgdb.NewProfilesQ(pgxtx.Exec(r.pool, ctx))
}

func (r Repository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return pgxtx.Transaction(r.pool, ctx, fn)
}
