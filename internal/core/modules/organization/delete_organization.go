package organization

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
)

func (s Service) DeleteOrganization(ctx context.Context, accountID, ID uuid.UUID) error {
	organization, err := s.GetOrganization(ctx, ID)
	if err != nil {
		return err
	}

	if organization.Status == models.OrganizationStatusSuspended {
		return errx.ErrorOrganizationIsSuspended.Raise(
			fmt.Errorf("organization is suspended"),
		)
	}

	initiator, err := s.getInitiator(ctx, accountID, organization.ID)
	if err != nil {
		return err
	}

	role, err := s.repo.GetMemberMaxRole(ctx, initiator.ID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get member max role: %w", err),
		)
	}

	if role.Head != true {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("only organization head can delete organization"),
		)
	}

	return s.repo.Transaction(ctx, func(ctx context.Context) error {
		err = s.repo.DeleteOrganization(ctx, ID)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to delete organization: %w", err),
			)
		}

		err = s.messenger.WriteOrganizationDeleted(ctx, organization)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to publish organization delete event: %w", err),
			)
		}

		return nil
	})
}
