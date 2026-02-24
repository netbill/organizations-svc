package organization

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/restkit/pagi"
)

type Module struct {
	repo      repo
	bucket    bucket
	messenger messenger
}

func New(repo repo, messenger messenger, bucket bucket) *Module {
	return &Module{
		repo:      repo,
		bucket:    bucket,
		messenger: messenger,
	}
}

type RepoUpdateOrganizationParams struct {
	Name   *string
	Icon   *string
	Banner *string
}

type repo interface {
	CreateOrganization(ctx context.Context, params CreateParams) (models.Organization, error)

	GetOrganizationByID(ctx context.Context, ID uuid.UUID) (models.Organization, error)
	GetOrganizations(
		ctx context.Context,
		filter FilterParams,
		limit, offset uint,
	) (pagi.Page[[]models.Organization], error)
	GetOrganizationsForUser(
		ctx context.Context,
		accountID uuid.UUID,
		limit, offset uint,
	) (pagi.Page[[]models.Organization], error)

	UpdateOrganization(
		ctx context.Context,
		ID uuid.UUID,
		params UpdateParams,
	) (models.Organization, error)
	UpdateOrganizationStatus(
		ctx context.Context,
		ID uuid.UUID,
		status string,
	) (models.Organization, error)

	DeleteOrganization(ctx context.Context, ID uuid.UUID) error

	GetMemberByAccountAndOrganization(
		ctx context.Context,
		accountID, organizationID uuid.UUID,
	) (models.Member, error)

	CreateMember(ctx context.Context, accountID, organizationID uuid.UUID) (models.Member, error)
	CreateMemberHead(ctx context.Context, accountID, organizationID uuid.UUID) (models.Member, error)

	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type messenger interface {
	WriteOrganizationCreated(ctx context.Context, organization models.Organization) error
	WriteOrganizationUpdated(ctx context.Context, organization models.Organization) error
	WriteOrganizationDeleted(ctx context.Context, organization models.Organization) error

	WriteOrgMemberCreated(ctx context.Context, member models.Member) error
}

type bucket interface {
	CreateOrganizationIconUploadMediaLinks(
		ctx context.Context,
		organizationID uuid.UUID,
	) (models.UploadMediaLink, error)

	ValidateOrganizationIcon(
		ctx context.Context,
		organizationID uuid.UUID,
		tempKey string,
	) error

	DeleteUploadOrganizationIcon(
		ctx context.Context,
		organizationID uuid.UUID,
		tempKey string,
	) error

	DeleteOrganizationIcon(
		ctx context.Context,
		organizationID uuid.UUID,
		finalKey string,
	) error

	UpdateOrganizationIcon(
		ctx context.Context,
		organizationID uuid.UUID,
		oldFinalKey *string,
		tempKey *string,
	) (*string, error)

	CreateOrganizationBannerUploadMediaLinks(
		ctx context.Context,
		organizationID uuid.UUID,
	) (models.UploadMediaLink, error)

	ValidateOrganizationBanner(
		ctx context.Context,
		organizationID uuid.UUID,
		tempKey string,
	) error

	DeleteUploadOrganizationBanner(
		ctx context.Context,
		organizationID uuid.UUID,
		tempKey string,
	) error

	DeleteOrganizationBanner(
		ctx context.Context,
		organizationID uuid.UUID,
		finalKey string,
	) error

	UpdateOrganizationBanner(
		ctx context.Context,
		organizationID uuid.UUID,
		oldFinalKey *string,
		tempKey *string,
	) (*string, error)
}
