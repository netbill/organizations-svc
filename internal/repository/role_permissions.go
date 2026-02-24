package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/models"
)

func (r *Repository) GetAllPermissions(ctx context.Context) ([]models.OrgRolePermission, error) {
	//TODO implement me
	panic("implement me")
}

func (r *Repository) UpdateRolePermissions(ctx context.Context, roleID uuid.UUID, permissions []uuid.UUID) ([]uuid.UUID, error) {
	//TODO implement me
	panic("implement me")
}

func (r *Repository) GetRolePermissions(
	ctx context.Context,
	roleID uuid.UUID,
) (models.OrgRolePermissionsWithDetailsForRole, error) {
	//TODO implement me
	panic("implement me")
}

func (r *Repository) CreateRolePermissionsRevision(ctx context.Context, roleID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (r *Repository) LockRolePermissionsRevision(ctx context.Context, roleID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (r *Repository) BumpRolePermissionsRevision(ctx context.Context, roleID uuid.UUID) (models.OrgRolePermissionsLinksRevision, error) {
	//TODO implement me
	panic("implement me")
}
