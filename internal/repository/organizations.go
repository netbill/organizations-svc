package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/internal/core/modules/organization"
	"github.com/netbill/restkit/pagi"
)

type OrganizationRow struct {
	ID     uuid.UUID `db:"id"`
	Status string    `db:"status"`
	Name   string    `db:"name"`

	IconKey   *string `db:"icon_key,omitempty"`
	BannerKey *string `db:"banner_key,omitempty"`

	Version int32 `db:"version"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (r OrganizationRow) IsNil() bool {
	return r.ID == uuid.Nil
}

func (r OrganizationRow) ToModel() models.Organization {
	return models.Organization{
		ID:        r.ID,
		Status:    r.Status,
		Name:      r.Name,
		IconKey:   r.IconKey,
		BannerKey: r.BannerKey,
		Version:   r.Version,
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
	UpdateStatus(status string) OrganizationsQ
	UpdateIconKey(icon *string) OrganizationsQ
	UpdateBannerKey(banner *string) OrganizationsQ

	Get(ctx context.Context) (OrganizationRow, error)
	Select(ctx context.Context) ([]OrganizationRow, error)

	UpdateOne(ctx context.Context) (OrganizationRow, error)

	Count(ctx context.Context) (uint, error)
	Page(limit, offset uint) OrganizationsQ

	Delete(ctx context.Context) error
}

func (r *Repository) CreateOrganization(
	ctx context.Context,
	params organization.CreateParams,
) (models.Organization, error) {
	row, err := r.OrganizationsSql.New().Insert(ctx, OrganizationRow{
		Name: params.Name,
	})
	if err != nil {
		return models.Organization{}, fmt.Errorf(
			"failed to create organization with name %s: %w",
			params.Name, err,
		)
	}

	return row.ToModel(), nil
}

func (r *Repository) GetOrganizationByID(
	ctx context.Context,
	ID uuid.UUID,
) (models.Organization, error) {
	row, err := r.OrganizationsSql.New().FilterByID(ID).Get(ctx)
	if err != nil {
		return models.Organization{}, fmt.Errorf("failed to get organization with ID %s: %w", ID, err)
	}
	if row.IsNil() {
		return models.Organization{}, errx.ErrorOrganizationNotExists.Raise(
			fmt.Errorf("organization with ID %s not found", ID),
		)
	}

	return row.ToModel(), nil
}

func (r *Repository) GetOrganizations(
	ctx context.Context,
	filter organization.FilterParams,
	limit, offset uint,
) (pagi.Page[[]models.Organization], error) {
	if limit == 0 {
		limit = 20
	}
	if limit > 1000 {
		limit = 1000
	}

	q := r.OrganizationsSql.New()
	if filter.Text != nil {
		q = q.FilterNameLike(*filter.Text)
	}
	if filter.Status != nil {
		q = q.FilterByStatus(*filter.Status)
	}

	rows, err := q.Page(limit, offset).Select(ctx)
	if err != nil {
		return pagi.Page[[]models.Organization]{}, fmt.Errorf("failed to get organizations, cause: %w", err)
	}

	total, err := q.Count(ctx)
	if err != nil {
		return pagi.Page[[]models.Organization]{}, fmt.Errorf("failed to count organizations, cause: %w", err)
	}

	organizations := make([]models.Organization, len(rows))
	for i, row := range rows {
		organizations[i] = row.ToModel()
	}

	return pagi.Page[[]models.Organization]{
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
) (pagi.Page[[]models.Organization], error) {
	if limit == 0 {
		limit = 10
	}

	row, err := r.OrganizationsSql.New().FilterByAccountID(accountID).Page(limit, offset).Select(ctx)
	if err != nil {
		return pagi.Page[[]models.Organization]{}, fmt.Errorf(
			"failed to get organizations for account ID %s, cause: %w", accountID, err,
		)
	}

	total, err := r.OrganizationsSql.New().FilterByAccountID(accountID).Count(ctx)
	if err != nil {
		return pagi.Page[[]models.Organization]{}, fmt.Errorf(
			"failed to count organizations for account ID %s, cause: %w", accountID, err,
		)
	}

	organizations := make([]models.Organization, len(row))
	for i, org := range row {
		organizations[i] = org.ToModel()
	}

	return pagi.Page[[]models.Organization]{
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
) (models.Organization, error) {
	q := r.OrganizationsSql.New().
		FilterByID(ID).
		UpdateName(params.Name).
		UpdateIconKey(params.IconKey).
		UpdateBannerKey(params.BannerKey)

	row, err := q.UpdateOne(ctx)
	if err != nil {
		return models.Organization{}, fmt.Errorf("failed to update organization with ID %s: %w", ID, err)
	}
	if row.IsNil() {
		return models.Organization{}, errx.ErrorOrganizationNotExists.Raise(
			fmt.Errorf("organization with ID %s not found", ID),
		)
	}

	return row.ToModel(), nil
}

func (r *Repository) UpdateOrganizationStatus(
	ctx context.Context,
	ID uuid.UUID,
	status string,
) (models.Organization, error) {
	row, err := r.OrganizationsSql.New().
		FilterByID(ID).
		UpdateStatus(status).
		UpdateOne(ctx)
	if err != nil {
		return models.Organization{}, fmt.Errorf("failed to update organization status with ID %s: %w", ID, err)
	}
	if row.IsNil() {
		return models.Organization{}, errx.ErrorOrganizationNotExists.Raise(
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
