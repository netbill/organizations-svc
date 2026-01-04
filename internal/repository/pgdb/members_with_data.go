package pgdb

import (
	"context"
	"database/sql"
	"fmt"

	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type MemberWithUserData struct {
	Member
	Username  string  `json:"username"`
	Official  bool    `json:"official"`
	Pseudonym *string `json:"pseudonym"`
}

func (mwd *MemberWithUserData) scan(row sq.RowScanner) error {
	err := row.Scan(
		&mwd.ID,
		&mwd.AccountID,
		&mwd.OrganizationID,
		&mwd.Position,
		&mwd.Label,
		&mwd.CreatedAt,
		&mwd.UpdatedAt,
		&mwd.Username,
		&mwd.Official,
		&mwd.Pseudonym,
	)
	if err != nil {
		return fmt.Errorf("scanning member with user data: %w", err)
	}
	return nil
}

func (q MembersQ) FilterByUsername(username string) MembersQ {
	q.selector = q.selector.Where(sq.Eq{"p.username": username})
	q.counter = q.counter.Where(sq.Eq{"p.username": username})
	return q
}

func (q MembersQ) FilterLikeUsername(username string) MembersQ {
	q.selector = q.selector.Where(sq.ILike{"p.username": "%" + username + "%"})
	q.counter = q.counter.Where(sq.ILike{"p.username": "%" + username + "%"})
	return q
}

func (q MembersQ) FilterLikePseudonym(pseudonym string) MembersQ {
	q.selector = q.selector.Where(sq.ILike{"p.pseudonym": "%" + pseudonym + "%"})
	q.counter = q.counter.Where(sq.ILike{"p.pseudonym": "%" + pseudonym + "%"})
	return q
}

func (q MembersQ) FilterBestMatch(term string) MembersQ {
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

func (q MembersQ) FilterRoleID(roleID uuid.UUID) MembersQ {
	query := sq.Expr(`
		EXISTS (
			SELECT 1
			FROM member_roles mr
			WHERE mr.member_id = m.id
				AND mr.role_id = ?
		)
	`, roleID)

	q.selector = q.selector.Where(query)
	q.counter = q.counter.Where(query)
	q.updater = q.updater.Where(query)
	q.deleter = q.deleter.Where(query)
	return q
}

func (q MembersQ) FilterByRoleRankUp(rankUp uint) MembersQ {
	query := sq.Expr(`
		EXISTS (
			SELECT 1
			FROM member_roles mr
			JOIN roles r ON r.id = mr.role_id
			WHERE mr.member_id = m.id
				AND r.rank >= ?
		)
	`, int(rankUp))

	q.selector = q.selector.Where(query)
	q.counter = q.counter.Where(query)
	q.updater = q.updater.Where(query)
	q.deleter = q.deleter.Where(query)
	return q
}

func (q MembersQ) FilterByRoleRankDown(rankDown uint) MembersQ {
	query := sq.Expr(`
		EXISTS (
			SELECT 1
			FROM member_roles mr
			JOIN roles r ON r.id = mr.role_id
			WHERE mr.member_id = m.id
				AND r.rank <= ?
		)
	`, int(rankDown))

	q.selector = q.selector.Where(query)
	q.counter = q.counter.Where(query)
	q.updater = q.updater.Where(query)
	q.deleter = q.deleter.Where(query)
	return q
}

func (q MembersQ) FilterByPermissionCode(code string) MembersQ {
	expr := sq.Expr(`
		EXISTS (
			SELECT 1
			FROM member_roles mr
			JOIN role_permissions rp ON rp.role_id = mr.role_id
			JOIN permissions perm ON perm.id = rp.permission_id
			WHERE mr.member_id = m.id
			  AND perm.code = ?
		)
	`, code)

	q.selector = q.selector.Where(expr)
	q.counter = q.counter.Where(expr)
	q.updater = q.updater.Where(expr)
	q.deleter = q.deleter.Where(expr)

	return q
}

func (q MembersQ) GetWithUserData(ctx context.Context) (MemberWithUserData, error) {
	q.selector = q.selector.
		Columns("p.username", "p.official", "p.pseudonym").
		Join("profiles p ON p.account_id = m.account_id")

	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return MemberWithUserData{}, fmt.Errorf("building select query for %s: %w", MembersTable, err)
	}

	var out MemberWithUserData
	if err = out.scan(q.db.QueryRowContext(ctx, query, args...)); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return MemberWithUserData{}, nil
		default:
			return MemberWithUserData{}, err
		}
	}

	return out, nil
}

