package pgdb

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/umisto/pgx"

	sq "github.com/Masterminds/squirrel"
)

const AgglomerationTable = "agglomerations"
const AgglomerationColumns = "id, status, name, icon, created_at, updated_at"

type Agglomeration struct {
	ID       uuid.UUID `json:"id"`
	Status   string    `json:"status"`
	Name     string    `json:"name"`
	Icon     string    `json:"icon"`
	MaxRoles uint      `json:"max_roles"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (a *Agglomeration) scan(row sq.RowScanner) error {
	err := row.Scan(
		&a.ID,
		&a.Status,
		&a.Name,
		&a.Icon,
		&a.MaxRoles,
		&a.CreatedAt,
		&a.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("scanning agglomeration: %w", err)
	}
	return nil
}

type AgglomerationsQ struct {
	db       pgx.DBTX
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewAgglomerationsQ(db pgx.DBTX) AgglomerationsQ {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return AgglomerationsQ{
		db:       db,
		selector: builder.Select(AgglomerationColumns).From(AgglomerationTable),
		inserter: builder.Insert(AgglomerationTable),
		updater:  builder.Update(AgglomerationTable),
		deleter:  builder.Delete(AgglomerationTable),
		counter:  builder.Select("COUNT(*) AS count").From(AgglomerationTable),
	}
}

type AgglomerationsQInsertInput struct {
	Name string
	Icon *string
}

func (q AgglomerationsQ) Insert(ctx context.Context, data AgglomerationsQInsertInput) (Agglomeration, error) {
	query, args, err := q.inserter.SetMap(map[string]interface{}{
		"name": data.Name,
		"icon": data.Icon,
	}).Suffix("RETURNING " + AgglomerationColumns).ToSql()
	if err != nil {
		return Agglomeration{}, fmt.Errorf("building insert query for %s: %w", AgglomerationTable, err)
	}

	var inserted Agglomeration
	err = inserted.scan(q.db.QueryRowContext(ctx, query, args...))
	if err != nil {
		return Agglomeration{}, err
	}

	return inserted, nil
}

func (q AgglomerationsQ) FilterByID(id uuid.UUID) AgglomerationsQ {
	q.selector = q.selector.Where(sq.Eq{"id": id})
	q.counter = q.counter.Where(sq.Eq{"id": id})
	q.updater = q.updater.Where(sq.Eq{"id": id})
	q.deleter = q.deleter.Where(sq.Eq{"id": id})
	return q
}

func (q AgglomerationsQ) FilterByStatus(status string) AgglomerationsQ {
	q.selector = q.selector.Where(sq.Eq{"status": status})
	q.counter = q.counter.Where(sq.Eq{"status": status})
	q.updater = q.updater.Where(sq.Eq{"status": status})
	q.deleter = q.deleter.Where(sq.Eq{"status": status})
	return q
}

func (q AgglomerationsQ) FilterNameLike(name string) AgglomerationsQ {
	q.selector = q.selector.Where(sq.Like{"name": "%" + name + "%"})
	q.counter = q.counter.Where(sq.Like{"name": "%" + name + "%"})
	return q
}

func (q AgglomerationsQ) OrderName(asc bool) AgglomerationsQ {
	if asc {
		q.selector = q.selector.OrderBy("name ASC", "id ASC")
	} else {
		q.selector = q.selector.OrderBy("name DESC", "id DESC")
	}
	return q
}

func (q AgglomerationsQ) Get(ctx context.Context) (Agglomeration, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return Agglomeration{}, fmt.Errorf("building select query for %s: %w", AgglomerationTable, err)
	}

	row := q.db.QueryRowContext(ctx, query, args...)

	var a Agglomeration
	if err = a.scan(row); err != nil {
		return Agglomeration{}, err
	}

	return a, nil

}

func (q AgglomerationsQ) Select(ctx context.Context) ([]Agglomeration, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", AgglomerationTable, err)
	}

	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing select query for %s: %w", AgglomerationTable, err)
	}
	defer rows.Close()

	var agglomerations []Agglomeration
	for rows.Next() {
		var agglomeration Agglomeration
		err = agglomeration.scan(rows)
		if err != nil {
			return nil, err
		}
		agglomerations = append(agglomerations, agglomeration)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return agglomerations, nil
}

func (q AgglomerationsQ) UpdateOne(ctx context.Context) (Agglomeration, error) {
	q.updater = q.updater.Set("updated_at", time.Now().UTC())

	query, args, err := q.updater.
		Suffix("RETURNING " + AgglomerationColumns).
		ToSql()
	if err != nil {
		return Agglomeration{}, fmt.Errorf("building update query for %s: %w", AgglomerationTable, err)
	}

	var updated Agglomeration
	if err = updated.scan(q.db.QueryRowContext(ctx, query, args...)); err != nil {
		return Agglomeration{}, err
	}

	return updated, nil
}

func (q AgglomerationsQ) UpdateMany(ctx context.Context) (int64, error) {
	q.updater = q.updater.Set("updated_at", time.Now().UTC())

	query, args, err := q.updater.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building update query for %s: %w", AgglomerationTable, err)
	}

	res, err := q.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("executing update query for %s: %w", AgglomerationTable, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("rows affected for %s: %w", AgglomerationTable, err)
	}

	return affected, nil
}

func (q AgglomerationsQ) UpdateName(name string) AgglomerationsQ {
	q.updater = q.updater.Set("name", name)
	return q
}

func (q AgglomerationsQ) UpdateIcon(icon string) AgglomerationsQ {
	q.updater = q.updater.Set("icon", icon)
	return q
}

func (q AgglomerationsQ) UpdateStatus(status string) AgglomerationsQ {
	q.updater = q.updater.Set("status", status)
	return q
}

func (q AgglomerationsQ) UpdateMaxRoles(maxRoles uint) AgglomerationsQ {
	q.updater = q.updater.Set("max_roles", maxRoles)
	return q
}

func (q AgglomerationsQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", AgglomerationTable, err)
	}

	_, err = q.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("executing delete query for %s: %w", AgglomerationTable, err)
	}

	return nil
}

func (q AgglomerationsQ) Page(limit, offset uint) AgglomerationsQ {
	q.selector = q.selector.Limit(uint64(limit)).Offset(uint64(offset))
	return q
}

func (q AgglomerationsQ) Count(ctx context.Context) (uint, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for %s: %w", AgglomerationTable, err)
	}

	row := q.db.QueryRowContext(ctx, query, args...)

	var count uint
	err = row.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("scanning count for %s: %w", AgglomerationTable, err)
	}

	return count, nil
}
