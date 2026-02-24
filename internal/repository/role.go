package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/internal/core/modules/role"
	"github.com/netbill/restkit/pagi"
)

type OrgRoleRow struct {
	ID             uuid.UUID `db:"id"`
	OrganizationID uuid.UUID `db:"organization_id"`

	Name        string `db:"name"`
	Description string `db:"description,omitempty"`
	Color       string `db:"color"`
	Version     int32  `db:"version"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (r OrgRoleRow) IsNil() bool {
	return r.ID == uuid.Nil
}

func (r OrgRoleRow) ToModel() models.Role {
	return models.Role{
		ID:             r.ID,
		OrganizationID: r.OrganizationID,
		Name:           r.Name,
		Description:    r.Description,
		Color:          r.Color,
		Version:        r.Version,
		CreatedAt:      r.CreatedAt,
		UpdatedAt:      r.UpdatedAt,
	}
}

type OrgRolesQ interface {
	New() OrgRolesQ

	Insert(ctx context.Context, input OrgRoleRow) (OrgRoleRow, error)

	Get(ctx context.Context) (OrgRoleRow, error)
	Select(ctx context.Context) ([]OrgRoleRow, error)
	Exists(ctx context.Context) (bool, error)

	UpdateOne(ctx context.Context) (OrgRoleRow, error)

	UpdateName(name string) OrgRolesQ
	UpdateDescription(description string) OrgRolesQ
	UpdateColor(color string) OrgRolesQ
	
	FilterByID(roleID uuid.UUID) OrgRolesQ
	FilterByOrganizationID(organizationID uuid.UUID) OrgRolesQ

	OrderByRank(ask bool) OrgRolesQ

	Delete(ctx context.Context) error
}

func (r *Repository) CreateRole(ctx context.Context, params role.CreateParams) (models.Role, error) {
	//TODO implement me
	panic("implement me")
}

func (r *Repository) GetRole(ctx context.Context, roleID uuid.UUID) (models.Role, error) {
	//TODO implement me
	panic("implement me")
}

func (r *Repository) GetRoles(
	ctx context.Context,
	filter role.FilterParams,
	limit, offset uint,
) (pagi.Page[[]models.Role], error) {
	//TODO implement me
	panic("implement me")
}

func (r *Repository) UpdateRole(
	ctx context.Context,
	roleID uuid.UUID,
	params role.UpdateParams,
) (models.Role, error) {
	//TODO implement me
	panic("implement me")
}

func (r *Repository) DeleteRole(ctx context.Context, roleID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}
