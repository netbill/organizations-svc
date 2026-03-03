package organization

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/models"
)

func (m *Module) Suspend(ctx context.Context, organizationID uuid.UUID, value bool) (models.Organization, error) {
	org, err := m.GetByID(ctx, organizationID)
	if err != nil {
		return models.Organization{}, err
	}

	if org.Status == models.OrganizationStatusSuspended && value {
		return org, nil
	}
	if org.Status != models.OrganizationStatusSuspended && !value {
		return org, nil
	}

	var newStatus string
	if value {
		newStatus = models.OrganizationStatusSuspended
	} else {
		newStatus = models.OrganizationStatusInactive
	}

	err = m.repo.Transaction(ctx, func(ctx context.Context) error {
		org, err = m.repo.UpdateOrganizationStatus(ctx, organizationID, newStatus)
		if err != nil {
			return err
		}

		if err = m.messenger.WriteOrganizationUpdated(ctx, org); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return models.Organization{}, err
	}

	return org, nil
}
