package pg

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/netbill/organizations-svc/internal/repository"
	"github.com/netbill/pgdbx"
)

const organizationMemberRoleTable = "organization_member_roles"
const organizationMemberRoleColumns = "member_id, role_id, created_at"

func scanOrganizationMemberRole(row sq.RowScanner) (r repository.OrganizationMemberRolesRow, err error) {
	err = row.Scan(&r.MemberID, &r.RoleID, &r.CreatedAt)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return repository.OrganizationMemberRolesRow{}, nil
	case err != nil:
		return repository.OrganizationMemberRolesRow{}, fmt.Errorf("scanning member_role: %w", err)
	}

	return r, nil
}

type orgMemberRoles struct {
	db       *pgdbx.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewOrgMemberRolesQ(db *pgdbx.DB) repository.OrgMemberRolesQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &orgMemberRoles{
		db:       db,
		selector: b.Select(organizationMemberRoleColumns).From(organizationMemberRoleTable),
		inserter: b.Insert(organizationMemberRoleTable),
		deleter:  b.Delete(organizationMemberRoleTable),
		counter:  b.Select("COUNT(*)").From(organizationMemberRoleTable),
	}
}

func (q *orgMemberRoles) New() repository.OrgMemberRolesQ {
	return NewOrgMemberRolesQ(q.db)
}

func (q *orgMemberRoles) Insert(ctx context.Context, data repository.OrganizationMemberRolesRow) (repository.OrganizationMemberRolesRow, error) {
	query, args, err := q.inserter.SetMap(map[string]any{
		"member_id": data.MemberID,
		"role_id":   data.RoleID,
	}).Suffix("RETURNING " + organizationMemberRoleColumns).ToSql()
	if err != nil {
		return repository.OrganizationMemberRolesRow{}, fmt.Errorf("building insert query for %s: %w", organizationMemberRoleTable, err)
	}

	return scanOrganizationMemberRole(q.db.QueryRow(ctx, query, args...))
}

func (q *orgMemberRoles) Get(ctx context.Context) (repository.OrganizationMemberRolesRow, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return repository.OrganizationMemberRolesRow{}, fmt.Errorf("building select query for %s: %w", organizationMemberRoleTable, err)
	}

	return scanOrganizationMemberRole(q.db.QueryRow(ctx, query, args...))
}

func (q *orgMemberRoles) Select(ctx context.Context) ([]repository.OrganizationMemberRolesRow, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", organizationMemberRoleTable, err)
	}

	rows, err := q.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing select query for %s: %w", organizationMemberRoleTable, err)
	}
	defer rows.Close()

	out := make([]repository.OrganizationMemberRolesRow, 0)
	for rows.Next() {
		r, err := scanOrganizationMemberRole(rows)
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

func (q *orgMemberRoles) FilterByMemberID(memberID uuid.UUID) repository.OrgMemberRolesQ {
	q.selector = q.selector.Where(sq.Eq{"member_id": memberID})
	q.counter = q.counter.Where(sq.Eq{"member_id": memberID})
	q.deleter = q.deleter.Where(sq.Eq{"member_id": memberID})
	return q
}

func (q *orgMemberRoles) FilterByRoleID(roleID uuid.UUID) repository.OrgMemberRolesQ {
	q.selector = q.selector.Where(sq.Eq{"role_id": roleID})
	q.counter = q.counter.Where(sq.Eq{"role_id": roleID})
	q.deleter = q.deleter.Where(sq.Eq{"role_id": roleID})
	return q
}

func (q *orgMemberRoles) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", organizationMemberRoleTable, err)
	}

	if _, err = q.db.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("executing delete query for %s: %w", organizationMemberRoleTable, err)
	}
	return nil
}

func (q *orgMemberRoles) Count(ctx context.Context) (uint, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for %s: %w", organizationMemberRoleTable, err)
	}

	var n uint
	if err = q.db.QueryRow(ctx, query, args...).Scan(&n); err != nil {
		return 0, fmt.Errorf("scanning count for %s: %w", organizationMemberRoleTable, err)
	}
	return n, nil
}

func (q *orgMemberRoles) Page(limit uint, offset uint) repository.OrgMemberRolesQ {
	q.selector = q.selector.Limit(uint64(limit)).Offset(uint64(offset))
	return q
}
