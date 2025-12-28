package pgdb

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/umisto/pgx"
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

func (q RolePermissionsQ) Insert(ctx context.Context, data RolePermission) (RolePermission, error) {
	query, args, err := q.inserter.SetMap(map[string]any{
		"role_id":       data.RoleID,
		"permission_id": data.PermissionID,
	}).Suffix("RETURNING " + RolePermissionsColumns).ToSql()
	if err != nil {
		return RolePermission{}, fmt.Errorf("building insert query for %s: %w", RolePermissionsTable, err)
	}

	var rp RolePermission
	if err = q.db.QueryRowContext(ctx, query, args...).Scan(&rp.RoleID, &rp.PermissionID); err != nil {
		return RolePermission{}, fmt.Errorf("executing insert query for %s: %w", RolePermissionsTable, err)
	}

	return rp, nil
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

func (q RolePermissionsQ) FilterByPermissionCode(code string) RolePermissionsQ {
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

func (q RolePermissionsQ) FilterByAgglomerationID(agglomerationID uuid.UUID) RolePermissionsQ {
	sub := sq.
		Select("id").
		From("roles").
		Where(sq.Eq{"agglomeration_id": agglomerationID})

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

func (q RolePermissionsQ) Count(ctx context.Context) (int64, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for %s: %w", RolePermissionsTable, err)
	}

	var n int64
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

//func (q RolePermissionsQ) CheckMemberHavePermissionInAgglomerationByCode(
//	ctx context.Context,
//	memberID uuid.UUID,
//	agglomerationID uuid.UUID,
//	code string,
//) (bool, error) {
//	const sqlq = `
//		SELECT EXISTS (
//			SELECT 1
//			FROM members m
//			JOIN member_roles mr ON mr.member_id = m.id
//			JOIN role_permissions rp ON rp.role_id = mr.role_id
//			JOIN permissions p ON p.id = rp.permission_id
//			WHERE m.id = $1
//				AND m.agglomeration_id = $2
//				AND p.code = $3
//		)
//	`
//
//	var ok bool
//	if err := q.db.QueryRowContext(ctx, sqlq, memberID, agglomerationID, code).Scan(&ok); err != nil {
//		return false, fmt.Errorf("scanning permission exists (member+code): %w", err)
//	}
//	return ok, nil
//}
//
//func (q RolePermissionsQ) CheckMemberHavePermissionInAgglomerationByID(
//	ctx context.Context,
//	memberID uuid.UUID,
//	agglomerationID uuid.UUID,
//	permissionID uuid.UUID,
//) (bool, error) {
//	const sqlq = `
//		SELECT EXISTS (
//			SELECT 1
//			FROM members m
//			JOIN member_roles mr ON mr.member_id = m.id
//			JOIN role_permissions rp ON rp.role_id = mr.role_id
//			WHERE m.id = $1
//				AND m.agglomeration_id = $2
//				AND rp.permission_id = $3
//		)
//	`
//
//	var ok bool
//	if err := q.db.QueryRowContext(ctx, sqlq, memberID, agglomerationID, permissionID).Scan(&ok); err != nil {
//		return false, fmt.Errorf("scanning permission exists (member+id): %w", err)
//	}
//	return ok, nil
//}
//
//func (q RolePermissionsQ) CheckMemberHavePermissionByCode(
//	ctx context.Context,
//	memberID uuid.UUID,
//	code string,
//) (bool, error) {
//	const sqlq = `
//		SELECT EXISTS (
//			SELECT 1
//			FROM member_roles mr
//			JOIN role_permissions rp ON rp.role_id = mr.role_id
//			JOIN permissions p ON p.id = rp.permission_id
//			WHERE mr.member_id = $1
//				AND p.code = $2
//		)
//	`
//
//	var ok bool
//	if err := q.db.QueryRowContext(ctx, sqlq, memberID, code).Scan(&ok); err != nil {
//		return false, fmt.Errorf("scanning permission exists (member+code): %w", err)
//	}
//	return ok, nil
//}
//
//func (q RolePermissionsQ) CheckMemberHavePermissionByID(
//	ctx context.Context,
//	memberID uuid.UUID,
//	permissionID uuid.UUID,
//) (bool, error) {
//	const sqlq = `
//		SELECT EXISTS (
//			SELECT 1
//			FROM member_roles mr
//			JOIN role_permissions rp ON rp.role_id = mr.role_id
//			WHERE mr.member_id = $1
//				AND rp.permission_id = $2
//		)
//	`
//
//	var ok bool
//	if err := q.db.QueryRowContext(ctx, sqlq, memberID, permissionID).Scan(&ok); err != nil {
//		return false, fmt.Errorf("scanning permission exists (member+id): %w", err)
//	}
//	return ok, nil
//}
//
//func (q RolePermissionsQ) CheckAccountHavePermissionByCode(
//	ctx context.Context,
//	accountID uuid.UUID,
//	agglomerationID uuid.UUID,
//	code string,
//) (bool, error) {
//	const sqlq = `
//		SELECT EXISTS (
//			SELECT 1
//			FROM members m
//			JOIN member_roles mr ON mr.member_id = m.id
//			JOIN role_permissions rp ON rp.role_id = mr.role_id
//			JOIN permissions p ON p.id = rp.permission_id
//			WHERE m.account_id = $1
//				AND m.agglomeration_id = $2
//				AND p.code = $3
//		)
//	`
//
//	var ok bool
//	if err := q.db.QueryRowContext(ctx, sqlq, accountID, agglomerationID, code).Scan(&ok); err != nil {
//		return false, fmt.Errorf("scanning permission exists (account+code): %w", err)
//	}
//	return ok, nil
//}
//
//func (q RolePermissionsQ) CheckAccountHavePermissionByID(
//	ctx context.Context,
//	accountID uuid.UUID,
//	agglomerationID uuid.UUID,
//	permissionID uuid.UUID,
//) (bool, error) {
//	const sqlq = `
//		SELECT EXISTS (
//			SELECT 1
//			FROM members m
//			JOIN member_roles mr ON mr.member_id = m.id
//			JOIN role_permissions rp ON rp.role_id = mr.role_id
//			WHERE m.account_id = $1
//				AND m.agglomeration_id = $2
//				AND rp.permission_id = $3
//		)
//	`
//
//	var ok bool
//	if err := q.db.QueryRowContext(ctx, sqlq, accountID, agglomerationID, permissionID).Scan(&ok); err != nil {
//		return false, fmt.Errorf("scanning permission exists (account+id): %w", err)
//	}
//	return ok, nil
//}
