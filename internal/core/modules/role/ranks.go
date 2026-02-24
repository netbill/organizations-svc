package role

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/orgperm"
)

func (m *Module) ReorderRanks(
	ctx context.Context,
	initiator models.AccountActor,
	organizationID uuid.UUID,
	order map[int32]uuid.UUID,
) error {
	member, err := m.getInitiator(ctx, initiator, organizationID)
	if err != nil {
		return err
	}

	ranks, err := m.repo.GetOrgRolesRanks(ctx, organizationID)
	if err != nil {
		return err
	}

	if len(order) != len(ranks) {
		return errx.ErrorInvalidOrgRoleRanksOrder.Raise(
			fmt.Errorf("order must contain all organization roles"),
		)
	}

	currentRoles := make(map[uuid.UUID]models.OrgRoleRank, len(ranks))
	for _, r := range ranks {
		currentRoles[r.RoleID] = r
	}

	for _, roleID := range order {
		if _, ok := currentRoles[roleID]; !ok {
			return errx.ErrorInvalidOrgRoleRanksOrder.Raise(
				fmt.Errorf("role %s does not belong to organization %s", roleID, organizationID),
			)
		}
	}

	if !member.Head {

		hasPermission, err := m.repo.CheckMemberHavePermission(
			ctx,
			member.ID,
			orgperm.RolesManageID,
		)
		if err != nil {
			return err
		}
		if !hasPermission {
			return errx.ErrorNotEnoughRights.Raise(
				fmt.Errorf("member %s does not have permission %s", member.ID, orgperm.RolesManageID),
			)
		}

		initiatorMaxRank, err := m.repo.GetMemberMaxRoleRank(ctx, member.ID)
		if err != nil {
			return fmt.Errorf("failed to get max role for member %s: %w", member.AccountID, err)
		}

		for _, r := range ranks {
			if r.Rank >= initiatorMaxRank {
				if _, touched := order[r.Rank]; touched {
					return errx.ErrorNotEnoughRights.Raise(
						fmt.Errorf(
							"member %s cannot manage role with rank %d",
							member.ID,
							r.Rank,
						),
					)
				}
			}
		}

		for newRank := range order {
			if newRank >= initiatorMaxRank {
				return errx.ErrorNotEnoughRights.Raise(
					fmt.Errorf(
						"member %s cannot assign rank %d",
						member.ID,
						newRank,
					),
				)
			}
		}
	}

	return m.repo.Transaction(ctx, func(ctx context.Context) error {
		if err = m.repo.LockOrgRoleRankRevision(ctx, organizationID); err != nil {
			return fmt.Errorf("failed to lock org roles for organization %s: %w", organizationID, err)
		}

		updatedRanks, err := m.repo.UpdateRolesRanks(ctx, organizationID, order)
		if err != nil {
			return err
		}

		revision, err := m.repo.BumpOrgRoleRankRevision(ctx, organizationID)
		if err != nil {
			return fmt.Errorf("failed to update org role ranks revision for organization %s: %w", organizationID, err)
		}

		if err = m.messenger.WriteOrgRolesRanksUpdated(ctx, organizationID, updatedRanks, revision); err != nil {
			return err
		}

		return nil
	})
}
