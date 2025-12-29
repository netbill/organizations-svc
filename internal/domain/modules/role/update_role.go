package role

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/errx"
	"github.com/umisto/cities-svc/internal/domain/models"
)

type UpdateParams struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Color       *string `json:"color"`
}

func (s Service) UpdateRole(ctx context.Context, roleID uuid.UUID, params UpdateParams) (role models.Role, err error) {
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

func (s Service) UpdateRoleByUser(
	ctx context.Context,
	accountID uuid.UUID,
	roleID uuid.UUID,
	params UpdateParams,
) (models.Role, error) {
	role, err := s.GetRole(ctx, roleID)
	if err != nil {
		return models.Role{}, err
	}

	if err = s.CheckPermissionsToManageRole(ctx, accountID, role.AgglomerationID, role.Rank); err != nil {
		return models.Role{}, err
	}

	return s.UpdateRole(ctx, roleID, params)
}

func (s Service) UpdateRoleRank(ctx context.Context, roleID uuid.UUID, newRank uint) (role models.Role, err error) {
	if err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		role, err = s.repo.UpdateRoleRank(ctx, roleID, newRank)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to update role rank: %w", err),
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

func (s Service) UpdateRoleRankByUser(
	ctx context.Context,
	accountID uuid.UUID,
	roleID uuid.UUID,
	newRank uint,
) (models.Role, error) {
	role, err := s.GetRole(ctx, roleID)
	if err != nil {
		return models.Role{}, err
	}

	if err = s.CheckPermissionsToManageRole(ctx, accountID, role.AgglomerationID, role.Rank); err != nil {
		return models.Role{}, err
	}

	if err = s.CheckPermissionsToManageRole(ctx, accountID, role.AgglomerationID, newRank); err != nil {
		return models.Role{}, err
	}

	return s.UpdateRoleRank(ctx, roleID, newRank)
}

func (s Service) UpdateRolesRanks(
	ctx context.Context,
	agglomerationID uuid.UUID,
	order map[uint]uuid.UUID,
) error {
	if err := s.repo.Transaction(ctx, func(ctx context.Context) error {
		if err := s.repo.UpdateRolesRanks(ctx, agglomerationID, order); err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to update roles ranks: %w", err),
			)
		}

		if err := s.messenger.WriteRoleRanksUpdated(ctx, agglomerationID, order); err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to send role ranks updated message: %w", err),
			)
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (s Service) UpdateRolesRanksByUser(
	ctx context.Context,
	accountID uuid.UUID,
	agglomerationID uuid.UUID,
	order map[uint]uuid.UUID,
) error {
	maxRole, err := s.repo.GetAccountMaxRoleInAgglomeration(ctx, accountID, agglomerationID)
	if err != nil {
		return err
	}

	rolesIDs := make(map[uuid.UUID]struct{})
	for _, roleID := range order {
		rolesIDs[roleID] = struct{}{}
	}

	RolesBefore, err := s.repo.FilterRoles(
		ctx,
		FilterParams{
			AgglomerationID: &agglomerationID,
		},
		0,
		0,
	)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to filter roles: %w", err),
		)
	}

	for _, role := range RolesBefore.Data {
		if _, ok := rolesIDs[role.ID]; ok {
			if err = s.CheckPermissionsToManageRole(ctx, accountID, agglomerationID, role.Rank); err != nil {
				return err
			}
			if role.Rank < maxRole.Rank {
				return errx.ErrorNotEnoughRights.Raise(
					fmt.Errorf("member %s with max role rank %d cannot manage role with rank %d",
						accountID, maxRole.Rank, role.Rank),
				)
			}
		}
	}

	for rank := range order {
		if err = s.CheckPermissionsToManageRole(ctx, accountID, agglomerationID, rank); err != nil {
			return err
		}
		if rank < maxRole.Rank {
			return errx.ErrorNotEnoughRights.Raise(
				fmt.Errorf("member %s with max role rank %d cannot manage role with rank %d",
					accountID, maxRole.Rank, rank),
			)
		}
	}

	return s.UpdateRolesRanks(ctx, agglomerationID, order)
}
