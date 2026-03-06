package media

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

func (u *Uploader) CreateOrganizationIconUploadMediaLinks(
	ctx context.Context,
	organizationID uuid.UUID,
) (models.UploadMediaLink, error) {
	key := createTempOrganizationIconKey(organizationID)

	uploadURL, getURL, err := u.s3.PresignPut(ctx, key, u.config.LinkTTL)
	if err != nil {
		return models.UploadMediaLink{}, fmt.Errorf("presigning put for organization organization icon: %w", err)
	}

	return models.UploadMediaLink{
		Key:        key,
		PreloadUrl: getURL,
		UploadURL:  uploadURL,
	}, nil
}

func (u *Uploader) UpdateOrganizationIcon(
	ctx context.Context,
	orgID uuid.UUID,
	tempKey string,
) (newKey string, err error) {
	if err = validateTempOrganizationIconKey(orgID, tempKey); err != nil {
		return "", err
	}

	out, err := u.s3.GetObjectRange(ctx, tempKey, 64*1024)
	switch {
	case errors.Is(err, awsx.ErrNotFound):
		return "", errx.ErrorOrganizationUploadedIconInvalid.Raise(
			fmt.Errorf("organization organization icon not found for key: %s", tempKey),
		)
	case err != nil:
		return "", fmt.Errorf("get object range for organization organization icon: %w", err)
	}
	defer out.Body.Close()

	if err = u.config.OrgIcon.Validate(out); err != nil {
		return "", errx.ErrorOrganizationUploadedIconInvalid.Raise(
			fmt.Errorf("validating organization organization icon: %w", err),
		)
	}

	finalKey := createOrganizationIconKey(orgID)

	if err = u.s3.CopyObject(ctx, tempKey, finalKey); err != nil {
		return "", fmt.Errorf("copying object for organization icon: %w", err)
	}

	return finalKey, nil
}

func (u *Uploader) DeleteUploadOrganizationIcon(
	ctx context.Context,
	organizationID uuid.UUID,
	key string,
) error {
	if err := validateTempOrganizationIconKey(organizationID, key); err != nil {
		return err
	}

	if err := u.s3.DeleteObject(ctx, key); err != nil {
		return fmt.Errorf("deleting temp organization organization icon object: %w", err)
	}

	return nil
}

func (u *Uploader) DeleteOrganizationIcon(
	ctx context.Context,
	organizationID uuid.UUID,
	key string,
) error {
	if err := validateFinalOrganizationIconKey(organizationID, key); err != nil {
		return err
	}

	if err := u.s3.DeleteObject(ctx, key); err != nil {
		return fmt.Errorf("deleting organization organization icon object: %w", err)
	}

	return nil
}

var tempOrganizationIconKeyRe = regexp.MustCompile(
	`^organization/icon/([0-9a-fA-F-]{36})/temp/([0-9a-fA-F-]{36})$`,
)

func createTempOrganizationIconKey(organizationID uuid.UUID) string {
	return fmt.Sprintf("organization/icon/%s/temp/%s", organizationID, uuid.New().String())
}

func validateTempOrganizationIconKey(organizationID uuid.UUID, key string) error {
	matches := tempOrganizationIconKeyRe.FindStringSubmatch(key)
	if matches == nil {
		return errx.ErrorOrganizationUploadedIconInvalid.Raise(
			fmt.Errorf("key %s does not match temp organization organization icon key pattern", key),
		)
	}

	if matches[1] != organizationID.String() {
		return errx.ErrorOrganizationUploadedIconInvalid.Raise(
			fmt.Errorf("key %s does not belong to organization organization %s", key, organizationID),
		)
	}

	return nil
}

var finalOrganizationIconKeyRe = regexp.MustCompile(
	`^organization/icon/([0-9a-fA-F-]{36})/([0-9a-fA-F-]{36})$`,
)

func createOrganizationIconKey(organizationID uuid.UUID) string {
	return fmt.Sprintf("organization/icon/%s/%s", organizationID, uuid.New().String())
}

func validateFinalOrganizationIconKey(organizationID uuid.UUID, key string) error {
	matches := finalOrganizationIconKeyRe.FindStringSubmatch(key)
	if matches == nil {
		return errx.ErrorOrganizationUploadedIconInvalid.Raise(
			fmt.Errorf("key %s does not match final organization organization icon key pattern", key),
		)
	}

	if matches[1] != organizationID.String() {
		return errx.ErrorOrganizationUploadedIconInvalid.Raise(
			fmt.Errorf("key %s does not belong to organization organization %s", key, organizationID),
		)
	}

	return nil
}
