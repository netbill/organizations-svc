package media

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/google/uuid"
	"github.com/netbill/awsx"
	"github.com/netbill/organizations-svc/internal/errx"
	"github.com/netbill/organizations-svc/internal/models"
)

var (
	tempOrganizationBannerKeyRe = regexp.MustCompile(
		`^organization/banner/([0-9a-fA-F-]{36})/temp/([0-9a-fA-F-]{36})$`,
	)

	finalOrganizationBannerKeyRe = regexp.MustCompile(
		`^organization/banner/([0-9a-fA-F-]{36})/([0-9a-fA-F-]{36})$`,
	)
)

func CreateTempOrganizationBannerKey(organizationID uuid.UUID) string {
	return fmt.Sprintf("organization/banner/%s/temp/%s", organizationID, uuid.New().String())
}

func validateFinalOrganizationBannerKey(organizationID uuid.UUID, key string) error {
	matches := finalOrganizationBannerKeyRe.FindStringSubmatch(key)
	if matches == nil {
		return errx.ErrorOrganizationUploadedBannerInvalid.Raise(
			fmt.Errorf("key %s does not match final organization organization banner key pattern", key),
		)
	}

	if matches[1] != organizationID.String() {
		return errx.ErrorOrganizationUploadedBannerInvalid.Raise(
			fmt.Errorf("key %s does not belong to organization organization %s", key, organizationID),
		)
	}

	return nil
}

func CreateOrganizationBannerKey(organizationID uuid.UUID) string {
	return fmt.Sprintf("organization/banner/%s/%s", organizationID, uuid.New().String())
}

func validateTempOrganizationBannerKey(organizationID uuid.UUID, key string) error {
	matches := tempOrganizationBannerKeyRe.FindStringSubmatch(key)
	if matches == nil {
		return errx.ErrorOrganizationUploadedBannerInvalid.Raise(
			fmt.Errorf("key %s does not match temp organization organization banner key pattern", key),
		)
	}

	if matches[1] != organizationID.String() {
		return errx.ErrorOrganizationUploadedBannerInvalid.Raise(
			fmt.Errorf("key %s does not belong to organization organization %s", key, organizationID),
		)
	}

	return nil
}

func (u *Uploader) CreateOrganizationBannerUploadMediaLinks(
	ctx context.Context,
	organizationID uuid.UUID,
) (models.UploadMediaLink, error) {
	key := CreateTempOrganizationBannerKey(organizationID)

	uploadURL, getURL, err := u.s3.PresignPut(ctx, key, u.config.LinkTTL)
	if err != nil {
		return models.UploadMediaLink{}, fmt.Errorf("presigning put for organization organization banner: %w", err)
	}

	return models.UploadMediaLink{
		Key:        key,
		PreloadUrl: getURL,
		UploadURL:  uploadURL,
	}, nil
}

func (u *Uploader) UpdateOrganizationBanner(
	ctx context.Context,
	organizationID uuid.UUID,
	key string,
) (string, error) {
	if err := validateTempOrganizationBannerKey(organizationID, key); err != nil {
		return "", err
	}

	out, err := u.s3.GetObjectRange(ctx, key, 64*1024)
	switch {
	case errors.Is(err, awsx.ErrNotFound):
		return "", errx.ErrorOrganizationUploadedBannerInvalid.Raise(
			fmt.Errorf("organization organization banner not found for key: %s", key),
		)
	case err != nil:
		return "", fmt.Errorf("get object range for organization organization banner: %w", err)
	}
	defer out.Body.Close()

	if err = u.config.OrgBanner.Validate(out); err != nil {
		return "", errx.ErrorOrganizationUploadedBannerInvalid.Raise(
			fmt.Errorf("validating organization organization banner: %w", err),
		)
	}

	finalKey := CreateOrganizationBannerKey(organizationID)

	if err := u.s3.CopyObject(ctx, key, finalKey); err != nil {
		return "", fmt.Errorf("copying object for organization banner: %w", err)
	}

	return finalKey, nil
}

func (u *Uploader) DeleteUploadOrganizationBanner(
	ctx context.Context,
	organizationID uuid.UUID,
	key string,
) error {
	if err := validateTempOrganizationBannerKey(organizationID, key); err != nil {
		return err
	}

	if err := u.s3.DeleteObject(ctx, key); err != nil {
		return fmt.Errorf("deleting temp organization organization banner object: %w", err)
	}

	return nil
}

func (u *Uploader) DeleteOrganizationBanner(
	ctx context.Context,
	organizationID uuid.UUID,
	key string,
) error {
	if err := validateFinalOrganizationBannerKey(organizationID, key); err != nil {
		return err
	}

	if err := u.s3.DeleteObject(ctx, key); err != nil {
		return fmt.Errorf("deleting organization organization banner object: %w", err)
	}

	return nil
}
