package pgdbsq

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/repository/pgdb"
)

const PermissionTable = "permissions"
const PermissionColumns = "id, code, description"

type Permission struct {
	ID          uuid.UUID `json:"id"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
}

func (p *Permission) scan(row sq.RowScanner) error {
	if err := row.Scan(&p.ID, &p.Code, &p.Description); err != nil {
		return fmt.Errorf("scanning permission: %w", err)
	}
	return nil
}

type PermissionQ struct {
	db       pgdb.DBTX
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewPermissionQ(db pgdb.DBTX) PermissionQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return PermissionQ{
		db:       db,
		selector: b.Select(PermissionColumns).From(PermissionTable),
		inserter: b.Insert(PermissionTable),
		updater:  b.Update(PermissionTable),
		deleter:  b.Delete(PermissionTable),
		counter:  b.Select("COUNT(*)").From(PermissionTable),
	}
}

func (q PermissionQ) New() PermissionQ { return NewPermissionQ(q.db) }

func (q PermissionQ) Insert(ctx context.Context, data Permission) (Permission, error) {
	query, args, err := q.inserter.SetMap(map[string]any{
		"id":          data.ID,
		"code":        data.Code,
		"description": data.Description,
	}).Suffix("RETURNING " + PermissionColumns).ToSql()
	if err != nil {
		return Permission{}, fmt.Errorf("building insert query for %s: %w", PermissionTable, err)
	}

	var out Permission
	if err = out.scan(q.db.QueryRowContext(ctx, query, args...)); err != nil {
		return Permission{}, err
	}
	return out, nil
}

func (q PermissionQ) Get(ctx context.Context) (Permission, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return Permission{}, fmt.Errorf("building select query for %s: %w", PermissionTable, err)
	}

	var out Permission
	if err = out.scan(q.db.QueryRowContext(ctx, query, args...)); err != nil {
		return Permission{}, err
	}
	return out, nil
}

func (q PermissionQ) Select(ctx context.Context) ([]Permission, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", PermissionTable, err)
	}

	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing select query for %s: %w", PermissionTable, err)
	}
	defer rows.Close()

	var out []Permission
	for rows.Next() {
		var p Permission
		if err = p.scan(rows); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (q PermissionQ) UpdateOne(ctx context.Context) (Permission, error) {
	query, args, err := q.updater.Suffix("RETURNING " + PermissionColumns).ToSql()
	if err != nil {
		return Permission{}, fmt.Errorf("building update query for %s: %w", PermissionTable, err)
	}

	var out Permission
	if err = out.scan(q.db.QueryRowContext(ctx, query, args...)); err != nil {
		return Permission{}, err
	}
	return out, nil
}

func (q PermissionQ) UpdateMany(ctx context.Context) (int64, error) {
	query, args, err := q.updater.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building update query for %s: %w", PermissionTable, err)
	}

	res, err := q.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("executing update query for %s: %w", PermissionTable, err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("rows affected for %s: %w", PermissionTable, err)
	}
	return n, nil
}

func (q PermissionQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", PermissionTable, err)
	}
	if _, err = q.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("executing delete query for %s: %w", PermissionTable, err)
	}
	return nil
}

func (q PermissionQ) Count(ctx context.Context) (int64, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for %s: %w", PermissionTable, err)
	}

	var n int64
	if err = q.db.QueryRowContext(ctx, query, args...).Scan(&n); err != nil {
		return 0, fmt.Errorf("scanning count for %s: %w", PermissionTable, err)
	}
	return n, nil
}

func (q PermissionQ) FilterByID(id uuid.UUID) PermissionQ {
	q.selector = q.selector.Where(sq.Eq{"id": id})
	q.counter = q.counter.Where(sq.Eq{"id": id})
	q.updater = q.updater.Where(sq.Eq{"id": id})
	q.deleter = q.deleter.Where(sq.Eq{"id": id})
	return q
}

func (q PermissionQ) FilterByCode(code string) PermissionQ {
	q.selector = q.selector.Where(sq.Eq{"code": code})
	q.counter = q.counter.Where(sq.Eq{"code": code})
	q.updater = q.updater.Where(sq.Eq{"code": code})
	q.deleter = q.deleter.Where(sq.Eq{"code": code})
	return q
}

func (q PermissionQ) UpdateCode(code string) PermissionQ {
	q.updater = q.updater.Set("code", code)
	return q
}

func (q PermissionQ) UpdateDescription(description string) PermissionQ {
	q.updater = q.updater.Set("description", description)
	return q
}
