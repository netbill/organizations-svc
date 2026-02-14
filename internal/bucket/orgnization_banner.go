package bucket

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"io"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
)

func CreateTempOrganizationBannerKey(organizationID, sessionID uuid.UUID) string {
	return fmt.Sprintf("organization/banner/%s/temp/%s", organizationID, sessionID)
}

func CreateOrganizationBannerKey(organizationID uuid.UUID) string {
	return fmt.Sprintf("organization/banner/%s", organizationID)
}

func (b Bucket) GeneratePreloadLinkForOrganizationBanner(
	ctx context.Context,
	organizationID, sessionID uuid.UUID,
) (string, string, error) {
	uploadBannerURL, getBannerURL, err := b.s3.PresignPut(
		ctx,
		CreateTempOrganizationBannerKey(organizationID, sessionID),
		b.config.Link.TTL,
	)
	if err != nil {
		return "", "", fmt.Errorf("presigning put for organization banner: %w", err)
	}

	return uploadBannerURL, getBannerURL, nil
}

func (b Bucket) UpdateOrganizationBanner(
	ctx context.Context,
	organizationID, sessionID uuid.UUID,
) (string, error) {
	finalKey := CreateOrganizationBannerKey(organizationID)
	tempKey := CreateTempOrganizationBannerKey(organizationID, sessionID)

	head, err := b.s3.HeadObject(ctx, tempKey)
	if err != nil {
		return "", fmt.Errorf("failed to head object for organization banner: %w", err)
	}

	if head.ContentLength == nil || *head.ContentLength == 0 {
		return "", errx.ErrorNoContentUploaded.Raise(
			fmt.Errorf("no content uploaded for organization banner in session %s", sessionID),
		)
	}

	rc, err := b.s3.GetObjectRange(ctx, tempKey, 2048)
	if err != nil {
		return "", fmt.Errorf("failed to get object range for organization banner: %w", err)
	}
	defer rc.Close()

	probe, err := io.ReadAll(rc)
	if err != nil {
		return "", fmt.Errorf("failed to read banner probe bytes: %w", err)
	}

	config, format, err := image.DecodeConfig(bytes.NewReader(probe))
	if err != nil {
		return "", fmt.Errorf("decode config: %w", err)
	}

	if b.config.Organization.Banner.MaxWidth > 0 && config.Width > b.config.Organization.Banner.MaxWidth {
		return "", errx.ErrorOrganizationBannerInvalid.Raise(
			fmt.Errorf("uploaded organization banner width %d exceeds the maximum allowed width", config.Width),
		)
	}
	if b.config.Organization.Banner.MaxHeight > 0 && config.Height > b.config.Organization.Banner.MaxHeight {
		return "", errx.ErrorOrganizationBannerInvalid.Raise(
			fmt.Errorf("uploaded organization banner height %d exceeds the maximum allowed height", config.Height),
		)
	}

	access := func(values []string, target string) bool {
		for _, v := range values {
			if v == target {
				return true
			}
		}
		return false
	}

	if !access(b.config.Organization.Banner.AllowedFormats, format) {
		return "", errx.ErrorOrganizationBannerInvalid.Raise(
			fmt.Errorf("uploaded organization banner format %s is not allowed", format),
		)
	}

	res, err := b.s3.CopyObject(ctx, tempKey, finalKey)
	if err != nil {
		return "", fmt.Errorf("failed to copy object for organization banner: %w", err)
	}

	return res, nil
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
