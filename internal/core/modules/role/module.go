package role

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/orgperm"
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
	//Role

	CreateRole(ctx context.Context, params CreateParams) (models.Role, error)
	GetRole(ctx context.Context, roleID uuid.UUID) (models.Role, error)
	UpdateRole(
		ctx context.Context,
		roleID uuid.UUID,
		params UpdateParams,
	) (models.Role, error)
	DeleteRole(
		ctx context.Context,
		roleID uuid.UUID,
	) error

	GetRoles(
		ctx context.Context,
		filter FilterParams,
		limit, offset uint,
	) (pagi.Page[[]models.Role], error)

	//Ranks

	GetRoleRank(ctx context.Context, roleID uuid.UUID) (models.OrgRoleRank, error)
	GetOrgRolesRanks(ctx context.Context, organizationID uuid.UUID) ([]models.OrgRoleRank, error)

	CreateRoleRank(ctx context.Context, roleID uuid.UUID, rank int32) ([]models.OrgRoleRank, error)
	DeleteRoleRank(ctx context.Context, roleID uuid.UUID) ([]models.OrgRoleRank, error)
	GetMemberMaxRoleRank(ctx context.Context, memberID uuid.UUID) (int32, error)

	LockOrgRoleRankRevision(ctx context.Context, organizationID uuid.UUID) error
	BumpOrgRoleRankRevision(ctx context.Context, organizationID uuid.UUID) (models.OrgRoleRanksRevision, error)

	UpdateRolesRanks(ctx context.Context, organizationID uuid.UUID, ranks map[int32]uuid.UUID) ([]models.OrgRoleRank, error)

	//Member roles

	GetMember(ctx context.Context, memberID uuid.UUID) (models.Member, error)
	GetMemberByAccountAndOrganization(
		ctx context.Context,
		accountID, organizationID uuid.UUID,
	) (models.Member, error)

	LockMemberRolesLinksRevision(ctx context.Context, memberID uuid.UUID) error
	BumpMemberRolesLinksRevision(ctx context.Context, memberID uuid.UUID) (models.OrgMemberRoleLinkRevision, error)

	RemoveMemberRole(
		ctx context.Context,
		memberID, roleID uuid.UUID,
	) ([]uuid.UUID, error)
	AddMemberRole(
		ctx context.Context,
		memberID, roleID uuid.UUID,
	) ([]uuid.UUID, error)

	//Permissions

	GetAllPermissions(ctx context.Context) ([]models.OrgRolePermission, error)

	GetRolePermissions(
		ctx context.Context,
		roleID uuid.UUID,
	) (models.OrgRolePermissionsWithDetailsForRole, error)

	UpdateRolePermissions(
		ctx context.Context,
		roleID uuid.UUID,
		permissions []uuid.UUID,
	) ([]uuid.UUID, error)

	CreateRolePermissionsRevision(
		ctx context.Context,
		roleID uuid.UUID,
	) error
	LockRolePermissionsRevision(
		ctx context.Context,
		roleID uuid.UUID,
	) error
	BumpRolePermissionsRevision(
		ctx context.Context,
		roleID uuid.UUID,
	) (models.OrgRolePermissionsLinksRevision, error)

	CheckMemberHavePermission(
		ctx context.Context,
		memberID uuid.UUID,
		permissionID uuid.UUID,
	) (bool, error)

	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type messenger interface {
	WriteOrgRoleCreated(
		ctx context.Context,
		role models.Role,
	) error
	WriteOrgRoleUpdated(
		ctx context.Context,
		role models.Role,
	) error
	WriteOrgRoleDeleted(
		ctx context.Context,
		role models.Role,
	) error

	WriteOrgRolePermissionsUpdated(
		ctx context.Context,
		role models.Role,
		permissions []uuid.UUID,
		revision models.OrgRolePermissionsLinksRevision,
	) error
	WriteOrgRolesRanksUpdated(
		ctx context.Context,
		organizationID uuid.UUID,
		ranks []models.OrgRoleRank,
		revision models.OrgRoleRanksRevision,
	) error
	WriteOrgMemberRolesUpdated(
		ctx context.Context,
		memberID uuid.UUID,
		roles []uuid.UUID,
		revision models.OrgMemberRoleLinkRevision,
	) error
}

func (m *Module) checkPermissionsToManageRole(
	ctx context.Context,
	initiator models.AccountActor,
	roleID uuid.UUID,
) (models.Role, error) {
	role, err := m.GetByID(ctx, roleID)
	if err != nil {
		return models.Role{}, err
	}

	member, err := m.getInitiator(ctx, initiator, role.OrganizationID)
	if err != nil {
		return models.Role{}, err
	}
	if member.Head {
		return role, nil
	}

	hasPermission, err := m.repo.CheckMemberHavePermission(ctx, member.ID, orgperm.RolesManageID)
	if err != nil {
		return models.Role{}, err
	}
	if !hasPermission {
		return models.Role{}, errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("member %s does not have permission %s", member.ID, orgperm.RolesManageID),
		)
	}

	initiatorMaxRank, err := m.repo.GetMemberMaxRoleRank(ctx, member.ID)
	if err != nil {
		if errors.Is(err, errx.ErrorRoleNotFound) {
			return models.Role{}, errx.ErrorNotEnoughRights.Raise(
				fmt.Errorf("member %s has no roles and cannot manage role %s", member.ID, role.ID),
			)
		}
		return models.Role{}, fmt.Errorf("failed to get max role for member %s: %w", member.AccountID, err)
	}

	roleRank, err := m.repo.GetRoleRank(ctx, role.ID)
	if err != nil {
		return models.Role{}, fmt.Errorf("failed to get ranks for organization %s: %w", role.OrganizationID, err)
	}

	if initiatorMaxRank < roleRank.Rank {
		return models.Role{}, errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("member %s with max role rank %d cannot manage role with rank %d",
				member.ID, initiatorMaxRank, roleRank.Rank,
			),
		)
	}

	return role, nil
}

func (m *Module) getInitiator(
	ctx context.Context,
	initiator models.AccountActor,
	organizationID uuid.UUID,
) (models.Member, error) {
	member, err := m.repo.GetMemberByAccountAndOrganization(ctx, initiator, organizationID)
	if errors.Is(err, errx.ErrorMemberNotFound) {
		return models.Member{}, errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf(
				"initiator member with account id %s and organization id %s not found: %w",
				initiator, organizationID, err,
			),
		)
	} else if err != nil {
		return models.Member{}, err
	}

	return member, nil
}

func (m *Module) getMember(
	ctx context.Context,
	memberID uuid.UUID,
) (models.Member, error) {
	member, err := m.repo.GetMember(ctx, memberID)
	if err != nil {
		return models.Member{}, err
	}

	return member, nil
}
