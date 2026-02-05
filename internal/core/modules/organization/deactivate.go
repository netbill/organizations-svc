package organization

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
)

func (m *Module) Deactivate(
	ctx context.Context,
	initiator models.Initiator,
	organizationID uuid.UUID,
) (models.Organization, error) {
	org, err := m.GetByID(ctx, organizationID)
	if err != nil {
		return models.Organization{}, err
	}

	member, err := m.repo.GetMemberByAccountAndOrganization(ctx, initiator.GetAccountID(), org.ID)
	if err != nil {
		return models.Organization{}, err
	}

	if !member.Head {
		return models.Organization{}, errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("only organization head member can activate organization, but member %s is not head", member.ID),
		)
	}

	if err = m.repo.Transaction(ctx, func(ctx context.Context) error {
		org, err = m.repo.UpdateOrganizationStatus(ctx, organizationID, models.OrganizationStatusInactive)
		if err != nil {
			return err
		}

		err = m.messenger.WriteOrganizationDeactivated(ctx, org)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return models.Organization{}, err
	}

	return org, nil
}
