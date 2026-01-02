package role

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/domain/errx"
	"github.com/netbill/organizations-svc/internal/domain/models"
)

func (s Service) GetMemberRoles(ctx context.Context, memberID uuid.UUID) ([]models.Role, error) {
	roles, err := s.repo.GetMemberRoles(ctx, memberID)
	if err != nil {
		return nil, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get member roles: %w", err),
		)
	}

	return roles, nil
}

func (s Service) GetMemberMaxRole(ctx context.Context, memberID uuid.UUID) (models.Role, error) {
	role, err := s.repo.GetMemberMaxRole(ctx, memberID)
	if err != nil {
		return models.Role{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get member max role: %w", err),
		)
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

	if err = s.checkPermissionsToManageRole(ctx, initiator.ID, role.Rank); err != nil {
		return err
	}

	return s.repo.Transaction(ctx, func(txCtx context.Context) error {
		if err = s.repo.AddMemberRole(ctx, memberID, roleID); err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("role Service MemberAddRoleByUser: repo AddMemberRole: %w", err),
			)
		}

		if err = s.messenger.WriteMemberRoleAdd(ctx, memberID, roleID); err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to write member role add event: %w", err),
			)
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

	if role.OrganizationID != member.OrganizationID {
		return errx.ErrorRoleNotFound.Raise(
			fmt.Errorf("role with id %s is not available in organization %s", role.ID, role.OrganizationID),
		)
	}

	if err = s.checkPermissionsToManageRole(ctx, initiator.ID, role.Rank); err != nil {
		return err
	}

	return s.repo.Transaction(ctx, func(txCtx context.Context) error {
		if err = s.repo.RemoveMemberRole(ctx, memberID, roleID); err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("role Service MemberRemoveRoleByUser: repo RemoveMemberRole: %w", err),
			)
		}

		if err = s.messenger.WriteMemberRoleRemove(ctx, memberID, roleID); err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to write member role remove event: %w", err),
			)
		}

		return nil
	})
}
