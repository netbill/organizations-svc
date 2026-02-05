package role

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/models"
)

func (m *Module) GetWithPermissions(
	ctx context.Context,
	initiator models.Initiator,
	roleID uuid.UUID,
) (models.Role, models.OrgRolePermissionDictWithDetails, error) {
	role, err := m.GetByID(ctx, roleID)
	if err != nil {
		return models.Role{}, models.OrgRolePermissionDictWithDetails{}, err
	}

	_, err = m.getInitiator(ctx, initiator, role.OrganizationID)
	if err != nil {
		return models.Role{}, models.OrgRolePermissionDictWithDetails{}, err
	}

	permissions, err := m.repo.GetRolePermissions(ctx, roleID)
	if err != nil {
		return models.Role{}, models.OrgRolePermissionDictWithDetails{}, err
	}

	return role, permissions, nil
}
