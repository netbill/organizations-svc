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

const organizationPermissionTable = "organization_role_permissions"
const organizationPermissionColumns = "id, code, description, deprecated_at, created_at, updated_at"

func scanOrganizationRolePermission(row sq.RowScanner) (p repository.OrganizationRolePermissionRow, err error) {
	err = row.Scan(
		&p.ID,
		&p.Code,
		&p.Description,
		&p.DeprecatedAt,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return repository.OrganizationRolePermissionRow{}, nil
	case err != nil:
		return repository.OrganizationRolePermissionRow{}, fmt.Errorf("scanning permission: %w", err)
	}
	return p, nil
}

type orgRolePermissions struct {
	db       *pgdbx.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewOrgRolePermissionsQ(db *pgdbx.DB) repository.OrgRolePermissionsQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	return &orgRolePermissions{
		db:       db,
		selector: b.Select(organizationPermissionColumns).From(organizationPermissionTable),
		inserter: b.Insert(organizationPermissionTable),
		updater:  b.Update(organizationPermissionTable),
		deleter:  b.Delete(organizationPermissionTable),
		counter:  b.Select("COUNT(*) AS count").From(organizationPermissionTable),
	}
}

func (q *orgRolePermissions) New() repository.OrgRolePermissionsQ {
	return NewOrgRolePermissionsQ(q.db)
}

func (q *orgRolePermissions) Insert(
	ctx context.Context,
	input repository.OrganizationRolePermissionRow,
) (repository.OrganizationRolePermissionRow, error) {
	if input.ID == uuid.Nil {
		input.ID = uuid.New()
	}
	if input.Code == "" {
		return repository.OrganizationRolePermissionRow{}, fmt.Errorf("missing code")
	}
	if input.Description == "" {
		return repository.OrganizationRolePermissionRow{}, fmt.Errorf("missing description")
	}

	query, args, err := q.inserter.
		SetMap(map[string]any{
			"id":            input.ID,
			"code":          input.Code,
			"description":   input.Description,
			"deprecated_at": input.DeprecatedAt,
		}).
		Suffix("RETURNING " + organizationPermissionColumns).
		ToSql()
	if err != nil {
		return repository.OrganizationRolePermissionRow{}, fmt.Errorf("building insert query for %s: %w", organizationPermissionTable, err)
	}

	return scanOrganizationRolePermission(q.db.QueryRow(ctx, query, args...))
}

func (q *orgRolePermissions) FilterByID(ids ...uuid.UUID) repository.OrgRolePermissionsQ {
	if len(ids) == 0 {
		return q
	}
	q.selector = q.selector.Where(sq.Eq{"id": ids})
	q.counter = q.counter.Where(sq.Eq{"id": ids})
	q.updater = q.updater.Where(sq.Eq{"id": ids})
	q.deleter = q.deleter.Where(sq.Eq{"id": ids})
	return q
}

func (q *orgRolePermissions) FilterByCode(codes ...string) repository.OrgRolePermissionsQ {
	if len(codes) == 0 {
		return q
	}
	q.selector = q.selector.Where(sq.Eq{"code": codes})
	q.counter = q.counter.Where(sq.Eq{"code": codes})
	q.updater = q.updater.Where(sq.Eq{"code": codes})
	q.deleter = q.deleter.Where(sq.Eq{"code": codes})
	return q
}

func (q *orgRolePermissions) FilterByDeprecated(deprecated bool) repository.OrgRolePermissionsQ {
	var cond sq.Sqlizer
	if deprecated {
		cond = sq.Expr("deprecated_at IS NOT NULL")
	} else {
		cond = sq.Expr("deprecated_at IS NULL")
	}

	q.selector = q.selector.Where(cond)
	q.counter = q.counter.Where(cond)
	q.updater = q.updater.Where(cond)
	q.deleter = q.deleter.Where(cond)
	return q
}

func (q *orgRolePermissions) UpdateDeprecatedAt(timestamp *time.Time) repository.OrgRolePermissionsQ {
	q.updater = q.updater.Set("deprecated_at", timestamp)
	q.updater = q.updater.Set("updated_at", sq.Expr("(now() AT TIME ZONE 'UTC')"))
	return q
}

func (q *orgRolePermissions) UpdateOne(ctx context.Context) (repository.OrganizationRolePermissionRow, error) {
	query, args, err := q.updater.
		Suffix("RETURNING " + organizationPermissionColumns).
		ToSql()
	if err != nil {
		return repository.OrganizationRolePermissionRow{}, fmt.Errorf("building update query for %s: %w", organizationPermissionTable, err)
	}

	return scanOrganizationRolePermission(q.db.QueryRow(ctx, query, args...))
}

func (q *orgRolePermissions) UpdateMany(ctx context.Context) (int64, error) {
	query, args, err := q.updater.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building update query for %s: %w", organizationPermissionTable, err)
	}

	res, err := q.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("executing update query for %s: %w", organizationPermissionTable, err)
	}

	return res.RowsAffected(), nil
}

func (q *orgRolePermissions) Select(ctx context.Context) ([]repository.OrganizationRolePermissionRow, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", organizationPermissionTable, err)
	}

	rows, err := q.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing select query for %s: %w", organizationPermissionTable, err)
	}
	defer rows.Close()

	out := make([]repository.OrganizationRolePermissionRow, 0)
	for rows.Next() {
		p, err := scanOrganizationRolePermission(rows)
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

func (q *orgRolePermissions) Get(ctx context.Context) (repository.OrganizationRolePermissionRow, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return repository.OrganizationRolePermissionRow{}, fmt.Errorf("building select query for %s: %w", organizationPermissionTable, err)
	}

	return scanOrganizationRolePermission(q.db.QueryRow(ctx, query, args...))
}

func (q *orgRolePermissions) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", organizationPermissionTable, err)
	}

	if _, err := q.db.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("executing delete query for %s: %w", organizationPermissionTable, err)
	}
	return nil
}

func (q *orgRolePermissions) Count(ctx context.Context) (uint, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for %s: %w", organizationPermissionTable, err)
	}

	var n uint
	if err = q.db.QueryRow(ctx, query, args...).Scan(&n); err != nil {
		return 0, fmt.Errorf("scanning count for %s: %w", organizationPermissionTable, err)
	}
	return n, nil
}

func (q *orgRolePermissions) Page(limit, offset uint) repository.OrgRolePermissionsQ {
	q.selector = q.selector.Limit(uint64(limit)).Offset(uint64(offset))
	return q
}
