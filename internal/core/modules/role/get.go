package role

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/restkit/pagi"
)

func (m *Module) GetByID(
	ctx context.Context,
	roleID uuid.UUID,
) (models.Role, error) {
	role, err := m.repo.GetRole(ctx, roleID)
	if err != nil {
		return models.Role{}, err
	}

	return role, nil
}

func (m *Module) GetWithPermissions(
	ctx context.Context,
	initiator models.AccountActor,
	roleID uuid.UUID,
) (models.Role, models.OrgRolePermissionsWithDetailsForRole, error) {
	role, err := m.GetByID(ctx, roleID)
	if err != nil {
		return models.Role{}, models.OrgRolePermissionsWithDetailsForRole{}, err
	}

	_, err = m.getInitiator(ctx, initiator, role.OrganizationID)
	if err != nil {
		return models.Role{}, models.OrgRolePermissionsWithDetailsForRole{}, err
	}

	permissions, err := m.repo.GetRolePermissions(ctx, roleID)
	if err != nil {
		return models.Role{}, models.OrgRolePermissionsWithDetailsForRole{}, err
	}

	return role, permissions, nil
}

type FilterParams struct {
	OrganizationID uuid.UUID
	RolesID        *[]uuid.UUID
	Name           *string
}

func (m *Module) GetList(
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
