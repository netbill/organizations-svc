package bucket

import (
	"context"
	"fmt"
	"io"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
)

func CreateTempOrganizationBannerKey(organizationID, sessionID uuid.UUID) string {
	return fmt.Sprintf("organization/banner/%s/temp/%s", organizationID, sessionID)
}

func CreateOrganizationBannerKey(organizationID uuid.UUID) string {
	return fmt.Sprintf("organization/banner/%s", organizationID)
}

func CreateTempOrganizationIconKey(organizationID, sessionID uuid.UUID) string {
	return fmt.Sprintf("organization/icon/%s/temp/%s", organizationID, sessionID)
}

func CreateOrganizationIconKey(organizationID uuid.UUID) string {
	return fmt.Sprintf("organization/icon/%s", organizationID)
}

func (b Bucket) GeneratePreloadLinkForOrganizationMedia(
	ctx context.Context,
	organizationID, sessionID uuid.UUID,
) (models.OrganizationUploadMediaLinks, error) {
	uploadBannerURL, getBannerURL, err := b.s3.PresignPut(
		ctx,
		CreateTempOrganizationBannerKey(organizationID, sessionID),
		b.tokensTTL.Org,
	)
	if err != nil {
		return models.OrganizationUploadMediaLinks{}, fmt.Errorf("presigning put for organization banner: %w", err)
	}

	uploadIconURL, getIconURL, err := b.s3.PresignPut(
		ctx,
		CreateTempOrganizationIconKey(organizationID, sessionID),
		b.tokensTTL.Org,
	)
	if err != nil {
		return models.OrganizationUploadMediaLinks{}, fmt.Errorf("presigning put for organization icon: %w", err)
	}

	return models.OrganizationUploadMediaLinks{
		BannerUploadURL: uploadBannerURL,
		BannerGetURL:    getBannerURL,
		IconUploadURL:   uploadIconURL,
		IconGetURL:      getIconURL,
	}, nil
}

func (b Bucket) AcceptUpdateOrganizationMedia(
	ctx context.Context,
	organizationID, sessionID uuid.UUID,
) (models.OrganizationMedia, error) {
	IconFinalKey := CreateOrganizationIconKey(organizationID)
	IconTempKey := CreateTempOrganizationIconKey(organizationID, sessionID)

	icon, size, err := b.s3.GetObjectRange(ctx, IconTempKey, 2048)
	if err != nil {
		return models.OrganizationMedia{}, fmt.Errorf("failed to get object range for organization icon: %w", err)
	}
	defer icon.Close()

	uploadIcon := false

	if size != 0 {
		uploadIcon = true
		if err = b.ValidateOrganizationIconUpload(icon, size); err != nil {
			return models.OrganizationMedia{}, err
		}
	}

	BannerFinalKey := CreateOrganizationBannerKey(organizationID)
	BannerTempKey := CreateTempOrganizationBannerKey(organizationID, sessionID)

	banner, size, err := b.s3.GetObjectRange(ctx, BannerTempKey, 2048)
	if err != nil {
		return models.OrganizationMedia{}, fmt.Errorf("failed to get object range for organization banner: %w", err)
	}
	defer banner.Close()

	uploadBanner := false

	if size != 0 {
		uploadBanner = true
		if err = b.ValidateOrganizationBannerUpload(banner, size); err != nil {
			return models.OrganizationMedia{}, err
		}
	}

	res := models.OrganizationMedia{}

	if uploadIcon {
		link, err := b.s3.CopyObject(ctx, IconTempKey, IconFinalKey)
		if err != nil {
			return models.OrganizationMedia{}, fmt.Errorf("failed to copy object for organization icon: %w", err)
		}
		res.Icon = &link
	}

	if uploadBanner {
		link, err := b.s3.CopyObject(ctx, BannerTempKey, BannerFinalKey)
		if err != nil {
			return models.OrganizationMedia{}, fmt.Errorf("failed to copy object for organization banner: %w", err)
		}
		res.Icon = &link
	}

	return res, nil
}

