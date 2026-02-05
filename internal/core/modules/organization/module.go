package organization

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/restkit/pagi"
)

type Module struct {
	repo      repo
	bucket    bucket
	token     token
	messenger messanger
}

func New(repo repo, messenger messanger, token token, bucket bucket) *Module {
	return &Module{
		repo:      repo,
		messenger: messenger,
		token:     token,
		bucket:    bucket,
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
	CreateMemberHead(ctx context.Context, accountID, organizationID uuid.UUID) (models.Member, error)

	GetRolePermissions(
		ctx context.Context,
		roleID uuid.UUID,
	) (models.OrgRolePermissionDictWithDetails, error)

	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type messanger interface {
	WriteOrganizationCreated(ctx context.Context, organization models.Organization) error

	WriteOrganizationActivated(ctx context.Context, organization models.Organization) error
	WriteOrganizationDeactivated(ctx context.Context, organization models.Organization) error

	WriteOrganizationUpdated(ctx context.Context, organization models.Organization) error

	WriteOrganizationDeleted(ctx context.Context, organization models.Organization) error

	WriteOrgMemberCreated(
		ctx context.Context,
		member models.Member,
	) error
}

type bucket interface {
	GeneratePreloadLinkForOrganizationMedia(
		ctx context.Context,
		organizationID, sessionID uuid.UUID,
	) (models.OrganizationUploadMediaLinks, error)

	AcceptUpdateOrganizationMedia(
		ctx context.Context,
		organizationID, sessionID uuid.UUID,
	) (models.OrganizationMedia, error)

	CancelUpdateOrganizationIcon(
		ctx context.Context,
		organizationID, sessionID uuid.UUID,
	) error

	CancelUpdateOrganizationBanner(
		ctx context.Context,
		organizationID, sessionID uuid.UUID,
	) error

	DeleteOrganizationIcon(
		ctx context.Context,
		organizationID uuid.UUID,
	) error

	DeleteOrganizationBanner(
		ctx context.Context,
		organizationID uuid.UUID,
	) error

	CleanOrganizationMediaSession(
		ctx context.Context,
		organizationID, sessionID uuid.UUID,
	) error
}

type token interface {
	NewUploadOrganizationMediaToken(
		accountID uuid.UUID,
		organizationID uuid.UUID,
		uploadSessionID uuid.UUID,
	) (string, error)
}

func (m *Module) chekPermissionForManageOrganization(
	ctx context.Context,
	initiator models.Initiator,
	organizationID uuid.UUID,
) error {
	member, err := m.repo.GetMemberByAccountAndOrganization(ctx, initiator.GetAccountID(), organizationID)
	if err != nil {
		if errors.Is(err, errx.ErrorMemberNotFound) {
			return errx.ErrorNotEnoughRights.Raise(
				fmt.Errorf("account is not a member of the organization"),
			)
		}
		return err
	}

	if member.Head {
		return nil
	}

	access, err := m.repo.CheckMemberHavePermission(
		ctx,
		member.ID,
		models.RolePermissionOrganizationUpdate,
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
