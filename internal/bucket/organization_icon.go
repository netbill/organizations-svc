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

func CreateTempOrganizationIconKey(organizationID, sessionID uuid.UUID) string {
	return fmt.Sprintf("organization/icon/%s/temp/%s", organizationID, sessionID)
}

func CreateOrganizationIconKey(organizationID uuid.UUID) string {
	return fmt.Sprintf("organization/icon/%s", organizationID)
}

func (b Bucket) GeneratePreloadLinkForOrganizationIcon(
	ctx context.Context,
	organizationID, sessionID uuid.UUID,
) (string, string, error) {
	uploadIconURL, getIconURL, err := b.s3.PresignPut(
		ctx,
		CreateTempOrganizationIconKey(organizationID, sessionID),
		b.config.Link.TTL,
	)
	if err != nil {
		return "", "", fmt.Errorf("presigning put for organization icon: %w", err)
	}

	return uploadIconURL, getIconURL, nil
}

func (b Bucket) UpdateOrganizationIcon(
	ctx context.Context,
	organizationID, sessionID uuid.UUID,
) (string, error) {
	finalKey := CreateOrganizationIconKey(organizationID)
	tempKey := CreateTempOrganizationIconKey(organizationID, sessionID)

	head, err := b.s3.HeadObject(ctx, tempKey)
	if err != nil {
		return "", fmt.Errorf("failed to head object for organization icon: %w", err)
	}

	if head.ContentLength == nil || *head.ContentLength == 0 {
		return "", errx.ErrorNoContentUploaded.Raise(
			fmt.Errorf("no content uploaded for organization icon in session %s", sessionID),
		)
	}

	rc, err := b.s3.GetObjectRange(ctx, tempKey, 2048)
	if err != nil {
		return "", fmt.Errorf("failed to get object range for organization icon: %w", err)
	}
	defer rc.Close()

	probe, err := io.ReadAll(rc)
	if err != nil {
		return "", fmt.Errorf("failed to read icon probe bytes: %w", err)
	}

	config, format, err := image.DecodeConfig(bytes.NewReader(probe))
	if err != nil {
		return "", fmt.Errorf("decode config: %w", err)
	}

	if b.config.Organization.Icon.MaxWidth > 0 && config.Width > b.config.Organization.Icon.MaxWidth {
		return "", errx.ErrorOrganizationIconInvalid.Raise(
			fmt.Errorf("uploaded organization icon width %d exceeds the maximum allowed width", config.Width),
		)
	}
	if b.config.Organization.Icon.MaxHeight > 0 && config.Height > b.config.Organization.Icon.MaxHeight {
		return "", errx.ErrorOrganizationIconInvalid.Raise(
			fmt.Errorf("uploaded organization icon height %d exceeds the maximum allowed height", config.Height),
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

	if !access(b.config.Organization.Icon.AllowedFormats, format) {
		return "", errx.ErrorOrganizationIconInvalid.Raise(
			fmt.Errorf("uploaded organization icon format %s is not allowed", format),
		)
	}

	res, err := b.s3.CopyObject(ctx, tempKey, finalKey)
	if err != nil {
		return "", fmt.Errorf("failed to copy object for organization icon: %w", err)
	}

	return res, nil
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
