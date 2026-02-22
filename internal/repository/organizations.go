package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/domain"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/modules/organization"
	"github.com/netbill/restkit/pagi"
)

type OrganizationRow struct {
	ID     uuid.UUID `db:"id"`
	Status string    `db:"status"`
	Name   string    `db:"name"`

	IconKey   *string `db:"icon_key,omitempty"`
	BannerKey *string `db:"banner_key,omitempty"`

	MaxRoles uint `db:"max_roles"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (r OrganizationRow) IsNil() bool {
	return r.ID == uuid.Nil
}

func (r OrganizationRow) ToModel() domain.Organization {
	return domain.Organization{
		ID:        r.ID,
		Status:    r.Status,
		Name:      r.Name,
		IconKey:   r.IconKey,
		BannerKey: r.BannerKey,
		MaxRoles:  r.MaxRoles,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

type OrganizationsQ interface {
	New() OrganizationsQ
	Insert(ctx context.Context, input OrganizationRow) (OrganizationRow, error)

	FilterByID(id ...uuid.UUID) OrganizationsQ
	FilterNameLike(name string) OrganizationsQ
	FilterByStatus(status string) OrganizationsQ
	FilterByAccountID(accountID uuid.UUID) OrganizationsQ

	UpdateName(name string) OrganizationsQ
	UpdateIconKey(icon *string) OrganizationsQ
	UpdateBannerKey(banner *string) OrganizationsQ
	UpdateStatus(status string) OrganizationsQ
	UpdateMaxRoles(maxRoles uint) OrganizationsQ

	Get(ctx context.Context) (OrganizationRow, error)
	Select(ctx context.Context) ([]OrganizationRow, error)

	UpdateOne(ctx context.Context) (OrganizationRow, error)
	UpdateMany(ctx context.Context) (int64, error)

	Count(ctx context.Context) (uint, error)
	Page(limit, offset uint) OrganizationsQ

	Delete(ctx context.Context) error
}

func (r *Repository) CreateOrganization(
	ctx context.Context,
	params organization.CreateParams,
) (domain.Organization, error) {
	row, err := r.OrganizationsSql.New().Insert(ctx, OrganizationRow{
		Name: params.Name,
	})
	if err != nil {
		return domain.Organization{}, fmt.Errorf(
			"failed to create organization with name %s: %w",
			params.Name, err,
		)
	}

	return row.ToModel(), nil
}

func (r *Repository) GetOrganizationByID(
	ctx context.Context,
	ID uuid.UUID,
) (domain.Organization, error) {
	row, err := r.OrganizationsSql.New().FilterByID(ID).Get(ctx)
	if err != nil {
		return domain.Organization{}, fmt.Errorf("failed to get organization with ID %s: %w", ID, err)
	}
	if row.IsNil() {
		return domain.Organization{}, errx.ErrorOrganizationNotFound.Raise(
			fmt.Errorf("organization with ID %s not found", ID),
		)
	}

	return row.ToModel(), nil
}

func (r *Repository) GetOrganizations(
	ctx context.Context,
	filter organization.FilterParams,
	limit, offset uint,
) (pagi.Page[[]domain.Organization], error) {
	if limit == 0 {
		limit = 10
	}

	q := r.OrganizationsSql.New()
	if filter.Name != nil {
		q = q.FilterNameLike(*filter.Name)
	}
	if filter.Status != nil {
		q = q.FilterByStatus(*filter.Status)
	}

	rows, err := q.Page(limit, offset).Select(ctx)
	if err != nil {
		return pagi.Page[[]domain.Organization]{}, fmt.Errorf("failed to get organizations, cause: %w", err)
	}

	total, err := q.Count(ctx)
	if err != nil {
		return pagi.Page[[]domain.Organization]{}, fmt.Errorf("failed to count organizations, cause: %w", err)
	}

	organizations := make([]domain.Organization, len(rows))
	for i, row := range rows {
		organizations[i] = row.ToModel()
	}

	return pagi.Page[[]domain.Organization]{
		Data:  organizations,
		Page:  uint(offset/limit) + 1,
		Size:  uint(len(organizations)),
		Total: total,
	}, nil

}

func (r *Repository) GetOrganizationsForUser(
	ctx context.Context,
	accountID uuid.UUID,
	limit, offset uint,
) (pagi.Page[[]domain.Organization], error) {
	if limit == 0 {
		limit = 10
	}

	row, err := r.OrganizationsSql.New().FilterByAccountID(accountID).Page(limit, offset).Select(ctx)
	if err != nil {
		return pagi.Page[[]domain.Organization]{}, fmt.Errorf(
			"failed to get organizations for account ID %s, cause: %w", accountID, err,
		)
	}

	total, err := r.OrganizationsSql.New().FilterByAccountID(accountID).Count(ctx)
	if err != nil {
		return pagi.Page[[]domain.Organization]{}, fmt.Errorf(
			"failed to count organizations for account ID %s, cause: %w", accountID, err,
		)
	}

	organizations := make([]domain.Organization, len(row))
	for i, org := range row {
		organizations[i] = org.ToModel()
	}

	return pagi.Page[[]domain.Organization]{
		Data:  organizations,
		Page:  uint(offset/limit) + 1,
		Size:  uint(len(organizations)),
		Total: total,
	}, nil
}

func (r *Repository) UpdateOrganization(
	ctx context.Context,
	ID uuid.UUID,
	params organization.UpdateParams,
) (domain.Organization, error) {
	q := r.OrganizationsSql.New().
		FilterByID(ID).
		UpdateName(params.Name).
		UpdateIconKey(params.IconKey).
		UpdateBannerKey(params.BannerKey)

	row, err := q.UpdateOne(ctx)
	if err != nil {
		return domain.Organization{}, fmt.Errorf("failed to update organization with ID %s: %w", ID, err)
	}
	if row.IsNil() {
		return domain.Organization{}, errx.ErrorOrganizationNotFound.Raise(
			fmt.Errorf("organization with ID %s not found", ID),
		)
	}

	return row.ToModel(), nil
}

func (r *Repository) UpdateOrganizationStatus(
	ctx context.Context,
	ID uuid.UUID,
	status string,
) (domain.Organization, error) {
	row, err := r.OrganizationsSql.New().
		FilterByID(ID).
		UpdateStatus(status).
		UpdateOne(ctx)
	if err != nil {
		return domain.Organization{}, fmt.Errorf("failed to update organization status with ID %s: %w", ID, err)
	}
	if row.IsNil() {
		return domain.Organization{}, errx.ErrorOrganizationNotFound.Raise(
			fmt.Errorf("organization with ID %s not found", ID),
		)
	}

	return row.ToModel(), nil
}

func (r *Repository) UpdateOrganizationMaxRoles(
	ctx context.Context,
	ID uuid.UUID,
	maxRoles uint,
) (domain.Organization, error) {
	row, err := r.OrganizationsSql.New().
		FilterByID(ID).
		UpdateMaxRoles(maxRoles).
		UpdateOne(ctx)
	if err != nil {
		return domain.Organization{}, fmt.Errorf("failed to update organization max roles with ID %s: %w", ID, err)
	}
	if row.IsNil() {
		return domain.Organization{}, errx.ErrorOrganizationNotFound.Raise(
			fmt.Errorf("organization with ID %s not found", ID),
		)
	}

	return row.ToModel(), nil
}

func (r *Repository) DeleteOrganization(ctx context.Context, ID uuid.UUID) error {
	err := r.OrganizationsSql.New().FilterByID(ID).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete organization with ID %s, cause: %w", ID, err)
	}

	return nil
}
