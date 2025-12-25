package pgdbsq

import (
	"context"
	"encoding/json"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type MemberRoleData struct {
	RoleID uuid.UUID `json:"role_id"`
	Head   bool      `json:"head"`
	Rank   uint      `json:"rank"`
	Name   string    `json:"name"`
}

type MemberWithUserData struct {
	Member
	Username  string           `json:"username"`
	Official  bool             `json:"official"`
	Pseudonym *string          `json:"pseudonym"`
	Roles     []MemberRoleData `json:"roles"`
}

func (mwd *MemberWithUserData) scan(row sq.RowScanner) error {
	var rolesRaw []byte

	err := row.Scan(
		&mwd.ID,
		&mwd.AccountID,
		&mwd.AgglomerationID,
		&mwd.Position,
		&mwd.Label,
		&mwd.CreatedAt,
		&mwd.UpdatedAt,
		&mwd.Username,
		&mwd.Official,
		&mwd.Pseudonym,
		&rolesRaw,
	)
	if err != nil {
		return fmt.Errorf("scanning member with user data: %w", err)
	}

	if len(rolesRaw) == 0 || string(rolesRaw) == "null" {
		mwd.Roles = nil
		return nil
	}

	if err = json.Unmarshal(rolesRaw, &mwd.Roles); err != nil {
		return fmt.Errorf("unmarshal roles: %w", err)
	}

	return nil
}

func (q MemberQ) WithUserData(roleLimit uint) MemberQ {
	q.selector = q.selector.
		Columns("p.username", "p.official", "p.pseudonym").
		Join("profiles p ON p.account_id = m.account_id").
		JoinClause(sq.Expr(
			`LEFT JOIN LATERAL (
				SELECT COALESCE(
					jsonb_agg(
						jsonb_build_object(
							'role_id', r.id,
							'head', r.head,
							'rank', r.rank,
							'name', r.name
						)
						ORDER BY r.rank ASC, r.id ASC
					),
					'[]'::jsonb
				) AS roles
				FROM member_roles mr
				JOIN roles r ON r.id = mr.role_id
				WHERE mr.member_id = m.id
				LIMIT ?
			) rr ON true`,
			uint64(roleLimit),
		)).
		Column("rr.roles")

	q.counter = q.counter.Join("profiles p ON p.account_id = m.account_id")

	return q
}

func (q MemberQ) FilterUsername(username string) MemberQ {
	q.selector = q.selector.Where(sq.Eq{"p.username": username})
	q.counter = q.counter.Where(sq.Eq{"p.username": username})
	return q
}

func (q MemberQ) FilterLikeUsername(username string) MemberQ {
	q.selector = q.selector.Where(sq.ILike{"p.username": "%" + username + "%"})
	q.counter = q.counter.Where(sq.ILike{"p.username": "%" + username + "%"})
	return q
}

func (q MemberQ) FilterLikePseudonym(pseudonym string) MemberQ {
	q.selector = q.selector.Where(sq.ILike{"p.pseudonym": "%" + pseudonym + "%"})
	q.counter = q.counter.Where(sq.ILike{"p.pseudonym": "%" + pseudonym + "%"})
	return q
}

func (q MemberQ) FilterBestMatch(term string) MemberQ {
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

func (q MemberQ) FilterRoleRank(rank uint) MemberQ {
	expr := sq.Expr(
		`EXISTS (
			SELECT 1
			FROM member_roles mr
			JOIN roles r ON r.id = mr.role_id
			WHERE mr.member_id = m.id
				AND r.rank = ?
		)`,
		int(rank),
	)

	q.selector = q.selector.Where(expr)
	q.counter = q.counter.Where(expr)
	q.updater = q.updater.Where(expr)
	q.deleter = q.deleter.Where(expr)
	return q
}

func (q MemberQ) GetWithUserData(ctx context.Context) (MemberWithUserData, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return MemberWithUserData{}, fmt.Errorf("building select query for %s: %w", MemberTable, err)
	}

	var out MemberWithUserData
	if err = out.scan(q.db.QueryRowContext(ctx, query, args...)); err != nil {
		return MemberWithUserData{}, err
	}

	return out, nil
}

func (q MemberQ) SelectWithUserData(ctx context.Context) ([]MemberWithUserData, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", MemberTable, err)
	}

	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing select query for %s: %w", MemberTable, err)
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
