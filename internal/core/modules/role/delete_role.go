package role

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/domain"
)

func (m *Module) Delete(
	ctx context.Context,
	initiator domain.AccountActor,
	roleID uuid.UUID,
) error {
	role, err := m.GetByID(ctx, roleID)
	if err != nil {
		return err
	}

	err = m.checkPermissionsToManageRole(ctx, initiator, role.OrganizationID, role.Rank)
	if err != nil {
		return err
	}

	return m.repo.Transaction(ctx, func(ctx context.Context) error {
		if err = m.repo.DeleteRole(ctx, roleID); err != nil {
			return err
		}

		if err = m.messenger.WriteOrgRoleDeleted(ctx, role); err != nil {
			return err
		}

		return nil
	})
}
