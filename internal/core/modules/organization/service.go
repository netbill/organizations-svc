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
	messenger messanger
}

func New(repo repo, messenger messanger) Service {
	return Service{
		repo:      repo,
		messenger: messenger,
	}
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
	WriteOrganizationSuspended(ctx context.Context, organization models.Organization) error

	WriteOrganizationUpdated(ctx context.Context, organization models.Organization) error

	WriteOrganizationDeleted(ctx context.Context, organization models.Organization) error

	WriteRoleCreated(ctx context.Context, role models.Role) error
	WriteRolePermissionsUpdated(
		ctx context.Context,
		RoleID uuid.UUID,
		permissions map[models.Permission]bool,
	) error
	WriteMemberCreated(
		ctx context.Context,
		member models.Member,
	) error
}

func (s Service) chekPermissionForManageOrganization(
	ctx context.Context,
	memberID uuid.UUID,
) error {
	access, err := s.repo.CheckMemberHavePermission(
		ctx,
		memberID,
		models.RolePermissionManageOrganization,
	)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to check initiator permissions: %w", err))
	}
	if !access {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("initiator has no access to activate organization"),
		)
	}

	return nil
}

func (s Service) getInitiator(ctx context.Context, accountID, organizationID uuid.UUID) (models.Member, error) {
	row, err := s.repo.GetMemberByAccountAndOrganization(ctx, accountID, organizationID)
	if err != nil {
		return models.Member{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get member with account id %s and organization id %s: %w",
				accountID, organizationID, err),
		)
	}
	if row.IsNil() {
		return models.Member{}, errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("member with account id %s and organization id %s not found", accountID, organizationID),
		)
	}

	return row, nil
}
