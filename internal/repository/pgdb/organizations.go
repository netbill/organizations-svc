package pgdb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/netbill/pgxtx"

	sq "github.com/Masterminds/squirrel"
)

const OrganizationTable = "organizations"
const OrganizationColumns = "id, status, name, icon, banner, max_roles, created_at, updated_at"

type Organization struct {
	ID     uuid.UUID `json:"id"`
	Status string    `json:"status"`
	Name   string    `json:"name"`

	Icon   pgtype.Text `json:"icon"`
	Banner pgtype.Text `json:"banner"`

	MaxRoles uint `json:"max_roles"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (a *Organization) scan(row sq.RowScanner) error {
	err := row.Scan(
		&a.ID,
		&a.Status,
		&a.Name,
		&a.Icon,
		&a.Banner,
		&a.MaxRoles,
		&a.CreatedAt,
		&a.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("scanning organization: %w", err)
	}
	return nil
}

type OrganizationsQ struct {
	db       pgxtx.DBTX
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewOrganizationsQ(db pgxtx.DBTX) OrganizationsQ {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return OrganizationsQ{
		db:       db,
		selector: builder.Select(OrganizationColumns).From(OrganizationTable),
		inserter: builder.Insert(OrganizationTable),
		updater:  builder.Update(OrganizationTable),
		deleter:  builder.Delete(OrganizationTable),
		counter:  builder.Select("COUNT(*)").From(OrganizationTable),
	}
}

type OrganizationsQInsertInput struct {
	Name string
}

func (q OrganizationsQ) Insert(ctx context.Context, data OrganizationsQInsertInput) (Organization, error) {
	query, args, err := q.inserter.SetMap(map[string]any{
		"name": data.Name,
	}).Suffix("RETURNING " + OrganizationColumns).ToSql()
	if err != nil {
		return Organization{}, fmt.Errorf("building insert query for %s: %w", OrganizationTable, err)
	}

	var inserted Organization
	err = inserted.scan(q.db.QueryRow(ctx, query, args...))
	if err != nil {
		return Organization{}, fmt.Errorf("executing insert query for %s: %w", OrganizationTable, err)
	}

	return inserted, nil
}

func (q OrganizationsQ) FilterByID(id uuid.UUID) OrganizationsQ {
	q.selector = q.selector.Where(sq.Eq{"id": id})
	q.counter = q.counter.Where(sq.Eq{"id": id})
	q.updater = q.updater.Where(sq.Eq{"id": id})
	q.deleter = q.deleter.Where(sq.Eq{"id": id})
	return q
}

func (q OrganizationsQ) FilterByStatus(status string) OrganizationsQ {
	q.selector = q.selector.Where(sq.Eq{"status": status})
	q.counter = q.counter.Where(sq.Eq{"status": status})
	q.updater = q.updater.Where(sq.Eq{"status": status})
	q.deleter = q.deleter.Where(sq.Eq{"status": status})
	return q
}

func (q OrganizationsQ) FilterByAccountID(accountID uuid.UUID) OrganizationsQ {
	sub := sq.
		Select("organization_id").
		From(OrganizationMembersTable).
		Where(sq.Eq{"account_id": accountID})

	subSQL, subArgs, err := sub.ToSql()
	if err != nil {
		q.selector = q.selector.Where(sq.Expr("1=0"))
		q.updater = q.updater.Where(sq.Expr("1=0"))
		q.deleter = q.deleter.Where(sq.Expr("1=0"))
		q.counter = q.counter.Where(sq.Expr("1=0"))
		return q
	}

	expr := sq.Expr("id IN ("+subSQL+")", subArgs...)

	q.selector = q.selector.Where(expr)
	q.updater = q.updater.Where(expr)
	q.deleter = q.deleter.Where(expr)
	q.counter = q.counter.Where(expr)

	return q
}

func (q OrganizationsQ) FilterNameLike(name string) OrganizationsQ {
	q.selector = q.selector.Where(sq.ILike{"name": "%" + name + "%"})
	q.counter = q.counter.Where(sq.ILike{"name": "%" + name + "%"})
	return q
}

func (q OrganizationsQ) Get(ctx context.Context) (Organization, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return Organization{}, fmt.Errorf("building select query for %s: %w", OrganizationTable, err)
	}

	var a Organization
	if err = a.scan(q.db.QueryRow(ctx, query, args...)); err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return Organization{}, pgx.ErrNoRows
		default:
			return Organization{}, err
		}
	}

	return a, nil
}

func (q OrganizationsQ) Select(ctx context.Context) ([]Organization, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", OrganizationTable, err)
	}

	rows, err := q.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing select query for %s: %w", OrganizationTable, err)
	}
	defer rows.Close()

	var organizations []Organization
	for rows.Next() {
		var organization Organization
		if err = organization.scan(rows); err != nil {
			return nil, err
		}
		organizations = append(organizations, organization)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return organizations, nil
}

func (q OrganizationsQ) UpdateOne(ctx context.Context) (Organization, error) {
	q.updater = q.updater.Set("updated_at", time.Now().UTC())

	query, args, err := q.updater.
		Suffix("RETURNING " + OrganizationColumns).
		ToSql()
	if err != nil {
		return Organization{}, fmt.Errorf("building update query for %s: %w", OrganizationTable, err)
	}

	var updated Organization
	if err = updated.scan(q.db.QueryRow(ctx, query, args...)); err != nil {
		return Organization{}, err
	}

	return updated, nil
}

func (q OrganizationsQ) UpdateMany(ctx context.Context) (int64, error) {
	q.updater = q.updater.Set("updated_at", time.Now().UTC())

	query, args, err := q.updater.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building update query for %s: %w", OrganizationTable, err)
	}

	res, err := q.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("executing update query for %s: %w", OrganizationTable, err)
	}

	return res.RowsAffected(), nil
}

func (q OrganizationsQ) UpdateName(name string) OrganizationsQ {
	q.updater = q.updater.Set("name", name)
	return q
}

func (q OrganizationsQ) UpdateStatus(status string) OrganizationsQ {
	q.updater = q.updater.Set("status", status)
	return q
}

func (q OrganizationsQ) UpdateIcon(icon pgtype.Text) OrganizationsQ {
	q.updater = q.updater.Set("icon", icon)
	return q
}

func (q OrganizationsQ) UpdateBanner(banner pgtype.Text) OrganizationsQ {
	q.updater = q.updater.Set("banner", banner)
	return q
}

func (q OrganizationsQ) UpdateMaxRoles(maxRoles uint) OrganizationsQ {
	q.updater = q.updater.Set("max_roles", maxRoles)
	return q
}

func (q OrganizationsQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", OrganizationTable, err)
	}

	if _, err = q.db.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("executing delete query for %s: %w", OrganizationTable, err)
	}

	return nil
}

func (q OrganizationsQ) Page(limit, offset uint) OrganizationsQ {
	q.selector = q.selector.Limit(uint64(limit)).Offset(uint64(offset))
	return q
}

func (q OrganizationsQ) Count(ctx context.Context) (uint, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for %s: %w", OrganizationTable, err)
	}

	var count uint
	if err = q.db.QueryRow(ctx, query, args...).Scan(&count); err != nil {
		return 0, fmt.Errorf("scanning count for %s: %w", OrganizationTable, err)
	}

	return count, nil
}
