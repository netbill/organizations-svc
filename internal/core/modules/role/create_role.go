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

func (s Service) CreateRole(
	ctx context.Context,
	accountID uuid.UUID,
	params CreateParams,
) (role models.Role, err error) {
	initiator, err := s.getInitiator(ctx, accountID, params.OrganizationID)
	if err != nil {
		return role, err
	}

	if err = s.checkPermissionsToManageRole(ctx, initiator.AccountID, params.Rank); err != nil {
		return models.Role{}, err
	}

	if err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		role, err = s.repo.CreateRole(ctx, params)
		if err != nil {
			return err
		}

		err = s.messenger.WriteOrgRoleCreated(ctx, role)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return models.Role{}, err
	}

	return role, nil
}
