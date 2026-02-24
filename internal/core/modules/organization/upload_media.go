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

	member, err := m.repo.GetMemberByAccountAndOrganization(ctx, actor, org.ID)
	if err != nil {
		return models.Organization{}, models.UploadOrgMediaLinks{}, err
	}
	if !member.Head {
		return models.Organization{}, models.UploadOrgMediaLinks{}, errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("only organization head member can activate organization, but member %s is not head", member.ID),
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

func (m *Module) DeleteOrgUploadIcon(
	ctx context.Context,
	actor models.AccountActor,
	organizationID uuid.UUID,
	key string,
) error {
	member, err := m.repo.GetMemberByAccountAndOrganization(ctx, actor, organizationID)
	if err != nil {
		return err
	}
	if !member.Head {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("only organization head member can activate organization, but member %s is not head", member.ID),
		)
	}

	err = m.bucket.DeleteUploadOrganizationIcon(ctx, organizationID, key)
	if err != nil {
		return fmt.Errorf("failed to delete upload media: %w", err)
	}

	return nil
}

func (m *Module) DeleteOrgUploadBanner(
	ctx context.Context,
	actor models.AccountActor,
	organizationID uuid.UUID,
	key string,
) error {
	member, err := m.repo.GetMemberByAccountAndOrganization(ctx, actor, organizationID)
	if err != nil {
		return err
	}
	if !member.Head {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("only organization head member can activate organization, but member %s is not head", member.ID),
		)
	}

	err = m.bucket.DeleteUploadOrganizationBanner(ctx, organizationID, key)
	if err != nil {
		return fmt.Errorf("failed to delete upload media: %w", err)
	}

	return nil
}
