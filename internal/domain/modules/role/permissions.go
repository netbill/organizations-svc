package role

import (
	"context"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/errx"
	"github.com/umisto/cities-svc/internal/domain/models"
)

func (s Service) GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]models.Permission, error) {
	return s.repo.GetRolePermissions(ctx, roleID)
}

func (s Service) SetRolePermissions(
	ctx context.Context,
	roleID uuid.UUID,
	permissions map[models.CodeRolePermission]bool,
) (perm []models.Permission, err error) {
	err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		err = s.repo.SetRolePermissions(ctx, roleID, permissions)
		if err != nil {
			return err
		}

		perm, err = s.repo.GetRolePermissions(ctx, roleID)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return perm, nil
}

func (s Service) SetRolePermissionsByUser(
	ctx context.Context,
	accountID, roleID uuid.UUID,
	permissions map[models.CodeRolePermission]bool,
) ([]models.Permission, error) {
	role, err := s.GetRole(ctx, roleID)
	if err != nil {
		return nil, err
	}

	maxRole, err := s.repo.GetAccountMaxRoleInAgglomeration(ctx, accountID, role.AgglomerationID)
	if err != nil {
		return nil, err
	}
	if err = s.CheckPermissionsToManageRole(ctx, accountID, roleID, role.Rank); err != nil {
		return nil, err
	}
	if role.Rank <= maxRole.Rank {
		return nil, errx.ErrorNotEnoughRights.Raise(
			nil,
		)
	}

	return s.SetRolePermissions(ctx, roleID, permissions)
}

func (s Service) GetAllPermissions(ctx context.Context) ([]models.Permission, error) {
	return s.repo.GetAllPermissions(ctx)
}
