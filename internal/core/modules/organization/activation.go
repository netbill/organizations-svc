package organization

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
)

func (m *Module) Activate(
	ctx context.Context,
	actor models.AccountActor,
	organizationID uuid.UUID,
) (models.Organization, error) {
	org, err := m.GetByID(ctx, organizationID)
	if err != nil {
		return models.Organization{}, err
	}

	member, err := m.repo.GetMemberByAccountAndOrganization(ctx, actor, org.ID)
	if err != nil {
		return models.Organization{}, err
	}
	if !member.Head {
		return models.Organization{}, errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("only organization head member can activate organization, but member %s is not head", member.ID),
		)
	}

	if org.Status == models.OrganizationStatusActive {
		return org, nil
	}

	err = m.repo.Transaction(ctx, func(ctx context.Context) error {
		org, err = m.repo.UpdateOrganizationStatus(ctx, organizationID, models.OrganizationStatusActive)
		if err != nil {
			return err
		}

		err = m.messenger.WriteOrganizationUpdated(ctx, org)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return models.Organization{}, err
	}

	return org, err
}

func (m *Module) Deactivate(
	ctx context.Context,
	actor models.AccountActor,
	organizationID uuid.UUID,
) (models.Organization, error) {
	org, err := m.GetByID(ctx, organizationID)
	if err != nil {
		return models.Organization{}, err
	}

	member, err := m.repo.GetMemberByAccountAndOrganization(ctx, actor, org.ID)
	if err != nil {
		return models.Organization{}, err
	}
	if !member.Head {
		return models.Organization{}, errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("only organization head member can activate organization, but member %s is not head", member.ID),
		)
	}

	if org.Status == models.OrganizationStatusInactive {
		return org, nil
	}

	if err = m.repo.Transaction(ctx, func(ctx context.Context) error {
		org, err = m.repo.UpdateOrganizationStatus(ctx, organizationID, models.OrganizationStatusInactive)
		if err != nil {
			return err
		}

		err = m.messenger.WriteOrganizationUpdated(ctx, org)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return models.Organization{}, err
	}

	return org, nil
}
