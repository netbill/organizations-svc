package organization

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
)

type UpdateParams struct {
	Name      string
	IconKey   *string
	BannerKey *string
}

func (m *Module) Update(
	ctx context.Context,
	actor models.AccountActor,
	organizationID uuid.UUID,
	params UpdateParams,
) (models.Organization, error) {
	org, err := m.GetByID(ctx, organizationID)
	if err != nil {
		return models.Organization{}, err
	}

	member, err := m.getInitiator(ctx, actor, org.ID)
	if err != nil {
		return models.Organization{}, err
	}
	if !member.Head {
		return models.Organization{}, errx.ErrorNotOrganizationHead.Raise(
			fmt.Errorf("only organization head member can activate organization, but member %s is not head", member.ID),
		)
	}

	if org.Status == models.OrganizationStatusSuspended {
		return models.Organization{}, errx.ErrorOrganizationIsSuspended.Raise(
			fmt.Errorf("organization %s is suspended", org.ID),
		)
	}

	upd := org.Name != params.Name

	if !ptrStrEq(params.IconKey, org.IconKey) {
		iconKey, err := m.updateOrganizationIcon(ctx, org, params)
		if err != nil {
			return models.Organization{}, fmt.Errorf("failed to validate organization icon: %w", err)
		}
		params.IconKey = iconKey
		upd = true
	}

	if !ptrStrEq(params.BannerKey, org.BannerKey) {
		bannerKey, err := m.updateOrganizationBanner(ctx, org, params)
		if err != nil {
			return models.Organization{}, fmt.Errorf("failed to validate organization banner: %w", err)
		}
		params.BannerKey = bannerKey
		upd = true
	}

	if !upd {
		return org, nil
	}

	if err = m.repo.Transaction(ctx, func(ctx context.Context) error {
		org, err = m.repo.UpdateOrganization(ctx, organizationID, params)
		if err != nil {
			return fmt.Errorf("failed to update organization: %w", err)
		}

		err = m.messenger.WriteOrganizationUpdated(ctx, org)
		if err != nil {
			return fmt.Errorf("failed to publish organization updated event: %w", err)
		}

		return nil
	}); err != nil {
		return models.Organization{}, err
	}

	return org, nil
}

func ptrStrEq(a, b *string) bool {
	return (a == nil && b == nil) || (a != nil && b != nil && *a == *b)
}
