package repository

import (
	"context"
	"database/sql"

	"github.com/umisto/cities-svc/internal/repository/pgdb"
	"github.com/umisto/pgx"
)

type Service struct {
	db *sql.DB
}

func New(db *sql.DB) Service {
	return Service{db: db}
}

func (s Service) agglomerationsQ() pgdb.AgglomerationsQ {
	return pgdb.NewAgglomerationsQ(s.db)
}

func (s Service) citiesQ() pgdb.CitiesQ {
	return pgdb.NewCitiesQ(s.db)
}

func (s Service) membersQ() pgdb.MembersQ {
	return pgdb.NewMembersQ(s.db)
}

func (s Service) memberRolesQ() pgdb.MemberRolesQ {
	return pgdb.NewMemberRolesQ(s.db)
}

func (s Service) rolesQ() pgdb.RolesQ {
	return pgdb.NewRolesQ(s.db)
}

func (s Service) rolePermissionsQ() pgdb.RolePermissionsQ {
	return pgdb.NewRolePermissionsQ(s.db)
}

func (s Service) permissionsQ() pgdb.PermissionsQ {
	return pgdb.NewPermissionsQ(s.db)
}

func (s Service) invitesQ() pgdb.InvitesQ {
	return pgdb.NewInvitesQ(s.db)
}

func (s Service) profilesQ() pgdb.ProfilesQ {
	return pgdb.NewProfilesQ(s.db)
}

func (s Service) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return pgx.Transaction(s.db, ctx, fn)
}
