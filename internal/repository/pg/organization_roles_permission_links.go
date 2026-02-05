package pg

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/netbill/pgdbx"

	"github.com/netbill/organizations-svc/internal/repository"
)

const organizationRolePermissionLinksTable = "organization_role_permission_links"
const organizationRolePermissionLinksColumns = "role_id, permission_code, created_at"

func scanOrganizationRolePermissionLink(row sq.RowScanner) (rp repository.OrganizationRolePermissionLinkRow, err error) {
	err = row.Scan(&rp.RoleID, &rp.PermissionCode, &rp.CreatedAt)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return repository.OrganizationRolePermissionLinkRow{}, nil
	case err != nil:
		return repository.OrganizationRolePermissionLinkRow{}, fmt.Errorf(
			"scanning organization role permission link: %w",
			err,
		)
	}
	return rp, nil
}

type orgRolePermissionLinks struct {
	db       *pgdbx.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewOrgRolePermissionLinksQ(db *pgdbx.DB) repository.OrgRolePermissionLinksQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &orgRolePermissionLinks{
		db:       db,
		selector: b.Select(organizationRolePermissionLinksColumns).From(organizationRolePermissionLinksTable),
		inserter: b.Insert(organizationRolePermissionLinksTable),
		deleter:  b.Delete(organizationRolePermissionLinksTable),
		counter:  b.Select("COUNT(*) AS count").From(organizationRolePermissionLinksTable),
	}
}

func (q *orgRolePermissionLinks) New() repository.OrgRolePermissionLinksQ {
	return NewOrgRolePermissionLinksQ(q.db)
}

func (q *orgRolePermissionLinks) Insert(
	ctx context.Context,
	roleID uuid.UUID,
	code ...string,
) ([]repository.OrganizationRolePermissionLinkRow, error) {
	codes := make([]string, 0, len(code))
	seen := make(map[string]struct{}, len(code))

	for _, c := range code {
		if c == "" {
			continue
		}
		if _, ok := seen[c]; ok {
			continue
		}
		seen[c] = struct{}{}
		codes = append(codes, c)
	}

	const sqlq = `
		WITH del AS (
			DELETE FROM organization_role_permission_links
			WHERE role_id = $1
		)
		INSERT INTO organization_role_permission_links (role_id, permission_code)
		SELECT $1, x.code
		FROM UNNEST($2::text[]) AS x(code)
		RETURNING role_id, permission_code, created_at
	`

	rows, err := q.db.Query(ctx, sqlq, roleID, codes)
	if err != nil {
		return nil, fmt.Errorf("upsert role permission links: %w", err)
	}
	defer rows.Close()

	out := make([]repository.OrganizationRolePermissionLinkRow, 0, len(codes))
	for rows.Next() {
		var r repository.OrganizationRolePermissionLinkRow
		if err := rows.Scan(&r.RoleID, &r.PermissionCode, &r.CreatedAt); err != nil {
			return nil, fmt.Errorf("scanning role permission link: %w", err)
		}
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (q *orgRolePermissionLinks) Get(ctx context.Context) (repository.OrganizationRolePermissionLinkRow, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return repository.OrganizationRolePermissionLinkRow{}, fmt.Errorf(
			"building select query for %s: %w",
			organizationRolePermissionLinksTable,
			err,
		)
	}

	return scanOrganizationRolePermissionLink(q.db.QueryRow(ctx, query, args...))
}

func (q *orgRolePermissionLinks) Select(
	ctx context.Context,
) ([]repository.OrganizationRolePermissionLinkRow, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf(
			"building select query for %s: %w",
			organizationRolePermissionLinksTable,
			err,
		)
	}

	rows, err := q.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf(
			"executing select query for %s: %w",
			organizationRolePermissionLinksTable,
			err,
		)
	}
	defer rows.Close()

	out := make([]repository.OrganizationRolePermissionLinkRow, 0)
	for rows.Next() {
		p, err := scanOrganizationRolePermissionLink(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (q *orgRolePermissionLinks) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf(
			"building delete query for %s: %w",
			organizationRolePermissionLinksTable,
			err,
		)
	}

	if _, err := q.db.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf(
			"executing delete query for %s: %w",
			organizationRolePermissionLinksTable,
			err,
		)
	}

	return nil
}

