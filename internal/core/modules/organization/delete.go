package organization

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/domain"
	"github.com/netbill/organizations-svc/internal/core/errx"
)

func (m *Module) Delete(
	ctx context.Context,
	initiator domain.AccountActor,
	organizationID uuid.UUID,
) error {
	organization, err := m.GetByID(ctx, organizationID)
	if err != nil {
		return err
	}

	member, err := m.repo.GetMemberByAccountAndOrganization(ctx, initiator, organization.ID)
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
