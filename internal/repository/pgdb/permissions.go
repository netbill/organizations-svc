package pgdb

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/umisto/pgx"
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

type PermissionsQ struct {
	db       pgx.DBTX
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewPermissionsQ(db pgx.DBTX) PermissionsQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return PermissionsQ{
		db:       db,
		selector: b.Select(PermissionColumns).From(PermissionTable),
		inserter: b.Insert(PermissionTable),
		updater:  b.Update(PermissionTable),
		deleter:  b.Delete(PermissionTable),
		counter:  b.Select("COUNT(*)").From(PermissionTable),
	}
}

func (q PermissionsQ) New() PermissionsQ { return NewPermissionsQ(q.db) }

func (q PermissionsQ) Insert(ctx context.Context, data Permission) (Permission, error) {
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

func (q PermissionsQ) Get(ctx context.Context) (Permission, error) {
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

func (q PermissionsQ) Select(ctx context.Context) ([]Permission, error) {
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

func (q PermissionsQ) UpdateOne(ctx context.Context) (Permission, error) {
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

func (q PermissionsQ) UpdateMany(ctx context.Context) (int64, error) {
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

func (q PermissionsQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", PermissionTable, err)
	}
	if _, err = q.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("executing delete query for %s: %w", PermissionTable, err)
	}
	return nil
}

func (q PermissionsQ) Count(ctx context.Context) (int64, error) {
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

func (q PermissionsQ) FilterByID(id uuid.UUID) PermissionsQ {
	q.selector = q.selector.Where(sq.Eq{"id": id})
	q.counter = q.counter.Where(sq.Eq{"id": id})
	q.updater = q.updater.Where(sq.Eq{"id": id})
	q.deleter = q.deleter.Where(sq.Eq{"id": id})
	return q
}

func (q PermissionsQ) FilterByCode(code string) PermissionsQ {
	q.selector = q.selector.Where(sq.Eq{"code": code})
	q.counter = q.counter.Where(sq.Eq{"code": code})
	q.updater = q.updater.Where(sq.Eq{"code": code})
	q.deleter = q.deleter.Where(sq.Eq{"code": code})
	return q
}

func (q PermissionsQ) FilterByRoleID(roleID uuid.UUID) PermissionsQ {
	q.selector = q.selector.
		Join("role_permissions rp ON rp.permission_id = permissions.id").
		Where(sq.Eq{"rp.role_id": roleID}).
		Distinct()

	q.counter = q.counter.
		Join("role_permissions rp ON rp.permission_id = permissions.id").
		Where(sq.Eq{"rp.role_id": roleID})

	return q
}

func (q PermissionsQ) FilterLikeDescription(description string) PermissionsQ {
	q.selector = q.selector.Where(sq.ILike{"description": "%" + description + "%"})
	q.counter = q.counter.Where(sq.ILike{"description": "%" + description + "%"})
	return q
}

func (q PermissionsQ) UpdateCode(code string) PermissionsQ {
	q.updater = q.updater.Set("code", code)
	return q
}

func (q PermissionsQ) UpdateDescription(description string) PermissionsQ {
	q.updater = q.updater.Set("description", description)
	return q
}

func (q PermissionsQ) Page(limit, offset uint) PermissionsQ {
	q.selector = q.selector.Limit(uint64(limit)).Offset(uint64(offset))
	return q
}
