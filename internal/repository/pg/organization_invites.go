package pg

import (
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/netbill/organizations-svc/internal/core/modules/invite"
	"github.com/netbill/pgdbx"
)

const OrganizationInviteTable = "organization_invites"

const OrganizationInviteColumns = "id, organization_id, account_id, status, updated_at, expires_at, created_at"

func scanOrganizationInvite(row sq.RowScanner) (r invite.OrgInviteRow, err error) {
	if err = row.Scan(
		&r.ID,
		&r.OrganizationID,
		&r.AccountID,
		&r.Status,
		&r.UpdatedAt,
		&r.ExpiresAt,
		&r.CreatedAt,
	); err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return invite.OrgInviteRow{}, nil
		default:
			return invite.OrgInviteRow{}, fmt.Errorf("scanning invite: %w", err)
		}
	}
	return r, nil
}

type orgInvites struct {
	db       *pgdbx.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewOrgInvitesQ(db *pgdbx.DB) invite.OrgInvitesQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return orgInvites{
		db:       db,
		selector: b.Select(OrganizationInviteColumns).From(OrganizationInviteTable),
		inserter: b.Insert(OrganizationInviteTable),
		updater:  b.Update(OrganizationInviteTable),
		deleter:  b.Delete(OrganizationInviteTable),
		counter:  b.Select("COUNT(*)").From(OrganizationInviteTable),
	}
}

func (q orgInvites) New() invite.OrgInvitesQ {
	return NewOrgInvitesQ(q.db)
}

func (q orgInvites) Insert(ctx context.Context, data invite.OrgInviteRow) (invite.OrgInviteRow, error) {
	query, args, err := q.inserter.SetMap(map[string]any{
		"organization_id": data.OrganizationID,
		"account_id":      data.AccountID,
		"expires_at":      data.ExpiresAt,
	}).Suffix("RETURNING " + OrganizationInviteColumns).ToSql()
	if err != nil {
		return invite.OrgInviteRow{}, fmt.Errorf("building insert query for %s: %w", OrganizationInviteTable, err)
	}

	return scanOrganizationInvite(q.db.QueryRow(ctx, query, args...))
}

func (q orgInvites) Get(ctx context.Context) (invite.OrgInviteRow, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return invite.OrgInviteRow{}, fmt.Errorf("building select query for %s: %w", OrganizationInviteTable, err)
	}

	return scanOrganizationInvite(q.db.QueryRow(ctx, query, args...))
}

func (q orgInvites) Select(ctx context.Context) ([]invite.OrgInviteRow, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", OrganizationInviteTable, err)
	}

	rows, err := q.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing select query for %s: %w", OrganizationInviteTable, err)
	}
	defer rows.Close()

	var out []invite.OrgInviteRow
	for rows.Next() {
		r, err := scanOrganizationInvite(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (q orgInvites) Exists(ctx context.Context) (bool, error) {
	subSQL, subArgs, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return false, err
	}
	sql := "SELECT EXISTS (" + subSQL + ")"

	var exists bool
	if err = q.db.QueryRow(ctx, sql, subArgs...).Scan(&exists); err != nil {
		return false, fmt.Errorf("executing exists query for %s: %w", OrganizationInviteTable, err)
	}
	return exists, nil
}

func (q orgInvites) FilterByID(id uuid.UUID) invite.OrgInvitesQ {
	q.selector = q.selector.Where(sq.Eq{"id": id})
	q.counter = q.counter.Where(sq.Eq{"id": id})
	q.updater = q.updater.Where(sq.Eq{"id": id})
	q.deleter = q.deleter.Where(sq.Eq{"id": id})
	return q
}

func (q orgInvites) FilterByOrganizationID(id uuid.UUID) invite.OrgInvitesQ {
	q.selector = q.selector.Where(sq.Eq{"organization_id": id})
	q.counter = q.counter.Where(sq.Eq{"organization_id": id})
	q.updater = q.updater.Where(sq.Eq{"organization_id": id})
	q.deleter = q.deleter.Where(sq.Eq{"organization_id": id})
	return q
}

func (q orgInvites) FilterByAccountID(id uuid.UUID) invite.OrgInvitesQ {
	q.selector = q.selector.Where(sq.Eq{"account_id": id})
	q.counter = q.counter.Where(sq.Eq{"account_id": id})
	q.updater = q.updater.Where(sq.Eq{"account_id": id})
	q.deleter = q.deleter.Where(sq.Eq{"account_id": id})
	return q
}

func (q orgInvites) FilterByStatus(status string) invite.OrgInvitesQ {
	q.selector = q.selector.Where(sq.Eq{"status": status})
	q.counter = q.counter.Where(sq.Eq{"status": status})
	q.updater = q.updater.Where(sq.Eq{"status": status})
	q.deleter = q.deleter.Where(sq.Eq{"status": status})
	return q
}

func (q orgInvites) FilterExpiresBefore(t time.Time) invite.OrgInvitesQ {
	q.selector = q.selector.Where(sq.Lt{"expires_at": t})
	q.counter = q.counter.Where(sq.Lt{"expires_at": t})
	q.updater = q.updater.Where(sq.Lt{"expires_at": t})
	q.deleter = q.deleter.Where(sq.Lt{"expires_at": t})
	return q
}

func (q orgInvites) FilterExpiresAfter(t time.Time) invite.OrgInvitesQ {
	q.selector = q.selector.Where(sq.GtOrEq{"expires_at": t})
	q.counter = q.counter.Where(sq.GtOrEq{"expires_at": t})
	q.updater = q.updater.Where(sq.GtOrEq{"expires_at": t})
	q.deleter = q.deleter.Where(sq.GtOrEq{"expires_at": t})
	return q
}

func (q orgInvites) UpdateOne(ctx context.Context) (invite.OrgInviteRow, error) {
	query, args, err := q.updater.
		Set("updated_at", time.Now().UTC()).
		Suffix("RETURNING " + OrganizationInviteColumns).ToSql()
	if err != nil {
		return invite.OrgInviteRow{}, fmt.Errorf("building update query for %s: %w", OrganizationInviteTable, err)
	}

	return scanOrganizationInvite(q.db.QueryRow(ctx, query, args...))
}

func (q orgInvites) UpdateStatus(status string) invite.OrgInvitesQ {
	q.updater = q.updater.Set("status", status)
	return q
}

func (q orgInvites) UpdateExpiresAt(t time.Time) invite.OrgInvitesQ {
	q.updater = q.updater.Set("expires_at", t)
	return q
}

func (q orgInvites) Count(ctx context.Context) (uint, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for %s: %w", OrganizationInviteTable, err)
	}

	var n uint
	if err = q.db.QueryRow(ctx, query, args...).Scan(&n); err != nil {
		return 0, fmt.Errorf("scanning count for %s: %w", OrganizationInviteTable, err)
	}
	return n, nil
}

func (q orgInvites) Page(limit uint, offset uint) invite.OrgInvitesQ {
	q.selector = q.selector.Limit(uint64(limit)).Offset(uint64(offset))
	return q
}

func (q orgInvites) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", OrganizationInviteTable, err)
	}

	if _, err = q.db.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("executing delete query for %s: %w", OrganizationInviteTable, err)
	}
	return nil
}
