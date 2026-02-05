package role

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/models"
)

func (m *Module) SetRolePermissions(
	ctx context.Context,
	initiator models.InitiatorData,
	roleID uuid.UUID,
	permissions models.OrgRolePermissionDict,
) (role models.Role, links models.OrgRolePermissionLinks, err error) {
	role, err = m.GetRole(ctx, roleID)
	if err != nil {
		return models.Role{}, models.OrgRolePermissionLinks{}, err
	}

	member, err := m.getInitiator(ctx, initiator, role.OrganizationID)
	if err != nil {
		return models.Role{}, models.OrgRolePermissionLinks{}, err
	}

	err = m.checkPermissionsToManageRole(ctx, initiator, member.OrganizationID, role.Rank)
	if err != nil {
		return models.Role{}, models.OrgRolePermissionLinks{}, err
	}

	err = m.repo.Transaction(ctx, func(ctx context.Context) error {
		links, err = m.repo.SetRolePermissions(ctx, roleID, permissions)
		if err != nil {
			return err
		}

		err = m.messenger.WriteOrgRolePermissionsUpdated(ctx, role, links)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return models.Role{}, models.OrgRolePermissionLinks{}, err
	}

	return role, links, nil
}

func (m *Module) GetAllPermissions(ctx context.Context) ([]models.OrgRolePermission, error) {
	res, err := m.repo.GetAllPermissions(ctx)
	if err != nil {
		return nil, err
	}

	return res, nil
}
