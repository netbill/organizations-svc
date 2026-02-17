package pg

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/netbill/organizations-svc/internal/repository"
)

func scanOrganizationMemberWithProfileData(row sq.RowScanner) (m repository.OrganizationMemberWithProfileDataRow, err error) {
	position := pgtype.Text{}
	label := pgtype.Text{}
	pseudonym := pgtype.Text{}
	icon := pgtype.Text{}

	err = row.Scan(
		&m.ID,
		&m.AccountID,
		&m.OrganizationID,
		&m.Head,
		&position,
		&label,
		&m.CreatedAt,
		&m.UpdatedAt,
		&m.Username,
		&m.Official,
		&pseudonym,
		&icon,
	)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return repository.OrganizationMemberWithProfileDataRow{}, nil
	case err != nil:
		return repository.OrganizationMemberWithProfileDataRow{}, fmt.Errorf("scanning member with user data: %w", err)
	}

	if position.Valid {
		m.Position = &position.String
	}
	if label.Valid {
		m.Label = &label.String
	}
	if pseudonym.Valid {
		m.Pseudonym = &pseudonym.String
	}
	if icon.Valid {
		m.Icon = &icon.String
	}

	return m, nil
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

func (q *orgMembers) FilterRoleID(roleID uuid.UUID) repository.OrgMembersQ {
	expr := sq.Expr(`
		EXISTS (
			SELECT 1
			FROM member_roles mr
			WHERE mr.member_id = m.id
				AND mr.role_id = ?
		)
	`, roleID)

	q.selector = q.selector.Where(expr)
	q.counter = q.counter.Where(expr)
	q.updater = q.updater.Where(expr)
	q.deleter = q.deleter.Where(expr)
	return q
}

func (q *orgMembers) FilterByRoleRankUp(rankUp uint) repository.OrgMembersQ {
	expr := sq.Expr(`
		EXISTS (
			SELECT 1
			FROM organization_member_roles mr
			JOIN organization_roles r ON r.id = mr.role_id
			WHERE mr.member_id = m.id
				AND r.rank >= ?
		)
	`, int(rankUp))

	q.selector = q.selector.Where(expr)
	q.counter = q.counter.Where(expr)
	q.updater = q.updater.Where(expr)
	q.deleter = q.deleter.Where(expr)
	return q
}

func (q *orgMembers) FilterByRoleRankDown(rankDown uint) repository.OrgMembersQ {
	expr := sq.Expr(`
		EXISTS (
			SELECT 1
			FROM organization_member_roles mr
			JOIN organization_roles r ON r.id = mr.role_id
			WHERE mr.member_id = m.id
				AND r.rank <= ?
		)
	`, int(rankDown))

	q.selector = q.selector.Where(expr)
	q.counter = q.counter.Where(expr)
	q.updater = q.updater.Where(expr)
	q.deleter = q.deleter.Where(expr)
	return q
}

func (q *orgMembers) FilterByPermissionID(permissionID uuid.UUID) repository.OrgMembersQ {
	expr := sq.Expr(`
		EXISTS (
			SELECT 1
			FROM organization_member_roles mr
			JOIN organization_role_permission_links rp ON rp.role_id = mr.role_id
			WHERE mr.member_id = m.id
			  AND rp.permission_id = ?
		)
	`, permissionID)

	q.selector = q.selector.Where(expr)
	q.counter = q.counter.Where(expr)
	q.updater = q.updater.Where(expr)
	q.deleter = q.deleter.Where(expr)
	return q
}

func (q *orgMembers) GetWithUserData(ctx context.Context) (repository.OrganizationMemberWithProfileDataRow, error) {
	q.selector = q.selector.
		Columns("p.username", "p.official", "p.pseudonym", "p.icon").
		Join("profiles p ON p.account_id = m.account_id")

	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return repository.OrganizationMemberWithProfileDataRow{}, fmt.Errorf("building select query for %s: %w", organizationMembersTable, err)
	}

	return scanOrganizationMemberWithProfileData(q.db.QueryRow(ctx, query, args...))
}

func (q *orgMembers) SelectWithUserData(ctx context.Context) ([]repository.OrganizationMemberWithProfileDataRow, error) {
	q.selector = q.selector.
		Columns("p.username", "p.official", "p.pseudonym", "p.icon").
		Join("profiles p ON p.account_id = m.account_id")

	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", organizationMembersTable, err)
	}

	rows, err := q.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing select query for %s: %w", organizationMembersTable, err)
	}
	defer rows.Close()

	out := make([]repository.OrganizationMemberWithProfileDataRow, 0)
	for rows.Next() {
		m, err := scanOrganizationMemberWithProfileData(rows)
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
