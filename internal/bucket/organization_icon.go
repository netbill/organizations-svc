package bucket

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

const (
	OrganizationIconMaxW = 512
	OrganizationIconMaxH = 512

	OrganizationIconContentLengthMax = 5 * 1024 * 1024 // 5 MB
)

func CreateTempOrganizationIconKey(organizationID, sessionID uuid.UUID) string {
	return fmt.Sprintf("organization/icon/%s/temp/%s", organizationID, sessionID)
}

func CreateOrganizationIconKey(organizationID uuid.UUID) string {
	return fmt.Sprintf("organization/icon/%s", organizationID)
}

var allowedOrganizationIconContentTypes = []string{
	"image/png",
	"image/jpeg",
	"image/jpg",
}

var allowedOrganizationIconExtensions = []string{
	"png",
	"jpeg",
	"jpg",
}

func (b Bucket) GeneratePreloadLinkForUpdateOrganizationIcon(
	ctx context.Context,
	organizationID, sessionID uuid.UUID,
) (uploadURL, getUrl string, error error) {
	uploadURL, getURL, err := b.s3.PresignPut(
		ctx,
		CreateTempOrganizationIconKey(organizationID, sessionID),
		OrganizationUploadTTL,
	)
	if err != nil {
		return "", "", fmt.Errorf("presigning put for organization icon: %w", err)
	}

	finalKey := CreateOrganizationIconKey(organizationID)
	tempKey := CreateTempOrganizationIconKey(organizationID, sessionID)

	_, err = b.s3.CopyObject(ctx, finalKey, tempKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to copy object for organization icon: %w", err)
	}

	return uploadURL, getURL, nil
}

func (b Bucket) GetContentLengthForOrganizationIcon(
	ctx context.Context,
	organizationID, sessionID uuid.UUID,
) (int64, error) {
	key := CreateTempOrganizationIconKey(organizationID, sessionID)
	_, contentLength, err := b.s3.GetObjectRange(ctx, key, probeBytes)
	if err != nil {
		return 0, fmt.Errorf("getting content length for organization icon: %w", err)
	}

	return contentLength, nil
}

func (b Bucket) ValidateOrganizationIconSize(
	ctx context.Context,
	organizationID, sessionID uuid.UUID,
) error {
	key := CreateTempOrganizationIconKey(organizationID, sessionID)
	return b.validateImgSize(
		ctx,
		key,
		OrganizationIconMaxW,
		OrganizationIconMaxH,
	)
}

func (b Bucket) AcceptUpdateOrganizationIcon(ctx context.Context, organizationID, sessionID uuid.UUID) (string, error) {
	tempKey := CreateTempOrganizationIconKey(organizationID, sessionID)
	finalKey := CreateOrganizationIconKey(organizationID)

	if err := b.validateImgObjet(
		ctx,
		tempKey,
		OrganizationIconContentLengthMax,
		allowedOrganizationIconContentTypes,
		allowedOrganizationIconExtensions,
		OrganizationIconMaxW,
		OrganizationIconMaxH,
	); err != nil {
		return "", err
	}

	res, err := b.s3.CopyObject(ctx, finalKey, tempKey)
	if err != nil {
		return "", fmt.Errorf("failed to copy object for organization icon: %w", err)
	}
	err = b.s3.DeleteObject(ctx, tempKey)
	if err != nil {
		return "", fmt.Errorf("failed to delete temp object for organization icon: %w", err)
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
