package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/internal/core/modules/role"
	"github.com/netbill/organizations-svc/internal/repository/pgdb"
	"github.com/netbill/pagi"
)

func (r Repository) CreateRole(ctx context.Context, params role.CreateParams) (models.Role, error) {
	row, err := r.orgRolesQ(ctx).Insert(ctx, pgdb.InsertRoleParams{
		OrganizationID: params.OrganizationID,
		Rank:           params.Rank,
		Name:           params.Name,
		Description:    params.Description,
		Color:          params.Color,
	})
	if err != nil {
		return models.Role{}, fmt.Errorf(
			"failed to create role for organization ID %s cause: %w",
			params.OrganizationID, err,
		)
	}

	return Role(row), nil
}

func (r Repository) CreateHeadRole(ctx context.Context, organizationID uuid.UUID) (models.Role, error) {
	row, err := r.orgRolesQ(ctx).Insert(ctx, pgdb.InsertRoleParams{
		OrganizationID: organizationID,
		Head:           true,
		Rank:           1,
		Name:           "Head",
		Description:    "Head role with all permissions",
		Color:          "#000000",
	})
	if err != nil {
		return models.Role{}, fmt.Errorf(
			"failed to create head role for organization ID %s cause: %w",
			organizationID, err,
		)
	}

	return Role(row), nil
}

func (r Repository) GetRole(ctx context.Context, roleID uuid.UUID) (models.Role, error) {
	row, err := r.orgRolesQ(ctx).FilterByID(roleID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return models.Role{}, errx.ErrorRoleNotFound.Raise(
				fmt.Errorf("role with ID %s not found, cause: %w", roleID, err),
			)
		default:
			return models.Role{}, fmt.Errorf("failed to get role with ID %s cause: %w", roleID, err)
		}
	}

	return Role(row), nil
}

func (r Repository) GetRoles(
	ctx context.Context,
	filter role.FilterParams,
	limit, offset uint,
) (pagi.Page[[]models.Role], error) {
	q := r.orgRolesQ(ctx)
	if filter.OrganizationID != nil {
		q = q.FilterByOrganizationID(*filter.OrganizationID)
	}
	if filter.RolesID != nil && len(*filter.RolesID) > 0 {
		q = q.FilterByID(*filter.RolesID...)
	}
	if filter.Head != nil {
		q = q.FilterHead(*filter.Head)
	}
	if filter.Rank != nil {
		q = q.FilterByRank(*filter.Rank)
	}
	if filter.Name != nil {
		q = q.FilterLikeName(*filter.Name)
	}

	if limit == 0 {
		limit = 10
	}

	rows, err := q.OrderByRoleRank(false).Page(limit, offset).Select(ctx)
	if err != nil {
		return pagi.Page[[]models.Role]{}, fmt.Errorf("failed to get roles, cause: %w", err)
	}

	total, err := q.Count(ctx)
	if err != nil {
		return pagi.Page[[]models.Role]{}, fmt.Errorf("failed to count roles, cause: %w", err)
	}

	collection := make([]models.Role, 0, len(rows))
	for _, row := range rows {
		collection = append(collection, Role(row))
	}

	return pagi.Page[[]models.Role]{
		Data:  collection,
		Total: total,
		Page:  uint(offset/limit) + 1,
		Size:  uint(len(collection)),
	}, nil
}

func (r Repository) UpdateRole(ctx context.Context, roleID uuid.UUID, params role.UpdateParams) (models.Role, error) {
	q := r.orgRolesQ(ctx).FilterByID(roleID)
	if params.Name != nil {
		q = q.UpdateName(*params.Name)
	}
	if params.Description != nil {
		q = q.UpdateDescription(*params.Description)
	}
	if params.Color != nil {
		q = q.UpdateColor(*params.Color)
	}

	row, err := q.UpdateOne(ctx)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return models.Role{}, errx.ErrorRoleNotFound.Raise(
			fmt.Errorf("role with ID %s not found, cause: %w", roleID, err),
		)
	case err != nil:
		return models.Role{}, fmt.Errorf("failed to update role with ID %s cause: %w", roleID, err)
	}

	return Role(row), nil
}

func (r Repository) UpdateRoleRank(ctx context.Context, roleID uuid.UUID, newRank uint) (models.Role, error) {
	row, err := r.orgRolesQ(ctx).UpdateRoleRank(ctx, roleID, newRank)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return models.Role{}, errx.ErrorRoleNotFound.Raise(
			fmt.Errorf("role with ID %s not found, cause: %w", roleID, err),
		)
	case err != nil:
		return models.Role{}, fmt.Errorf("failed to update rank for role with ID %s cause: %w", roleID, err)
	}

	return Role(row), nil
}

func (r Repository) UpdateRolesRanks(
	ctx context.Context,
	organizationID uuid.UUID,
	order map[uuid.UUID]uint,
) error {
	_, err := r.orgRolesQ(ctx).UpdateRolesRanks(ctx, organizationID, order)
	if err != nil {
		return fmt.Errorf("failed to update roles ranks for organization ID %s cause: %w", organizationID, err)
	}

	return nil
}

func (r Repository) DeleteRole(ctx context.Context, roleID uuid.UUID) error {
	err := r.orgRolesQ(ctx).DeleteAndShiftRanks(ctx, roleID)
	if err != nil {
		return fmt.Errorf("failed to delete role with ID %s cause: %w", roleID, err)
	}

	return nil
}

func (r Repository) GetMemberMaxRole(
	ctx context.Context,
	memberID uuid.UUID,
) (models.Role, error) {
	res, err := r.orgRolesQ(ctx).
		FilterByMemberID(memberID).
		OrderByRoleRank(false). // DESC => max
		Get(ctx)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return models.Role{}, errx.ErrorRoleNotFound.Raise(
			fmt.Errorf("no roles found for member ID %s, cause: %w", memberID, err),
		)
	case err != nil:
		return models.Role{}, fmt.Errorf("failed to get max role for member ID %s, cause: %w", memberID, err)
	}

	return Role(res), nil
}

func Role(r pgdb.OrganizationRole) models.Role {
	return models.Role{
		ID:             r.ID,
		OrganizationID: r.OrganizationID,
		Head:           r.Head,
		Rank:           r.Rank,
		Name:           r.Name,
		Description:    r.Description,
		Color:          r.Color,
		CreatedAt:      r.CreatedAt,
		UpdatedAt:      r.UpdatedAt,
	}
}
