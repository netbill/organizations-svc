package role

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
)

type UpdateParams struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Color       *string `json:"color"`
}

func (m *Module) UpdateRole(
	ctx context.Context,
	accountID uuid.UUID,
	roleID uuid.UUID,
	params UpdateParams,
) (models.Role, error) {
	role, err := m.GetRole(ctx, roleID)
	if err != nil {
		return models.Role{}, err
	}

	initiator, err := m.getInitiator(ctx, accountID, role.OrganizationID)
	if err != nil {
		return models.Role{}, err
	}

	if err = m.checkPermissionsToManageRole(ctx, initiator.AccountID, role.Rank); err != nil {
		return models.Role{}, err
	}

	if err = m.repo.Transaction(ctx, func(ctx context.Context) error {
		role, err = m.repo.UpdateRole(ctx, roleID, params)
		if err != nil {
			return err
		}

		if err = m.messenger.WriteOrgRoleUpdated(ctx, role); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return models.Role{}, err
	}

	return role, nil
}

func (m *Module) UpdateRolesRanks(
	ctx context.Context,
	accountID uuid.UUID,
	organizationID uuid.UUID,
	order map[uuid.UUID]uint,
) error {
	initiator, err := m.getInitiator(ctx, accountID, organizationID)
	if err != nil {
		return err
	}

	maxRole, err := m.repo.GetMemberMaxRole(ctx, initiator.ID)
	if err != nil {
		return err
	}

	rolesIDs := make(map[uuid.UUID]struct{}, len(order))
	for roleID := range order {
		rolesIDs[roleID] = struct{}{}
	}

	rankToRole := make(map[uint]uuid.UUID, len(order))
	for roleID, newRank := range order {
		if prevRoleID, ok := rankToRole[newRank]; ok && prevRoleID != roleID {
			return err
		}
		rankToRole[newRank] = roleID
	}

	hasPermission, err := m.repo.CheckMemberHavePermission(
		ctx,
		initiator.AccountID,
		models.RolePermissionManageRoles,
	)
	if err != nil {
		return err
	}
	if !hasPermission {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("member %s does not have permission %s", initiator.AccountID, models.RolePermissionManageRoles),
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

		if role.Head {
			return errx.ErrorCannotUpdateHeadRoleRank.Raise(
				fmt.Errorf("cannot update rank of head role %s", role.ID),
			)
		}

		if role.Rank >= maxRole.Rank {
			return errx.ErrorNotEnoughRights.Raise(
				fmt.Errorf("member %s with max role rank %d cannot manage role with rank %d",
					accountID, maxRole.Rank, role.Rank),
			)
		}
	}

	for _, newRank := range order {
		if newRank >= maxRole.Rank {
			return errx.ErrorNotEnoughRights.Raise(
				fmt.Errorf("member %s with max role rank %d cannot manage role with rank %d",
					accountID, maxRole.Rank, newRank),
			)
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
