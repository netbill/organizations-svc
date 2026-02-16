package role

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/orgperm"
)

func (m *Module) UpdateRanks(
	ctx context.Context,
	initiator models.AccountActor,
	organizationID uuid.UUID,
	order map[uuid.UUID]uint,
) error {
	member, err := m.getInitiator(ctx, initiator, organizationID)
	if err != nil {
		return err
	}

	if !member.Head {
		maxRole, err := m.repo.GetMemberMaxRole(ctx, member.ID)
		if err != nil {
			return err
		}

		rolesIDs := make(map[uuid.UUID]struct{}, len(order))
		for roleID := range order {
			rolesIDs[roleID] = struct{}{}
		}

		hasPermission, err := m.repo.CheckMemberHavePermission(ctx, member.AccountID, orgperm.RolesManageID)
		if err != nil {
			return err
		}
		if !hasPermission {
			return errx.ErrorNotEnoughRights.Raise(
				fmt.Errorf("member %s does not have permission %s", initiator, orgperm.RolesManageID),
			)
		}

		rolesBefore, err := m.repo.GetRoles(ctx, FilterParams{
			OrganizationID: &organizationID,
		}, 1000, 0)
		if err != nil {
			return err
		}

		for _, role := range rolesBefore.Data {
			if _, ok := rolesIDs[role.ID]; !ok {
				continue
			}

			if role.Rank >= maxRole.Rank {
				return errx.ErrorNotEnoughRights.Raise(
					fmt.Errorf(
						"member %s with max role rank %d cannot manage role with rank %d",
						initiator, maxRole.Rank, role.Rank,
					),
				)
			}
		}

		for _, newRank := range order {
			if newRank >= maxRole.Rank {
				return errx.ErrorNotEnoughRights.Raise(
					fmt.Errorf(
						"member %s with max role rank %d cannot manage role with rank %d",
						initiator, maxRole.Rank, newRank,
					),
				)
			}
		}
	}

	return m.repo.Transaction(ctx, func(ctx context.Context) error {
		if err = m.repo.UpdateRolesRanks(ctx, organizationID, order); err != nil {
			return err
		}

		if err = m.messenger.WriteOrgRolesRanksUpdated(ctx, organizationID, order); err != nil {
			return err
		}

		return nil
	})
}
