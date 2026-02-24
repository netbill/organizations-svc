package pg

import (
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/netbill/organizations-svc/internal/repository"
	"github.com/netbill/pgdbx"
)

const organizationRolesTable = "organization_roles"

const organizationRolesColumns = "id, organization_id, name, description, color, version, created_at, updated_at"
const organizationRolesColumnsR = "r.id, r.organization_id, r.name, r.description, r.color, r.version, r.created_at, r.updated_at"

func scanOrgRole(row sq.RowScanner) (res repository.OrgRoleRow, err error) {
	err = row.Scan(
		&res.ID,
		&res.OrganizationID,
		&res.Name,
		&res.Description,
		&res.Color,
		&res.Version,
		&res.CreatedAt,
		&res.UpdatedAt,
	)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return repository.OrgRoleRow{}, nil
	case err != nil:
		return repository.OrgRoleRow{}, fmt.Errorf("scanning org role: %w", err)
	}
	return res, nil
}

type orgRoles struct {
	db       *pgdbx.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewOrgRolesQ(db *pgdbx.DB) repository.OrgRolesQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &orgRoles{
		db:       db,
		selector: b.Select(organizationRolesColumnsR).From(organizationRolesTable + " r"),
		inserter: b.Insert(organizationRolesTable),
		updater:  b.Update(organizationRolesTable + " r"),
		deleter:  b.Delete(organizationRolesTable + " r"),
		counter:  b.Select("COUNT(*)").From(organizationRolesTable + " r"),
	}
}

func (q *orgRoles) New() repository.OrgRolesQ {
	return NewOrgRolesQ(q.db)
}

func (q *orgRoles) Insert(ctx context.Context, input repository.OrgRoleRow) (repository.OrgRoleRow, error) {
	query, args, err := q.inserter.SetMap(map[string]any{
		"organization_id": input.OrganizationID,
		"name":            input.Name,
		"description":     input.Description,
		"color":           input.Color,
	}).Suffix("RETURNING " + organizationRolesColumns).ToSql()
	if err != nil {
		return repository.OrgRoleRow{}, fmt.Errorf("building insert query for %s: %w", organizationRolesTable, err)
	}

	return scanOrgRole(q.db.QueryRow(ctx, query, args...))
}

func (q *orgRoles) Get(ctx context.Context) (repository.OrgRoleRow, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return repository.OrgRoleRow{}, fmt.Errorf("building select query for %s: %w", organizationRolesTable, err)
	}

	return scanOrgRole(q.db.QueryRow(ctx, query, args...))
}

func (q *orgRoles) Select(ctx context.Context) ([]repository.OrgRoleRow, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", organizationRolesTable, err)
	}

	rows, err := q.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing select query for %s: %w", organizationRolesTable, err)
	}
	defer rows.Close()

	out := make([]repository.OrgRoleRow, 0)
	for rows.Next() {
		r, err := scanOrgRole(rows)
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

func (q *orgRoles) Exists(ctx context.Context) (bool, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return false, fmt.Errorf("building exists query for %s: %w", organizationRolesTable, err)
	}

	var count uint
	if err = q.db.QueryRow(ctx, query, args...).Scan(&count); err != nil {
		return false, fmt.Errorf("scanning exists for %s: %w", organizationRolesTable, err)
	}

	return count > 0, nil
}

func (q *orgRoles) UpdateOne(ctx context.Context) (repository.OrgRoleRow, error) {
	q.updater = q.updater.
		Set("updated_at", time.Now().UTC()).
		Set("version", sq.Expr("version + 1"))

	query, args, err := q.updater.Suffix("RETURNING " + organizationRolesColumns).ToSql()
	if err != nil {
		return repository.OrgRoleRow{}, fmt.Errorf("building update query for %s: %w", organizationRolesTable, err)
	}

	return scanOrgRole(q.db.QueryRow(ctx, query, args...))
}

func (q *orgRoles) UpdateName(name string) repository.OrgRolesQ {
	q.updater = q.updater.Set("name", name)
	return q
}

func (q *orgRoles) UpdateDescription(description string) repository.OrgRolesQ {
	q.updater = q.updater.Set("description", description)
	return q
}

func (q *orgRoles) UpdateColor(color string) repository.OrgRolesQ {
	q.updater = q.updater.Set("color", color)
	return q
}

func (q *orgRoles) FilterByID(roleID uuid.UUID) repository.OrgRolesQ {
	q.selector = q.selector.Where(sq.Eq{"r.id": roleID})
	q.counter = q.counter.Where(sq.Eq{"r.id": roleID})
	q.updater = q.updater.Where(sq.Eq{"r.id": roleID})
	q.deleter = q.deleter.Where(sq.Eq{"r.id": roleID})
	return q
}

func (q *orgRoles) FilterByOrganizationID(organizationID uuid.UUID) repository.OrgRolesQ {
	q.selector = q.selector.Where(sq.Eq{"r.organization_id": organizationID})
	q.counter = q.counter.Where(sq.Eq{"r.organization_id": organizationID})
	q.updater = q.updater.Where(sq.Eq{"r.organization_id": organizationID})
	q.deleter = q.deleter.Where(sq.Eq{"r.organization_id": organizationID})
	return q
}

func (q *orgRoles) OrderByRank(ask bool) repository.OrgRolesQ {
	dir := "DESC"
	if ask {
		dir = "ASC"
	}

	q.selector = q.selector.
		Join("organization_role_ranks rr ON rr.role_id = r.id").
		OrderBy("rr.rank " + dir)

	return q
}

func (q *orgRoles) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", organizationRolesTable, err)
	}

	if _, err = q.db.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("executing delete query for %s: %w", organizationRolesTable, err)
	}

	return nil
}
