package role

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/domain/errx"
	"github.com/netbill/organizations-svc/internal/domain/models"
)

func (s Service) SetRolePermissions(
	ctx context.Context,
	accountID, roleID uuid.UUID,
	permissions map[string]bool,
) (role models.Role, perm map[models.Permission]bool, err error) {
	role, err = s.GetRole(ctx, roleID)
	if err != nil {
		return models.Role{}, nil, err
	}

	initiator, err := s.getInitiator(ctx, accountID, role.OrganizationID)
	if err != nil {
		return models.Role{}, nil, err
	}

	if err = s.checkPermissionsToManageRole(ctx, initiator.ID, role.Rank); err != nil {
		return models.Role{}, nil, err
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

		err = s.messenger.WriteRolePermissionsUpdated(ctx, roleID, perm)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to send role permissions updated message: %w", err),
			)
		}

		return nil
	})
	if err != nil {
		return models.Role{}, nil, err
	}

	return role, perm, nil
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
