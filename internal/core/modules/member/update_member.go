package member

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
)

type UpdateParams struct {
	Position *string
	Label    *string
}

func (s Service) UpdateMember(
	ctx context.Context,
	accountID, memberID uuid.UUID,
	params UpdateParams,
) (models.Member, error) {
	member, err := s.GetMemberByID(ctx, memberID)
	if err != nil {
		return models.Member{}, err
	}

	initiator, err := s.GetInitiatorMember(ctx, accountID, member.OrganizationID)
	if err != nil {
		return models.Member{}, err
	}

	hasPermission, err := s.repo.CheckMemberHavePermission(
		ctx,
		initiator.ID,
		models.RolePermissionManageMembers,
	)
	if err != nil {
		return models.Member{}, err
	}
	if !hasPermission {
		return models.Member{}, errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("initiator member %s has no manage members permission", initiator.ID),
		)
	}

	firstMaxRole, err := s.repo.GetMemberMaxRole(ctx, initiator.ID)
	if err != nil {
		return models.Member{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get max role for member %s: %w", initiator.ID, err),
		)
	}
	if firstMaxRole.Head == false {
		secMaxRole, err := s.repo.GetMemberMaxRole(ctx, member.ID)
		if err != nil {
			return models.Member{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get max role for member %s: %w", member.ID, err),
			)
		}

		if firstMaxRole.Rank < secMaxRole.Rank {
			return models.Member{}, errx.ErrorNotEnoughRights.Raise(
				fmt.Errorf(
					"member %s with rank %d cannot manage member %s with rank %d",
					initiator.ID,
					firstMaxRole.Rank,
					member.ID,
					secMaxRole.Rank,
				),
			)
		}
	}

	err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		member, err = s.repo.UpdateMember(ctx, memberID, params)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to update member %s: %w", memberID, err),
			)
		}

		if err = s.messenger.WriteMemberUpdated(ctx, member); err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to send member updated message for member %s: %w", memberID, err),
			)
		}
		return nil
	})

	return member, err
}
