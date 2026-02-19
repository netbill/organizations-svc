package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/internal/core/modules/role"
	"github.com/netbill/restkit/pagi"
)

type OrganizationRoleRow struct {
	ID             uuid.UUID `db:"id"`
	OrganizationID uuid.UUID `db:"organization_id"`
	Rank           uint      `db:"rank"`
	Name           string    `db:"name"`
	Description    string    `db:"description"`
	Color          string    `db:"color"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (r OrganizationRoleRow) IsNil() bool {
	return r.ID == uuid.Nil
}

func (r OrganizationRoleRow) ToModel() models.Role {
	return models.Role{
		ID:             r.ID,
		OrganizationID: r.OrganizationID,
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

type OrganizationRolePermissionRow struct {
	ID           uuid.UUID  `db:"id"`
	Code         string     `db:"code"`
	Description  string     `db:"description"`
	DeprecatedAt *time.Time `db:"deprecated_at"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
}

func (r OrganizationRolePermissionRow) IsNil() bool {
	return r.ID == uuid.Nil
}

type OrgRolePermissionsQ interface {
	New() OrgRolePermissionsQ

	Insert(ctx context.Context, input OrganizationRolePermissionRow) (OrganizationRolePermissionRow, error)

	FilterByID(ids ...uuid.UUID) OrgRolePermissionsQ
	FilterByCode(codes ...string) OrgRolePermissionsQ
	FilterByDeprecated(deprecated bool) OrgRolePermissionsQ

	UpdateOne(ctx context.Context) (OrganizationRolePermissionRow, error)
	UpdateMany(ctx context.Context) (int64, error)

	UpdateDeprecatedAt(timestamp *time.Time) OrgRolePermissionsQ

	Select(ctx context.Context) ([]OrganizationRolePermissionRow, error)
	Get(ctx context.Context) (OrganizationRolePermissionRow, error)

	Delete(ctx context.Context) error
	Count(ctx context.Context) (uint, error)
	Page(limit, offset uint) OrgRolePermissionsQ
}

type OrganizationRolePermissionLinkRow struct {
	RoleID       uuid.UUID `db:"role_id"`
	PermissionID uuid.UUID `db:"permission_id"`
}

func (r OrganizationRolePermissionLinkRow) IsNil() bool {
	return r.RoleID == uuid.Nil
}

type OrgRolePermissionLinksQ interface {
	New() OrgRolePermissionLinksQ

	Upsert(
		ctx context.Context,
		roleID uuid.UUID,
		permissions map[uuid.UUID]bool,
	) ([]OrganizationRolePermissionLinkRow, error)

	Select(ctx context.Context) ([]OrganizationRolePermissionLinkRow, error)
	Get(ctx context.Context) (OrganizationRolePermissionLinkRow, error)

	FilterByAccountID(accountID uuid.UUID) OrgRolePermissionLinksQ
	FilterByOrganizationID(organizationID uuid.UUID) OrgRolePermissionLinksQ
	FilterByMemberID(memberID uuid.UUID) OrgRolePermissionLinksQ
	FilterByRoleID(roleID uuid.UUID) OrgRolePermissionLinksQ
	FilterByPermissionID(permissionIDs ...uuid.UUID) OrgRolePermissionLinksQ

	Delete(ctx context.Context) error

	Count(ctx context.Context) (uint, error)
	Page(limit, offset uint) OrgRolePermissionLinksQ
	Exists(ctx context.Context) (bool, error)
}

