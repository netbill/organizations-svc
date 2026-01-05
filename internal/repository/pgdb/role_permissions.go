package pgdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/netbill/pgx"
)

const PermissionTable = "organization_role_permission"
const PermissionColumns = "id, code, description"

type RolePermission struct {
	ID          uuid.UUID `json:"id"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
}

func (p *RolePermission) scan(row sq.RowScanner) error {
	if err := row.Scan(&p.ID, &p.Code, &p.Description); err != nil {
		return fmt.Errorf("scanning permission: %w", err)
	}
	return nil
}

type RolePermissionsQ struct {
	db       pgx.DBTX
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewPermissionsQ(db pgx.DBTX) RolePermissionsQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return RolePermissionsQ{
		db:       db,
		selector: b.Select(PermissionColumns).From(PermissionTable),
		inserter: b.Insert(PermissionTable),
		updater:  b.Update(PermissionTable),
		deleter:  b.Delete(PermissionTable),
		counter:  b.Select("COUNT(*)").From(PermissionTable),
	}
}

func (q RolePermissionsQ) Insert(ctx context.Context, data RolePermission) (RolePermission, error) {
	query, args, err := q.inserter.SetMap(map[string]any{
		"id":          data.ID,
		"code":        data.Code,
		"description": data.Description,
	}).Suffix("RETURNING " + PermissionColumns).ToSql()
	if err != nil {
		return RolePermission{}, fmt.Errorf("building insert query for %s: %w", PermissionTable, err)
	}

	var out RolePermission
	if err = out.scan(q.db.QueryRowContext(ctx, query, args...)); err != nil {
		return RolePermission{}, err
	}
	return out, nil
}

func (q RolePermissionsQ) Get(ctx context.Context) (RolePermission, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return RolePermission{}, fmt.Errorf("building select query for %s: %w", PermissionTable, err)
	}

	var out RolePermission
	if err = out.scan(q.db.QueryRowContext(ctx, query, args...)); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return RolePermission{}, nil
		default:
			return RolePermission{}, err
		}
	}
	return out, nil
}

func (q RolePermissionsQ) Select(ctx context.Context) ([]RolePermission, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", PermissionTable, err)
	}

	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing select query for %s: %w", PermissionTable, err)
	}
	defer rows.Close()

	var out []RolePermission
	for rows.Next() {
		var p RolePermission
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

func (q RolePermissionsQ) UpdateOne(ctx context.Context) (RolePermission, error) {
	query, args, err := q.updater.Suffix("RETURNING " + PermissionColumns).ToSql()
	if err != nil {
		return RolePermission{}, fmt.Errorf("building update query for %s: %w", PermissionTable, err)
	}

	var out RolePermission
	if err = out.scan(q.db.QueryRowContext(ctx, query, args...)); err != nil {
		return RolePermission{}, err
	}
	return out, nil
}

func (q RolePermissionsQ) UpdateMany(ctx context.Context) (int64, error) {
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

func (q RolePermissionsQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", PermissionTable, err)
	}
	if _, err = q.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("executing delete query for %s: %w", PermissionTable, err)
	}
	return nil
}

func (q RolePermissionsQ) FilterByID(id uuid.UUID) RolePermissionsQ {
	q.selector = q.selector.Where(sq.Eq{"id": id})
	q.counter = q.counter.Where(sq.Eq{"id": id})
	q.updater = q.updater.Where(sq.Eq{"id": id})
	q.deleter = q.deleter.Where(sq.Eq{"id": id})
	return q
}

func (q RolePermissionsQ) FilterByCode(code ...string) RolePermissionsQ {
	q.selector = q.selector.Where(sq.Eq{"code": code})
	q.counter = q.counter.Where(sq.Eq{"code": code})
	q.updater = q.updater.Where(sq.Eq{"code": code})
	q.deleter = q.deleter.Where(sq.Eq{"code": code})
	return q
}

func (q RolePermissionsQ) FilterByRoleID(roleID uuid.UUID) RolePermissionsQ {
	q.selector = q.selector.
		Join("organization_role_permission_links rp ON rp.permission_id = role_permissions.id").
		Where(sq.Eq{"rp.role_id": roleID}).
		Distinct()

	q.counter = q.counter.
		Join("organization_role_permission_links rp ON rp.permission_id = role_permissions.id").
		Where(sq.Eq{"rp.role_id": roleID})

	return q
}

func (q RolePermissionsQ) FilterLikeDescription(description string) RolePermissionsQ {
	q.selector = q.selector.Where(sq.ILike{"description": "%" + description + "%"})
	q.counter = q.counter.Where(sq.ILike{"description": "%" + description + "%"})
	return q
}

func (q RolePermissionsQ) GetForRole(
	ctx context.Context,
	roleID uuid.UUID,
) (map[RolePermission]bool, error) {

	const sqlq = `
		SELECT
			p.id,
			p.code,
			p.description,
			(rp.permission_id IS NOT NULL) AS enabled
		FROM organization_role_permissions p
		LEFT JOIN organization_role_permission_links rp
			ON rp.permission_id = p.id
			AND rp.role_id = $1
		ORDER BY p.code
	`

	rows, err := q.db.QueryContext(ctx, sqlq, roleID)
	if err != nil {
		return nil, fmt.Errorf("query organization_role_permissions for role: %w", err)
	}
	defer rows.Close()

	out := make(map[RolePermission]bool)

	for rows.Next() {
		var p RolePermission
		var enabled bool

		if err := rows.Scan(
			&p.ID,
			&p.Code,
			&p.Description,
			&enabled,
		); err != nil {
			return nil, fmt.Errorf("scanning permission for role: %w", err)
		}

		out[p] = enabled
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}
