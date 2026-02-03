package role

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/models"
)

type CreateParams struct {
	OrganizationID uuid.UUID `json:"organization_id"`
	Rank           uint      `json:"rank"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Color          string    `json:"color"`
}

func (m *Module) CreateRole(
	ctx context.Context,
	initiator models.InitiatorData,
	params CreateParams,
) (role models.Role, err error) {
	err = m.checkPermissionsToManageRole(ctx, initiator, params.OrganizationID, params.Rank)
	if err != nil {
		return models.Role{}, err
	}

	if err = m.repo.Transaction(ctx, func(ctx context.Context) error {
		role, err = m.repo.CreateRole(ctx, params)
		if err != nil {
			return err
		}

		err = m.messenger.WriteOrgRoleCreated(ctx, role)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return models.Role{}, err
	}

	return role, nil
}
