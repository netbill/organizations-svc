package organization

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
)

func (s Service) StartUpdateSession(
	ctx context.Context,
	accountID, organizationID uuid.UUID,
) (models.OrganizationUploadMediaLinks, error) {
	org, err := s.GetOrganization(ctx, organizationID)
	if err != nil {
		return models.OrganizationUploadMediaLinks{}, err
	}

	if err = s.chekPermissionForManageOrganization(ctx, accountID, org.ID); err != nil {
		return models.OrganizationUploadMediaLinks{}, err
	}

	uploadSessionID := uuid.New()

	iconUploadLink, iconGetLink, err := s.bucket.GeneratePreloadLinkForUpdateOrganizationIcon(
		ctx,
		organizationID,
		uploadSessionID,
	)
	if err != nil {
		return models.OrganizationUploadMediaLinks{}, fmt.Errorf(
			"failed to generate preload link for update organization icon: %w", err,
		)
	}

	bannerUploadLink, bannerGetLink, err := s.bucket.GeneratePreloadLinkForUpdateOrganizationBanner(
		ctx,
		organizationID,
		uploadSessionID,
	)
	if err != nil {
		return models.OrganizationUploadMediaLinks{}, fmt.Errorf(
			"failed to generate preload link for update organization banner: %w", err,
		)
	}

	uploadToken, err := s.token.NewUploadOrganizationMediaToken(
		accountID,
		organizationID,
		uploadSessionID,
	)
	if err != nil {
		return models.OrganizationUploadMediaLinks{}, fmt.Errorf(
			"failed to generate upload token for organization media: %w", err,
		)
	}

	return models.OrganizationUploadMediaLinks{
		IconUploadLink:   iconUploadLink,
		IconGetLink:      iconGetLink,
		BannerUploadLink: bannerUploadLink,
		BannerGetLink:    bannerGetLink,
		UploadToken:      uploadToken,
	}, nil
}

type UpdateParams struct {
	Name   string
	icon   *string
	banner *string
}

func (p *UpdateParams) Icon() *string {
	return p.icon
}

func (p *UpdateParams) Banner() *string {
	return p.banner
}

func (p *UpdateParams) SetIcon(icon string) *UpdateParams {
	p.icon = &icon
	return p
}

func (p *UpdateParams) SetBanner(s string) *UpdateParams {
	p.banner = &s
	return p
}

func (s Service) UpdateOrganization(
	ctx context.Context,
	accountID uuid.UUID,
	organizationID uuid.UUID,
	uploadSessionID uuid.UUID,
	params UpdateParams,
) (models.Organization, error) {
	org, err := s.GetOrganization(ctx, organizationID)
	if err != nil {
		return models.Organization{}, err
	}

	if err = s.chekPermissionForManageOrganization(ctx, accountID, org.ID); err != nil {
		return models.Organization{}, err
	}

	anyUpdated := false

	icon, banner, err := s.updateOrganizationMedia(
		ctx,
		organizationID,
		uploadSessionID,
	)
	if err != nil {
		return models.Organization{}, fmt.Errorf("failed to update organization media: %w", err)
	}

	params.icon = icon
	if params.Name != nil && *params.Name != org.Name {
		anyUpdated = true
	}

	params.icon = org.Icon
	if icon != org.Icon {
		anyUpdated = true
		params.icon = icon
	}

	params.banner = org.Banner
	if banner != org.Banner {
		anyUpdated = true
		params.banner = banner
	}

	if anyUpdated {
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
	}

	return org, nil
}

func (s Service) updateOrganizationMedia(
	ctx context.Context,
	organizationID uuid.UUID,
	uploadSessionID uuid.UUID,
) (icon *string, banner *string, err error) {
	iconCL, err := s.bucket.GetContentLengthForOrganizationIcon(
		ctx,
		organizationID,
		uploadSessionID,
	)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"failed to get content length for organization icon: %w", err,
		)
	}

	if iconCL > 0 {
		validate, err := s.bucket.ValidateOrganizationIconSize(
			ctx,
			organizationID,
			uploadSessionID,
		)
		if err != nil {
			return nil, nil, fmt.Errorf(
				"failed to validate organization icon size: %w", err,
			)
		}
		if !validate {
			return nil, nil, errx.ErrorOrganizationIconTooLarge.Raise(err)
		}

		img, err := s.bucket.GetLoadedOrganizationIcon(
			ctx,
			organizationID,
			uploadSessionID,
		)
		if err != nil {
			return nil, nil, fmt.Errorf(
				"failed to get object for organization icon: %w", err,
			)
		}

		validate, err = s.bucket.ValidateOrganizationIconFormat(ctx, img)
		if err != nil {
			return nil, nil, fmt.Errorf(
				"failed to validate organization icon format: %w", err,
			)
		}
		if !validate {
			return nil, nil, errx.ErrorOrganizationIconContentFormatNotAllowed.Raise(err)
		}

		validate, err = s.bucket.ValidateOrganizationIconContentType(
			ctx,
			organizationID,
			uploadSessionID,
		)
		if err != nil {
			return nil, nil, fmt.Errorf(
				"failed to validate organization icon content type: %w", err,
			)
		}
		if !validate {
			return nil, nil, errx.ErrorOrganizationBannerContentTypeNotAllowed.Raise(err)
		}
	}

	avatarCL, err := s.bucket.GetContentLengthForOrganizationBanner(
		ctx,
		organizationID,
		uploadSessionID,
	)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"failed to get content length for organization banner: %w", err,
		)
	}

	if avatarCL > 0 {
		validate, err := s.bucket.ValidateOrganizationBannerSize(
			ctx,
			organizationID,
			uploadSessionID,
		)
		if err != nil {
			return nil, nil, fmt.Errorf(
				"failed to validate organization banner size: %w", err,
			)
		}
		if !validate {
			return nil, nil, errx.ErrorOrganizationBannerTooLarge.Raise(err)
		}

		img, err := s.bucket.GetLoadedOrganizationBanner(
			ctx,
			organizationID,
			uploadSessionID,
		)
		if err != nil {
			return nil, nil, fmt.Errorf(
				"failed to get object for organization banner: %w", err,
			)
		}

		validate, err = s.bucket.ValidateOrganizationBannerFormat(ctx, img)
		if err != nil {
			return nil, nil, fmt.Errorf(
				"failed to validate organization banner format: %w", err,
			)
		}
		if !validate {
			return nil, nil, errx.ErrorOrganizationBannerContentFormatNotAllowed.Raise(err)
		}

		validate, err = s.bucket.ValidateOrganizationBannerContentType(
			ctx,
			organizationID,
			uploadSessionID,
		)
		if err != nil {
			return nil, nil, fmt.Errorf(
				"failed to validate organization banner content type: %w", err,
			)
		}
		if !validate {
			return nil, nil, errx.ErrorOrganizationBannerContentTypeNotAllowed.Raise(err)
		}
	}

	icon, banner, err = s.bucket.AcceptUpdateOrganizationMedia(
		ctx,
		organizationID,
		uploadSessionID,
	)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"failed to accept update organization media in bucket: %w", err,
		)
	}

	return icon, banner, err
}
