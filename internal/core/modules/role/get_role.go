package role

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/restkit/pagi"
)

func (m *Module) GetRole(ctx context.Context, roleID uuid.UUID) (models.Role, error) {
	role, err := m.repo.GetRole(ctx, roleID)
	if err != nil {
		return models.Role{}, err
	}

	return role, nil
}

func (m *Module) GetRoleWithPermissions(ctx context.Context, accountID, roleID uuid.UUID) (models.Role, map[models.Permission]bool, error) {
	role, err := m.GetRole(ctx, roleID)
	if err != nil {
		return models.Role{}, nil, err
	}

	_, err = m.getInitiator(ctx, accountID, role.OrganizationID)
	if err != nil {
		return models.Role{}, nil, err
	}

	permissions, err := m.repo.GetRolePermissions(ctx, roleID)
	if err != nil {
		return models.Role{}, nil, err
	}

	return role, permissions, nil
}

type FilterParams struct {
	OrganizationID *uuid.UUID
	RolesID        *[]uuid.UUID
	Head           *bool
	Rank           *int
	Name           *string
}

func (m *Module) GetRoles(
	ctx context.Context,
	params FilterParams,
	limit, offset uint,
) (pagi.Page[[]models.Role], error) {
	res, err := m.repo.GetRoles(ctx, params, limit, offset)
	if err != nil {
		return pagi.Page[[]models.Role]{}, err
	}

	return res, nil
}
