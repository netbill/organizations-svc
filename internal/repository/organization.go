package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/organization"
	"github.com/netbill/organizations-svc/internal/errx"
	"github.com/netbill/organizations-svc/internal/models"
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

type OrganizationRepo struct {
	query OrganizationsQ
}

func NewOrganizationRepo(organizationsSql OrganizationsQ) *OrganizationRepo {
	return &OrganizationRepo{
		query: organizationsSql,
	}
}

func (r *OrganizationRepo) Create(
	ctx context.Context,
	params organization.CreateParams,
) (models.Organization, error) {
	row, err := r.query.New().Insert(ctx, OrganizationRow{
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

func (r *OrganizationRepo) Get(
	ctx context.Context,
	ID uuid.UUID,
) (models.Organization, error) {
	row, err := r.query.New().FilterByID(ID).Get(ctx)
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

func (r *OrganizationRepo) GetListByIds(
	ctx context.Context,
	IDs []uuid.UUID,
) ([]models.Organization, error) {
	rows, err := r.query.New().FilterByID(IDs...).Select(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get organizations by IDs: %w", err)
	}

	organizations := make([]models.Organization, len(rows))
	for i, row := range rows {
		organizations[i] = row.ToModel()
	}

	return organizations, nil
}

func (r *OrganizationRepo) GetList(
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

	q := r.query.New()
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

func (r *OrganizationRepo) GetForAccountAndOrg(
	ctx context.Context,
	accountID uuid.UUID,
	limit, offset uint,
) (pagi.Page[[]models.Organization], error) {
	if limit == 0 {
		limit = 10
	}

	row, err := r.query.New().FilterByAccountID(accountID).Page(limit, offset).Select(ctx)
	if err != nil {
		return pagi.Page[[]models.Organization]{}, fmt.Errorf(
			"failed to get organizations for account ID %s, cause: %w", accountID, err,
		)
	}

	total, err := r.query.New().FilterByAccountID(accountID).Count(ctx)
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

func (r *OrganizationRepo) Update(
	ctx context.Context,
	ID uuid.UUID,
	params organization.UpdateParams,
) (models.Organization, error) {
	q := r.query.New().FilterByID(ID)

	if params.Name != nil {
		q = q.UpdateName(*params.Name)
	}
	if params.IconKey != nil {
		q = q.UpdateIconKey(params.IconKey)
	}
	if params.BannerKey != nil {
		q = q.UpdateBannerKey(params.BannerKey)
	}

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

func (r *OrganizationRepo) UpdateStatus(
	ctx context.Context,
	ID uuid.UUID,
	status string,
) (models.Organization, error) {
	row, err := r.query.New().
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

func (r *OrganizationRepo) Delete(ctx context.Context, ID uuid.UUID) error {
	err := r.query.New().FilterByID(ID).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete organization with ID %s, cause: %w", ID, err)
	}

	return nil
}
