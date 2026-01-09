package repository

import (
	"context"
	"database/sql"

	"github.com/netbill/organizations-svc/internal/repository/pgdb"
	"github.com/netbill/pgx"
	replicaspg "github.com/netbill/replicas/pgdb"
)

type Service struct {
	db *sql.DB
}

func New(db *sql.DB) Service {
	return Service{db: db}
}

func (s Service) organizationsQ(ctx context.Context) pgdb.OrganizationsQ {
	return pgdb.NewOrganizationsQ(pgx.Exec(s.db, ctx))
}

func (s Service) membersQ(ctx context.Context) pgdb.OrganizationMembersQ {
	return pgdb.NewOrganizationMembersQ(pgx.Exec(s.db, ctx))
}

func (s Service) memberRolesQ(ctx context.Context) pgdb.OrganizationMemberRolesQ {
	return pgdb.NewOrganizationMemberRolesQ(pgx.Exec(s.db, ctx))
}

func (s Service) rolesQ(ctx context.Context) pgdb.OrganizationRolesQ {
	return pgdb.NewRolesQ(pgx.Exec(s.db, ctx))
}

func (s Service) rolePermissionLinksQ(ctx context.Context) pgdb.OrganizationRolePermissionLinksQ {
	return pgdb.NewRolePermissionsQ(pgx.Exec(s.db, ctx))
}

func (s Service) rolePermissionsQ(ctx context.Context) pgdb.OrganizationRolePermissionsQ {
	return pgdb.NewOrganizationPermissionsQ(pgx.Exec(s.db, ctx))
}

func (s Service) invitesQ(ctx context.Context) pgdb.OrganizationInvitesQ {
	return pgdb.NewOrganizationInvitesQ(pgx.Exec(s.db, ctx))
}

func (s Service) profilesQ(ctx context.Context) replicaspg.ProfilesQ {
	return replicaspg.NewProfilesQ(pgx.Exec(s.db, ctx))
}

func (s Service) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return pgx.Transaction(s.db, ctx, fn)
}
