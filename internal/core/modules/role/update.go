package role

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/models"
)

type UpdateParams struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

func (m *Module) Update(
	ctx context.Context,
	initiator models.AccountActor,
	roleID uuid.UUID,
	params UpdateParams,
) (models.Role, error) {
	role, err := m.checkPermissionsToManageRole(ctx, initiator, roleID)
	if err != nil {
		return models.Role{}, err
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
		return models.Role{}, err
	}

	return role, nil
}
