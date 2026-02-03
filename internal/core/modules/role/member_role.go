package role

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
)

func (m *Module) GetMemberRoles(
	ctx context.Context,
	memberID uuid.UUID,
) ([]models.Role, error) {
	roles, err := m.repo.GetMemberRoles(ctx, memberID)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (m *Module) GetMemberMaxRole(
	ctx context.Context,
	memberID uuid.UUID,
) (models.Role, error) {
	role, err := m.repo.GetMemberMaxRole(ctx, memberID)
	if err != nil {
		return models.Role{}, err
	}

	return role, nil
}

type MemberAddRoleParams struct {
	MemberID uuid.UUID
	RoleID   uuid.UUID
}

func (m *Module) MemberAddRole(
	ctx context.Context,
	initiator models.InitiatorData,
	memberID, roleID uuid.UUID,
) error {
	member, err := m.getMember(ctx, memberID)
	if err != nil {
		return err
	}

	role, err := m.GetRole(ctx, roleID)
	if err != nil {
		return err
	}

	if role.OrganizationID != member.OrganizationID {
		return errx.ErrorRoleNotFound.Raise(
			fmt.Errorf("role with id %s is not available in organization %s", role.ID, role.OrganizationID),
		)
	}

	err = m.checkPermissionsToManageRole(ctx, initiator, role.OrganizationID, role.Rank)
	if err != nil {
		return err
	}

	return m.repo.Transaction(ctx, func(ctx context.Context) error {
		if err = m.repo.AddMemberRole(ctx, memberID, roleID); err != nil {
			return err
		}

		if err = m.messenger.WriteOrgMemberRoleAdd(ctx, memberID, roleID); err != nil {
			return err
		}

		return nil
	})
}

type MemberRemoveRoleParams struct {
	MemberID uuid.UUID
	RoleID   uuid.UUID
}

func (m *Module) MemberRemoveRole(
	ctx context.Context,
	initiator models.InitiatorData,
	memberID, roleID uuid.UUID,
) error {
	member, err := m.getMember(ctx, memberID)
	if err != nil {
		return err
	}

	role, err := m.GetRole(ctx, roleID)
	if err != nil {
		return err
	}

	if role.OrganizationID != member.OrganizationID {
		return errx.ErrorRoleNotFound.Raise(
			fmt.Errorf("role with id %s is not available in organization %s", role.ID, role.OrganizationID),
		)
	}

	err = m.checkPermissionsToManageRole(ctx, initiator, role.OrganizationID, role.Rank)
	if err != nil {
		return err
	}

	return m.repo.Transaction(ctx, func(ctx context.Context) error {
		if err = m.repo.RemoveMemberRole(ctx, memberID, roleID); err != nil {
			return err
		}

		if err = m.messenger.WriteOrgMemberRoleRemove(ctx, memberID, roleID); err != nil {
			return err
		}

		return nil
	})
}
