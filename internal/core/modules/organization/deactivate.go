package organization

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/domain"
	"github.com/netbill/organizations-svc/internal/core/errx"
)

func (m *Module) Deactivate(
	ctx context.Context,
	initiator domain.AccountActor,
	organizationID uuid.UUID,
) (domain.Organization, error) {
	org, err := m.GetByID(ctx, organizationID)
	if err != nil {
		return domain.Organization{}, err
	}

	member, err := m.repo.GetMemberByAccountAndOrganization(ctx, initiator, org.ID)
	if err != nil {
		return domain.Organization{}, err
	}

	if !member.Head {
		return domain.Organization{}, errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("only organization head member can activate organization, but member %s is not head", member.ID),
		)
	}

	if err = m.repo.Transaction(ctx, func(ctx context.Context) error {
		org, err = m.repo.UpdateOrganizationStatus(ctx, organizationID, domain.OrganizationStatusInactive)
		if err != nil {
			return err
		}

		err = m.messenger.WriteOrganizationDeactivated(ctx, org)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return domain.Organization{}, err
	}

	return org, nil
}
