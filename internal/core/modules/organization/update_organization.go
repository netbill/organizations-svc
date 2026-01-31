package organization

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/models"
)

func (s Service) OpenUpdateOrganizationSession(
	ctx context.Context,
	accountID, organizationID uuid.UUID,
) (models.Organization, models.UpdateOrganizationMedia, error) {
	org, err := s.GetOrganization(ctx, organizationID)
	if err != nil {
		return models.Organization{}, models.UpdateOrganizationMedia{}, err
	}

	if err = s.chekPermissionForManageOrganization(ctx, accountID, org.ID); err != nil {
		return models.Organization{}, models.UpdateOrganizationMedia{}, err
	}

	uploadSessionID := uuid.New()
	links, err := s.bucket.GeneratePreloadLinkForOrganizationMedia(ctx, org.ID, uploadSessionID)
	if err != nil {
		return models.Organization{}, models.UpdateOrganizationMedia{}, err
	}

	uploadToken, err := s.token.NewUploadOrganizationMediaToken(
		accountID,
		organizationID,
		uploadSessionID,
	)
	if err != nil {
		return models.Organization{}, models.UpdateOrganizationMedia{}, err
	}

	return org, models.UpdateOrganizationMedia{
		Links: models.OrganizationUploadMediaLinks{
			IconUploadURL:   links.IconUploadURL,
			IconGetURL:      links.IconGetURL,
			BannerUploadURL: links.BannerUploadURL,
			BannerGetURL:    links.BannerGetURL,
		},
		UploadSessionID: uploadSessionID,
		UploadToken:     uploadToken,
	}, nil
}

type UpdateParams struct {
	Name  string
	Media UpdateMediaParams
}

type UpdateMediaParams struct {
	UploadSessionID uuid.UUID

	DeletedIcon   bool
	icon          *string
	DeletedBanner bool
	banner        *string
}

func (p UpdateParams) GetUpdatedIcon() *string {
	if p.Media.DeletedIcon {
		return nil
	}
	return p.Media.icon
}

func (p UpdateParams) GetUpdatedBanner() *string {
	if p.Media.DeletedBanner {
		return nil
	}
	return p.Media.banner
}

func (s Service) UpdateOrganization(
	ctx context.Context,
	accountID uuid.UUID,
	organizationID uuid.UUID,
	params UpdateParams,
) (models.Organization, error) {
	org, err := s.GetOrganization(ctx, organizationID)
	if err != nil {
		return models.Organization{}, err
	}

	if err = s.chekPermissionForManageOrganization(ctx, accountID, org.ID); err != nil {
		return models.Organization{}, err
	}

	params.Media.icon = org.Icon
	params.Media.banner = org.Banner

	if params.Media.DeletedIcon == true {
		if err = s.bucket.DeleteOrganizationIcon(
			ctx,
			organizationID,
		); err != nil {
			return models.Organization{}, err
		}

		params.Media.icon = nil
	}

	if params.Media.DeletedBanner == true {
		if err = s.bucket.DeleteOrganizationBanner(
			ctx,
			organizationID,
		); err != nil {
			return models.Organization{}, err
		}

		params.Media.banner = nil
	}

	if !(params.Media.DeletedBanner == params.Media.DeletedIcon == true) {
		links, err := s.bucket.AcceptUpdateOrganizationMedia(
			ctx,
			organizationID,
			params.Media.UploadSessionID,
		)
		if err != nil {
			return models.Organization{}, err
		}

		params.Media.icon = links.Icon
		params.Media.banner = links.Banner
	}

	err = s.bucket.CleanOrganizationMediaSession(
		ctx,
		organizationID,
		params.Media.UploadSessionID,
	)
	if err != nil {
		return models.Organization{}, err
	}

	if err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		org, err = s.repo.UpdateOrganization(ctx, organizationID, params)
		if err != nil {
			return fmt.Errorf("failed to update organization: %w", err)
		}

		err = s.messenger.WriteOrganizationUpdated(ctx, org)
		if err != nil {
			return fmt.Errorf("failed to publish organization updated event: %w", err)
		}

		return nil
	}); err != nil {
		return models.Organization{}, err
	}

	return org, nil
}

func (s Service) DeleteUpdateOrganizationIconInSession(
	ctx context.Context,
	accountID, organizationID, uploadSessionID uuid.UUID,
) error {
	org, err := s.GetOrganization(ctx, organizationID)
	if err != nil {
		return err
	}

	if err = s.chekPermissionForManageOrganization(ctx, accountID, org.ID); err != nil {
		return err
	}

	return s.bucket.CancelUpdateOrganizationIcon(
		ctx,
		organizationID,
		uploadSessionID,
	)
}

func (s Service) DeleteUpdateOrganizationBannerInSession(
	ctx context.Context,
	accountID, organizationID, uploadSessionID uuid.UUID,
) error {
	org, err := s.GetOrganization(ctx, organizationID)
	if err != nil {
		return err
	}

	if err = s.chekPermissionForManageOrganization(ctx, accountID, org.ID); err != nil {
		return err
	}

	return s.bucket.CancelUpdateOrganizationBanner(
		ctx,
		organizationID,
		uploadSessionID,
	)
}
