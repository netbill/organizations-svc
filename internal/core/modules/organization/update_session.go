package organization

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
)

func (m *Module) OpenUpdateSession(
	ctx context.Context,
	initiator models.AccountActor,
	organizationID uuid.UUID,
) (models.Organization, models.UpdateOrganizationMedia, error) {
	org, err := m.GetByID(ctx, organizationID)
	if err != nil {
		return models.Organization{}, models.UpdateOrganizationMedia{}, err
	}

	err = m.chekPermissionForManageOrganization(ctx, initiator, org.ID)
	if err != nil {
		return models.Organization{}, models.UpdateOrganizationMedia{}, err
	}

	uploadSessionID := uuid.New()
	uploadIconLink, getIconLink, err := m.bucket.GeneratePreloadLinkForOrganizationIcon(ctx, org.ID, uploadSessionID)
	if err != nil {
		return models.Organization{}, models.UpdateOrganizationMedia{}, err
	}

	uploadBannerLink, getBannerLink, err := m.bucket.GeneratePreloadLinkForOrganizationBanner(ctx, org.ID, uploadSessionID)
	if err != nil {
		return models.Organization{}, models.UpdateOrganizationMedia{}, err
	}

	uploadToken, err := m.token.GenerateUploadOrganizationMediaToken(
		initiator,
		organizationID,
		uploadSessionID,
	)
	if err != nil {
		return models.Organization{}, models.UpdateOrganizationMedia{}, err
	}

	return org, models.UpdateOrganizationMedia{
		Links: models.OrganizationUploadMediaLinks{
			IconUploadURL:   uploadIconLink,
			IconGetURL:      getIconLink,
			BannerUploadURL: uploadBannerLink,
			BannerGetURL:    getBannerLink,
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

func (m *Module) UpdateWithSession(
	ctx context.Context,
	initiator models.AccountActor,
	scope models.UploadScope,
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

	if params.Media.DeletedIcon {
		if err := m.bucket.DeleteOrganizationIcon(ctx, org.ID); err != nil {
			return models.Organization{}, err
		}
	} else {
		key, err := m.bucket.UpdateOrganizationIcon(ctx, org.ID, scope)
		switch {
		case errors.Is(err, errx.ErrorNoContentUploaded):
			params.Media.icon = org.Icon
		case err != nil:
			return models.Organization{}, err
		default:
			params.Media.banner = &key
		}
	}

	if params.Media.DeletedBanner {
		if err := m.bucket.DeleteOrganizationBanner(ctx, org.ID); err != nil {
			return models.Organization{}, err
		}
	} else {
		key, err := m.bucket.UpdateOrganizationBanner(ctx, org.ID, scope)
		switch {
		case errors.Is(err, errx.ErrorNoContentUploaded):
			params.Media.banner = org.Banner
		case err != nil:
			return models.Organization{}, err
		default:
			params.Media.banner = &key
		}
	}

	params.Media.icon = org.Icon
	params.Media.banner = org.Banner

	err = m.bucket.CleanOrganizationMediaSession(ctx, organizationID, scope)
	if err != nil {
		return models.Organization{}, err
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

func (m *Module) DeleteUpdateIconInSession(
	ctx context.Context,
	initiator models.AccountActor,
	organizationID, uploadSessionID uuid.UUID,
) error {
	err := m.chekPermissionForManageOrganization(ctx, initiator, organizationID)
	if err != nil {
		return err
	}

	return m.bucket.CancelUpdateOrganizationIcon(
		ctx,
		organizationID,
		uploadSessionID,
	)
}

func (m *Module) DeleteUpdateBannerInSession(
	ctx context.Context,
	initiator models.AccountActor,
	organizationID, uploadSessionID uuid.UUID,
) error {
	err := m.chekPermissionForManageOrganization(ctx, initiator, organizationID)
	if err != nil {
		return err
	}

	return m.bucket.CancelUpdateOrganizationBanner(
		ctx,
		organizationID,
		uploadSessionID,
	)
}

func (m *Module) CancelUpdateSession(
	ctx context.Context,
	initiator models.AccountActor,
	uploadSessionID uuid.UUID,
	organizationID uuid.UUID,
) error {
	err := m.chekPermissionForManageOrganization(ctx, initiator, organizationID)
	if err != nil {
		return err
	}

	return m.bucket.CleanOrganizationMediaSession(
		ctx,
		organizationID,
		uploadSessionID,
	)
}
