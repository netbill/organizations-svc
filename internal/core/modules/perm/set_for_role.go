package perm

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/models"
)

type SetForRole struct {
	OrgUpdate     bool `json:"org_update"`
	RolesManage   bool `json:"roles_manage"`
	InviteManage  bool `json:"invite_manage"`
	MembersDelete bool `json:"members_delete"`
	MemberUpdate  bool `json:"member_update"`
	PlacesCreate  bool `json:"places_create"`
	PlacesDelete  bool `json:"places_delete"`
	PlacesUpdate  bool `json:"places_update"`
}

func (m *Module) SetForRole(
	ctx context.Context,
	initiator models.AccountActor,
	roleID uuid.UUID,
	params SetForRole,
) (role models.Role, links models.OrgRolePermissionsWithDetailsForRole, err error) {
	role, err = m.repo.GetRole(ctx, roleID)
	if err != nil {
		return models.Role{}, nil, err
	}

	member, err := m.getInitiator(ctx, initiator, role.OrganizationID)
	if err != nil {
		return models.Role{}, nil, err
	}

	err = m.checkPermissionsToManageRole(ctx, initiator, member.OrganizationID, role.Rank)
	if err != nil {
		return models.Role{}, nil, err
	}

	err = m.repo.Transaction(ctx, func(ctx context.Context) error {
		err = m.repo.SetRolePermissions(ctx, roleID, params)
		if err != nil {
			return err
		}

		err = m.messenger.WriteOrgRolePermissionsUpdated(ctx, role, params)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return models.Role{}, nil, err
	}

	links, err = m.repo.GetRolePermissions(ctx, roleID)
	if err != nil {
		return models.Role{}, nil, err
	}

	return role, links, nil
}
