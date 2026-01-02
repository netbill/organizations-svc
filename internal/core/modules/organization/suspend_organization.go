package organization

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
)

func (s Service) SuspendOrganization(ctx context.Context, ID uuid.UUID) (agglo models.Organization, err error) {
	if err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		agglo, err = s.repo.UpdateOrganizationStatus(ctx, ID, models.OrganizationStatusSuspended)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to suspend organization: %w", err))
		}

		err = s.messenger.WriteOrganizationSuspended(ctx, agglo)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to publish organization suspend event: %w", err),
			)
		}

		return nil
	}); err != nil {
		return models.Organization{}, err
	}

	return agglo, nil
}
