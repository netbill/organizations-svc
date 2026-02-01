package role

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
)

func (m *Module) DeleteRole(ctx context.Context, accountID, roleID uuid.UUID) error {
	role, err := m.GetRole(ctx, roleID)
	if err != nil {
		return err
	}

	initiator, err := m.getInitiator(ctx, accountID, role.OrganizationID)
	if err != nil {
		return err
	}

	if err = m.checkPermissionsToManageRole(ctx, initiator.AccountID, role.Rank); err != nil {
		return err
	}

	if role.Head {
		return errx.ErrorCannotDeleteHeadRole.Raise(
			fmt.Errorf("cannot delete head role"),
		)
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
