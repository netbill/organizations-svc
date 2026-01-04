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
	return Service{db: db}
}

func (s Service) organizationsQ(ctx context.Context) pgdb.OrganizationsQ {
	return pgdb.NewOrganizationsQ(pgx.Exec(s.db, ctx))
}

func (s Service) membersQ(ctx context.Context) pgdb.MembersQ {
	return pgdb.NewMembersQ(pgx.Exec(s.db, ctx))
}

func (s Service) memberRolesQ(ctx context.Context) pgdb.MemberRolesQ {
	return pgdb.NewMemberRolesQ(pgx.Exec(s.db, ctx))
}

func (s Service) rolesQ(ctx context.Context) pgdb.RolesQ {
	return pgdb.NewRolesQ(pgx.Exec(s.db, ctx))
}

func (s Service) rolePermissionsQ(ctx context.Context) pgdb.RolePermissionsQ {
	return pgdb.NewRolePermissionsQ(pgx.Exec(s.db, ctx))
}
func (s Service) permissionsQ(ctx context.Context) pgdb.PermissionsQ {
	return pgdb.NewPermissionsQ(pgx.Exec(s.db, ctx))
}

func (s Service) invitesQ(ctx context.Context) pgdb.InvitesQ {
	return pgdb.NewInvitesQ(pgx.Exec(s.db, ctx))
}

func (s Service) profilesQ(ctx context.Context) pgdb.ProfilesQ {
	return pgdb.NewProfilesQ(pgx.Exec(s.db, ctx))
}

func (s Service) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return pgx.Transaction(s.db, ctx, fn)
}
