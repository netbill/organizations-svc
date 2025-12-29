package role

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/errx"
	"github.com/umisto/cities-svc/internal/domain/models"
)

func (s Service) GetRole(ctx context.Context, roleID uuid.UUID) (models.Role, error) {
	role, err := s.repo.GetRole(ctx, roleID)
	if err != nil {
		return models.Role{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get role: %w", err),
		)
	}
	if role.IsNil() {
		return models.Role{}, errx.ErrorRoleNotFound.Raise(
			fmt.Errorf("role not found: %s", roleID),
		)
	}

	return role, nil
}

func (s Service) GetRoleWithPermissions(ctx context.Context, roleID uuid.UUID) (models.Role, []models.Permission, error) {
	role, err := s.GetRole(ctx, roleID)
	if err != nil {
		return models.Role{}, nil, err
	}

	permissions, err := s.repo.GetRolePermissions(ctx, roleID)
	if err != nil {
		return models.Role{}, nil, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get role permissions: %w", err),
		)
	}

	return role, permissions, nil
}
