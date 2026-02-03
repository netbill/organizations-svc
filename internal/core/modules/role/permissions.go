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
	permissions map[string]bool,
) (role models.Role, perm map[models.Permission]bool, err error) {
	role, err = m.GetRole(ctx, roleID)
	if err != nil {
		return models.Role{}, nil, err
	}

	member, err := m.getInitiator(ctx, initiator, role.OrganizationID)
	if err != nil {
		return models.Role{}, nil, err
	}

	if err = m.checkPermissionsToManageRole(ctx, initiator, member.OrganizationID, role.Rank); err != nil {
		return models.Role{}, nil, err
	}

	err = m.repo.Transaction(ctx, func(ctx context.Context) error {
		err = m.repo.SetRolePermissions(ctx, roleID, permissions)
		if err != nil {
			return err
		}

		perm, err = m.repo.GetRolePermissions(ctx, roleID)
		if err != nil {
			return err
		}

		err = m.messenger.WriteOrgRolePermissionsUpdated(ctx, role, perm)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return models.Role{}, nil, err
	}

	return role, perm, nil
}

func (m *Module) GetAllPermissions(ctx context.Context) ([]models.Permission, error) {
	res, err := m.repo.GetAllPermissions(ctx)
	if err != nil {
		return nil, err
	}

	return res, nil
}
