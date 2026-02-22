package role

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/domain"
)

func (m *Module) GetWithPermissions(
	ctx context.Context,
	initiator domain.AccountActor,
	roleID uuid.UUID,
) (domain.Role, domain.OrgRolePermissionsWithDetailsForRole, error) {
	role, err := m.GetByID(ctx, roleID)
	if err != nil {
		return domain.Role{}, domain.OrgRolePermissionsWithDetailsForRole{}, err
	}

	_, err = m.getInitiator(ctx, initiator, role.OrganizationID)
	if err != nil {
		return domain.Role{}, domain.OrgRolePermissionsWithDetailsForRole{}, err
	}

	permissions, err := m.repo.GetRolePermissions(ctx, roleID)
	if err != nil {
		return domain.Role{}, domain.OrgRolePermissionsWithDetailsForRole{}, err
	}

	return role, permissions, nil
}
