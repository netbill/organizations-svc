package organization

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
)

func (m *Module) Delete(
	ctx context.Context,
	actor models.AccountActor,
	organizationID uuid.UUID,
) error {
	organization, err := m.GetByID(ctx, organizationID)
	if errors.Is(err, errx.ErrorOrganizationNotFound) {
		buried, err := m.repo.OrganizationIsBuried(ctx, organizationID)
		if err != nil {
			return err
		}
		if buried {
			return errx.ErrorOrganizationDeleted.Raise(
				fmt.Errorf("organization with id %s is already deleted", organizationID),
			)
		}
	}
	if err != nil {
		return err
	}

	member, err := m.getInitiator(ctx, actor, organization.ID)
	if err != nil {
		return err
	}
	if !member.Head {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("initiator member %s is not head of organization %s", member.ID, organization.ID),
		)
	}

	if organization.Status == models.OrganizationStatusSuspended {
		return errx.ErrorOrganizationIsSuspended.Raise(
			fmt.Errorf("organization %s is suspended", organization.ID),
		)
	}

	places, err := m.repo.GetPlaceExistsForOrganization(ctx, organization.ID)
	if err != nil {
		return fmt.Errorf("failed to get places for organization: %w", err)
	}
	if places {
		return errx.ErrorOrganizationHavePlace.Raise(
			fmt.Errorf("organization %s has places, cannot be deleted", organization.ID),
		)
	}

	return m.repo.Transaction(ctx, func(ctx context.Context) error {
		err = m.repo.BuryOrganization(ctx, organizationID)
		if err != nil {
			return fmt.Errorf("failed to bury organization: %w", err)
		}

		err = m.repo.DeleteOrganization(ctx, organizationID)
		if err != nil {
			return fmt.Errorf("failed to delete organization: %w", err)
		}

		err = m.messenger.WriteOrganizationDeleted(ctx, organization)
		if err != nil {
			return fmt.Errorf("failed to publish organization delete event: %w", err)
		}

		return nil
	})
}
