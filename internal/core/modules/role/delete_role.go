package role

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
)

func (s Service) DeleteRole(ctx context.Context, accountID, roleID uuid.UUID) error {
	role, err := s.GetRole(ctx, roleID)
	if err != nil {
		return err
	}

	initiator, err := s.getInitiator(ctx, accountID, role.OrganizationID)
	if err != nil {
		return err
	}

	if err = s.checkPermissionsToManageRole(ctx, initiator.AccountID, role.Rank); err != nil {
		return err
	}

	if role.Head {
		return errx.ErrorCannotDeleteHeadRole.Raise(
			fmt.Errorf("cannot delete head role"),
		)
	}

	return s.repo.Transaction(ctx, func(ctx context.Context) error {
		if err = s.repo.DeleteRole(ctx, roleID); err != nil {
			return err
		}

		if err = s.messenger.WriteOrgRoleDeleted(ctx, role); err != nil {
			return err
		}

		return nil
	})
}
