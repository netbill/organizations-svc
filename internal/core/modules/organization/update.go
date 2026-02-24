package organization

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/models"
)

type UpdateParams struct {
	Name      string
	IconKey   *string
	BannerKey *string
}

func (m *Module) Update(
	ctx context.Context,
	initiator models.AccountActor,
	organizationID uuid.UUID,
	params UpdateParams,
) (models.Organization, error) {
	org, err := m.GetByID(ctx, organizationID)
	if err != nil {
		return models.Organization{}, err
	}

	err = m.chekPermissionForManageOrganization(ctx, initiator, org.ID)
	if err != nil {
		return models.Organization{}, err
	}

	if params.IconKey != nil {
		err = m.bucket.ValidateOrganizationIcon(ctx, organizationID, *params.IconKey)
		if err != nil {
			return models.Organization{}, fmt.Errorf("failed to validate organization icon: %w", err)
		}
	}

	if params.BannerKey != nil {
		err = m.bucket.ValidateOrganizationBanner(ctx, organizationID, *params.BannerKey)
		if err != nil {
			return models.Organization{}, fmt.Errorf("failed to validate organization banner: %w", err)
		}
	}

	avatarKey, err := m.bucket.UpdateOrganizationIcon(ctx, organizationID, org.IconKey, params.IconKey)
	if err != nil {
		return models.Organization{}, fmt.Errorf("failed to update organization icon: %w", err)
	}
	params.IconKey = avatarKey

	bannerKey, err := m.bucket.UpdateOrganizationBanner(ctx, organizationID, org.BannerKey, params.BannerKey)
	if err != nil {
		return models.Organization{}, fmt.Errorf("failed to update organization banner: %w", err)
	}
	params.BannerKey = bannerKey

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
