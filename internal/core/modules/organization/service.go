package organization

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/pagi"
)

type Service struct {
	repo      repo
	bucket    bucket
	token     token
	messenger messanger
}

func New(repo repo, messenger messanger) Service {
	return Service{
		repo:      repo,
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
	UpdateOrganizationStatus(ctx context.Context, ID uuid.UUID, status string) (models.Organization, error)
	UpdateOrganizationMaxRoles(ctx context.Context, ID uuid.UUID, maxRoles uint) (models.Organization, error)

	DeleteOrganization(ctx context.Context, ID uuid.UUID) error

	CheckMemberHavePermission(
		ctx context.Context,
		memberID uuid.UUID,
		permissionCode string,
	) (bool, error)
	GetMemberByAccountAndOrganization(
		ctx context.Context,
		accountID, organizationID uuid.UUID,
	) (models.Member, error)

	GetMemberMaxRole(ctx context.Context, memberID uuid.UUID) (models.Role, error)

	CreateMember(ctx context.Context, accountID, organizationID uuid.UUID) (models.Member, error)
	CreateHeadRole(ctx context.Context, organizationID uuid.UUID) (models.Role, error)
	AddMemberRole(ctx context.Context, memberID, roleID uuid.UUID) error

	GetRolePermissions(ctx context.Context, roleID uuid.UUID) (map[models.Permission]bool, error)

	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type messanger interface {
	WriteOrganizationCreated(ctx context.Context, organization models.Organization) error

	WriteOrganizationActivated(ctx context.Context, organization models.Organization) error
	WriteOrganizationDeactivated(ctx context.Context, organization models.Organization) error

	WriteOrganizationUpdated(ctx context.Context, organization models.Organization) error

	WriteOrganizationDeleted(ctx context.Context, organization models.Organization) error

	WriteOrgRoleCreated(ctx context.Context, role models.Role) error
	WriteOrgRolePermissionsUpdated(
		ctx context.Context,
		role models.Role,
		permissions map[models.Permission]bool,
	) error
	WriteOrgMemberCreated(
		ctx context.Context,
		member models.Member,
	) error
}

type bucket interface {
	GeneratePreloadLinkForUpdateOrganizationIcon(
		ctx context.Context,
		organizationID, sessionID uuid.UUID,
	) (uploadLink, getLink string, err error)

	GeneratePreloadLinkForUpdateOrganizationBanner(
		ctx context.Context,
		organizationID, sessionID uuid.UUID,
	) (uploadLink, getLink string, err error)

	GetContentLengthForOrganizationIcon(
		ctx context.Context,
		organizationID, sessionID uuid.UUID,
	) (int64, error)

	GetContentLengthForOrganizationBanner(
		ctx context.Context,
		organizationID, sessionID uuid.UUID,
	) (int64, error)

	ValidateOrganizationIconSize(
		ctx context.Context,
		organizationID, sessionID uuid.UUID,
	) (bool, error)

	ValidateOrganizationBannerSize(
		ctx context.Context,
		organizationID, sessionID uuid.UUID,
	) (bool, error)

	GetLoadedOrganizationIcon(
		ctx context.Context,
		organizationID, sessionID uuid.UUID,
	) (data []byte, err error)

	GetLoadedOrganizationBanner(
		ctx context.Context,
		organizationID, sessionID uuid.UUID,
	) (data []byte, err error)

	ValidateOrganizationIconFormat(
		ctx context.Context,
		imageData []byte,
	) (bool, error)

	ValidateOrganizationBannerFormat(
		ctx context.Context,
		imageData []byte,
	) (bool, error)

	ValidateOrganizationIconContentType(
		ctx context.Context,
		organizationID, sessionID uuid.UUID,
	) (bool, error)

	ValidateOrganizationBannerContentType(
		ctx context.Context,
		organizationID, sessionID uuid.UUID,
	) (bool, error)

	AcceptUpdateOrganizationMedia(
		ctx context.Context,
		organizationID, sessionID uuid.UUID,
	) (icon *string, banner *string, err error)
}

type token interface {
	NewUploadOrganizationMediaToken(
		accountID uuid.UUID,
		organizationID uuid.UUID,
		uploadSessionID uuid.UUID,
	) (string, error)
}

func (s Service) chekPermissionForManageOrganization(
	ctx context.Context,
	accountID, organizationID uuid.UUID,
) error {
	member, err := s.repo.GetMemberByAccountAndOrganization(ctx, accountID, organizationID)
	if err != nil {
		return fmt.Errorf(
			"failed to get member with account id %s and organization id %s: %w",
			accountID, organizationID, err,
		)
	}
	if member.IsNil() {
		return fmt.Errorf(
			"member with account id %s and organization id %s not found", accountID, organizationID,
		)
	}

	access, err := s.repo.CheckMemberHavePermission(
		ctx,
		member.ID,
		models.RolePermissionManageOrganization,
	)
	if err != nil {
		return err
	}
	if !access {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("initiator has no access to activate organization"),
		)
	}

	return nil
}
