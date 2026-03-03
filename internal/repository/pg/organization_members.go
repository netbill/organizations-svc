package pg

import (
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/netbill/organizations-svc/internal/repository"
	"github.com/netbill/pgdbx"
)

const organizationMembersTable = "organization_members"
const organizationMemberColumns = "id, account_id, organization_id, head, position, label, version, created_at, updated_at"
const organizationMemberColumnsM = "m.id, m.account_id, m.organization_id, m.head, m.position, m.label, m.version, m.created_at, m.updated_at"

func scanOrganizationMember(row sq.RowScanner) (m repository.OrganizationMemberRow, err error) {
	position := pgtype.Text{}
	label := pgtype.Text{}

	err = row.Scan(
		&m.ID,
		&m.AccountID,
		&m.OrganizationID,
		&m.Head,
		&position,
		&label,
		&m.Version,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return repository.OrganizationMemberRow{}, nil
	case err != nil:
		return repository.OrganizationMemberRow{}, fmt.Errorf("scanning member: %w", err)
	}

	if position.Valid {
		m.Position = &position.String
	}
	if label.Valid {
		m.Label = &label.String
	}

	return m, nil
}

type orgMembers struct {
	db       *pgdbx.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewOrgMembersQ(db *pgdbx.DB) repository.OrgMembersQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	return &orgMembers{
		db:       db,
		selector: b.Select(organizationMemberColumnsM).From(organizationMembersTable + " m"),
		inserter: b.Insert(organizationMembersTable),
		updater:  b.Update(organizationMembersTable + " m"),
		deleter:  b.Delete(organizationMembersTable + " m"),
		counter:  b.Select("COUNT(*)").From(organizationMembersTable + " m"),
	}
}

func (q *orgMembers) New() repository.OrgMembersQ {
	return NewOrgMembersQ(q.db)
}

func (q *orgMembers) Insert(ctx context.Context, data repository.OrganizationMemberRow) (repository.OrganizationMemberRow, error) {
	query, args, err := q.inserter.SetMap(map[string]any{
		"account_id":      data.AccountID,
		"organization_id": data.OrganizationID,
		"head":            data.Head,
		"position":        data.Position,
		"label":           data.Label,
	}).Suffix("RETURNING " + organizationMemberColumns).ToSql()
	if err != nil {
		return repository.OrganizationMemberRow{}, fmt.Errorf("building insert query for %s: %w", organizationMembersTable, err)
	}

	return scanOrganizationMember(q.db.QueryRow(ctx, query, args...))
}

func (q *orgMembers) Exists(ctx context.Context) (bool, error) {
	existsQ := q.selector.
		Columns("1").
		RemoveLimit().
		RemoveOffset().
		Prefix("SELECT EXISTS (").
		Suffix(") AS exists").
		Limit(1)

	query, args, err := existsQ.ToSql()
	if err != nil {
		return false, fmt.Errorf("building exists query for %s: %w", organizationMembersTable, err)
	}

	var ok bool
	if err = q.db.QueryRow(ctx, query, args...).Scan(&ok); err != nil {
		return false, fmt.Errorf("scanning exists for %s: %w", organizationMembersTable, err)
	}

	return ok, nil
}

func (q *orgMembers) Get(ctx context.Context) (repository.OrganizationMemberRow, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return repository.OrganizationMemberRow{}, fmt.Errorf("building select query for %s: %w", organizationMembersTable, err)
	}

	return scanOrganizationMember(q.db.QueryRow(ctx, query, args...))
}

