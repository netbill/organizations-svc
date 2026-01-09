package repository

import (
	"context"
	"database/sql"

	"github.com/netbill/organizations-svc/internal/repository/pgdb"
	"github.com/netbill/pgx"
)

type Service struct {
	db *sql.DB
}

func New(db *sql.DB) Service {
	return Service{
		db: db,
	}
}

func (s Service) organizationsQ(ctx context.Context) pgdb.OrganizationsQ {
	return pgdb.NewOrganizationsQ(pgx.Exec(s.db, ctx))
}

func (s Service) membersQ(ctx context.Context) pgdb.OrgMembersQ {
	return pgdb.NewOrgMembersQ(pgx.Exec(s.db, ctx))
}

func (s Service) memberRolesQ(ctx context.Context) pgdb.OrgMemberRolesQ {
	return pgdb.NewOrgMemberRolesQ(pgx.Exec(s.db, ctx))
}

func (s Service) rolesQ(ctx context.Context) pgdb.OrgRolesQ {
	return pgdb.NewOrgRolesQ(pgx.Exec(s.db, ctx))
}

func (s Service) rolePermissionLinksQ(ctx context.Context) pgdb.OrgRolePermissionLinksQ {
	return pgdb.NewOrgRolePermissionLinksQ(pgx.Exec(s.db, ctx))
}

func (s Service) rolePermissionsQ(ctx context.Context) pgdb.OrgRolePermissionsQ {
	return pgdb.NewOrgRolePermissionsQ(pgx.Exec(s.db, ctx))
}

func (s Service) invitesQ(ctx context.Context) pgdb.OrgInvitesQ {
	return pgdb.NewOrgInvitesQ(pgx.Exec(s.db, ctx))
}

func (s Service) profilesQ(ctx context.Context) pgdb.ProfilesQ {
	return pgdb.NewProfilesQ(pgx.Exec(s.db, ctx))
}

func (s Service) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return pgx.Transaction(s.db, ctx, fn)
}
