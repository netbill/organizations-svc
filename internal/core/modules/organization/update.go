package organization

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/models"
)

func (m *Module) OpenUpdateSession(
	ctx context.Context,
	initiator models.Initiator,
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
	links, err := m.bucket.GeneratePreloadLinkForOrganizationMedia(ctx, org.ID, uploadSessionID)
	if err != nil {
		return models.Organization{}, models.UpdateOrganizationMedia{}, err
	}

	uploadToken, err := m.token.NewUploadOrganizationMediaToken(
		initiator.GetAccountID(),
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

func (m *Module) UpdateWithSession(
	ctx context.Context,
	initiator models.Initiator,
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

	params.Media.icon = org.Icon
	params.Media.banner = org.Banner

	if params.Media.DeletedIcon == true {
		if err = m.bucket.DeleteOrganizationIcon(
			ctx,
			organizationID,
		); err != nil {
			return models.Organization{}, err
		}

		params.Media.icon = nil
	}

	if params.Media.DeletedBanner == true {
		if err = m.bucket.DeleteOrganizationBanner(
			ctx,
			organizationID,
		); err != nil {
			return models.Organization{}, err
		}

		params.Media.banner = nil
	}

	if !(params.Media.DeletedBanner == params.Media.DeletedIcon == true) {
		links, err := m.bucket.AcceptUpdateOrganizationMedia(
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

	err = m.bucket.CleanOrganizationMediaSession(
		ctx,
		organizationID,
		params.Media.UploadSessionID,
	)
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
	initiator models.Initiator,
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
	initiator models.Initiator,
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
	initiator models.Initiator,
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
