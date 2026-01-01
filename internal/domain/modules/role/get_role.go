package role

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/umisto/agglomerations-svc/internal/domain/errx"
	"github.com/umisto/agglomerations-svc/internal/domain/models"
	"github.com/umisto/pagi"
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

func (s Service) GetRoleWithPermissions(ctx context.Context, accountID, roleID uuid.UUID) (models.Role, []models.Permission, error) {
	role, err := s.GetRole(ctx, roleID)
	if err != nil {
		return models.Role{}, nil, err
	}

	if err = s.CheckPermissionsToManageRole(ctx, accountID, role.AgglomerationID, role.Rank); err != nil {
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

type FilterParams struct {
	AgglomerationID *uuid.UUID
	RolesID         *[]uuid.UUID
	Head            *bool
	Rank            *int
	Name            *string
}

func (s Service) GetRoles(
	ctx context.Context,
	params FilterParams,
	offset uint,
	limit uint,
) (pagi.Page[[]models.Role], error) {
	res, err := s.repo.GetRoles(ctx, params, offset, limit)
	if err != nil {
		return pagi.Page[[]models.Role]{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to filter roles: %w", err),
		)
	}

	return res, nil
}
