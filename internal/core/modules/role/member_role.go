package role

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
)

func (s Service) GetMemberRoles(ctx context.Context, memberID uuid.UUID) ([]models.Role, error) {
	roles, err := s.repo.GetMemberRoles(ctx, memberID)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (s Service) GetMemberMaxRole(ctx context.Context, memberID uuid.UUID) (models.Role, error) {
	role, err := s.repo.GetMemberMaxRole(ctx, memberID)
	if err != nil {
		return models.Role{}, err
	}

	return role, nil
}

func (s Service) MemberAddRole(
	ctx context.Context,
	accountID, memberID, roleID uuid.UUID,
) error {
	member, err := s.getMember(ctx, memberID)
	if err != nil {
		return err
	}

	initiator, err := s.getInitiator(ctx, accountID, member.OrganizationID)
	if err != nil {
		return err
	}

	role, err := s.GetRole(ctx, roleID)
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

	if err = s.checkPermissionsToManageRole(ctx, initiator.AccountID, role.Rank); err != nil {
		return err
	}

	return s.repo.Transaction(ctx, func(ctx context.Context) error {
		if err = s.repo.AddMemberRole(ctx, memberID, roleID); err != nil {
			return err
		}

		if err = s.messenger.WriteOrgMemberRoleAdd(ctx, memberID, roleID); err != nil {
			return err
		}

		return nil
	})
}

func (s Service) RemoveMemberRole(
	ctx context.Context,
	accountID, memberID, roleID uuid.UUID,
) error {
	member, err := s.getMember(ctx, memberID)
	if err != nil {
		return err
	}

	initiator, err := s.getInitiator(ctx, accountID, member.OrganizationID)
	if err != nil {
		return err
	}

	role, err := s.GetRole(ctx, roleID)
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

	if err = s.checkPermissionsToManageRole(ctx, initiator.AccountID, role.Rank); err != nil {
		return err
	}

	return s.repo.Transaction(ctx, func(ctx context.Context) error {
		if err = s.repo.RemoveMemberRole(ctx, memberID, roleID); err != nil {
			return err
		}

		if err = s.messenger.WriteOrgMemberRoleRemove(ctx, memberID, roleID); err != nil {
			return err
		}

		return nil
	})
}
