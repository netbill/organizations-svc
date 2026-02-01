package repository

import (
	"context"
)

type Repository struct {
	Transactioner
	OrganizationsSql          OrganizationsQ
	OrgMembersSql             OrgMembersQ
	OrgMemberRolesSql         OrgMemberRolesQ
	OrgRolesSql               OrgRolesQ
	OrgRolePermissionLinksSql OrgRolePermissionLinksQ
	OrgRolePermissionsSql     OrgRolePermissionsQ
	OrgInvitesSql             OrgInvitesQ
	ProfilesSql               ProfilesQ
}

type Transactioner interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

func (r *Repository) organizationsQ() OrganizationsQ {
	return r.OrganizationsSql.New()
}

func (r *Repository) orgMembersQ() OrgMembersQ {
	return r.OrgMembersSql.New()
}

func (r *Repository) orgMemberRolesQ() OrgMemberRolesQ {
	return r.OrgMemberRolesSql.New()
}

func (r *Repository) orgRolesQ() OrgRolesQ {
	return r.OrgRolesSql.New()
}

func (r *Repository) orgRolePermissionLinksQ() OrgRolePermissionLinksQ {
	return r.OrgRolePermissionLinksSql.New()
}

func (r *Repository) orgRolePermissionsQ() OrgRolePermissionsQ {
	return r.OrgRolePermissionsSql.New()
}

func (r *Repository) orgInvitesQ() OrgInvitesQ {
	return r.OrgInvitesSql.New()
}

func (r *Repository) profilesQ() ProfilesQ {
	return r.ProfilesSql.New()
}