func (b Bucket) ValidateOrganizationIconUpload(
	icon io.Reader,
	size int64,
) error {
	probe, err := io.ReadAll(icon)
	if err != nil {
		return fmt.Errorf("failed to read icon probe bytes: %w", err)
	}

	valid, err := b.OrgIconValidator.ValidateImageSize(uint(size))
	if err != nil {
		return fmt.Errorf("failed to validate organization icon image size: %w", err)
	}
	if !valid {
		return errx.ErrorOrganizationIconTooLarge.Raise(
			fmt.Errorf("uploaded organization icon size %d exceeds the maximum allowed size", size),
		)
	}

	valid, err = b.OrgIconValidator.ValidateImageContentType(probe)
	if err != nil {
		return fmt.Errorf("failed to validate organization icon content type: %w", err)
	}
	if !valid {
		return errx.ErrorOrganizationIconContentTypeNotAllowed.Raise(
			fmt.Errorf("uploaded organization icon content type is not allowed"),
		)
	}

	valid, err = b.OrgIconValidator.ValidateImageFormat(probe)
	if err != nil {
		return fmt.Errorf("failed to validate organization icon image format: %w", err)
	}
	if !valid {
		return errx.ErrorOrganizationIconContentFormatNotAllowed.Raise(
			fmt.Errorf("uploaded organization icon format is not allowed"),
		)
	}

	valid, err = b.OrgIconValidator.ValidateImageResolution(probe)
	if err != nil {
		return fmt.Errorf("failed to validate organization icon image resolution: %w", err)
	}
	if !valid {
		return errx.ErrorOrganizationIconResolutionNotAllowed.Raise(
			fmt.Errorf("uploaded organization icon has invalid image resolution"),
		)
	}

	return nil
}

func (b Bucket) ValidateOrganizationBannerUpload(
	banner io.Reader,
	size int64,
) error {
	probe, err := io.ReadAll(banner)
	if err != nil {
		return fmt.Errorf("failed to read banner probe bytes: %w", err)
	}

	valid, err := b.OrgBannerValidator.ValidateImageSize(uint(size))
	if err != nil {
		return fmt.Errorf("failed to validate organization banner image size: %w", err)
	}
	if !valid {
		return errx.ErrorOrganizationBannerTooLarge.Raise(
			fmt.Errorf("uploaded organization banner size %d exceeds the maximum allowed size", size),
		)
	}

	valid, err = b.OrgBannerValidator.ValidateImageContentType(probe)
	if err != nil {
		return fmt.Errorf("failed to validate organization banner content type: %w", err)
	}
	if !valid {
		return errx.ErrorOrganizationBannerContentTypeNotAllowed.Raise(
			fmt.Errorf("uploaded organization banner content type is not allowed"),
		)
	}

	valid, err = b.OrgBannerValidator.ValidateImageFormat(probe)
	if err != nil {
		return fmt.Errorf("failed to validate organization banner image format: %w", err)
	}
	if !valid {
		return errx.ErrorOrganizationBannerContentFormatNotAllowed.Raise(
			fmt.Errorf("uploaded organization banner format is not allowed"),
		)
	}

	valid, err = b.OrgBannerValidator.ValidateImageResolution(probe)
	if err != nil {
		return fmt.Errorf("failed to validate organization banner image resolution: %w", err)
	}
	if !valid {
		return errx.ErrorOrganizationBannerResolutionNotAllowed.Raise(
			fmt.Errorf("uploaded organization banner has invalid image resolution"),
		)
	}

	return nil
}

func (b Bucket) CancelUpdateOrganizationIcon(
	ctx context.Context,
	organizationID, sessionID uuid.UUID,
) error {
	key := CreateTempOrganizationIconKey(organizationID, sessionID)

	if err := b.s3.DeleteObject(ctx, key); err != nil {
		return fmt.Errorf("deleting temp organization icon object: %w", err)
	}

	return nil
}

func (b Bucket) CancelUpdateOrganizationBanner(
	ctx context.Context,
	organizationID, sessionID uuid.UUID,
) error {
	key := CreateTempOrganizationBannerKey(organizationID, sessionID)

	if err := b.s3.DeleteObject(ctx, key); err != nil {
		return fmt.Errorf("deleting temp organization banner object: %w", err)
	}

	return nil
}

func (b Bucket) DeleteOrganizationIcon(
	ctx context.Context,
	organizationID uuid.UUID,
) error {
	key := CreateOrganizationIconKey(organizationID)
	if err := b.s3.DeleteObject(ctx, key); err != nil {
		return fmt.Errorf("deleting organization icon object: %w", err)
	}

	return nil
}

func (b Bucket) DeleteOrganizationBanner(
	ctx context.Context,
	organizationID uuid.UUID,
) error {
	key := CreateOrganizationBannerKey(organizationID)
	if err := b.s3.DeleteObject(ctx, key); err != nil {
		return fmt.Errorf("deleting organization banner object: %w", err)
	}

	return nil
}

func (b Bucket) CleanOrganizationMediaSession(
	ctx context.Context,
	organizationID, sessionID uuid.UUID,
) error {
	err := b.s3.DeleteObject(ctx, CreateTempOrganizationIconKey(organizationID, sessionID))
	if err != nil {
		return fmt.Errorf(
			"failed to delete temp object for organization icon: %w", err,
		)
	}

	err = b.s3.DeleteObject(ctx, CreateTempOrganizationBannerKey(organizationID, sessionID))
	if err != nil {
		return fmt.Errorf(
			"failed to delete temp object for organization banner: %w", err,
		)
	}

	return nil
}
