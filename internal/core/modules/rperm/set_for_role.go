package rperm

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/models"
)

func (m *Module) SetForRole(
	ctx context.Context,
	initiator models.Initiator,
	roleID uuid.UUID,
	permissions models.OrgRolePermissionAccess,
) (role models.Role, links models.OrgRolePermissionDictWithDetails, err error) {
	role, err = m.repo.GetRole(ctx, roleID)
	if err != nil {
		return models.Role{}, models.OrgRolePermissionDictWithDetails{}, err
	}

	member, err := m.getInitiator(ctx, initiator, role.OrganizationID)
	if err != nil {
		return models.Role{}, models.OrgRolePermissionDictWithDetails{}, err
	}

	err = m.checkPermissionsToManageRole(ctx, initiator, member.OrganizationID, role.Rank)
	if err != nil {
		return models.Role{}, models.OrgRolePermissionDictWithDetails{}, err
	}

	err = m.repo.Transaction(ctx, func(ctx context.Context) error {
		links, err = m.repo.SetRolePermissions(ctx, roleID, permissions)
		if err != nil {
			return err
		}

		err = m.messenger.WriteOrgRolePermissionsUpdated(ctx, role, permissions)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return models.Role{}, models.OrgRolePermissionDictWithDetails{}, err
	}

	return role, links, nil
}
