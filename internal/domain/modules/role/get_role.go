package role

import (
	"context"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/entity"
)

func (s Service) GetRole(ctx context.Context, roleID uuid.UUID) (entity.Role, error) {
	return s.repo.GetRole(ctx, roleID)
}

func (s Service) GetRoleWithPermissions(ctx context.Context, roleID uuid.UUID) (entity.Role, []entity.Permission, error) {
	role, err := s.repo.GetRole(ctx, roleID)
	if err != nil {
		return entity.Role{}, nil, err
	}

	permissions, err := s.repo.GetRolePermissions(ctx, roleID)
	if err != nil {
		return entity.Role{}, nil, err
	}

	return role, permissions, nil
}
