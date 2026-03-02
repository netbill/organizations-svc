package pg

import (
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/netbill/organizations-svc/internal/repository"
	"github.com/netbill/pgdbx"
)

const placesTable = "places"
const placesColumns = `id, organization_id, source_created_at, replica_created_at`
const placesColumnsP = `p.id, p.organization_id, p.source_created_at, p.replica_created_at`

func scanPlace(row sq.RowScanner) (p repository.PlaceRow, err error) {
	err = row.Scan(
		&p.ID,
		&p.OrganizationID,
		&p.SourceCreatedAt,
		&p.ReplicaCreatedAt,
	)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return repository.PlaceRow{}, nil
	case err != nil:
		return repository.PlaceRow{}, fmt.Errorf("scanning place: %w", err)
	}
	return p, nil
}

type places struct {
	db       *pgdbx.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewPlacesQ(db *pgdbx.DB) repository.PlacesQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &places{
		db:       db,
		selector: b.Select(placesColumnsP).From(placesTable + " p"),
		inserter: b.Insert(placesTable),
		updater:  b.Update(placesTable + " p"),
		deleter:  b.Delete(placesTable + " p"),
		counter:  b.Select("COUNT(*)").From(placesTable + " p"),
	}
}

func (q *places) New() repository.PlacesQ {
	return NewPlacesQ(q.db)
}

func (q *places) Insert(ctx context.Context, data repository.PlaceRow) error {
	query, args, err := q.inserter.SetMap(map[string]any{
		"id":                data.ID,
		"organization_id":   data.OrganizationID,
		"source_created_at": data.SourceCreatedAt.UTC(),
	}).ToSql()
	if err != nil {
		return fmt.Errorf("building insert query for %s: %w", placesTable, err)
	}

	if _, err = q.db.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("executing insert query for %s: %w", placesTable, err)
	}

	return nil
}

func (q *places) Get(ctx context.Context) (repository.PlaceRow, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return repository.PlaceRow{}, fmt.Errorf("building select query for %s: %w", placesTable, err)
	}

	return scanPlace(q.db.QueryRow(ctx, query, args...))
}

func (q *places) Select(ctx context.Context) ([]repository.PlaceRow, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", placesTable, err)
	}

	rows, err := q.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing select query for %s: %w", placesTable, err)
	}
	defer rows.Close()

	out := make([]repository.PlaceRow, 0)
	for rows.Next() {
		p, err := scanPlace(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (q *places) Exists(ctx context.Context) (bool, error) {
	subSQL, subArgs, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return false, fmt.Errorf("building exists query for %s: %w", placesTable, err)
	}

	var ok bool
	if err = q.db.QueryRow(ctx, "SELECT EXISTS ("+subSQL+")", subArgs...).Scan(&ok); err != nil {
		return false, fmt.Errorf("scanning exists for %s: %w", placesTable, err)
	}

	return ok, nil
}

func (q *places) UpdateOne(ctx context.Context) error {
	q.updater = q.updater.Set("replica_updated_at", time.Now().UTC())

	query, args, err := q.updater.ToSql()
	if err != nil {
		return fmt.Errorf("building update query for %s: %w", placesTable, err)
	}

	if _, err = q.db.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("executing update query for %s: %w", placesTable, err)
	}

	return nil
}

func (q *places) FilterByID(id ...uuid.UUID) repository.PlacesQ {
	q.selector = q.selector.Where(sq.Eq{"p.id": id})
	q.counter = q.counter.Where(sq.Eq{"p.id": id})
	q.updater = q.updater.Where(sq.Eq{"p.id": id})
	q.deleter = q.deleter.Where(sq.Eq{"p.id": id})
	return q
}

func (q *places) FilterByOrganizationID(organizationID ...uuid.UUID) repository.PlacesQ {
	q.selector = q.selector.Where(sq.Eq{"p.organization_id": organizationID})
	q.counter = q.counter.Where(sq.Eq{"p.organization_id": organizationID})
	q.updater = q.updater.Where(sq.Eq{"p.organization_id": organizationID})
	q.deleter = q.deleter.Where(sq.Eq{"p.organization_id": organizationID})
	return q
}

func (q *places) UpdateClassID(classID uuid.UUID) repository.PlacesQ {
	q.updater = q.updater.Set("class_id", classID)
	return q
}

func (q *places) UpdateName(name string) repository.PlacesQ {
	q.updater = q.updater.Set("name", name)
	return q
}

func (q *places) UpdateAddress(address string) repository.PlacesQ {
	q.updater = q.updater.Set("address", address)
	return q
}

func (q *places) UpdateStatus(status string) repository.PlacesQ {
	q.updater = q.updater.Set("status", status)
	return q
}

func (q *places) UpdateVerified(verified bool) repository.PlacesQ {
	q.updater = q.updater.Set("verified", verified)
	return q
}

func (q *places) UpdateDescription(description *string) repository.PlacesQ {
	q.updater = q.updater.Set("description", description)
	return q
}

func (q *places) UpdateIconKey(icon *string) repository.PlacesQ {
	q.updater = q.updater.Set("icon_key", icon)
	return q
}

func (q *places) UpdateBannerKey(banner *string) repository.PlacesQ {
	q.updater = q.updater.Set("banner_key", banner)
	return q
}

func (q *places) UpdateWebsite(website *string) repository.PlacesQ {
	q.updater = q.updater.Set("website", website)
	return q
}

func (q *places) UpdatePhone(phone *string) repository.PlacesQ {
	q.updater = q.updater.Set("phone", phone)
	return q
}

func (q *places) UpdateVersion(v int32) repository.PlacesQ {
	q.updater = q.updater.Set("version", v)
	return q
}

func (q *places) UpdateSourceUpdatedAt(v time.Time) repository.PlacesQ {
	q.updater = q.updater.Set("source_updated_at", v.UTC())
	return q
}

func (q *places) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", placesTable, err)
	}

	if _, err = q.db.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("executing delete query for %s: %w", placesTable, err)
	}

	return nil
}
