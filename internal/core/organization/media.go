package organization

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/models"
)

type media interface {
	CreateOrganizationIconUploadMediaLinks(
		ctx context.Context,
		organizationID uuid.UUID,
	) (models.UploadMediaLink, error)

	UpdateOrganizationIcon(
		ctx context.Context,
		orgID uuid.UUID,
		tempKey string,
	) (newKey string, err error)

	DeleteOrganizationIcon(
		ctx context.Context,
		organizationID uuid.UUID,
		key string,
	) error

	DeleteUploadOrganizationIcon(
		ctx context.Context,
		orgID uuid.UUID,
		key string,
	) error

	CreateOrganizationBannerUploadMediaLinks(
		ctx context.Context,
		organizationID uuid.UUID,
	) (models.UploadMediaLink, error)

	UpdateOrganizationBanner(
		ctx context.Context,
		organizationID uuid.UUID,
		key string,
	) (string, error)

	DeleteOrganizationBanner(
		ctx context.Context,
		organizationID uuid.UUID,
		key string,
	) error

	DeleteUploadOrganizationBanner(
		ctx context.Context,
		orgID uuid.UUID,
		key string,
	) error
}

func (s *Service) CreateUploadMediaLinks(
	ctx context.Context,
	actor models.AccountActor,
	organizationID uuid.UUID,
) (models.Organization, models.UploadOrgMediaLinks, error) {
	_, err := s.AuthorizeOrgHead(ctx, actor, organizationID)
	if err != nil {
		return models.Organization{}, models.UploadOrgMediaLinks{}, err
	}

	org, err := s.ValidateOrg(ctx, organizationID)
	if err != nil {
		return models.Organization{}, models.UploadOrgMediaLinks{}, err
	}

	iconLinks, err := s.media.CreateOrganizationIconUploadMediaLinks(ctx, organizationID)
	if err != nil {
		return models.Organization{}, models.UploadOrgMediaLinks{}, err
	}

	bannerLinks, err := s.media.CreateOrganizationBannerUploadMediaLinks(ctx, organizationID)
	if err != nil {
		return models.Organization{}, models.UploadOrgMediaLinks{}, err
	}

	return org, models.UploadOrgMediaLinks{
		Icon:   iconLinks,
		Banner: bannerLinks,
	}, nil
}

type DeleteUploaderParams struct {
	Icon   *string
	Banner *string
}

func (s *Service) DeleteUploadMedia(
	ctx context.Context,
	actor models.AccountActor,
	organizationID uuid.UUID,
	params DeleteUploaderParams,
) error {
	_, err := s.AuthorizeOrgHead(ctx, actor, organizationID)
	if err != nil {
		return err
	}

	_, err = s.ValidateOrg(ctx, organizationID)
	if err != nil {
		return err
	}

	if params.Icon != nil {
		if err = s.media.DeleteUploadOrganizationIcon(ctx, organizationID, *params.Icon); err != nil {
			return err
		}
	}

	if params.Banner != nil {
		if err = s.media.DeleteUploadOrganizationBanner(ctx, organizationID, *params.Banner); err != nil {
			return err
		}
	}

	return nil
}
