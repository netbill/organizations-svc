package role

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/umisto/agglomerations-svc/internal/domain/errx"
	"github.com/umisto/agglomerations-svc/internal/domain/models"
)

func (s Service) GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]models.Permission, error) {
	permission, err := s.repo.GetRolePermissions(ctx, roleID)
	if err != nil {
		return nil, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get role permissions: %w", err),
		)
	}

	return permission, nil
}

func (s Service) SetRolePermissions(
	ctx context.Context,
	accountID, roleID uuid.UUID,
	permissions map[models.CodeRolePermission]bool,
) (perm []models.Permission, err error) {
	role, err := s.GetRole(ctx, roleID)
	if err != nil {
		return nil, err
	}

	maxRole, err := s.repo.GetAccountMaxRoleInAgglomeration(ctx, accountID, role.AgglomerationID)
	if err != nil {
		return nil, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get account max role in agglomeration: %w", err),
		)
	}

	if err = s.CheckPermissionsToManageRole(ctx, accountID, roleID, role.Rank); err != nil {
		return nil, err
	}
	if role.Rank <= maxRole.Rank {
		return nil, errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("account does not have enough rights to set permissions for this role"),
		)
	}

	err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		err = s.repo.SetRolePermissions(ctx, roleID, permissions)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to set role permissions: %w", err),
			)
		}

		perm, err = s.repo.GetRolePermissions(ctx, roleID)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get role permissions after setting: %w", err),
			)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return perm, nil
}

func (s Service) GetAllPermissions(ctx context.Context) ([]models.Permission, error) {
	res, err := s.repo.GetAllPermissions(ctx)
	if err != nil {
		return nil, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get all permissions: %w", err),
		)
	}

	return res, nil
}