func (r *Repository) CreateRole(ctx context.Context, params role.CreateParams) (models.Role, error) {
	row, err := r.OrgRolesSql.New().Insert(ctx, OrganizationRoleRow{
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

func (r *Repository) GetRole(ctx context.Context, roleID uuid.UUID) (models.Role, error) {
	row, err := r.OrgRolesSql.New().FilterByID(roleID).Get(ctx)
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

func (r *Repository) GetRoles(
	ctx context.Context,
	filter role.FilterParams,
	limit, offset uint,
) (pagi.Page[[]models.Role], error) {
	q := r.OrgRolesSql.New()
	if filter.OrganizationID != nil {
		q = q.FilterByOrganizationID(*filter.OrganizationID)
	}
	if filter.RolesID != nil && len(*filter.RolesID) > 0 {
		q = q.FilterByID(*filter.RolesID...)
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

func (r *Repository) UpdateRole(
	ctx context.Context,
	roleID uuid.UUID,
	params role.UpdateParams,
) (models.Role, error) {
	q := r.OrgRolesSql.New().FilterByID(roleID)
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

func (r *Repository) UpdateRoleRank(
	ctx context.Context,
	roleID uuid.UUID,
	newRank uint,
) (models.Role, error) {
	row, err := r.OrgRolesSql.New().UpdateRoleRank(ctx, roleID, newRank)
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

func (r *Repository) UpdateRolesRanks(
	ctx context.Context,
	organizationID uuid.UUID,
	order map[uuid.UUID]uint,
) error {
	_, err := r.OrgRolesSql.New().UpdateRolesRanks(ctx, organizationID, order)
	if err != nil {
		return fmt.Errorf("failed to update roles ranks, cause: %w", err)
	}

	return nil
}

func (r *Repository) DeleteRole(ctx context.Context, roleID uuid.UUID) error {
	err := r.OrgRolesSql.New().DeleteAndShiftRanks(ctx, roleID)
	if err != nil {
		return fmt.Errorf("failed to delete role, cause: %w", err)
	}

	return nil
}

func (r *Repository) GetRolePermissions(
	ctx context.Context,
	roleID uuid.UUID,
) (models.OrgRolePermissionsWithDetailsForRole, error) {
	dict, err := r.OrgRolePermissionsSql.New().Select(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions dict: %w", err)
	}

	links, err := r.OrgRolePermissionLinksSql.New().
		FilterByRoleID(roleID).
		Select(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get role permission links: %w", err)
	}

	enabled := make(map[uuid.UUID]struct{}, len(links))
	for i := range links {
		enabled[links[i].PermissionID] = struct{}{}
	}

	out := make(models.OrgRolePermissionsWithDetailsForRole, len(dict))
	for i := range dict {
		p := dict[i]

		_, ok := enabled[p.ID]
		out[p.ID] = models.OrgRolePermissionDetails{
			Code:        p.Code,
			Description: p.Description,
			Enabled:     ok,
		}
	}

	return out, nil
}

func (r *Repository) GetAllPermissions(ctx context.Context) ([]models.OrgRolePermission, error) {
	permissions, err := r.OrgRolePermissionsSql.New().Select(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all permissions, cause: %w", err)
	}

	result := make([]models.OrgRolePermission, len(permissions))
	for i, perm := range permissions {
		result[i] = models.OrgRolePermission{
			ID:          perm.ID,
			Code:        perm.Code,
			Description: perm.Description,
		}
	}

	return result, nil
}

func (r *Repository) SetRolePermissions(
	ctx context.Context,
	roleID uuid.UUID,
	params role.SetPermissions,
) (models.OrgRolePermissionsWithDetailsForRole, error) {
	dict, err := r.OrgRolePermissionsSql.New().FilterByDeprecated(false).Select(ctx)
	if err != nil {
		return nil, fmt.Errorf("select params dict: %w", err)
	}

	rows, err := r.OrgRolePermissionLinksSql.New().Upsert(ctx, roleID, params)
	if err != nil {
		return nil, fmt.Errorf("set role params: %w", err)
	}

	enabled := make(map[uuid.UUID]struct{}, len(rows))
	for i := range rows {
		enabled[rows[i].PermissionID] = struct{}{}
	}

	out := make(models.OrgRolePermissionsWithDetailsForRole, len(dict))
	for i := range dict {
		p := dict[i]
		_, ok := enabled[p.ID]
		out[p.ID] = models.OrgRolePermissionDetails{
			Code:        p.Code,
			Description: p.Description,
			Enabled:     ok,
		}
	}

	return out, nil
}
