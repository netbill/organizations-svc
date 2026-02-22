package organization

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/domain"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/orgperm"
	"github.com/netbill/restkit/pagi"
)

type Module struct {
	repo      repo
	bucket    bucket
	messenger messanger
}

func New(repo repo, messenger messanger, bucket bucket) *Module {
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
	CreateOrganization(ctx context.Context, params CreateParams) (domain.Organization, error)

	GetOrganizationByID(ctx context.Context, ID uuid.UUID) (domain.Organization, error)
	GetOrganizations(
		ctx context.Context,
		filter FilterParams,
		limit, offset uint,
	) (pagi.Page[[]domain.Organization], error)
	GetOrganizationsForUser(
		ctx context.Context,
		accountID uuid.UUID,
		limit, offset uint,
	) (pagi.Page[[]domain.Organization], error)

	UpdateOrganization(
		ctx context.Context,
		ID uuid.UUID,
		params UpdateParams,
	) (domain.Organization, error)
	UpdateOrganizationStatus(ctx context.Context, ID uuid.UUID, status string) (domain.Organization, error)
	UpdateOrganizationMaxRoles(ctx context.Context, ID uuid.UUID, maxRoles uint) (domain.Organization, error)

	DeleteOrganization(ctx context.Context, ID uuid.UUID) error

	CheckMemberHavePermission(
		ctx context.Context,
		memberID uuid.UUID,
		permissionID uuid.UUID,
	) (bool, error)
	GetMemberByAccountAndOrganization(
		ctx context.Context,
		accountID, organizationID uuid.UUID,
	) (domain.Member, error)

	GetMemberMaxRole(ctx context.Context, memberID uuid.UUID) (domain.Role, error)

	CreateMember(ctx context.Context, accountID, organizationID uuid.UUID) (domain.Member, error)
	CreateMemberHead(ctx context.Context, accountID, organizationID uuid.UUID) (domain.Member, error)

	GetRolePermissions(
		ctx context.Context,
		roleID uuid.UUID,
	) (domain.OrgRolePermissionsWithDetailsForRole, error)

	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type messanger interface {
	WriteOrganizationCreated(ctx context.Context, organization domain.Organization) error

	WriteOrganizationActivated(ctx context.Context, organization domain.Organization) error
	WriteOrganizationDeactivated(ctx context.Context, organization domain.Organization) error

	WriteOrganizationUpdated(ctx context.Context, organization domain.Organization) error

	WriteOrganizationDeleted(ctx context.Context, organization domain.Organization) error

	WriteOrgMemberCreated(ctx context.Context, member domain.Member) error
}

type bucket interface {
	CreateOrganizationIconUploadMediaLinks(
		ctx context.Context,
		organizationID uuid.UUID,
	) (domain.UploadMediaLink, error)

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
	) (domain.UploadMediaLink, error)

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

func (m *Module) chekPermissionForManageOrganization(
	ctx context.Context,
	initiator domain.AccountActor,
	organizationID uuid.UUID,
) error {
	member, err := m.repo.GetMemberByAccountAndOrganization(ctx, initiator, organizationID)
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

	access, err := m.repo.CheckMemberHavePermission(ctx, member.ID, orgperm.OrganizationUpdateID)
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
