package core

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/models"
)

func (m *OrganizationModule) CreateUploadMediaLinks(
	ctx context.Context,
	actor models.AccountActor,
	organizationID uuid.UUID,
) (models.Organization, models.UploadOrgMediaLinks, error) {
	_, err := m.authorizeOrgHead(ctx, actor, organizationID)
	if err != nil {
		return models.Organization{}, models.UploadOrgMediaLinks{}, err
	}

	org, err := m.validateOrg(ctx, organizationID)
	if err != nil {
		return models.Organization{}, models.UploadOrgMediaLinks{}, err
	}

	iconLinks, err := m.media.CreateOrganizationIconUploadMediaLinks(ctx, organizationID)
	if err != nil {
		return models.Organization{}, models.UploadOrgMediaLinks{}, err
	}

	bannerLinks, err := m.media.CreateOrganizationBannerUploadMediaLinks(ctx, organizationID)
	if err != nil {
		return models.Organization{}, models.UploadOrgMediaLinks{}, err
	}

	return org, models.UploadOrgMediaLinks{
		Icon:   iconLinks,
		Banner: bannerLinks,
	}, nil
}

type DeleteUploadOrgMediaParams struct {
	Icon   *string
	Banner *string
}

func (m *OrganizationModule) DeleteUploadMedia(
	ctx context.Context,
	actor models.AccountActor,
	organizationID uuid.UUID,
	params DeleteUploadOrgMediaParams,
) error {
	_, err := m.authorizeOrgHead(ctx, actor, organizationID)
	if err != nil {
		return err
	}

	_, err = m.validateOrg(ctx, organizationID)
	if err != nil {
		return err
	}

	if params.Icon != nil {
		if err = m.media.DeleteUploadOrganizationIcon(ctx, organizationID, *params.Icon); err != nil {
			return err
		}
	}

	if params.Banner != nil {
		if err = m.media.DeleteUploadOrganizationBanner(ctx, organizationID, *params.Banner); err != nil {
			return err
		}
	}

	return nil
}
