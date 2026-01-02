package pgdb

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/netbill/pgx"
)

const RolePermissionsTable = "role_permissions"
const RolePermissionsColumns = "role_id, permission_id"

type RolePermissionsQ struct {
	db       pgx.DBTX
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

type RolePermission struct {
	RoleID       uuid.UUID `json:"role_id"`
	PermissionID uuid.UUID `json:"permission_id"`
}

func NewRolePermissionsQ(db pgx.DBTX) RolePermissionsQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return RolePermissionsQ{
		db:       db,
		selector: b.Select(RolePermissionsColumns).From(RolePermissionsTable),
		inserter: b.Insert(RolePermissionsTable),
		deleter:  b.Delete(RolePermissionsTable),
		counter:  b.Select("COUNT(*)").From(RolePermissionsTable),
	}
}

func (q RolePermissionsQ) Insert(ctx context.Context, data ...RolePermission) error {
	if len(data) == 0 {
		return nil
	}

	ins := q.inserter.Columns("role_id", "permission_id")

	for _, rp := range data {
		ins = ins.Values(rp.RoleID, rp.PermissionID)
	}

	// если не хочешь падать на дублях — раскомментируй
	// ins = ins.Suffix("ON CONFLICT DO NOTHING")

	query, args, err := ins.ToSql()
	if err != nil {
		return fmt.Errorf("building insert query for %s: %w", RolePermissionsTable, err)
	}

	if _, err := q.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("executing insert query for %s: %w", RolePermissionsTable, err)
	}

	return nil
}

func (q RolePermissionsQ) Get(ctx context.Context) (RolePermission, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return RolePermission{}, fmt.Errorf("building select query for %s: %w", RolePermissionsTable, err)
	}

	var rp RolePermission
	if err = q.db.QueryRowContext(ctx, query, args...).Scan(&rp.RoleID, &rp.PermissionID); err != nil {
		return RolePermission{}, fmt.Errorf("executing select query for %s: %w", RolePermissionsTable, err)
	}

	return rp, nil
}

func (q RolePermissionsQ) Select(ctx context.Context) ([]RolePermission, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", RolePermissionsTable, err)
	}

	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing select query for %s: %w", RolePermissionsTable, err)
	}
	defer rows.Close()

	var rps []RolePermission
	for rows.Next() {
		var rp RolePermission
		if err := rows.Scan(&rp.RoleID, &rp.PermissionID); err != nil {
			return nil, fmt.Errorf("scanning row for %s: %w", RolePermissionsTable, err)
		}
		rps = append(rps, rp)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows for %s: %w", RolePermissionsTable, err)
	}

	return rps, nil
}

func (q RolePermissionsQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", RolePermissionsTable, err)
	}

	_, err = q.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("executing delete query for %s: %w", RolePermissionsTable, err)
	}

	return nil
}

func (q RolePermissionsQ) FilterByRoleID(roleID uuid.UUID) RolePermissionsQ {
	q.selector = q.selector.Where(sq.Eq{"role_id": roleID})
	q.deleter = q.deleter.Where(sq.Eq{"role_id": roleID})
	q.counter = q.counter.Where(sq.Eq{"role_id": roleID})
	return q
}

func (q RolePermissionsQ) FilterByPermissionID(permissionID uuid.UUID) RolePermissionsQ {
	q.selector = q.selector.Where(sq.Eq{"permission_id": permissionID})
	q.deleter = q.deleter.Where(sq.Eq{"permission_id": permissionID})
	q.counter = q.counter.Where(sq.Eq{"permission_id": permissionID})
	return q
}

