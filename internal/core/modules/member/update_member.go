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

func (m *Module) UpdateMember(
	ctx context.Context,
	accountID, memberID uuid.UUID,
	params UpdateParams,
) (models.Member, error) {
	member, err := m.GetMemberByID(ctx, memberID)
	if err != nil {
		return models.Member{}, err
	}

	initiator, err := m.GetInitiatorMember(ctx, accountID, member.OrganizationID)
	if err != nil {
		return models.Member{}, err
	}

	hasPermission, err := m.repo.CheckMemberHavePermission(
		ctx,
		initiator.AccountID,
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

	firstMaxRole, err := m.repo.GetMemberMaxRole(ctx, initiator.ID)
	if err != nil {
		return models.Member{}, fmt.Errorf("failed to get max role for member %s: %w", initiator.AccountID, err)
	}
	if firstMaxRole.Head == false {
		secMaxRole, err := m.repo.GetMemberMaxRole(ctx, member.ID)
		if err != nil {
			return models.Member{}, fmt.Errorf("failed to get max role for member %s: %w", member.ID, err)
		}

		if firstMaxRole.Rank < secMaxRole.Rank {
			return models.Member{}, errx.ErrorNotEnoughRights.Raise(
				fmt.Errorf(
					"member %s with rank %d cannot manage member %s with rank %d",
					initiator.AccountID,
					firstMaxRole.Rank,
					member.ID,
					secMaxRole.Rank,
				),
			)
		}
	}

	err = m.repo.Transaction(ctx, func(ctx context.Context) error {
		member, err = m.repo.UpdateMember(ctx, memberID, params)
		if err != nil {
			return fmt.Errorf("failed to update member %s: %w", memberID, err)
		}

		if err = m.messenger.WriteOrgMemberUpdated(ctx, member); err != nil {
			return fmt.Errorf("failed to send member updated message for member %s: %w", memberID, err)
		}
		return nil
	})

	return member, err
}
