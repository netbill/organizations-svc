package role

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/errx"
)

func (s Service) DeleteRole(ctx context.Context, roleID uuid.UUID) error {
	return s.repo.Transaction(ctx, func(ctx context.Context) error {
		if err := s.repo.DeleteRole(ctx, roleID); err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to delete role: %w", err),
			)
		}

		if err := s.messenger.WriteRoleDeleted(ctx, roleID); err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to send role deleted message: %w", err),
			)
		}

		return nil
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
