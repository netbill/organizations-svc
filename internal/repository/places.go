package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/internal/core/modules/place"
)

type PlaceRow struct {
	ID               uuid.UUID `db:"id"`
	OrganizationID   uuid.UUID `db:"organization_id"`
	SourceCreatedAt  time.Time `db:"source_created_at"`
	ReplicaCreatedAt time.Time `db:"replica_created_at"`
}

func (r PlaceRow) IsNil() bool {
	return r.ID == uuid.Nil
}

func (r PlaceRow) ToModel() models.Place {
	return models.Place{
		ID:             r.ID,
		OrganizationID: r.OrganizationID,
		CreatedAt:      r.SourceCreatedAt,
	}
}

type PlacesQ interface {
	New() PlacesQ
	Insert(ctx context.Context, input PlaceRow) error

	Get(ctx context.Context) (PlaceRow, error)
	Select(ctx context.Context) ([]PlaceRow, error)
	Exists(ctx context.Context) (bool, error)

	FilterByID(id ...uuid.UUID) PlacesQ
	FilterByOrganizationID(organizationID ...uuid.UUID) PlacesQ

	Delete(ctx context.Context) error
}

func (r *Repository) CreatePlace(ctx context.Context, params place.CreateParams) error {
	return r.PlacesSql.New().Insert(ctx, PlaceRow{
		ID:              params.ID,
		OrganizationID:  params.OrganizationID,
		SourceCreatedAt: params.CreatedAt,
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
		return models.Place{}, errx.ErrorPlaceNotExists.Raise(
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

func (r *Repository) PlaceExists(ctx context.Context, id uuid.UUID) (bool, error) {
	return r.PlacesSql.New().FilterByID(id).Exists(ctx)
}

func (r *Repository) DeletePlaceByID(ctx context.Context, id uuid.UUID) error {
	return r.PlacesSql.New().FilterByID(id).Delete(ctx)
}
