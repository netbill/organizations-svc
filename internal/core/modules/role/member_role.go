package role

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
)

func (m *Module) GetMemberRoles(ctx context.Context, memberID uuid.UUID) ([]models.Role, error) {
	roles, err := m.repo.GetMemberRoles(ctx, memberID)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (m *Module) GetMemberMaxRole(ctx context.Context, memberID uuid.UUID) (models.Role, error) {
	role, err := m.repo.GetMemberMaxRole(ctx, memberID)
	if err != nil {
		return models.Role{}, err
	}

	return role, nil
}

func (m *Module) MemberAddRole(
	ctx context.Context,
	accountID, memberID, roleID uuid.UUID,
) error {
	member, err := m.getMember(ctx, memberID)
	if err != nil {
		return err
	}

	initiator, err := m.getInitiator(ctx, accountID, member.OrganizationID)
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

	if role.Head {
		return errx.ErrorCannotAddHeadRoleToMember.Raise(
			fmt.Errorf("cannot add head role %s to member %s", role.ID, member.ID),
		)
	}

	if err = m.checkPermissionsToManageRole(ctx, initiator.AccountID, role.Rank); err != nil {
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

func (m *Module) RemoveMemberRole(
	ctx context.Context,
	accountID, memberID, roleID uuid.UUID,
) error {
	member, err := m.getMember(ctx, memberID)
	if err != nil {
		return err
	}

	initiator, err := m.getInitiator(ctx, accountID, member.OrganizationID)
	if err != nil {
		return err
	}

	role, err := m.GetRole(ctx, roleID)
	if err != nil {
		return err
	}

	if role.Head {
		return errx.ErrorCannotRemoveHeadRoleFromMember.Raise(
			fmt.Errorf("cannot remove head role %s from member %s", role.ID, member.ID),
		)
	}

	if role.OrganizationID != member.OrganizationID {
		return errx.ErrorRoleNotFound.Raise(
			fmt.Errorf("role with id %s is not available in organization %s", role.ID, role.OrganizationID),
		)
	}

	if err = m.checkPermissionsToManageRole(ctx, initiator.AccountID, role.Rank); err != nil {
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
