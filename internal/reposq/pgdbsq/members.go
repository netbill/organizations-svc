package pgdbsq

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/repository/pgdb"

	sq "github.com/Masterminds/squirrel"
)

const MemberTable = "members"

const MemberColumns = "id, account_id, agglomeration_id, position, label, created_at, updated_at"
const MemberColumnsM = "m.id, m.account_id, m.agglomeration_id, m.position, m.label, m.created_at, m.updated_at"

type Member struct {
	ID              uuid.UUID `json:"id"`
	AccountID       uuid.UUID `json:"account_id"`
	AgglomerationID uuid.UUID `json:"agglomeration_id"`
	Position        *string   `json:"position"`
	Label           *string   `json:"label"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (m *Member) scan(row sq.RowScanner) error {
	err := row.Scan(
		&m.ID,
		&m.AccountID,
		&m.AgglomerationID,
		&m.Position,
		&m.Label,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("scanning member: %w", err)
	}
	return nil
}

type MemberQ struct {
	db       pgdb.DBTX
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewMemberQ(db pgdb.DBTX) MemberQ {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return MemberQ{
		db:       db,
		selector: builder.Select(MemberColumnsM).From(MemberTable + " m"),
		inserter: builder.Insert(MemberTable),
		updater:  builder.Update(MemberTable + " m"),
		deleter:  builder.Delete(MemberTable + " m"),
		counter:  builder.Select("COUNT(*)").From(MemberTable + " m"),
	}
}

func (q MemberQ) New() MemberQ {
	return NewMemberQ(q.db)
}

func (q MemberQ) Insert(ctx context.Context, data Member) (Member, error) {
	query, args, err := q.inserter.SetMap(map[string]interface{}{
		"account_id":       data.AccountID,
		"agglomeration_id": data.AgglomerationID,
		"position":         data.Position,
		"label":            data.Label,
	}).Suffix("RETURNING " + MemberColumns).ToSql()
	if err != nil {
		return Member{}, fmt.Errorf("building insert query for %s: %w", MemberTable, err)
	}

	var inserted Member
	err = inserted.scan(q.db.QueryRowContext(ctx, query, args...))
	if err != nil {
		return Member{}, err
	}

	return inserted, nil
}

func (q MemberQ) Get(ctx context.Context) (Member, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return Member{}, fmt.Errorf("building select query for %s: %w", MemberTable, err)
	}

	var m Member
	err = m.scan(q.db.QueryRowContext(ctx, query, args...))
	if err != nil {
		return Member{}, err
	}

	return m, nil
}

func (q MemberQ) Select(ctx context.Context) ([]Member, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", MemberTable, err)
	}

	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing select query for %s: %w", MemberTable, err)
	}
	defer rows.Close()

	var out []Member
	for rows.Next() {
		var m Member
		if err = m.scan(rows); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (q MemberQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", MemberTable, err)
	}

	_, err = q.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("executing delete query for %s: %w", MemberTable, err)
	}

	return nil
}

func (q MemberQ) Count(ctx context.Context) (int64, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for %s: %w", MemberTable, err)
	}

	var count int64
	err = q.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("scanning count for %s: %w", MemberTable, err)
	}

	return count, nil
}

func (q MemberQ) UpdateOne(ctx context.Context) (Member, error) {
	q.updater = q.updater.Set("updated_at", time.Now().UTC())

	query, args, err := q.updater.Suffix("RETURNING " + MemberColumns).ToSql()
	if err != nil {
		return Member{}, fmt.Errorf("building update query for %s: %w", MemberTable, err)
	}

	var updated Member
	err = updated.scan(q.db.QueryRowContext(ctx, query, args...))
	if err != nil {
		return Member{}, err
	}

	return updated, nil
}

func (q MemberQ) UpdateMany(ctx context.Context) (int64, error) {
	q.updater = q.updater.Set("updated_at", time.Now().UTC())

	query, args, err := q.updater.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building update query for %s: %w", MemberTable, err)
	}

	res, err := q.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("executing update query for %s: %w", MemberTable, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("rows affected for %s: %w", MemberTable, err)
	}

	return affected, nil
}

func (q MemberQ) FilterByID(id uuid.UUID) MemberQ {
	q.selector = q.selector.Where(sq.Eq{"m.id": id})
	q.counter = q.counter.Where(sq.Eq{"m.id": id})
	q.updater = q.updater.Where(sq.Eq{"m.id": id})
	q.deleter = q.deleter.Where(sq.Eq{"m.id": id})
	return q
}

func (q MemberQ) FilterByAccountID(accountID uuid.UUID) MemberQ {
	q.selector = q.selector.Where(sq.Eq{"m.account_id": accountID})
	q.counter = q.counter.Where(sq.Eq{"m.account_id": accountID})
	q.updater = q.updater.Where(sq.Eq{"m.account_id": accountID})
	q.deleter = q.deleter.Where(sq.Eq{"m.account_id": accountID})
	return q
}

func (q MemberQ) FilterByAgglomerationID(agglomerationID uuid.UUID) MemberQ {
	q.selector = q.selector.Where(sq.Eq{"m.agglomeration_id": agglomerationID})
	q.counter = q.counter.Where(sq.Eq{"m.agglomeration_id": agglomerationID})
	q.updater = q.updater.Where(sq.Eq{"m.agglomeration_id": agglomerationID})
	q.deleter = q.deleter.Where(sq.Eq{"m.agglomeration_id": agglomerationID})
	return q
}

func (q MemberQ) FilterLikePosition(position string) MemberQ {
	q.selector = q.selector.Where(sq.ILike{"m.position": "%" + position + "%"})
	q.counter = q.counter.Where(sq.ILike{"m.position": "%" + position + "%"})
	q.updater = q.updater.Where(sq.ILike{"m.position": "%" + position + "%"})
	q.deleter = q.deleter.Where(sq.ILike{"m.position": "%" + position + "%"})
	return q
}

func (q MemberQ) FilterLikeLabel(label string) MemberQ {
	q.selector = q.selector.Where(sq.ILike{"m.label": "%" + label + "%"})
	q.counter = q.counter.Where(sq.ILike{"m.label": "%" + label + "%"})
	q.updater = q.updater.Where(sq.ILike{"m.label": "%" + label + "%"})
	q.deleter = q.deleter.Where(sq.ILike{"m.label": "%" + label + "%"})
	return q
}

func (q MemberQ) UpdatePosition(position *string) MemberQ {
	q.updater = q.updater.Set("position", position)
	return q
}

func (q MemberQ) UpdateLabel(label *string) MemberQ {
	q.updater = q.updater.Set("label", label)
	return q
}

func (q MemberQ) CursorCreatedAt(limit uint, asc bool, createdAt time.Time, id uuid.UUID) MemberQ {
	if asc {
		q.selector = q.selector.OrderBy("m.created_at ASC", "m.id ASC")
	} else {
		q.selector = q.selector.OrderBy("m.created_at DESC", "m.id DESC")
	}

	q.selector = q.selector.Limit(uint64(limit))

	if asc {
		q.selector = q.selector.Where(sq.Expr("(m.created_at, m.id) > (?, ?)", createdAt, id))
	} else {
		q.selector = q.selector.Where(sq.Expr("(m.created_at, m.id) < (?, ?)", createdAt, id))
	}

	return q
}
