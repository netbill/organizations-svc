package bucket

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

const (
	OrganizationBannerMaxW = 2048
	OrganizationBannerMaxH = 1048

	OrganizationBannerContentLengthMax = 5 * 1024 * 1024 // 5 MB
)

func CreateTempOrganizationBannerKey(organizationID, sessionID uuid.UUID) string {
	return fmt.Sprintf("organization/banner/%s/temp/%s", organizationID, sessionID)
}

func CreateOrganizationBannerKey(organizationID uuid.UUID) string {
	return fmt.Sprintf("organization/banner/%s", organizationID)
}

var allowedOrganizationBannerContentTypes = []string{
	"image/png",
	"image/jpeg",
	"image/jpg",
}

var allowedOrganizationBannerExtensions = []string{
	"png",
	"jpeg",
	"jpg",
}

func (b Bucket) GeneratePreloadLinkForUpdateOrganizationBanner(
	ctx context.Context,
	organizationID, sessionID uuid.UUID,
) (uploadURL, getUrl string, error error) {
	uploadURL, getURL, err := b.s3.PresignPut(
		ctx,
		CreateTempOrganizationBannerKey(organizationID, sessionID),
		OrganizationUploadTTL,
	)
	if err != nil {
		return "", "", fmt.Errorf("presigning put for organization banner: %w", err)
	}

	finalKey := CreateOrganizationBannerKey(organizationID)
	tempKey := CreateTempOrganizationBannerKey(organizationID, sessionID)

	_, err = b.s3.CopyObject(ctx, finalKey, tempKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to copy object for organization banner: %w", err)
	}

	return uploadURL, getURL, nil
}

func (b Bucket) GetContentLengthForOrganizationBanner(
	ctx context.Context,
	organizationID, sessionID uuid.UUID,
) (int64, error) {
	key := CreateTempOrganizationBannerKey(organizationID, sessionID)
	_, contentLength, err := b.s3.GetObjectRange(ctx, key, probeBytes)
	if err != nil {
		return 0, fmt.Errorf("getting content length for organization banner: %w", err)
	}

	return contentLength, nil
}

func (b Bucket) AcceptUpdateOrganizationBanner(ctx context.Context, organizationID, sessionID uuid.UUID) (string, error) {
	tempKey := CreateTempOrganizationBannerKey(organizationID, sessionID)
	finalKey := CreateOrganizationBannerKey(organizationID)

	if err := b.validateImgObjet(
		ctx,
		tempKey,
		OrganizationBannerContentLengthMax,
		allowedOrganizationBannerContentTypes,
		allowedOrganizationBannerExtensions,
		OrganizationBannerMaxW,
		OrganizationBannerMaxH,
	); err != nil {
		return "", err
	}

	res, err := b.s3.CopyObject(ctx, tempKey, finalKey)
	if err != nil {
		return "", fmt.Errorf("failed to copy object for organization banner: %w", err)
	}

	err = b.s3.DeleteObject(ctx, tempKey)
	if err != nil {
		return "", fmt.Errorf("failed to delete temp object for organization banner: %w", err)
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
