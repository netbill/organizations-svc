package media

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/models"
)

type Uploader struct {
	bucket Bucket
}

type Bucket interface {
	CreateOrganizationIconUploadMediaLinks(ctx context.Context, organizationID uuid.UUID) (models.UploadMediaLink, error)
	CreateOrganizationBannerUploadMediaLinks(ctx context.Context, organizationID uuid.UUID) (models.UploadMediaLink, error)

	ValidateOrganizationIcon(ctx context.Context, organizationID uuid.UUID, key string) error
	ValidateOrganizationBanner(ctx context.Context, organizationID uuid.UUID, key string) error

	UpdateOrganizationIcon(ctx context.Context, organizationID uuid.UUID, key string) (string, error)
	UpdateOrganizationBanner(ctx context.Context, organizationID uuid.UUID, key string) (string, error)

	DeleteOrganizationIcon(ctx context.Context, organizationID uuid.UUID, key string) error
	DeleteOrganizationBanner(ctx context.Context, organizationID uuid.UUID, key string) error
	DeleteUploadOrganizationIcon(ctx context.Context, organizationID uuid.UUID, key string) error
	DeleteUploadOrganizationBanner(ctx context.Context, organizationID uuid.UUID, key string) error
}

func (m *Uploader) CreateOrganizationIconUploadMediaLinks(
	ctx context.Context,
	orgID uuid.UUID,
) (models.UploadMediaLink, error) {
	return m.bucket.CreateOrganizationIconUploadMediaLinks(ctx, orgID)
}

func (m *Uploader) CreateOrganizationBannerUploadMediaLinks(
	ctx context.Context,
	orgID uuid.UUID,
) (models.UploadMediaLink, error) {
	return m.bucket.CreateOrganizationBannerUploadMediaLinks(ctx, orgID)
}

func (m *Uploader) UpdateOrganizationIcon(
	ctx context.Context,
	orgID uuid.UUID,
	oldKey *string,
	tempKey *string,
) (newKey *string, err error) {
	if ptrStrEq(oldKey, tempKey) {
		return oldKey, nil
	}

	if tempKey != nil {
		if err = m.bucket.ValidateOrganizationIcon(ctx, orgID, *tempKey); err != nil {
			return nil, fmt.Errorf("failed to validate organization icon: %w", err)
		}

		iconKey, err := m.bucket.UpdateOrganizationIcon(ctx, orgID, *tempKey)
		if err != nil {
			return nil, fmt.Errorf("failed to update organization icon: %w", err)
		}

		newKey = &iconKey
	}

	if oldKey != nil {
		if err = m.bucket.DeleteOrganizationIcon(ctx, orgID, *oldKey); err != nil {
			return nil, fmt.Errorf("failed to delete organization icon: %w", err)
		}
	}

	return newKey, nil
}

func (m *Uploader) UpdateOrganizationBanner(
	ctx context.Context,
	orgID uuid.UUID,
	oldKey *string,
	tempKey *string,
) (newKey *string, err error) {
	if ptrStrEq(oldKey, tempKey) {
		return oldKey, nil
	}

	if tempKey != nil {
		if err = m.bucket.ValidateOrganizationBanner(ctx, orgID, *tempKey); err != nil {
			return nil, fmt.Errorf("failed to validate organization banner: %w", err)
		}

		key, err := m.bucket.UpdateOrganizationBanner(ctx, orgID, *tempKey)
		if err != nil {
			return nil, fmt.Errorf("failed to update organization banner: %w", err)
		}

		newKey = &key
	}

	if oldKey != nil {
		if err = m.bucket.DeleteOrganizationBanner(ctx, orgID, *oldKey); err != nil {
			return nil, fmt.Errorf("failed to delete organization banner: %w", err)
		}
	}

	return newKey, nil
}

func (m *Uploader) DeleteUploadOrganizationIcon(
	ctx context.Context,
	orgID uuid.UUID,
	key string,
) error {
	return m.bucket.DeleteUploadOrganizationIcon(ctx, orgID, key)
}

func (m *Uploader) DeleteUploadOrganizationBanner(
	ctx context.Context,
	orgID uuid.UUID,
	key string,
) error {
	return m.bucket.DeleteUploadOrganizationBanner(ctx, orgID, key)
}

func ptrStrEq(a, b *string) bool {
	return (a == nil && b == nil) || (a != nil && b != nil && *a == *b)
}