func (q MembersQ) SelectWithUserData(ctx context.Context) ([]MemberWithUserData, error) {
	q.selector = q.selector.
		Columns("p.username", "p.official", "p.pseudonym").
		Join("profiles p ON p.account_id = m.account_id")

	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", MembersTable, err)
	}

	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing select query for %s: %w", MembersTable, err)
	}
	defer rows.Close()

	var out []MemberWithUserData
	for rows.Next() {
		var mwd MemberWithUserData
		if err = mwd.scan(rows); err != nil {
			return nil, err
		}
		out = append(out, mwd)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

type MemberRoleData struct {
	RoleID uuid.UUID `json:"role_id"`
	Head   bool      `json:"head"`
	Rank   uint      `json:"rank"`
	Name   string    `json:"name"`
}

type MemberWithRoleDataRow struct {
	MemberWithUserData
	roles []MemberRoleData
}

func (q MembersQ) SelectWithRolesData(ctx context.Context, roleLimit uint) ([]MemberWithRoleDataRow, error) {
	q.selector = q.selector.
		Columns("p.username", "p.official", "p.pseudonym").
		Join("profiles p ON p.account_id = m.account_id").
		LeftJoin("member_roles mr ON mr.member_id = m.id").
		LeftJoin("roles r ON r.id = mr.role_id").
		Columns("r.id", "r.head", "r.rank", "r.name").
		OrderBy("m.id ASC", "r.rank ASC", "r.id ASC")

	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", MembersTable, err)
	}

	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing select query for %s: %w", MembersTable, err)
	}
	defer rows.Close()

	out := make([]MemberWithRoleDataRow, 0)
	idx := make(map[uuid.UUID]int)

	for rows.Next() {
		var mwd MemberWithUserData

		var roleID uuid.NullUUID
		var head sql.NullBool
		var rank sql.NullInt64
		var name sql.NullString

		if err = rows.Scan(
			&mwd.ID,
			&mwd.AccountID,
			&mwd.OrganizationID,
			&mwd.Position,
			&mwd.Label,
			&mwd.CreatedAt,
			&mwd.UpdatedAt,
			&mwd.Username,
			&mwd.Official,
			&mwd.Pseudonym,
			&roleID,
			&head,
			&rank,
			&name,
		); err != nil {
			return nil, fmt.Errorf("scanning member with role data: %w", err)
		}

		i, ok := idx[mwd.ID]
		if !ok {
			out = append(out, MemberWithRoleDataRow{MemberWithUserData: mwd})
			i = len(out) - 1
			idx[mwd.ID] = i
		}

		if !roleID.Valid {
			continue
		}

		if roleLimit > 0 && uint(len(out[i].roles)) >= roleLimit {
			continue
		}

		out[i].roles = append(out[i].roles, MemberRoleData{
			RoleID: roleID.UUID,
			Head:   head.Valid && head.Bool,
			Rank:   uint(rank.Int64),
			Name:   name.String,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (q MembersQ) CanInteract(ctx context.Context, firstMemberID, secondMemberID uuid.UUID) (bool, error) {
	const sqlq = `
		SELECT
			(m1.organization_id = m2.organization_id)
			AND (COALESCE(r1.min_rank, 2147483647) < COALESCE(r2.min_rank, 2147483647)) AS can
		FROM members m1
		JOIN members m2 ON m2.id = $2
		LEFT JOIN LATERAL (
			SELECT MIN(r.rank) AS min_rank
			FROM member_roles mr
			JOIN roles r ON r.id = mr.role_id
			WHERE mr.member_id = m1.id
		) r1 ON true
		LEFT JOIN LATERAL (
			SELECT MIN(r.rank) AS min_rank
			FROM member_roles mr
			JOIN roles r ON r.id = mr.role_id
			WHERE mr.member_id = m2.id
		) r2 ON true
		WHERE m1.id = $1
		LIMIT 1
	`

	var ok bool
	if err := q.db.QueryRowContext(ctx, sqlq, firstMemberID, secondMemberID).Scan(&ok); err != nil {
		return false, fmt.Errorf("scanning can_interact: %w", err)
	}
	return ok, nil
}
