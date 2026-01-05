package role

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
	messenger messenger
}

func New(repo repo, messenger messenger) Service {
	return Service{
		repo:      repo,
		messenger: messenger,
	}
}

type repo interface {
	CreateRole(ctx context.Context, params CreateParams) (models.Role, error)

	GetRole(ctx context.Context, roleID uuid.UUID) (models.Role, error)
	GetRoles(
		ctx context.Context,
		filter FilterParams,
		limit, offset uint,
	) (pagi.Page[[]models.Role], error)

	UpdateRole(ctx context.Context, roleID uuid.UUID, params UpdateParams) (models.Role, error)
	UpdateRolesRanks(
		ctx context.Context,
		organizationID uuid.UUID,
		order map[uuid.UUID]uint,
	) error

	DeleteRole(ctx context.Context, roleID uuid.UUID) error

	GetMember(ctx context.Context, memberID uuid.UUID) (models.Member, error)
	GetMemberByAccountAndOrganization(
		ctx context.Context,
		accountID, organizationID uuid.UUID,
	) (models.Member, error)

	GetRolePermissions(ctx context.Context, roleID uuid.UUID) (map[models.Permission]bool, error)
	SetRolePermissions(
		ctx context.Context,
		roleID uuid.UUID,
		permissions map[string]bool,
	) error

	GetAllPermissions(ctx context.Context) ([]models.Permission, error)

	GetMemberMaxRole(ctx context.Context, memberID uuid.UUID) (models.Role, error)

	GetMemberRoles(ctx context.Context, memberID uuid.UUID) ([]models.Role, error)
	RemoveMemberRole(ctx context.Context, memberID, roleID uuid.UUID) error
	AddMemberRole(ctx context.Context, memberID, roleID uuid.UUID) error

	CheckMemberHavePermission(
		ctx context.Context,
		memberID uuid.UUID,
		permissionCode string,
	) (bool, error)

	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type messenger interface {
	WriteRoleCreated(ctx context.Context, role models.Role) error
	WriteRoleUpdated(ctx context.Context, role models.Role) error
	WriteRoleDeleted(ctx context.Context, role models.Role) error

	WriteRolesRanksUpdated(
		ctx context.Context,
		organizationID uuid.UUID,
		order map[uuid.UUID]uint,
	) error
	WriteRolePermissionsUpdated(
		ctx context.Context,
		roleID uuid.UUID,
		permissions map[models.Permission]bool,
	) error

	WriteMemberRoleAdd(
		ctx context.Context,
		memberID uuid.UUID,
		roleID uuid.UUID,
	) error
	WriteMemberRoleRemove(
		ctx context.Context,
		memberID uuid.UUID,
		roleID uuid.UUID,
	) error
}

func (s Service) checkPermissionsToManageRole(
	ctx context.Context,
	memberID uuid.UUID,
	rank uint,
) error {
	hasPermission, err := s.repo.CheckMemberHavePermission(
		ctx,
		memberID,
		models.RolePermissionManageRoles,
	)
	if err != nil {
		return err
	}
	if !hasPermission {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("member %s does not have permission %s", memberID, models.RolePermissionManageRoles),
		)
	}

	maxRole, err := s.repo.GetMemberMaxRole(ctx, memberID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get max role for member %s, cause: %w",
				memberID, err),
		)
	}
	if maxRole.IsNil() {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("member %s has no roles assigned", memberID),
		)
	}

	if maxRole.Rank < rank {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("member %s with max role rank %d cannot manage role with rank %d",
				memberID, maxRole.Rank, rank),
		)
	}

	return nil
}

func (s Service) getInitiator(
	ctx context.Context,
	accountID, organizationID uuid.UUID,
) (models.Member, error) {
	initiator, err := s.repo.GetMemberByAccountAndOrganization(ctx, accountID, organizationID)
	if err != nil {
		return models.Member{}, err
	}
	if initiator.IsNil() {
		return models.Member{}, errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("initiator member with account id %s and organization id %s not found: %w",
				accountID, organizationID, err),
		)
	}

	return initiator, nil
}

func (s Service) getMember(
	ctx context.Context,
	memberID uuid.UUID,
) (models.Member, error) {
	initiator, err := s.repo.GetMember(ctx, memberID)
	if err != nil {
		return models.Member{}, err
	}
	if initiator.IsNil() {
		return models.Member{}, errx.ErrorMemberNotFound.Raise(
			fmt.Errorf("initiator member with id %s not found", memberID),
		)
	}

	return initiator, nil
}
