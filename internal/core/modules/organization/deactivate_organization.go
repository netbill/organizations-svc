package organization

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
)

func (s Service) DeactivateOrganization(
	ctx context.Context,
	accountID,
	organizationID uuid.UUID,
) (models.Organization, error) {
	org, err := s.GetOrganization(ctx, organizationID)
	if err != nil {
		return models.Organization{}, err
	}

	if org.Status == models.OrganizationStatusSuspended {
		return models.Organization{}, errx.ErrorOrganizationIsSuspended.Raise(
			fmt.Errorf("organization is not suspended"),
		)
	}

	initiator, err := s.getInitiator(ctx, accountID, organizationID)
	if err != nil {
		return models.Organization{}, err
	}

	err = s.chekPermissionForManageOrganization(ctx, initiator.ID)
	if err != nil {
		return models.Organization{}, err
	}

	if err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		org, err = s.repo.UpdateOrganizationStatus(ctx, organizationID, models.OrganizationStatusInactive)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to deactivate organization: %w", err),
			)
		}

		err = s.messenger.WriteOrganizationDeactivated(ctx, org)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to publish organization deactivate event: %w", err))
		}

		return nil
	}); err != nil {
		return models.Organization{}, err
	}

	return org, nil
}
