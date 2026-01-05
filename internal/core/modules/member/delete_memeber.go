package member

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
)

func (s Service) DeleteMember(ctx context.Context, accountID, memberID uuid.UUID) error {
	member, err := s.GetMemberByID(ctx, memberID)
	if err != nil {
		return err
	}

	initiator, err := s.GetInitiatorMember(ctx, accountID, member.OrganizationID)
	if err != nil {
		return err
	}

	hasPermission, err := s.repo.CheckMemberHavePermission(
		ctx,
		initiator.ID,
		models.RolePermissionManageMembers,
	)
	if err != nil {
		return err
	}
	if !hasPermission {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("initiator member %s has no manage members permission", initiator.ID),
		)
	}

	firstMaxRole, err := s.repo.GetMemberMaxRole(ctx, initiator.ID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get max role for member %s: %w", initiator.ID, err),
		)
	}

	secMaxRole, err := s.repo.GetMemberMaxRole(ctx, member.ID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get max role for member %s: %w", member.ID, err),
		)
	}
	if secMaxRole.Head {
		return errx.ErrorCannotDeleteOrganizationHeadMember.Raise(
			fmt.Errorf("cannot delete organization head member %s", member.ID),
		)
	}

	if firstMaxRole.Rank < secMaxRole.Rank {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf(
				"member %s with rank %d cannot manage member %s with rank %d",
				initiator.ID,
				firstMaxRole.Rank,
				member.ID,
				secMaxRole.Rank,
			),
		)
	}

	return s.repo.Transaction(ctx, func(ctx context.Context) error {
		err = s.repo.DeleteMember(ctx, memberID)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to delete member %s: %w", memberID, err),
			)
		}

		if err = s.messenger.WriteMemberDeleted(ctx, member); err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to send member deleted message for member %s: %w", memberID, err),
			)
		}

		return nil
	})
}
