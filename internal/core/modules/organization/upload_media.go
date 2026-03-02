package organization

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
)

func (m *Module) CreateOrgUploadMediaLinks(
	ctx context.Context,
	actor models.AccountActor,
	organizationID uuid.UUID,
) (models.Organization, models.UploadOrgMediaLinks, error) {
	org, err := m.GetByID(ctx, organizationID)
	if err != nil {
		return models.Organization{}, models.UploadOrgMediaLinks{}, err
	}

	member, err := m.getInitiator(ctx, actor, organizationID)
	if err != nil {
		return models.Organization{}, models.UploadOrgMediaLinks{}, err
	}
	if !member.Head {
		return models.Organization{}, models.UploadOrgMediaLinks{}, errx.ErrorNotOrganizationHead.Raise(
			fmt.Errorf("only organization head member can activate organization, but member %s is not head", member.ID),
		)
	}

	if org.Status == models.OrganizationStatusSuspended {
		return models.Organization{}, models.UploadOrgMediaLinks{}, errx.ErrorOrganizationIsSuspended.Raise(
			fmt.Errorf("organization with id %s is suspended", organizationID),
		)
	}

	iconLinks, err := m.bucket.CreateOrganizationIconUploadMediaLinks(ctx, organizationID)
	if err != nil {
		return models.Organization{}, models.UploadOrgMediaLinks{}, fmt.Errorf("failed to create upload media links: %w", err)
	}

	bannerLinks, err := m.bucket.CreateOrganizationBannerUploadMediaLinks(ctx, organizationID)
	if err != nil {
		return models.Organization{}, models.UploadOrgMediaLinks{}, fmt.Errorf("failed to create upload media links: %w", err)
	}

	return org, models.UploadOrgMediaLinks{
		Icon:   iconLinks,
		Banner: bannerLinks,
	}, nil
}

func (m *Module) updateOrganizationIcon(
	ctx context.Context,
	organization models.Organization,
	params UpdateParams,
) (newKey *string, err error) {
	if params.IconKey != nil {
		if err = m.bucket.ValidateOrganizationIcon(ctx, organization.ID, *params.IconKey); err != nil {
			return nil, fmt.Errorf("failed to validate organization icon: %w", err)
		}

		iconKey, err := m.bucket.UpdateOrganizationIcon(ctx, organization.ID, *params.IconKey)
		if err != nil {
			return nil, fmt.Errorf("failed to update organization icon: %w", err)
		}

		if err = m.bucket.DeleteOrganizationIcon(ctx, organization.ID, iconKey); err != nil {
			return nil, fmt.Errorf("failed to delete organization icon: %w", err)
		}

		newKey = &iconKey
	}

	if organization.IconKey != nil {
		if err = m.bucket.DeleteOrganizationIcon(ctx, organization.ID, *organization.IconKey); err != nil {
			return nil, fmt.Errorf("failed to delete organization icon: %w", err)
		}
	}

	return newKey, nil
}

func (m *Module) DeleteOrgUploadIcon(
	ctx context.Context,
	actor models.AccountActor,
	organizationID uuid.UUID,
	key string,
) error {
	org, err := m.GetByID(ctx, organizationID)
	if err != nil {
		return err
	}

	if org.Status == models.OrganizationStatusSuspended {
		return errx.ErrorOrganizationIsSuspended.Raise(
			fmt.Errorf("organization with id %s is suspended", organizationID),
		)
	}

	member, err := m.getInitiator(ctx, actor, organizationID)
	if err != nil {
		return err
	}
	if !member.Head {
		return errx.ErrorNotOrganizationHead.Raise(
			fmt.Errorf("only organization head member can activate organization, but member %s is not head", member.ID),
		)
	}

	err = m.bucket.DeleteUploadOrganizationIcon(ctx, organizationID, key)
	if err != nil {
		return fmt.Errorf("failed to delete upload media: %w", err)
	}

	return nil
}

func (m *Module) updateOrganizationBanner(
	ctx context.Context,
	organization models.Organization,
	params UpdateParams,
) (newKey *string, err error) {
	if params.BannerKey != nil {
		if err = m.bucket.ValidateOrganizationBanner(ctx, organization.ID, *params.BannerKey); err != nil {
			return nil, fmt.Errorf("failed to validate organization banner: %w", err)
		}

		key, err := m.bucket.UpdateOrganizationBanner(ctx, organization.ID, *params.BannerKey)
		if err != nil {
			return nil, fmt.Errorf("failed to update organization banner: %w", err)
		}

		if err = m.bucket.DeleteOrganizationBanner(ctx, organization.ID, key); err != nil {
			return nil, fmt.Errorf("failed to delete organization banner: %w", err)
		}

		newKey = &key
	}

	if organization.BannerKey != nil {
		if err = m.bucket.DeleteOrganizationBanner(ctx, organization.ID, *organization.BannerKey); err != nil {
			return nil, fmt.Errorf("failed to delete organization banner: %w", err)
		}
	}

	return newKey, nil
}

func (m *Module) DeleteOrgUploadBanner(
	ctx context.Context,
	actor models.AccountActor,
	organizationID uuid.UUID,
	key string,
) error {
	org, err := m.GetByID(ctx, organizationID)
	if err != nil {
		return err
	}

	if org.Status == models.OrganizationStatusSuspended {
		return errx.ErrorOrganizationIsSuspended.Raise(
			fmt.Errorf("organization with id %s is suspended", organizationID),
		)
	}

	member, err := m.getInitiator(ctx, actor, organizationID)
	if err != nil {
		return err
	}
	if !member.Head {
		return errx.ErrorNotOrganizationHead.Raise(
			fmt.Errorf("only organization head member can activate organization, but member %s is not head", member.ID),
		)
	}

	err = m.bucket.DeleteUploadOrganizationBanner(ctx, organizationID, key)
	if err != nil {
		return fmt.Errorf("failed to delete upload media: %w", err)
	}

	return nil
}
