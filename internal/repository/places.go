package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/internal/core/modules/place"
	"github.com/paulmach/orb"
)

type PlaceRow struct {
	ID             uuid.UUID `db:"id"`
	ClassID        uuid.UUID `db:"class_id"`
	OrganizationID uuid.UUID `db:"organization_id"`

	Status   string    `db:"status"`
	Verified bool      `db:"verified"`
	Point    orb.Point `db:"point"`
	Address  string    `db:"address"`
	Name     string    `db:"name"`

	Description *string `db:"description"`
	IconKey     *string `db:"icon_key"`
	BannerKey   *string `db:"banner_key"`
	Website     *string `db:"website"`
	Phone       *string `db:"phone"`

	Version          int32     `db:"version"`
	SourceCreatedAt  time.Time `db:"source_created_at"`
	SourceUpdatedAt  time.Time `db:"source_updated_at"`
	ReplicaCreatedAt time.Time `db:"replica_created_at"`
	ReplicaUpdatedAt time.Time `db:"replica_updated_at"`
}

func (r PlaceRow) IsNil() bool {
	return r.ID == uuid.Nil
}

func (r PlaceRow) ToModel() models.Place {
	return models.Place{
		ID:             r.ID,
		ClassID:        r.ClassID,
		OrganizationID: r.OrganizationID,
		Status:         r.Status,
		Verified:       r.Verified,
		Point:          r.Point,
		Address:        r.Address,
		Name:           r.Name,
		Description:    r.Description,
		IconKey:        r.IconKey,
		BannerKey:      r.BannerKey,
		Website:        r.Website,
		Phone:          r.Phone,
		Version:        r.Version,
		CreatedAt:      r.SourceCreatedAt,
		UpdatedAt:      r.SourceUpdatedAt,
	}
}

type PlacesQ interface {
	New() PlacesQ
	Insert(ctx context.Context, input PlaceRow) error

	Get(ctx context.Context) (PlaceRow, error)
	Select(ctx context.Context) ([]PlaceRow, error)
	Exists(ctx context.Context) (bool, error)

	UpdateOne(ctx context.Context) error

	UpdateClassID(classID uuid.UUID) PlacesQ
	UpdateName(name string) PlacesQ
	UpdateAddress(address string) PlacesQ
	UpdateStatus(status string) PlacesQ
	UpdateVerified(verified bool) PlacesQ
	UpdateDescription(description *string) PlacesQ
	UpdateIconKey(icon *string) PlacesQ
	UpdateBannerKey(banner *string) PlacesQ
	UpdateWebsite(website *string) PlacesQ
	UpdatePhone(phone *string) PlacesQ
	UpdateVersion(v int32) PlacesQ
	UpdateSourceUpdatedAt(v time.Time) PlacesQ

	FilterByID(id ...uuid.UUID) PlacesQ
	FilterByOrganizationID(organizationID ...uuid.UUID) PlacesQ

	Delete(ctx context.Context) error
}

func (r *Repository) CreatePlace(ctx context.Context, params place.CreateParams) error {
	return r.PlacesSql.New().Insert(ctx, PlaceRow{
		ID:              params.ID,
		ClassID:         params.ClassID,
		OrganizationID:  params.OrganizationID,
		Status:          params.Status,
		Verified:        params.Verified,
		Point:           params.Point,
		Address:         params.Address,
		Name:            params.Name,
		Description:     params.Description,
		IconKey:         params.IconKey,
		BannerKey:       params.BannerKey,
		Website:         params.Website,
		Phone:           params.Phone,
		Version:         1,
		SourceCreatedAt: params.CreatedAt,
		SourceUpdatedAt: params.CreatedAt,
	})
}

func (r *Repository) GetPlaceExistsForOrganization(ctx context.Context, organizationID uuid.UUID) (bool, error) {
	return r.PlacesSql.New().FilterByOrganizationID(organizationID).Exists(ctx)
}

func (r *Repository) GetPlaceByID(ctx context.Context, id uuid.UUID) (models.Place, error) {
	row, err := r.PlacesSql.New().FilterByID(id).Get(ctx)
	if err != nil {
		return models.Place{}, err
	}
	if row.IsNil() {
		return models.Place{}, errx.ErrorPlaceNotFound.Raise(
			fmt.Errorf("place with id %s not found", id),
		)
	}

	return row.ToModel(), nil
}

func (r *Repository) GetPlacesByIDs(ctx context.Context, ids []uuid.UUID) ([]models.Place, error) {
	rows, err := r.PlacesSql.New().FilterByID(ids...).Select(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]models.Place, 0, len(rows))
	for _, row := range rows {
		res = append(res, row.ToModel())
	}

	return res, nil
}

func (r *Repository) UpdatePlaceByID(ctx context.Context, id uuid.UUID, params place.UpdateParams) error {
	return r.PlacesSql.New().
		FilterByID(id).
		UpdateClassID(params.ClassID).
		UpdateName(params.Name).
		UpdateAddress(params.Address).
		UpdateStatus(params.Status).
		UpdateVerified(params.Verified).
		UpdateDescription(params.Description).
		UpdateIconKey(params.IconKey).
		UpdateBannerKey(params.BannerKey).
		UpdateWebsite(params.Website).
		UpdatePhone(params.Phone).
		UpdateVersion(params.Version).
		UpdateSourceUpdatedAt(params.UpdatedAt).
		UpdateOne(ctx)
}

func (r *Repository) DeletePlaceByID(ctx context.Context, id uuid.UUID) error {
	return r.PlacesSql.New().FilterByID(id).Delete(ctx)
}
