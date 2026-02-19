package boot

import (
	"github.com/netbill/organizations-svc/internal/repository"
	"github.com/netbill/organizations-svc/internal/repository/pg"
	"github.com/netbill/pgdbx"
)

func (c *Config) NewRepository(db *pgdbx.DB) *repository.Repository {
	return &repository.Repository{
		OrgInvitesSql:             pg.NewOrgInvitesQ(db),
		OrgMemberRolesSql:         pg.NewOrgMemberRolesQ(db),
		OrgMembersSql:             pg.NewOrgMembersQ(db),
		OrganizationsSql:          pg.NewOrganizationsQ(db),
		OrgRolePermissionsSql:     pg.NewOrgRolePermissionsQ(db),
		OrgRolePermissionLinksSql: pg.NewOrgRolePermissionLinksQ(db),
		OrgRolesSql:               pg.NewOrgRolesQ(db),
		ProfilesSql:               pg.NewProfilesQ(db),
		Transactioner:             pg.NewTransaction(db),
	}
}
