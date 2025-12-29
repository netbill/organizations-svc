package role

import (
	"context"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/entity"
)

func (s Service) GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]entity.Permission, error) {
	return s.repo.GetRolePermissions(ctx, roleID)
}

func (s Service) SetRolePermissions(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) ([]entity.Permission, error) {
	return s.repo.SetRolePermissions(ctx, roleID, permissionIDs)
}