func (q *orgRolePermissionLinks) FilterByRoleID(roleID uuid.UUID) repository.OrgRolePermissionLinksQ {
	q.selector = q.selector.Where(sq.Eq{"role_id": roleID})
	q.deleter = q.deleter.Where(sq.Eq{"role_id": roleID})
	q.counter = q.counter.Where(sq.Eq{"role_id": roleID})
	return q
}

func (q *orgRolePermissionLinks) FilterByPermissionCode(
	code ...string,
) repository.OrgRolePermissionLinksQ {
	q.selector = q.selector.Where(sq.Eq{"permission_code": code})
	q.deleter = q.deleter.Where(sq.Eq{"permission_code": code})
	q.counter = q.counter.Where(sq.Eq{"permission_code": code})
	return q
}

func (q *orgRolePermissionLinks) FilterByAccountID(
	accountID uuid.UUID,
) repository.OrgRolePermissionLinksQ {
	sub := sq.
		Select("DISTINCT mr.role_id").
		From("organization_members m").
		Join("organization_member_roles mr ON mr.member_id = m.id").
		Where(sq.Eq{"m.account_id": accountID})

	subSQL, subArgs, err := sub.ToSql()
	if err != nil {
		ex := sq.Expr("1=0")
		q.selector = q.selector.Where(ex)
		q.deleter = q.deleter.Where(ex)
		q.counter = q.counter.Where(ex)
		return q
	}

	expr := sq.Expr("role_id IN ("+subSQL+")", subArgs...)
	q.selector = q.selector.Where(expr)
	q.deleter = q.deleter.Where(expr)
	q.counter = q.counter.Where(expr)
	return q
}

func (q *orgRolePermissionLinks) FilterByOrganizationID(
	organizationID uuid.UUID,
) repository.OrgRolePermissionLinksQ {
	sub := sq.
		Select("r.id").
		From("organization_roles r").
		Where(sq.Eq{"r.organization_id": organizationID})

	subSQL, subArgs, err := sub.ToSql()
	if err != nil {
		ex := sq.Expr("1=0")
		q.selector = q.selector.Where(ex)
		q.deleter = q.deleter.Where(ex)
		q.counter = q.counter.Where(ex)
		return q
	}

	expr := sq.Expr("role_id IN ("+subSQL+")", subArgs...)
	q.selector = q.selector.Where(expr)
	q.deleter = q.deleter.Where(expr)
	q.counter = q.counter.Where(expr)
	return q
}

func (q *orgRolePermissionLinks) FilterByMemberID(
	memberID uuid.UUID,
) repository.OrgRolePermissionLinksQ {
	sub := sq.
		Select("mr.role_id").
		From("organization_member_roles mr").
		Where(sq.Eq{"mr.member_id": memberID})

	subSQL, subArgs, err := sub.ToSql()
	if err != nil {
		ex := sq.Expr("1=0")
		q.selector = q.selector.Where(ex)
		q.deleter = q.deleter.Where(ex)
		q.counter = q.counter.Where(ex)
		return q
	}

	expr := sq.Expr("role_id IN ("+subSQL+")", subArgs...)
	q.selector = q.selector.Where(expr)
	q.deleter = q.deleter.Where(expr)
	q.counter = q.counter.Where(expr)
	return q
}

func (q *orgRolePermissionLinks) Count(ctx context.Context) (uint, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf(
			"building count query for %s: %w",
			organizationRolePermissionLinksTable,
			err,
		)
	}

	var n uint
	if err = q.db.QueryRow(ctx, query, args...).Scan(&n); err != nil {
		return 0, fmt.Errorf(
			"scanning count for %s: %w",
			organizationRolePermissionLinksTable,
			err,
		)
	}
	return n, nil
}

func (q *orgRolePermissionLinks) Page(
	limit, offset uint,
) repository.OrgRolePermissionLinksQ {
	q.selector = q.selector.Limit(uint64(limit)).Offset(uint64(offset))
	return q
}

func (q *orgRolePermissionLinks) Exists(ctx context.Context) (bool, error) {
	query, args, err := q.selector.
		Columns("1").
		Limit(1).
		ToSql()
	if err != nil {
		return false, err
	}

	var one int
	if err = q.db.QueryRow(ctx, query, args...).Scan(&one); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
