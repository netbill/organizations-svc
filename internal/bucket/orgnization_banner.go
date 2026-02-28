package bucket

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/google/uuid"
	"github.com/netbill/awsx"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
)

func CreateTempOrganizationBannerKey(organizationID uuid.UUID) string {
	return fmt.Sprintf("organization/banner/%s/temp/%s", organizationID, uuid.New().String())
}

func CreateOrganizationBannerKey(organizationID uuid.UUID) string {
	return fmt.Sprintf("organization/banner/%s/%s", organizationID, uuid.New().String())
}

func (s *Storage) CreateOrganizationBannerUploadMediaLinks(
	ctx context.Context,
	organizationID uuid.UUID,
) (models.UploadMediaLink, error) {
	key := CreateTempOrganizationBannerKey(organizationID)

	uploadURL, getURL, err := s.s3.PresignPut(
		ctx,
		key,
		s.config.LinkTTL,
	)
	if err != nil {
		return models.UploadMediaLink{}, fmt.Errorf("presigning put for organization organization banner: %w", err)
	}

	return models.UploadMediaLink{
		Key:        key,
		PreloadUrl: getURL,
		UploadURL:  uploadURL,
	}, nil
}

func (s *Storage) ValidateOrganizationBanner(
	ctx context.Context,
	organizationID uuid.UUID,
	key string,
) error {
	err := validateTempOrganizationBannerKey(organizationID, key)
	if err != nil {
		return err
	}

	out, err := s.s3.GetObjectRange(ctx, key, 64*1024)
	switch {
	case errors.Is(err, awsx.ErrNotFound):
		return errx.ErrorNoContentUploaded.Raise(
			fmt.Errorf("organization organization banner not found for key: %s", key),
		)
	case err != nil:
		return fmt.Errorf("get object range for organization organization banner: %w", err)
	}
	defer out.Body.Close()

	if err = s.config.OrgBanner.Validate(out); err != nil {
		switch {
		case errors.Is(err, awsx.ErrorNoContentUploaded):
			return errx.ErrorNoContentUploaded.Raise(err)
		case errors.Is(err, awsx.ErrorSizeExceedsMax):
			return errx.ErrorOrganizationBannerContentIsExceedsMax.Raise(err)
		case errors.Is(err, awsx.ErrorResolutionIsInvalid):
			return errx.ErrorOrganizationBannerResolutionIsInvalid.Raise(err)
		case errors.Is(err, awsx.ErrorFormatNotAllowed):
			return errx.ErrorOrganizationBannerFormatIsNotAllowed.Raise(err)
		default:
			return fmt.Errorf("validate organization banner content: %w", err)
		}
	}

	return nil
}

func (s *Storage) DeleteUploadOrganizationBanner(
	ctx context.Context,
	organizationID uuid.UUID,
	key string,
) error {
	if err := validateTempOrganizationBannerKey(organizationID, key); err != nil {
		return err
	}

	if err := s.s3.DeleteObject(ctx, key); err != nil {
		return fmt.Errorf("deleting temp organization organization banner object: %w", err)
	}

	return nil
}

func (s *Storage) DeleteOrganizationBanner(
	ctx context.Context,
	organizationID uuid.UUID,
	key string,
) error {
	if err := validateFinalOrganizationBannerKey(organizationID, key); err != nil {
		return err
	}

	if err := s.s3.DeleteObject(ctx, key); err != nil {
		return fmt.Errorf("deleting organization organization banner object: %w", err)
	}

	return nil
}

func (s *Storage) UpdateOrganizationBanner(
	ctx context.Context,
	organizationID uuid.UUID,
	key string,
) (string, error) {
	if err := validateTempOrganizationBannerKey(organizationID, key); err != nil {
		return "", err
	}

	finalKey := CreateOrganizationBannerKey(organizationID)

	if err := s.s3.CopyObject(ctx, key, finalKey); err != nil {
		return "", fmt.Errorf("copying object for organization banner: %w", err)
	}

	return finalKey, nil
}

var (
	tempOrganizationBannerKeyRe = regexp.MustCompile(
		`^organization/banner/([0-9a-fA-F-]{36})/temp/([0-9a-fA-F-]{36})$`,
	)

	finalOrganizationBannerKeyRe = regexp.MustCompile(
		`^organization/banner/([0-9a-fA-F-]{36})/([0-9a-fA-F-]{36})$`,
	)
)

func validateTempOrganizationBannerKey(organizationID uuid.UUID, key string) error {
	if key == "" {
		return errx.ErrorOrganizationBannerKeyIsInvalid.Raise(fmt.Errorf("empty key"))
	}

	matches := tempOrganizationBannerKeyRe.FindStringSubmatch(key)
	if matches == nil {
		return errx.ErrorOrganizationBannerKeyIsInvalid.Raise(fmt.Errorf("key %s does not match temp organization organization banner key pattern", key))
	}

	if matches[1] != organizationID.String() {
		return errx.ErrorOrganizationBannerKeyIsInvalid.Raise(fmt.Errorf("key %s does not belong to organization organization %s", key, organizationID))
	}

	return nil
}

func validateFinalOrganizationBannerKey(organizationID uuid.UUID, key string) error {
	if key == "" {
		return errx.ErrorOrganizationBannerKeyIsInvalid.Raise(fmt.Errorf("empty key"))
	}

	matches := finalOrganizationBannerKeyRe.FindStringSubmatch(key)
	if matches == nil {
		return errx.ErrorOrganizationBannerKeyIsInvalid.Raise(fmt.Errorf("key %s does not match final organization organization banner key pattern", key))
	}

	if matches[1] != organizationID.String() {
		return errx.ErrorOrganizationBannerKeyIsInvalid.Raise(fmt.Errorf("key %s does not belong to organization organization %s", key, organizationID))
	}

	return nil
}
