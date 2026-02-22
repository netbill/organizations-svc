package bucket

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/google/uuid"
	"github.com/netbill/awsx"
	"github.com/netbill/organizations-svc/internal/core/domain"
	"github.com/netbill/organizations-svc/internal/core/errx"
)

func CreateTempOrganizationBannerKey(organizationID uuid.UUID) string {
	return fmt.Sprintf("organization/banner/%s/temp/%s", organizationID, uuid.New())
}

func CreateFinalOrganizationBannerKey(organizationID uuid.UUID) string {
	return fmt.Sprintf("organization/banner/%s/%s", organizationID, uuid.New())
}

func (s Storage) CreateOrganizationBannerUploadMediaLinks(
	ctx context.Context,
	organizationID uuid.UUID,
) (domain.UploadMediaLink, error) {
	key := CreateTempOrganizationBannerKey(organizationID)

	uploadLink, getLink, err := s.s3.PresignPut(ctx, key, s.config.LinkTTL)
	if err != nil {
		return domain.UploadMediaLink{}, fmt.Errorf("presign put object for organization banner: %w", err)
	}

	return domain.UploadMediaLink{
		Key:        key,
		PreloadUrl: getLink,
		UploadURL:  uploadLink,
	}, nil
}

func (s Storage) ValidateOrganizationBanner(
	ctx context.Context,
	organizationID uuid.UUID,
	tempKey string,
) error {
	if err := validateTempOrganizationBannerKey(organizationID, tempKey); err != nil {
		return err
	}

	out, err := s.s3.GetObjectRange(ctx, tempKey, 64*1024)
	switch {
	case errors.Is(err, awsx.ErrNotFound):
		return errx.ErrorNoContentUploaded.Raise(
			fmt.Errorf("organization banner not found for key: %s", tempKey),
		)
	case err != nil:
		return fmt.Errorf("get object range for organization banner: %w", err)
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

func (s Storage) DeleteUploadOrganizationBanner(
	ctx context.Context,
	organizationID uuid.UUID,
	tempKey string,
) error {
	if err := validateTempOrganizationBannerKey(organizationID, tempKey); err != nil {
		return err
	}

	if err := s.s3.DeleteObject(ctx, tempKey); err != nil {
		return fmt.Errorf("delete temp organization banner: %w", err)
	}

	return nil
}

func (s Storage) DeleteOrganizationBanner(
	ctx context.Context,
	organizationID uuid.UUID,
	finalKey string,
) error {
	if err := validateFinalOrganizationBannerKey(organizationID, finalKey); err != nil {
		return err
	}

	if err := s.s3.DeleteObject(ctx, finalKey); err != nil {
		return fmt.Errorf("delete organization banner: %w", err)
	}

	return nil
}

func (s Storage) UpdateOrganizationBanner(
	ctx context.Context,
	organizationID uuid.UUID,
	oldFinalKey *string,
	tempKey *string,
) (*string, error) {
	if ptrStrEq(oldFinalKey, tempKey) {
		return oldFinalKey, nil
	}

	if tempKey == nil {
		return nil, s.DeleteOrganizationBanner(ctx, organizationID, *oldFinalKey)
	}

	//if err := s.ValidateOrganizationBanner(ctx, organizationID, *tempKey); err != nil {
	//	return nil, err
	//}

	finalKey := CreateFinalOrganizationBannerKey(organizationID)

	if err := s.s3.CopyObject(ctx, *tempKey, finalKey); err != nil {
		return nil, fmt.Errorf("copy object for organization banner: %w", err)
	}

	if err := s.s3.DeleteObject(ctx, *tempKey); err != nil {
		return nil, fmt.Errorf("delete temp organization banner: %w", err)
	}

	if oldFinalKey != nil {
		if err := s.DeleteOrganizationBanner(ctx, organizationID, *oldFinalKey); err != nil {
			return nil, err
		}
	}

	return &finalKey, nil
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
		return errx.ErrorOrganizationBannerKeyIsInvalid.Raise(fmt.Errorf("invalid key format"))
	}

	if matches[1] != organizationID.String() {
		return errx.ErrorOrganizationBannerKeyIsInvalid.Raise(fmt.Errorf("key does not belong to the organization"))
	}

	return nil
}

func validateFinalOrganizationBannerKey(organizationID uuid.UUID, key string) error {
	if key == "" {
		return errx.ErrorOrganizationBannerKeyIsInvalid.Raise(fmt.Errorf("empty key"))
	}

	matches := finalOrganizationBannerKeyRe.FindStringSubmatch(key)
	if matches == nil {
		return errx.ErrorOrganizationBannerKeyIsInvalid.Raise(fmt.Errorf("invalid key format"))
	}

	if matches[1] != organizationID.String() {
		return errx.ErrorOrganizationBannerKeyIsInvalid.Raise(fmt.Errorf("key does not belong to the organization"))
	}

	if tempOrganizationBannerKeyRe.MatchString(key) {
		return errx.ErrorOrganizationBannerKeyIsInvalid.Raise(fmt.Errorf("final key cannot be temp key"))
	}

	return nil
}