func (q *orgMembers) Select(ctx context.Context) ([]repository.OrganizationMemberRow, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", organizationMembersTable, err)
	}

	rows, err := q.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing select query for %s: %w", organizationMembersTable, err)
	}
	defer rows.Close()

	out := make([]repository.OrganizationMemberRow, 0)
	for rows.Next() {
		m, err := scanOrganizationMember(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, m)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (q *orgMembers) FilterByID(id uuid.UUID) repository.OrgMembersQ {
	q.selector = q.selector.Where(sq.Eq{"m.id": id})
	q.counter = q.counter.Where(sq.Eq{"m.id": id})
	q.updater = q.updater.Where(sq.Eq{"m.id": id})
	q.deleter = q.deleter.Where(sq.Eq{"m.id": id})
	return q
}

func (q *orgMembers) FilterByAccountID(accountID uuid.UUID) repository.OrgMembersQ {
	q.selector = q.selector.Where(sq.Eq{"m.account_id": accountID})
	q.counter = q.counter.Where(sq.Eq{"m.account_id": accountID})
	q.updater = q.updater.Where(sq.Eq{"m.account_id": accountID})
	q.deleter = q.deleter.Where(sq.Eq{"m.account_id": accountID})
	return q
}

func (q *orgMembers) FilterByOrganizationID(organizationID ...uuid.UUID) repository.OrgMembersQ {
	q.selector = q.selector.Where(sq.Eq{"m.organization_id": organizationID})
	q.counter = q.counter.Where(sq.Eq{"m.organization_id": organizationID})
	q.updater = q.updater.Where(sq.Eq{"m.organization_id": organizationID})
	q.deleter = q.deleter.Where(sq.Eq{"m.organization_id": organizationID})
	return q
}

func (q *orgMembers) FilterLikePosition(position string) repository.OrgMembersQ {
	q.selector = q.selector.Where(sq.ILike{"m.position": "%" + position + "%"})
	q.counter = q.counter.Where(sq.ILike{"m.position": "%" + position + "%"})
	q.updater = q.updater.Where(sq.ILike{"m.position": "%" + position + "%"})
	q.deleter = q.deleter.Where(sq.ILike{"m.position": "%" + position + "%"})
	return q
}

func (q *orgMembers) FilterLikeLabel(label string) repository.OrgMembersQ {
	q.selector = q.selector.Where(sq.ILike{"m.label": "%" + label + "%"})
	q.counter = q.counter.Where(sq.ILike{"m.label": "%" + label + "%"})
	q.updater = q.updater.Where(sq.ILike{"m.label": "%" + label + "%"})
	q.deleter = q.deleter.Where(sq.ILike{"m.label": "%" + label + "%"})
	return q
}

func (q *orgMembers) FilterByHead(head bool) repository.OrgMembersQ {
	q.selector = q.selector.Where(sq.Eq{"m.head": head})
	q.counter = q.counter.Where(sq.Eq{"m.head": head})
	q.updater = q.updater.Where(sq.Eq{"m.head": head})
	q.deleter = q.deleter.Where(sq.Eq{"m.head": head})
	return q
}

func (q *orgMembers) FilterByUsername(username string) repository.OrgMembersQ {
	q.selector = q.selector.Where(sq.Eq{"p.username": username})
	q.counter = q.counter.Where(sq.Eq{"p.username": username})
	return q
}

func (q *orgMembers) FilterLikeUsername(username string) repository.OrgMembersQ {
	q.selector = q.selector.Where(sq.ILike{"p.username": "%" + username + "%"})
	q.counter = q.counter.Where(sq.ILike{"p.username": "%" + username + "%"})
	return q
}

func (q *orgMembers) FilterLikePseudonym(pseudonym string) repository.OrgMembersQ {
	q.selector = q.selector.Where(sq.ILike{"p.pseudonym": "%" + pseudonym + "%"})
	q.counter = q.counter.Where(sq.ILike{"p.pseudonym": "%" + pseudonym + "%"})
	return q
}

func (q *orgMembers) FilterBestMatch(term string) repository.OrgMembersQ {
	like := "%" + term + "%"
	prefix := term + "%"

	q.selector = q.selector.Where(sq.Or{
		sq.ILike{"p.username": like},
		sq.ILike{"p.pseudonym": like},
	})
	q.counter = q.counter.Where(sq.Or{
		sq.ILike{"p.username": like},
		sq.ILike{"p.pseudonym": like},
	})

	q.selector = q.selector.OrderByClause(sq.Expr(
		`CASE
			WHEN lower(p.username) = lower(?) THEN 0
			WHEN lower(p.pseudonym) = lower(?) THEN 1
			WHEN lower(p.username) LIKE lower(?) THEN 2
			WHEN lower(p.pseudonym) LIKE lower(?) THEN 3
			WHEN lower(p.username) LIKE lower(?) THEN 4
			WHEN lower(p.pseudonym) LIKE lower(?) THEN 5
			ELSE 6
		END`,
		term, term,
		prefix, prefix,
		like, like,
	))

	q.selector = q.selector.OrderBy("p.username ASC", "m.id ASC")
	return q
}

func (q *orgMembers) UpdateOne(ctx context.Context) (repository.OrganizationMemberRow, error) {
	q.updater = q.updater.
		Set("updated_at", time.Now().UTC()).
		Set("version", sq.Expr("version + 1"))

	query, args, err := q.updater.Suffix("RETURNING " + organizationMemberColumns).ToSql()
	if err != nil {
		return repository.OrganizationMemberRow{}, fmt.Errorf("building update query for %s: %w", organizationMembersTable, err)
	}

	return scanOrganizationMember(q.db.QueryRow(ctx, query, args...))
}

func (q *orgMembers) UpdatePosition(v *string) repository.OrgMembersQ {
	q.updater = q.updater.Set("position", v)
	return q
}

func (q *orgMembers) UpdateLabel(v *string) repository.OrgMembersQ {
	q.updater = q.updater.Set("label", v)
	return q
}

func (q *orgMembers) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", organizationMembersTable, err)
	}

	_, err = q.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("executing delete query for %s: %w", organizationMembersTable, err)
	}

	return nil
}

func (q *orgMembers) Count(ctx context.Context) (uint, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for %s: %w", organizationMembersTable, err)
	}

	var count uint
	if err = q.db.QueryRow(ctx, query, args...).Scan(&count); err != nil {
		return 0, fmt.Errorf("scanning count for %s: %w", organizationMembersTable, err)
	}

	return count, nil
}

func (q *orgMembers) Page(limit uint, offset uint) repository.OrgMembersQ {
	q.selector = q.selector.Limit(uint64(limit)).Offset(uint64(offset))
	return q
}
