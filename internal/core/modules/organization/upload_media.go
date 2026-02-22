package organization

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/domain"
)

func (m *Module) CreateOrgUploadMediaLinks(
	ctx context.Context,
	actor domain.AccountActor,
	organizationID uuid.UUID,
) (domain.Organization, domain.UploadOrgMediaLinks, error) {
	org, err := m.GetByID(ctx, organizationID)
	if err != nil {
		return domain.Organization{}, domain.UploadOrgMediaLinks{}, err
	}

	err = m.chekPermissionForManageOrganization(ctx, actor, org.ID)
	if err != nil {
		return domain.Organization{}, domain.UploadOrgMediaLinks{}, err
	}

	iconLinks, err := m.bucket.CreateOrganizationIconUploadMediaLinks(ctx, organizationID)
	if err != nil {
		return domain.Organization{}, domain.UploadOrgMediaLinks{}, fmt.Errorf("failed to create upload media links: %w", err)
	}

	bannerLinks, err := m.bucket.CreateOrganizationBannerUploadMediaLinks(ctx, organizationID)
	if err != nil {
		return domain.Organization{}, domain.UploadOrgMediaLinks{}, fmt.Errorf("failed to create upload media links: %w", err)
	}

	return org, domain.UploadOrgMediaLinks{
		Icon:   iconLinks,
		Banner: bannerLinks,
	}, nil
}

func (m *Module) DeleteOrgUploadIcon(
	ctx context.Context,
	actor domain.AccountActor,
	organizationID uuid.UUID,
	key string,
) error {
	err := m.chekPermissionForManageOrganization(ctx, actor, organizationID)
	if err != nil {
		return err
	}

	err = m.bucket.DeleteUploadOrganizationIcon(ctx, organizationID, key)
	if err != nil {
		return fmt.Errorf("failed to delete upload media: %w", err)
	}

	return nil
}

func (m *Module) DeleteOrgUploadBanner(
	ctx context.Context,
	actor domain.AccountActor,
	organizationID uuid.UUID,
	key string,
) error {
	err := m.chekPermissionForManageOrganization(ctx, actor, organizationID)
	if err != nil {
		return err
	}

	err = m.bucket.DeleteUploadOrganizationBanner(ctx, organizationID, key)
	if err != nil {
		return fmt.Errorf("failed to delete upload media: %w", err)
	}

	return nil
}
