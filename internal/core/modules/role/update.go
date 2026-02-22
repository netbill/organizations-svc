package role

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/domain"
)

type UpdateParams struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Color       *string `json:"color"`
}

func (m *Module) Update(
	ctx context.Context,
	initiator domain.AccountActor,
	roleID uuid.UUID,
	params UpdateParams,
) (domain.Role, error) {
	role, err := m.GetByID(ctx, roleID)
	if err != nil {
		return domain.Role{}, err
	}

	member, err := m.getInitiator(ctx, initiator, role.OrganizationID)
	if err != nil {
		return domain.Role{}, err
	}

	if !member.Head {
		if err = m.checkPermissionsToManageRole(ctx, initiator, member.OrganizationID, role.Rank); err != nil {
			return domain.Role{}, err
		}
	}

	if err = m.repo.Transaction(ctx, func(ctx context.Context) error {
		role, err = m.repo.UpdateRole(ctx, roleID, params)
		if err != nil {
			return err
		}

		if err = m.messenger.WriteOrgRoleUpdated(ctx, role); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return domain.Role{}, err
	}

	return role, nil
}
