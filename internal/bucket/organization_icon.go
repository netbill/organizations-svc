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
	return fmt.Sprintf("organization/icon/%s/temp/%s", organizationID, uuid.New().String())
}

func CreateOrganizationIconKey(organizationID uuid.UUID) string {
	return fmt.Sprintf("organization/icon/%s/%s", organizationID, uuid.New().String())
}

func (s *Storage) CreateOrganizationIconUploadMediaLinks(
	ctx context.Context,
	organizationID uuid.UUID,
) (models.UploadMediaLink, error) {
	key := CreateTempOrganizationIconKey(organizationID)

	uploadURL, getURL, err := s.s3.PresignPut(
		ctx,
		key,
		s.config.LinkTTL,
	)
	if err != nil {
		return models.UploadMediaLink{}, fmt.Errorf("presigning put for organization organization icon: %w", err)
	}

	return models.UploadMediaLink{
		Key:        key,
		PreloadUrl: getURL,
		UploadURL:  uploadURL,
	}, nil
}

func (s *Storage) ValidateOrganizationIcon(
	ctx context.Context,
	organizationID uuid.UUID,
	key string,
) error {
	err := validateTempOrganizationIconKey(organizationID, key)
	if err != nil {
		return err
	}

	out, err := s.s3.GetObjectRange(ctx, key, 64*1024)
	switch {
	case errors.Is(err, awsx.ErrNotFound):
		return errx.ErrorNoContentUploaded.Raise(
			fmt.Errorf("organization organization icon not found for key: %s", key),
		)
	case err != nil:
		return fmt.Errorf("get object range for organization organization icon: %w", err)
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
	key string,
) error {
	if err := validateTempOrganizationIconKey(organizationID, key); err != nil {
		return err
	}

	if err := s.s3.DeleteObject(ctx, key); err != nil {
		return fmt.Errorf("deleting temp organization organization icon object: %w", err)
	}

	return nil
}

func (s *Storage) DeleteOrganizationIcon(
	ctx context.Context,
	organizationID uuid.UUID,
	key string,
) error {
	if err := validateFinalOrganizationIconKey(organizationID, key); err != nil {
		return err
	}

	if err := s.s3.DeleteObject(ctx, key); err != nil {
		return fmt.Errorf("deleting organization organization icon object: %w", err)
	}

	return nil
}

func (s *Storage) UpdateOrganizationIcon(
	ctx context.Context,
	organizationID uuid.UUID,
	key string,
) (string, error) {
	if err := validateTempOrganizationIconKey(organizationID, key); err != nil {
		return "", err
	}

	finalKey := CreateOrganizationIconKey(organizationID)

	if err := s.s3.CopyObject(ctx, key, finalKey); err != nil {
		return "", fmt.Errorf("copying object for organization icon: %w", err)
	}

	return finalKey, nil
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
		return errx.ErrorOrganizationIconKeyIsInvalid.Raise(fmt.Errorf("key %s does not match temp organization organization icon key pattern", key))
	}

	if matches[1] != organizationID.String() {
		return errx.ErrorOrganizationIconKeyIsInvalid.Raise(fmt.Errorf("key %s does not belong to organization organization %s", key, organizationID))
	}

	return nil
}

func validateFinalOrganizationIconKey(organizationID uuid.UUID, key string) error {
	if key == "" {
		return errx.ErrorOrganizationIconKeyIsInvalid.Raise(fmt.Errorf("empty key"))
	}

	matches := finalOrganizationIconKeyRe.FindStringSubmatch(key)
	if matches == nil {
		return errx.ErrorOrganizationIconKeyIsInvalid.Raise(fmt.Errorf("key %s does not match final organization organization icon key pattern", key))
	}

	if matches[1] != organizationID.String() {
		return errx.ErrorOrganizationIconKeyIsInvalid.Raise(fmt.Errorf("key %s does not belong to organization organization %s", key, organizationID))
	}

	return nil
}
