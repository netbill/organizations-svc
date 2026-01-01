package controller

import (
	"context"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/models"
	"github.com/umisto/cities-svc/internal/domain/modules/role"
	"github.com/umisto/logium"
	"github.com/umisto/pagi"
)

type Role interface {
	CreateRole(ctx context.Context, params role.CreateParams) (models.Role, error)
	CreateRoleByUser(
		ctx context.Context,
		initiatorID uuid.UUID,
		params role.CreateParams,
	) (models.Role, error)

	GetRole(ctx context.Context, roleID uuid.UUID) (models.Role, error)
	GetRoleWithPermissions(ctx context.Context, roleID uuid.UUID) (models.Role, []models.Permission, error)
	FilterRoles(
		ctx context.Context,
		params role.FilterParams,
		offset uint,
		limit uint,
	) (pagi.Page[[]models.Role], error)

	UpdateRoleByUser(
		ctx context.Context,
		accountID uuid.UUID,
		roleID uuid.UUID,
		params role.UpdateParams,
	) (models.Role, error)

	UpdateRolesRanksByUser(
		ctx context.Context,
		accountID uuid.UUID,
		agglomerationID uuid.UUID,
		order map[uuid.UUID]uint,
	) error

	DeleteRole(ctx context.Context, roleID uuid.UUID) error
	DeleteRoleByUser(ctx context.Context, accountID, roleID uuid.UUID) error

	GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]models.Permission, error)
	SetRolePermissions(
		ctx context.Context,
		roleID uuid.UUID,
		permissions map[models.CodeRolePermission]bool,
	) ([]models.Permission, error)
	SetRolePermissionsByUser(
		ctx context.Context,
		accountID, roleID uuid.UUID,
		permissions map[models.CodeRolePermission]bool,
	) ([]models.Permission, error)

	GetAllPermissions(ctx context.Context) ([]models.Permission, error)
}

type RoleController struct {
	domain Role
	log    logium.Logger
}
