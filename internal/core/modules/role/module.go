package role

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/restkit/pagi"
)

type Module struct {
	repo      repo
	messenger messenger
}

func New(repo repo, messenger messenger) *Module {
	return &Module{
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

	UpdateRole(
		ctx context.Context,
		roleID uuid.UUID,
		params UpdateParams,
	) (models.Role, error)
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

	GetRolePermissions(
		ctx context.Context,
		roleID uuid.UUID,
	) (models.OrgRolePermissionLinks, error)
	SetRolePermissions(
		ctx context.Context,
		roleID uuid.UUID,
		perms models.OrgRolePermissionDict,
	) (models.OrgRolePermissionLinks, error)

	GetAllPermissions(ctx context.Context) ([]models.OrgRolePermission, error)

	GetMemberMaxRole(ctx context.Context, memberID uuid.UUID) (models.Role, error)

	GetMemberRoles(ctx context.Context, memberID uuid.UUID) ([]models.Role, error)
	RemoveMemberRole(
		ctx context.Context,
		memberID, roleID uuid.UUID,
	) error
	AddMemberRole(
		ctx context.Context,
		memberID, roleID uuid.UUID,
	) (models.OrgMemberRolesLink, error)

	CheckMemberHavePermission(
		ctx context.Context,
		memberID uuid.UUID,
		permissionCode string,
	) (bool, error)

	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type messenger interface {
	WriteOrgRoleCreated(ctx context.Context, role models.Role) error
	WriteOrgRoleUpdated(ctx context.Context, role models.Role) error
	WriteOrgRoleDeleted(ctx context.Context, role models.Role) error

	WriteOrgRolesRanksUpdated(
		ctx context.Context,
		organizationID uuid.UUID,
		order map[uuid.UUID]uint,
	) error
	WriteOrgRolePermissionsUpdated(
		ctx context.Context,
		role models.Role,
		permissions models.OrgRolePermissionLinks,
	) error

	WriteOrgMemberRoleAdd(
		ctx context.Context,
		link models.OrgMemberRolesLink,
	) error
	WriteOrgMemberRoleRemove(
		ctx context.Context,
		memberID uuid.UUID,
		roleID uuid.UUID,
	) error
}

func (m *Module) checkPermissionsToManageRole(
	ctx context.Context,
	initiator models.InitiatorData,
	organizationID uuid.UUID,
	rank uint,
) error {
	member, err := m.getInitiator(ctx, initiator, organizationID)
	if err != nil {
		return err
	}

	hasPermission, err := m.repo.CheckMemberHavePermission(
		ctx,
		member.ID,
		models.RolePermissionManageRoles,
	)
	if err != nil {
		return err
	}
	if !hasPermission {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("member %s does not have permission %s", member.ID, models.RolePermissionManageRoles),
		)
	}

	maxRole, err := m.repo.GetMemberMaxRole(ctx, member.ID)
	if err != nil {
		return err
	}
	if maxRole.IsNil() {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("member %s has no roles assigned", member.ID),
		)
	}

	if maxRole.Rank < rank {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("member %s with max role rank %d cannot manage role with rank %d",
				member.ID, maxRole.Rank, rank,
			),
		)
	}

	return nil
}

func (m *Module) getInitiator(
	ctx context.Context,
	initiator models.InitiatorData,
	organizationID uuid.UUID,
) (models.Member, error) {
	member, err := m.repo.GetMemberByAccountAndOrganization(ctx, initiator.AccountID, organizationID)
	if err != nil {
		return models.Member{}, err
	}
	if member.IsNil() {
		return models.Member{}, errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf(
				"initiator member with account id %s and organization id %s not found: %w",
				initiator.AccountID, organizationID, err,
			),
		)
	}

	return member, nil
}

func (m *Module) getMember(
	ctx context.Context,
	memberID uuid.UUID,
) (models.Member, error) {
	initiator, err := m.repo.GetMember(ctx, memberID)
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
