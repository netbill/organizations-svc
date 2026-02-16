package perm

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/orgperm"
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
	GetMember(ctx context.Context, memberID uuid.UUID) (models.Member, error)
	GetMemberByAccountAndOrganization(
		ctx context.Context,
		accountID, organizationID uuid.UUID,
	) (models.Member, error)

	GetRole(ctx context.Context, roleID uuid.UUID) (models.Role, error)
	GetMemberMaxRole(ctx context.Context, memberID uuid.UUID) (models.Role, error)

	GetAllPermissions(ctx context.Context) ([]models.OrgRolePermission, error)
	GetRolePermissions(
		ctx context.Context,
		roleID uuid.UUID,
	) (models.OrgRolePermissionsWithDetailsForRole, error)
	SetRolePermissions(
		ctx context.Context,
		roleID uuid.UUID,
		params SetForRole,
	) error

	CheckMemberHavePermission(
		ctx context.Context,
		memberID uuid.UUID,
		permissionID uuid.UUID,
	) (bool, error)

	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type messenger interface {
	WriteOrgRolePermissionsUpdated(
		ctx context.Context,
		role models.Role,
		permissions SetForRole,
	) error
}

func (m *Module) checkPermissionsToManageRole(
	ctx context.Context,
	initiator models.AccountActor,
	organizationID uuid.UUID,
	rank uint,
) error {
	member, err := m.getInitiator(ctx, initiator, organizationID)
	if err != nil {
		return err
	}

	hasPermission, err := m.repo.CheckMemberHavePermission(ctx, member.ID, orgperm.RolesManageID)
	if err != nil {
		return err
	}
	if !hasPermission {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("member %s does not have permission %s", member.ID, orgperm.RolesManageCode),
		)
	}

	maxRole, err := m.repo.GetMemberMaxRole(ctx, member.ID)
	if err != nil {
		if errors.Is(err, errx.ErrorRoleNotFound) {
			return errx.ErrorNotEnoughRights.Raise(
				fmt.Errorf("member %s has no roles assigned: %w", member.ID, err),
			)
		}
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
	initiator models.AccountActor,
	organizationID uuid.UUID,
) (models.Member, error) {
	member, err := m.repo.GetMemberByAccountAndOrganization(ctx, initiator, organizationID)
	if err != nil {
		if errors.Is(err, errx.ErrorMemberNotFound) {
			return models.Member{}, errx.ErrorNotEnoughRights.Raise(
				fmt.Errorf(
					"initiator member with account id %s and organization id %s not found: %w",
					initiator, organizationID, err,
				),
			)
		}
	}

	return member, nil
}
