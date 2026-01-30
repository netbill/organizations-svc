package organization

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
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
	if org.IsNil() {
		return models.Organization{}, errx.ErrorOrganizationNotFound.Raise(
			fmt.Errorf("organization %s is not found", organizationID),
		)
	}

	err = s.chekPermissionForManageOrganization(ctx, accountID, org.ID)
	if err != nil {
		return models.Organization{}, err
	}

	err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		org, err = s.repo.UpdateOrganizationStatus(ctx, organizationID, models.OrganizationStatusActive)
		if err != nil {
			return fmt.Errorf("failed to activate organization, cause: %w", err)
		}

		err = s.messenger.WriteOrganizationActivated(ctx, org)
		if err != nil {
			return fmt.Errorf("failed to publish organization activated event, cause: %w", err)
		}

		return nil
	})
	if err != nil {
		return models.Organization{}, err
	}

	return org, err
}
