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

func CreateTempOrganizationIconKey(organizationID uuid.UUID) string {
	return fmt.Sprintf("organization/icon/%s/temp/%s", organizationID, uuid.New())
}

func CreateFinalOrganizationIconKey(organizationID uuid.UUID) string {
	return fmt.Sprintf("organization/icon/%s/%s", organizationID, uuid.New())
}

func (s *Storage) CreateOrganizationIconUploadMediaLinks(
	ctx context.Context,
	organizationID uuid.UUID,
) (models.UploadMediaLink, error) {
	key := CreateTempOrganizationIconKey(organizationID)

	uploadLink, getLink, err := s.s3.PresignPut(ctx, key, s.config.LinkTTL)
	if err != nil {
		return models.UploadMediaLink{}, fmt.Errorf("presign put object for organization icon: %w", err)
	}

	return models.UploadMediaLink{
		Key:        key,
		PreloadUrl: getLink,
		UploadURL:  uploadLink,
	}, nil
}

func (s *Storage) ValidateOrganizationIcon(
	ctx context.Context,
	organizationID uuid.UUID,
	tempKey string,
) error {
	if err := validateTempOrganizationIconKey(organizationID, tempKey); err != nil {
		return err
	}

	out, err := s.s3.GetObjectRange(ctx, tempKey, 64*1024)
	switch {
	case errors.Is(err, awsx.ErrNotFound):
		return errx.ErrorNoContentUploaded.Raise(
			fmt.Errorf("organization icon not found for key: %s", tempKey),
		)
	case err != nil:
		return fmt.Errorf("get object range for organization icon: %w", err)
	}
	defer out.Body.Close()

	if err = s.config.OrgIcon.Validate(out); err != nil {
		switch {
		case errors.Is(err, awsx.ErrorNoContentUploaded):
			return errx.ErrorNoContentUploaded.Raise(err)
		case errors.Is(err, awsx.ErrorSizeExceedsMax):
			return errx.ErrorOrganizationIconContentIsExceedsMax.Raise(err)
		case errors.Is(err, awsx.ErrorResolutionIsInvalid):
			return errx.ErrorOrganizationIconResolutionIsInvalid.Raise(err)
		case errors.Is(err, awsx.ErrorFormatNotAllowed):
			return errx.ErrorOrganizationIconFormatIsNotAllowed.Raise(err)
		default:
			return fmt.Errorf("validate organization icon content: %w", err)
		}
	}

	return nil
}

func (s *Storage) DeleteUploadOrganizationIcon(
	ctx context.Context,
	organizationID uuid.UUID,
	tempKey string,
) error {
	if err := validateTempOrganizationIconKey(organizationID, tempKey); err != nil {
		return err
	}

	if err := s.s3.DeleteObject(ctx, tempKey); err != nil {
		return fmt.Errorf("delete temp organization icon: %w", err)
	}

	return nil
}

func (s *Storage) DeleteOrganizationIcon(
	ctx context.Context,
	organizationID uuid.UUID,
	finalKey string,
) error {
	if err := validateFinalOrganizationIconKey(organizationID, finalKey); err != nil {
		return err
	}

	if err := s.s3.DeleteObject(ctx, finalKey); err != nil {
		return fmt.Errorf("delete organization icon: %w", err)
	}

	return nil
}

func (s *Storage) UpdateOrganizationIcon(
	ctx context.Context,
	organizationID uuid.UUID,
	oldFinalKey *string,
	tempKey *string,
) (*string, error) {
	if ptrStrEq(oldFinalKey, tempKey) {
		return oldFinalKey, nil
	}

	if tempKey == nil {
		return nil, s.DeleteOrganizationIcon(ctx, organizationID, *oldFinalKey)
	}

	//if err := s.ValidateOrganizationIcon(ctx, organizationID, *tempKey); err != nil {
	//	return nil, err
	//}

	finalKey := CreateFinalOrganizationIconKey(organizationID)

	if err := s.s3.CopyObject(ctx, *tempKey, finalKey); err != nil {
		return nil, fmt.Errorf("copy object for organization icon: %w", err)
	}

	if err := s.s3.DeleteObject(ctx, *tempKey); err != nil {
		return nil, fmt.Errorf("delete temp organization icon: %w", err)
	}

	if oldFinalKey != nil {
		if err := s.DeleteOrganizationIcon(ctx, organizationID, *oldFinalKey); err != nil {
			return nil, err
		}
	}

	return &finalKey, nil
}

var (
	tempOrganizationIconKeyRe = regexp.MustCompile(
		`^organization/icon/([0-9a-fA-F-]{36})/temp/([0-9a-fA-F-]{36})$`,
	)
	finalOrganizationIconKeyRe = regexp.MustCompile(
		`^organization/icon/([0-9a-fA-F-]{36})/([0-9a-fA-F-]{36})$`,
	)
)

func validateTempOrganizationIconKey(organizationID uuid.UUID, key string) error {
	if key == "" {
		return errx.ErrorOrganizationIconKeyIsInvalid.Raise(fmt.Errorf("empty key"))
	}

	matches := tempOrganizationIconKeyRe.FindStringSubmatch(key)
	if matches == nil {
		return errx.ErrorOrganizationIconKeyIsInvalid.Raise(fmt.Errorf("invalid key format"))
	}

	if matches[1] != organizationID.String() {
		return errx.ErrorOrganizationIconKeyIsInvalid.Raise(fmt.Errorf("key does not belong to the organization"))
	}

	return nil
}

func validateFinalOrganizationIconKey(organizationID uuid.UUID, key string) error {
	if key == "" {
		return errx.ErrorOrganizationIconKeyIsInvalid.Raise(fmt.Errorf("empty key"))
	}

	matches := finalOrganizationIconKeyRe.FindStringSubmatch(key)
	if matches == nil {
		return errx.ErrorOrganizationIconKeyIsInvalid.Raise(fmt.Errorf("invalid key format"))
	}

	if matches[1] != organizationID.String() {
		return errx.ErrorOrganizationIconKeyIsInvalid.Raise(fmt.Errorf("key does not belong to the organization"))
	}

	if tempOrganizationIconKeyRe.MatchString(key) {
		return errx.ErrorOrganizationIconKeyIsInvalid.Raise(fmt.Errorf("final key cannot be temp key"))
	}

	return nil
}
