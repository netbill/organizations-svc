package role

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/models"
)

func (m *Module) Delete(
	ctx context.Context,
	initiator models.AccountActor,
	roleID uuid.UUID,
) error {
	role, err := m.checkPermissionsToManageRole(ctx, initiator, roleID)
	if err != nil {
		return err
	}

	return m.repo.Transaction(ctx, func(ctx context.Context) error {
		err = m.repo.LockOrgRoleRankRevision(ctx, role.OrganizationID)
		if err != nil {
			return err
		}

		ranks, err := m.repo.DeleteRoleRank(ctx, roleID)
		if err != nil {
			return err
		}

		revision, err := m.repo.BumpOrgRoleRankRevision(ctx, role.OrganizationID)
		if err != nil {
			return err
		}

		if err = m.repo.DeleteRole(ctx, roleID); err != nil {
			return err
		}

		if err = m.messenger.WriteOrgRolesRanksUpdated(ctx, role.OrganizationID, ranks, revision); err != nil {
			return err
		}

		if err = m.messenger.WriteOrgRoleDeleted(ctx, role); err != nil {
			return err
		}

		return nil
	})
}
