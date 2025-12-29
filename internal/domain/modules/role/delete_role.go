package role

import (
	"context"

	"github.com/google/uuid"
)

func (s Service) DeleteRole(ctx context.Context, roleID uuid.UUID) error {
	return s.repo.Transaction(ctx, func(ctx context.Context) error {
		if err := s.repo.DeleteRole(ctx, roleID); err != nil {
			return err
		}

		return s.messenger.WriteRoleDeleted(ctx, roleID)
	})
}

func (s Service) DeleteRoleByUser(ctx context.Context, accountID, roleID uuid.UUID) error {
	role, err := s.GetRole(ctx, roleID)
	if err != nil {
		return err
	}

	if err = s.CheckPermissionsToManageRole(ctx, accountID, role.AgglomerationID, role.Rank); err != nil {
		return err
	}

	return s.DeleteRole(ctx, role.ID)
}
