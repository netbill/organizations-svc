package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/internal/core/modules/role"
	"github.com/netbill/pagi"
)

type OrganizationRoleRow struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Head           bool      `json:"head"`
	Rank           uint      `json:"rank"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Color          string    `json:"color"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (r OrganizationRoleRow) IsNil() bool {
	return r.ID == uuid.Nil
}

func (r OrganizationRoleRow) ToModel() models.Role {
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

type OrgRolesQ interface {
	New() OrgRolesQ
	Insert(ctx context.Context, input OrganizationRoleRow) (OrganizationRoleRow, error)

	FilterByID(id ...uuid.UUID) OrgRolesQ
	FilterByOrganizationID(organizationID uuid.UUID) OrgRolesQ
	FilterByAccountID(accountID uuid.UUID) OrgRolesQ
	FilterByMemberID(memberID uuid.UUID) OrgRolesQ
	FilterHead(head bool) OrgRolesQ
	FilterByRank(rank int) OrgRolesQ
	FilterLikeName(name string) OrgRolesQ

	Get(ctx context.Context) (OrganizationRoleRow, error)
	Select(ctx context.Context) ([]OrganizationRoleRow, error)

	UpdateMany(ctx context.Context) (int64, error)
	UpdateOne(ctx context.Context) (OrganizationRoleRow, error)

	UpdateName(name string) OrgRolesQ
	UpdateDescription(description string) OrgRolesQ
	UpdateColor(color string) OrgRolesQ

	OrderByRoleRank(asc bool) OrgRolesQ
	Page(limit, offset uint) OrgRolesQ

	DeleteAndShiftRanks(ctx context.Context, roleID uuid.UUID) error
	UpdateRoleRank(ctx context.Context, roleID uuid.UUID, newRank uint) (OrganizationRoleRow, error)
	UpdateRolesRanks(
		ctx context.Context,
		organizationID uuid.UUID,
		order map[uuid.UUID]uint,
	) ([]OrganizationRoleRow, error)

	Count(ctx context.Context) (uint, error)
}

func (r Repository) CreateRole(ctx context.Context, params role.CreateParams) (models.Role, error) {
	row, err := r.orgRolesQ().Insert(ctx, OrganizationRoleRow{
		OrganizationID: params.OrganizationID,
		Rank:           params.Rank,
		Name:           params.Name,
		Description:    params.Description,
		Color:          params.Color,
	})
	if err != nil {
		return models.Role{}, fmt.Errorf("failed to create role, cause: %w", err)
	}

	return row.ToModel(), nil
}

func (r Repository) CreateHeadRole(ctx context.Context, organizationID uuid.UUID) (models.Role, error) {
	row, err := r.orgRolesQ().Insert(ctx, OrganizationRoleRow{
		OrganizationID: organizationID,
		Head:           true,
		Rank:           1,
		Name:           "Head",
		Description:    "Head role with all permissions",
		Color:          "#000000",
	})
	if err != nil {
		return models.Role{}, fmt.Errorf("failed to create head role, cause: %w", err)
	}

	return row.ToModel(), nil
}

func (r Repository) GetRole(ctx context.Context, roleID uuid.UUID) (models.Role, error) {
	row, err := r.orgRolesQ().FilterByID(roleID).Get(ctx)
	if err != nil {
		return models.Role{}, fmt.Errorf("failed to get role, cause: %w", err)
	}
	if row.IsNil() {
		return models.Role{}, errx.ErrorRoleNotFound.Raise(
			fmt.Errorf("role with ID %s not found", roleID),
		)
	}

	return row.ToModel(), nil
}

func (r Repository) GetRoles(
	ctx context.Context,
	filter role.FilterParams,
	limit, offset uint,
) (pagi.Page[[]models.Role], error) {
	q := r.orgRolesQ()
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
		collection = append(collection, row.ToModel())
	}

	return pagi.Page[[]models.Role]{
		Data:  collection,
		Total: total,
		Page:  uint(offset/limit) + 1,
		Size:  uint(len(collection)),
	}, nil
}

func (r Repository) UpdateRole(ctx context.Context, roleID uuid.UUID, params role.UpdateParams) (models.Role, error) {
	q := r.orgRolesQ().FilterByID(roleID)
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
	if err != nil {
		return models.Role{}, fmt.Errorf("failed to update role, cause: %w", err)
	}
	if row.IsNil() {
		return models.Role{}, errx.ErrorRoleNotFound.Raise(
			fmt.Errorf("role with ID %s not found", roleID),
		)
	}

	return row.ToModel(), nil
}

func (r Repository) UpdateRoleRank(ctx context.Context, roleID uuid.UUID, newRank uint) (models.Role, error) {
	row, err := r.orgRolesQ().UpdateRoleRank(ctx, roleID, newRank)
	if err != nil {
		return models.Role{}, fmt.Errorf("failed to update role rank, cause: %w", err)
	}
	if row.IsNil() {
		return models.Role{}, errx.ErrorRoleNotFound.Raise(
			fmt.Errorf("role with ID %s not found", roleID),
		)
	}

	return row.ToModel(), nil
}

func (r Repository) UpdateRolesRanks(
	ctx context.Context,
	organizationID uuid.UUID,
	order map[uuid.UUID]uint,
) error {
	_, err := r.orgRolesQ().UpdateRolesRanks(ctx, organizationID, order)
	if err != nil {
		return fmt.Errorf("failed to update roles ranks, cause: %w", err)
	}

	return nil
}

func (r Repository) DeleteRole(ctx context.Context, roleID uuid.UUID) error {
	err := r.orgRolesQ().DeleteAndShiftRanks(ctx, roleID)
	if err != nil {
		return fmt.Errorf("failed to delete role, cause: %w", err)
	}

	return nil
}

func (r Repository) GetMemberMaxRole(
	ctx context.Context,
	memberID uuid.UUID,
) (models.Role, error) {
	row, err := r.orgRolesQ().
		FilterByMemberID(memberID).
		OrderByRoleRank(false). // DESC => max
		Get(ctx)
	if err != nil {
		return models.Role{}, fmt.Errorf("failed to get member max role, cause: %w", err)
	}
	if row.IsNil() {
		return models.Role{}, errx.ErrorRoleNotFound.Raise(
			fmt.Errorf("no roles found for member with ID %s", memberID),
		)
	}

	return row.ToModel(), nil
}
