package pgdb

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/umisto/pgx"
)

const InviteTable = "invites"
const InviteColumns = "id, agglomeration_id, account_id, status, expires_at, created_at"

type Invite struct {
	ID              uuid.UUID `json:"id"`
	AgglomerationID uuid.UUID `json:"agglomeration_id"`
	AccountID       uuid.UUID `json:"account_id,omitempty"`
	Status          string    `json:"status"`
	ExpiresAt       time.Time `json:"expires_at"`
	CreatedAt       time.Time `json:"created_at"`
}

func (i *Invite) scan(row sq.RowScanner) error {
	if err := row.Scan(
		&i.ID,
		&i.AgglomerationID,
		&i.AccountID,
		&i.Status,
		&i.ExpiresAt,
		&i.CreatedAt,
	); err != nil {
		return fmt.Errorf("scanning invite: %w", err)
	}
	return nil
}

type InvitesQ struct {
	db       pgx.DBTX
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewInvitesQ(db pgx.DBTX) InvitesQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return InvitesQ{
		db:       db,
		selector: b.Select(InviteColumns).From(InviteTable),
		inserter: b.Insert(InviteTable),
		updater:  b.Update(InviteTable),
		deleter:  b.Delete(InviteTable),
		counter:  b.Select("COUNT(*)").From(InviteTable),
	}
}

type InsertInviteParams struct {
	AgglomerationID uuid.UUID
	AccountID       uuid.UUID
	ExpiresAt       time.Time
}

func (q InvitesQ) Insert(ctx context.Context, data InsertInviteParams) (Invite, error) {
	query, args, err := q.inserter.SetMap(map[string]any{
		"agglomeration_id": data.AgglomerationID,
		"account_id":       data.AccountID,
		"expires_at":       data.ExpiresAt,
	}).Suffix("RETURNING " + InviteColumns).ToSql()
	if err != nil {
		return Invite{}, fmt.Errorf("building insert query for %s: %w", InviteTable, err)
	}

	var out Invite
	if err = out.scan(q.db.QueryRowContext(ctx, query, args...)); err != nil {
		return Invite{}, err
	}
	return out, nil
}

func (q InvitesQ) Get(ctx context.Context) (Invite, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return Invite{}, fmt.Errorf("building select query for %s: %w", InviteTable, err)
	}

	var out Invite
	if err = out.scan(q.db.QueryRowContext(ctx, query, args...)); err != nil {
		return Invite{}, err
	}
	return out, nil
}

func (q InvitesQ) Select(ctx context.Context) ([]Invite, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", InviteTable, err)
	}

	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing select query for %s: %w", InviteTable, err)
	}
	defer rows.Close()

	var out []Invite
	for rows.Next() {
		var i Invite
		if err = i.scan(rows); err != nil {
			return nil, err
		}
		out = append(out, i)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (q InvitesQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", InviteTable, err)
	}

	if _, err = q.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("executing delete query for %s: %w", InviteTable, err)
	}
	return nil
}

func (q InvitesQ) Count(ctx context.Context) (int64, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for %s: %w", InviteTable, err)
	}

	var n int64
	if err = q.db.QueryRowContext(ctx, query, args...).Scan(&n); err != nil {
		return 0, fmt.Errorf("scanning count for %s: %w", InviteTable, err)
	}
	return n, nil
}

func (q InvitesQ) UpdateOne(ctx context.Context) (Invite, error) {
	query, args, err := q.updater.Suffix("RETURNING " + InviteColumns).ToSql()
	if err != nil {
		return Invite{}, fmt.Errorf("building update query for %s: %w", InviteTable, err)
	}

	var out Invite
	if err = out.scan(q.db.QueryRowContext(ctx, query, args...)); err != nil {
		return Invite{}, err
	}
	return out, nil
}

func (q InvitesQ) UpdateMany(ctx context.Context) (int64, error) {
	query, args, err := q.updater.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building update query for %s: %w", InviteTable, err)
	}

	res, err := q.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("executing update query for %s: %w", InviteTable, err)
	}

	aff, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("rows affected for %s: %w", InviteTable, err)
	}

	return aff, nil
}

func (q InvitesQ) FilterByID(id uuid.UUID) InvitesQ {
	q.selector = q.selector.Where(sq.Eq{"id": id})
	q.counter = q.counter.Where(sq.Eq{"id": id})
	q.updater = q.updater.Where(sq.Eq{"id": id})
	q.deleter = q.deleter.Where(sq.Eq{"id": id})
	return q
}

func (q InvitesQ) FilterByAgglomerationID(id uuid.UUID) InvitesQ {
	q.selector = q.selector.Where(sq.Eq{"agglomeration_id": id})
	q.counter = q.counter.Where(sq.Eq{"agglomeration_id": id})
	q.updater = q.updater.Where(sq.Eq{"agglomeration_id": id})
	q.deleter = q.deleter.Where(sq.Eq{"agglomeration_id": id})
	return q
}

func (q InvitesQ) FilterByAccountID(id uuid.UUID) InvitesQ {
	q.selector = q.selector.Where(sq.Eq{"account_id": id})
	q.counter = q.counter.Where(sq.Eq{"account_id": id})
	q.updater = q.updater.Where(sq.Eq{"account_id": id})
	q.deleter = q.deleter.Where(sq.Eq{"account_id": id})
	return q
}

func (q InvitesQ) FilterByStatus(status string) InvitesQ {
	q.selector = q.selector.Where(sq.Eq{"status": status})
	q.counter = q.counter.Where(sq.Eq{"status": status})
	q.updater = q.updater.Where(sq.Eq{"status": status})
	q.deleter = q.deleter.Where(sq.Eq{"status": status})
	return q
}

func (q InvitesQ) FilterExpiresBefore(t time.Time) InvitesQ {
	q.selector = q.selector.Where(sq.Lt{"expires_at": t})
	q.counter = q.counter.Where(sq.Lt{"expires_at": t})
	q.updater = q.updater.Where(sq.Lt{"expires_at": t})
	q.deleter = q.deleter.Where(sq.Lt{"expires_at": t})
	return q
}

func (q InvitesQ) FilterExpiresAfter(t time.Time) InvitesQ {
	q.selector = q.selector.Where(sq.GtOrEq{"expires_at": t})
	q.counter = q.counter.Where(sq.GtOrEq{"expires_at": t})
	q.updater = q.updater.Where(sq.GtOrEq{"expires_at": t})
	q.deleter = q.deleter.Where(sq.GtOrEq{"expires_at": t})
	return q
}

func (q InvitesQ) UpdateStatus(status string) InvitesQ {
	q.updater = q.updater.Set("status", status)
	return q
}

func (q InvitesQ) UpdateExpiresAt(t time.Time) InvitesQ {
	q.updater = q.updater.Set("expires_at", t)
	return q
}

func (q InvitesQ) CursorCreatedAt(limit uint, asc bool, createdAt time.Time, id uuid.UUID) InvitesQ {
	if asc {
		q.selector = q.selector.OrderBy("created_at ASC", "id ASC")
	} else {
		q.selector = q.selector.OrderBy("created_at DESC", "id DESC")
	}

	q.selector = q.selector.Limit(uint64(limit))

	if asc {
		q.selector = q.selector.Where(sq.Expr("(created_at, id) > (?, ?)", createdAt, id))
	} else {
		q.selector = q.selector.Where(sq.Expr("(created_at, id) < (?, ?)", createdAt, id))
	}

	return q
}
