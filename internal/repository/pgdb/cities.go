package pgdb

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
	"github.com/umisto/pgx"

	sq "github.com/Masterminds/squirrel"
)

const CityTable = "cities"

const CityColumns = "id, agglomeration_id, status, slug, name, icon, banner, point, created_at, updated_at"
const CityColumnsC = "c.id, c.agglomeration_id, c.status, c.slug, c.name, c.icon, c.banner, c.point, c.created_at, c.updated_at"

type City struct {
	ID              uuid.UUID  `json:"id"`
	AgglomerationID *uuid.UUID `json:"agglomeration_id"`
	Slug            *string    `json:"slug"`
	Name            string     `json:"name"`
	Icon            *string    `json:"icon"`
	Banner          *string    `json:"banner"`
	Point           orb.Point  `json:"point"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (c *City) scan(row sq.RowScanner) error {
	err := row.Scan(
		&c.ID,
		&c.AgglomerationID,
		&c.Slug,
		&c.Name,
		&c.Icon,
		&c.Banner,
		&c.Point,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("scanning city: %w", err)
	}
	return nil
}

type CitiesQ struct {
	db       pgx.DBTX
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewCitiesQ(db pgx.DBTX) CitiesQ {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return CitiesQ{
		db:       db,
		selector: builder.Select(CityColumnsC).From(CityTable + " c"),
		inserter: builder.Insert(CityTable),
		updater:  builder.Update(CityTable + " c"),
		deleter:  builder.Delete(CityTable + " c"),
		counter:  builder.Select("COUNT(*)").From(CityTable + " c"),
	}
}

type CityInsertParams struct {
	AgglomerationID *uuid.UUID
	Slug            *string
	Name            string
	Icon            *string
	Banner          *string
	Point           orb.Point
}

func (q CitiesQ) Insert(ctx context.Context, data CityInsertParams) (City, error) {
	query, args, err := q.inserter.SetMap(map[string]interface{}{
		"agglomeration_id": data.AgglomerationID,
		"slug":             data.Slug,
		"name":             data.Name,
		"icon":             data.Icon,
		"banner":           data.Banner,
		"point":            data.Point,
	}).Suffix("RETURNING " + CityColumns).ToSql()
	if err != nil {
		return City{}, fmt.Errorf("building insert query for %s: %w", CityTable, err)
	}

	var inserted City
	err = inserted.scan(q.db.QueryRowContext(ctx, query, args...))
	if err != nil {
		return City{}, err
	}

	return inserted, nil
}
func (q CitiesQ) FilterByID(id uuid.UUID) CitiesQ {
	q.selector = q.selector.Where(sq.Eq{"c.id": id})
	q.counter = q.counter.Where(sq.Eq{"c.id": id})
	q.updater = q.updater.Where(sq.Eq{"c.id": id})
	q.deleter = q.deleter.Where(sq.Eq{"c.id": id})
	return q
}

func (q CitiesQ) FilterByAgglomerationID(id uuid.UUID) CitiesQ {
	q.selector = q.selector.Where(sq.Eq{"c.agglomeration_id": id})
	q.counter = q.counter.Where(sq.Eq{"c.agglomeration_id": id})
	q.updater = q.updater.Where(sq.Eq{"c.agglomeration_id": id})
	q.deleter = q.deleter.Where(sq.Eq{"c.agglomeration_id": id})
	return q
}

func (q CitiesQ) FilterBySlug(slug string) CitiesQ {
	q.selector = q.selector.Where(sq.Eq{"c.slug": slug})
	q.counter = q.counter.Where(sq.Eq{"c.slug": slug})
	q.updater = q.updater.Where(sq.Eq{"c.slug": slug})
	q.deleter = q.deleter.Where(sq.Eq{"c.slug": slug})
	return q
}

func (q CitiesQ) FilterLikeName(name string) CitiesQ {
	q.selector = q.selector.Where(sq.ILike{"c.name": "%" + name + "%"})
	q.counter = q.counter.Where(sq.ILike{"c.name": "%" + name + "%"})
	q.updater = q.updater.Where(sq.ILike{"c.name": "%" + name + "%"})
	q.deleter = q.deleter.Where(sq.ILike{"c.name": "%" + name + "%"})
	return q
}

func (q CitiesQ) Get(ctx context.Context) (City, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return City{}, fmt.Errorf("building select query for %s: %w", CityTable, err)
	}

	var c City
	err = c.scan(q.db.QueryRowContext(ctx, query, args...))
	if err != nil {
		return City{}, err
	}

	return c, nil
}

func (q CitiesQ) Select(ctx context.Context) ([]City, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", CityTable, err)
	}

	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing select query for %s: %w", CityTable, err)
	}
	defer rows.Close()

	var out []City
	for rows.Next() {
		var c City
		if err = c.scan(rows); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (q CitiesQ) UpdateOne(ctx context.Context) (City, error) {
	q.updater = q.updater.Set("updated_at", time.Now().UTC())

	query, args, err := q.updater.Suffix("RETURNING " + CityColumns).ToSql()
	if err != nil {
		return City{}, fmt.Errorf("building update query for %s: %w", CityTable, err)
	}

	var updated City
	if err = updated.scan(q.db.QueryRowContext(ctx, query, args...)); err != nil {
		return City{}, err
	}
	return updated, nil
}

func (q CitiesQ) UpdateMany(ctx context.Context) (int64, error) {
	q.updater = q.updater.Set("updated_at", time.Now().UTC())

	query, args, err := q.updater.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building update query for %s: %w", CityTable, err)
	}

	res, err := q.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("executing update query for %s: %w", CityTable, err)
	}

	aff, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("rows affected for %s: %w", CityTable, err)
	}

	return aff, nil
}

func (q CitiesQ) UpdateAgglomerationID(id uuid.NullUUID) CitiesQ {
	q.updater = q.updater.Set("agglomeration_id", id)
	return q
}

func (q CitiesQ) UpdateStatus(status string) CitiesQ {
	q.updater = q.updater.Set("status", status)
	return q
}

func (q CitiesQ) UpdateSlug(slug sql.NullString) CitiesQ {
	q.updater = q.updater.Set("slug", slug)
	return q
}

func (q CitiesQ) UpdateName(name string) CitiesQ {
	q.updater = q.updater.Set("name", name)
	return q
}

func (q CitiesQ) UpdateIcon(icon sql.NullString) CitiesQ {
	q.updater = q.updater.Set("icon", icon)
	return q
}

func (q CitiesQ) UpdateBanner(banner sql.NullString) CitiesQ {
	q.updater = q.updater.Set("banner", banner)
	return q
}

func (q CitiesQ) UpdatePoint(point any) CitiesQ {
	q.updater = q.updater.Set("point", point)
	return q
}

func (q CitiesQ) OrderName(asc bool) CitiesQ {
	if asc {
		q.selector = q.selector.OrderBy("c.name ASC", "c.id ASC")
	} else {
		q.selector = q.selector.OrderBy("c.name DESC", "c.id DESC")
	}
	return q
}

func (q CitiesQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", CityTable, err)
	}

	_, err = q.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("executing delete query for %s: %w", CityTable, err)
	}

	return nil
}

func (q CitiesQ) Page(limit, offset uint) CitiesQ {
	q.selector = q.selector.Limit(uint64(limit)).Offset(uint64(offset))
	return q
}

func (q CitiesQ) Count(ctx context.Context) (int64, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for %s: %w", CityTable, err)
	}

	var count int64
	err = q.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("scanning count for %s: %w", CityTable, err)
	}

	return count, nil
}
