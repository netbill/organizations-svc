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

func (s Service) UpdateRole(
	ctx context.Context,
	accountID uuid.UUID,
	roleID uuid.UUID,
	params UpdateParams,
) (models.Role, error) {
	role, err := s.GetRole(ctx, roleID)
	if err != nil {
		return models.Role{}, err
	}

	initiator, err := s.getInitiator(ctx, accountID, role.OrganizationID)
	if err != nil {
		return models.Role{}, err
	}

	if err = s.checkPermissionsToManageRole(ctx, initiator.ID, role.Rank); err != nil {
		return models.Role{}, err
	}

	if err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		role, err = s.repo.UpdateRole(ctx, roleID, params)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to update role: %w", err),
			)
		}

		if err = s.messenger.WriteRoleUpdated(ctx, role); err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to send role updated message: %w", err),
			)
		}

		return nil
	}); err != nil {
		return models.Role{}, err
	}

	return role, nil
}

func (s Service) UpdateRolesRanks(
	ctx context.Context,
	accountID uuid.UUID,
	organizationID uuid.UUID,
	order map[uuid.UUID]uint,
) error {
	initiator, err := s.getInitiator(ctx, accountID, organizationID)
	if err != nil {
		return err
	}

	maxRole, err := s.repo.GetMemberMaxRole(ctx, initiator.ID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get account max role in organization: %w", err),
		)
	}

	rolesIDs := make(map[uuid.UUID]struct{}, len(order))
	for roleID := range order {
		rolesIDs[roleID] = struct{}{}
	}

	rankToRole := make(map[uint]uuid.UUID, len(order))
	for roleID, newRank := range order {
		if prevRoleID, ok := rankToRole[newRank]; ok && prevRoleID != roleID {
			return errx.ErrorInvalidInput.Raise(
				fmt.Errorf("duplicate rank %d for roles %s and %s", newRank, prevRoleID, roleID),
			)
		}
		rankToRole[newRank] = roleID
	}

	hasPermission, err := s.repo.CheckMemberHavePermission(
		ctx,
		initiator.ID,
		models.RolePermissionManageRoles,
	)
	if err != nil {
		return err
	}
	if !hasPermission {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("member %s does not have permission %s", initiator.ID, models.RolePermissionManageRoles),
		)
	}

	rolesBefore, err := s.repo.GetRoles(ctx, FilterParams{
		OrganizationID: &organizationID,
	}, 1000, 0)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to filter roles: %w", err),
		)
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

	return s.repo.Transaction(ctx, func(ctx context.Context) error {
		if err = s.repo.UpdateRolesRanks(ctx, organizationID, order); err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to update roles ranks: %w", err),
			)
		}

		if err = s.messenger.WriteRolesRanksUpdated(ctx, organizationID, order); err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to send role ranks updated message: %w", err),
			)
		}

		return nil
	})
}
