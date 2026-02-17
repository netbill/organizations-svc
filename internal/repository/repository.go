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