func (q RolePermissionsQ) FilterByPermissionCode(code ...string) RolePermissionsQ {
	sub := sq.
		Select("id").
		From(PermissionTable).
		Where(sq.Eq{"code": code})

	subSQL, subArgs, err := sub.ToSql()
	if err != nil {
		q.selector = q.selector.Where(sq.Expr("1=0"))
		q.deleter = q.deleter.Where(sq.Expr("1=0"))
		q.counter = q.counter.Where(sq.Expr("1=0"))
		return q
	}

	expr := sq.Expr("permission_id IN ("+subSQL+")", subArgs...)
	q.selector = q.selector.Where(expr)
	q.deleter = q.deleter.Where(expr)
	q.counter = q.counter.Where(expr)

	return q
}

func (q RolePermissionsQ) FilterByAccountID(accountID uuid.UUID) RolePermissionsQ {
	sub := sq.
		Select("DISTINCT mr.role_id").
		From("members m").
		Join("member_roles mr ON mr.member_id = m.id").
		Where(sq.Eq{"m.account_id": accountID})

	subSQL, subArgs, err := sub.ToSql()
	if err != nil {
		q.selector = q.selector.Where(sq.Expr("1=0"))
		q.deleter = q.deleter.Where(sq.Expr("1=0"))
		q.counter = q.counter.Where(sq.Expr("1=0"))
		return q
	}

	expr := sq.Expr("role_id IN ("+subSQL+")", subArgs...)
	q.selector = q.selector.Where(expr)
	q.deleter = q.deleter.Where(expr)
	q.counter = q.counter.Where(expr)

	return q
}

func (q RolePermissionsQ) FilterByOrganizationID(organizationID uuid.UUID) RolePermissionsQ {
	sub := sq.
		Select("id").
		From("roles").
		Where(sq.Eq{"organization_id": organizationID})

	subSQL, subArgs, err := sub.ToSql()
	if err != nil {
		q.selector = q.selector.Where(sq.Expr("1=0"))
		q.deleter = q.deleter.Where(sq.Expr("1=0"))
		q.counter = q.counter.Where(sq.Expr("1=0"))
		return q
	}

	expr := sq.Expr("role_id IN ("+subSQL+")", subArgs...)
	q.selector = q.selector.Where(expr)
	q.deleter = q.deleter.Where(expr)
	q.counter = q.counter.Where(expr)

	return q
}

func (q RolePermissionsQ) FilterByMemberID(memberID uuid.UUID) RolePermissionsQ {
	sub := sq.
		Select("mr.role_id").
		From("member_roles mr").
		Where(sq.Eq{"mr.member_id": memberID})

	subSQL, subArgs, err := sub.ToSql()
	if err != nil {
		q.selector = q.selector.Where(sq.Expr("1=0"))
		q.deleter = q.deleter.Where(sq.Expr("1=0"))
		q.counter = q.counter.Where(sq.Expr("1=0"))
		return q
	}

	expr := sq.Expr("role_id IN ("+subSQL+")", subArgs...)
	q.selector = q.selector.Where(expr)
	q.deleter = q.deleter.Where(expr)
	q.counter = q.counter.Where(expr)

	return q
}

func (q RolePermissionsQ) Count(ctx context.Context) (uint, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for %s: %w", RolePermissionsTable, err)
	}

	var n uint
	if err = q.db.QueryRowContext(ctx, query, args...).Scan(&n); err != nil {
		return 0, fmt.Errorf("scanning count for %s: %w", RolePermissionsTable, err)
	}
	return n, nil
}

func (q RolePermissionsQ) Page(limit, offset uint) RolePermissionsQ {
	q.selector = q.selector.Limit(uint64(limit)).Offset(uint64(offset))
	return q
}

func (q RolePermissionsQ) Exists(ctx context.Context) (bool, error) {
	subSQL, subArgs, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return false, fmt.Errorf("building exists query for %s: %w", RolePermissionsTable, err)
	}

	sqlq := "SELECT EXISTS (" + subSQL + ")"

	var ok bool
	if err := q.db.QueryRowContext(ctx, sqlq, subArgs...).Scan(&ok); err != nil {
		return false, fmt.Errorf("scanning exists for %s: %w", RolePermissionsTable, err)
	}
	return ok, nil
}
