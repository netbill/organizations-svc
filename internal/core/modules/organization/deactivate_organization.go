package organization

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/models"
)

func (m *Module) DeactivateOrganization(
	ctx context.Context,
	accountID,
	organizationID uuid.UUID,
) (models.Organization, error) {
	org, err := m.GetOrganization(ctx, organizationID)
	if err != nil {
		return models.Organization{}, err
	}

	err = m.chekPermissionForManageOrganization(ctx, accountID, org.ID)
	if err != nil {
		return models.Organization{}, err
	}

	if err = m.repo.Transaction(ctx, func(ctx context.Context) error {
		org, err = m.repo.UpdateOrganizationStatus(ctx, organizationID, models.OrganizationStatusInactive)
		if err != nil {
			return err
		}

		//TODO clean organization
		err = m.messenger.WriteOrganizationDeactivated(ctx, org)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return models.Organization{}, err
	}

	return org, nil
}
