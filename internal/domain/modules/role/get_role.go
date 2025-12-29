package role

import (
	"context"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/models"
)

func (s Service) GetRole(ctx context.Context, roleID uuid.UUID) (models.Role, error) {
	return s.repo.GetRole(ctx, roleID)
}

func (s Service) GetRoleWithPermissions(ctx context.Context, roleID uuid.UUID) (models.Role, []models.Permission, error) {
	role, err := s.repo.GetRole(ctx, roleID)
	if err != nil {
		return models.Role{}, nil, err
	}

	permissions, err := s.repo.GetRolePermissions(ctx, roleID)
	if err != nil {
		return models.Role{}, nil, err
	}

	return role, permissions, nil
}
