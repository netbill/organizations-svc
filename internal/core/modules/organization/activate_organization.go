package organization

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/models"
)

func (s Service) ActivateOrganization(
	ctx context.Context,
	accountID, organizationID uuid.UUID,
) (models.Organization, error) {
	org, err := s.GetOrganization(ctx, organizationID)
	if err != nil {
		return models.Organization{}, err
	}

	err = s.chekPermissionForManageOrganization(ctx, accountID, org.ID)
	if err != nil {
		return models.Organization{}, err
	}

	err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		org, err = s.repo.UpdateOrganizationStatus(ctx, organizationID, models.OrganizationStatusActive)
		if err != nil {
			return err
		}

		err = s.messenger.WriteOrganizationActivated(ctx, org)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return models.Organization{}, err
	}

	return org, err
}
