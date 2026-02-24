package role

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/models"
)

func (m *Module) GetAllPermissions(ctx context.Context) ([]models.OrgRolePermission, error) {
	res, err := m.repo.GetAllPermissions(ctx)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (m *Module) UpdatePermissions(
	ctx context.Context,
	initiator models.AccountActor,
	roleID uuid.UUID,
	permissions []uuid.UUID,
) (role models.Role, links models.OrgRolePermissionsWithDetailsForRole, err error) {
	role, err = m.checkPermissionsToManageRole(ctx, initiator, roleID)
	if err != nil {
		return models.Role{}, models.OrgRolePermissionsWithDetailsForRole{}, err
	}

	err = m.repo.Transaction(ctx, func(ctx context.Context) error {
		if err = m.repo.LockRolePermissionsRevision(ctx, roleID); err != nil {
			return err
		}

		permissions, err = m.repo.UpdateRolePermissions(ctx, roleID, permissions)
		if err != nil {
			return err
		}

		revision, err := m.repo.BumpRolePermissionsRevision(ctx, role.ID)
		if err != nil {
			return err
		}

		if err = m.messenger.WriteOrgRolePermissionsUpdated(ctx, role, permissions, revision); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return models.Role{}, models.OrgRolePermissionsWithDetailsForRole{}, err
	}

	links, err = m.repo.GetRolePermissions(ctx, roleID)
	if err != nil {
		return models.Role{}, models.OrgRolePermissionsWithDetailsForRole{}, err
	}

	return role, links, nil
}
