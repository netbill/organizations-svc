package organization

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/domain/errx"
	"github.com/netbill/organizations-svc/internal/domain/models"
)

func (s Service) ActivateOrganization(
	ctx context.Context,
	accountID, organizationID uuid.UUID,
) (models.Organization, error) {
	agglo, err := s.GetOrganization(ctx, organizationID)
	if err != nil {
		return models.Organization{}, err
	}

	if agglo.Status == models.OrganizationStatusSuspended {
		return models.Organization{}, errx.ErrorOrganizationIsSuspended.Raise(
			fmt.Errorf("organization is suspended"),
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

	err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		agglo, err = s.repo.UpdateOrganizationStatus(ctx, organizationID, models.OrganizationStatusActive)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to activate organization: %w", err))
		}

		err = s.messenger.WriteOrganizationActivated(ctx, agglo)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to publish organization activated event: %w", err))
		}

		return nil
	})
	if err != nil {
		return models.Organization{}, err
	}

	return agglo, err
}
