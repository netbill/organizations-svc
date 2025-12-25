package pgdbsq

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/repository/pgdb"
)

const MemberRoleTable = "member_roles"
const MemberRoleColumns = "member_id, role_id"

type MemberRole struct {
	MemberID uuid.UUID `json:"member_id"`
	RoleID   uuid.UUID `json:"role_id"`
}

func (mr *MemberRole) scan(row sq.RowScanner) error {
	if err := row.Scan(&mr.MemberID, &mr.RoleID); err != nil {
		return fmt.Errorf("scanning member_role: %w", err)
	}
	return nil
}

type MemberRoleQ struct {
	db       pgdb.DBTX
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewMemberRoleQ(db pgdb.DBTX) MemberRoleQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return MemberRoleQ{
		db:       db,
		selector: b.Select(MemberRoleColumns).From(MemberRoleTable),
		inserter: b.Insert(MemberRoleTable),
		deleter:  b.Delete(MemberRoleTable),
		counter:  b.Select("COUNT(*)").From(MemberRoleTable),
	}
}

func (q MemberRoleQ) New() MemberRoleQ { return NewMemberRoleQ(q.db) }

func (q MemberRoleQ) Insert(ctx context.Context, data MemberRole) (MemberRole, error) {
	query, args, err := q.inserter.SetMap(map[string]any{
		"member_id": data.MemberID,
		"role_id":   data.RoleID,
	}).Suffix("RETURNING " + MemberRoleColumns).ToSql()
	if err != nil {
		return MemberRole{}, fmt.Errorf("building insert query for %s: %w", MemberRoleTable, err)
	}

	var out MemberRole
	if err = out.scan(q.db.QueryRowContext(ctx, query, args...)); err != nil {
		return MemberRole{}, err
	}
	return out, nil
}

// удобно, чтобы "добавить если нет" не падало на PK конфликте
func (q MemberRoleQ) InsertIgnore(ctx context.Context, data MemberRole) error {
	query, args, err := q.inserter.SetMap(map[string]any{
		"member_id": data.MemberID,
		"role_id":   data.RoleID,
	}).Suffix("ON CONFLICT (member_id, role_id) DO NOTHING").ToSql()
	if err != nil {
		return fmt.Errorf("building insert ignore query for %s: %w", MemberRoleTable, err)
	}

	if _, err = q.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("executing insert ignore query for %s: %w", MemberRoleTable, err)
	}
	return nil
}

func (q MemberRoleQ) Get(ctx context.Context) (MemberRole, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return MemberRole{}, fmt.Errorf("building select query for %s: %w", MemberRoleTable, err)
	}

	var out MemberRole
	if err = out.scan(q.db.QueryRowContext(ctx, query, args...)); err != nil {
		return MemberRole{}, err
	}
	return out, nil
}

func (q MemberRoleQ) Select(ctx context.Context) ([]MemberRole, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", MemberRoleTable, err)
	}

	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing select query for %s: %w", MemberRoleTable, err)
	}
	defer rows.Close()

	var out []MemberRole
	for rows.Next() {
		var mr MemberRole
		if err = mr.scan(rows); err != nil {
			return nil, err
		}
		out = append(out, mr)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (q MemberRoleQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", MemberRoleTable, err)
	}
	if _, err = q.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("executing delete query for %s: %w", MemberRoleTable, err)
	}
	return nil
}

func (q MemberRoleQ) Count(ctx context.Context) (int64, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for %s: %w", MemberRoleTable, err)
	}

	var n int64
	if err = q.db.QueryRowContext(ctx, query, args...).Scan(&n); err != nil {
		return 0, fmt.Errorf("scanning count for %s: %w", MemberRoleTable, err)
	}
	return n, nil
}

func (q MemberRoleQ) FilterMemberID(memberID uuid.UUID) MemberRoleQ {
	q.selector = q.selector.Where(sq.Eq{"member_id": memberID})
	q.counter = q.counter.Where(sq.Eq{"member_id": memberID})
	q.deleter = q.deleter.Where(sq.Eq{"member_id": memberID})
	return q
}

func (q MemberRoleQ) FilterRoleID(roleID uuid.UUID) MemberRoleQ {
	q.selector = q.selector.Where(sq.Eq{"role_id": roleID})
	q.counter = q.counter.Where(sq.Eq{"role_id": roleID})
	q.deleter = q.deleter.Where(sq.Eq{"role_id": roleID})
	return q
}
