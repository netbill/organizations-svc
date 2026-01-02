package organization

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
)

type UpdateParams struct {
	Name *string `json:"name,omitempty"`
	Icon *string `json:"icon,omitempty"`
}

func (s Service) UpdateOrganization(
	ctx context.Context,
	accountID, organizationID uuid.UUID,
	params UpdateParams,
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

	initiator, err := s.getInitiator(ctx, accountID, agglo.ID)
	if err != nil {
		return models.Organization{}, err
	}

	if err = s.chekPermissionForManageOrganization(
		ctx,
		initiator.ID,
	); err != nil {
		return models.Organization{}, err
	}

	if err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		agglo, err = s.repo.UpdateOrganization(ctx, organizationID, params)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to update organization: %w", err),
			)
		}

		err = s.messenger.WriteOrganizationUpdated(ctx, agglo)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to publish organization updated event: %w", err),
			)
		}

		return nil
	}); err != nil {
		return models.Organization{}, err
	}

	return agglo, nil
}
