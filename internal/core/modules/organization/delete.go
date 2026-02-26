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
		return nil
	} else if err != nil {
		return err
	}

	member, err := m.repo.GetMemberByAccountAndOrganization(ctx, actor, organization.ID)
	if err != nil {
		return err
	}
	if !member.Head {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("initiator member %s is not head of organization %s", member.ID, organization.ID),
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
