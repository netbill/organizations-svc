package bucket

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	OrganizationIconMaxW          = 512
	OrganizationIconMaxH          = 512
	OrganizationIconProbeMaxBytes = int64(512 * 1024)

	OrganizationContentLengthMax               = 5 * 1024 * 1024 // 5 MB
	OrganizationIconUploadTTL    time.Duration = 1 * time.Hour
)

func CreateTempOrganizationIconKey(organizationID, sessionID string) string {
	return fmt.Sprintf("organization/icon/%s/temp/%s", organizationID, sessionID)
}

func CreateOrganizationIconKey(organizationID string) string {
	return fmt.Sprintf("organization/icon/%s", organizationID)
}

var allowedOrganizationIconContentTypes = []string{
	"image/png",
	"image/jpeg",
	"image/jpg",
}

func isAllowedDetectedContentType(ct string) bool {
	for _, allowedCT := range allowedOrganizationIconContentTypes {
		if ct == allowedCT {
			return true
		}
	}
	return false
}

var allowedOrganizationIconExtensions = []string{
	"png",
	"jpeg",
	"jpg",
}

func isAllowedImageFormat(format string) bool {
	for _, allowedExt := range allowedOrganizationIconExtensions {
		if format == allowedExt {
			return true
		}
	}
	return false
}

func (r Bucket) GetPreloadLinkForUpdateOrganizationIcon(
	ctx context.Context,
	organizationID, sessionID string,
) (uploadURL, getUrl string, error error) {
	uploadURL, getURL, err := r.s3.PresignPut(
		ctx,
		CreateTempOrganizationIconKey(organizationID, sessionID),
		OrganizationIconUploadTTL,
	)
	if err != nil {
		return "", "", fmt.Errorf("presigning put for organization icon: %w", err)
	}

	return uploadURL, getURL, nil
}

func (b Bucket) AcceptUpdateOrganizationIcon(ctx context.Context, organizationID, sessionID uuid.UUID) (string, error) {

}

func (r Bucket) CancelUpdateOrganizationIcon(
	ctx context.Context,
	organizationID, sessionID string,
) error {
	key := CreateTempOrganizationIconKey(organizationID, sessionID)
	if err := r.s3.DeleteObject(ctx, key); err != nil {
		return fmt.Errorf("deleting temp organization icon object: %w", err)
	}

	return nil
}

func (r Bucket) DeleteOrganizationIcon(
	ctx context.Context,
	organizationID string,
) error {
	key := CreateOrganizationIconKey(organizationID)
	if err := r.s3.DeleteObject(ctx, key); err != nil {
		return fmt.Errorf("deleting organization icon object: %w", err)
	}

	return nil
}
