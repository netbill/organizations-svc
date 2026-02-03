package organization

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
)

func (m *Module) DeleteOrganization(
	ctx context.Context,
	initiator models.InitiatorData,
	organizationID uuid.UUID,
) error {
	organization, err := m.GetOrganization(ctx, organizationID)
	if err != nil {
		return err
	}

	member, err := m.repo.GetMemberByAccountAndOrganization(ctx, initiator.AccountID, organization.ID)
	if err != nil {
		return err
	}

	if !member.Head {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("initiator member %s is not head of organization %s", member.ID, organization.ID),
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
