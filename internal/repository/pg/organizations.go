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

const organizationsTable = "organizations"
const organizationsColumns = "id, status, name, icon_key, banner_key, max_roles, created_at, updated_at"

func scanOrganization(row sq.RowScanner) (o repository.OrganizationRow, err error) {
	var iconKey pgtype.Text
	var bannerKey pgtype.Text
	var maxRoles int32

	err = row.Scan(
		&o.ID,
		&o.Status,
		&o.Name,
		&iconKey,
		&bannerKey,
		&maxRoles,
		&o.CreatedAt,
		&o.UpdatedAt,
	)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return repository.OrganizationRow{}, nil
	case err != nil:
		return repository.OrganizationRow{}, fmt.Errorf("scanning organization: %w", err)
	}

	if iconKey.Valid {
		o.IconKey = &iconKey.String
	}
	if bannerKey.Valid {
		o.BannerKey = &bannerKey.String
	}

	// max_roles в БД обычно int, в модели тебе нужен uint
	if maxRoles < 0 {
		maxRoles = 0
	}
	o.MaxRoles = uint(maxRoles)

	return o, nil
}

type organizations struct {
	db       *pgdbx.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewOrganizationsQ(db *pgdbx.DB) repository.OrganizationsQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &organizations{
		db:       db,
		selector: b.Select(organizationsColumns).From(organizationsTable),
		inserter: b.Insert(organizationsTable),
		updater:  b.Update(organizationsTable),
		deleter:  b.Delete(organizationsTable),
		counter:  b.Select("COUNT(*)").From(organizationsTable),
	}
}

func (q *organizations) New() repository.OrganizationsQ {
	return NewOrganizationsQ(q.db)
}

func (q *organizations) Insert(ctx context.Context, data repository.OrganizationRow) (repository.OrganizationRow, error) {
	query, args, err := q.inserter.SetMap(map[string]any{
		"icon_key":   data.IconKey,
		"banner_key": data.BannerKey,
		"name":       data.Name,
	}).Suffix("RETURNING " + organizationsColumns).ToSql()
	if err != nil {
		return repository.OrganizationRow{}, fmt.Errorf("building insert query for %s: %w", organizationsTable, err)
	}

	return scanOrganization(q.db.QueryRow(ctx, query, args...))
}

func (q *organizations) FilterByID(id ...uuid.UUID) repository.OrganizationsQ {
	q.selector = q.selector.Where(sq.Eq{"id": id})
	q.counter = q.counter.Where(sq.Eq{"id": id})
	q.updater = q.updater.Where(sq.Eq{"id": id})
	q.deleter = q.deleter.Where(sq.Eq{"id": id})
	return q
}

func (q *organizations) FilterByStatus(status string) repository.OrganizationsQ {
	q.selector = q.selector.Where(sq.Eq{"status": status})
	q.counter = q.counter.Where(sq.Eq{"status": status})
	q.updater = q.updater.Where(sq.Eq{"status": status})
	q.deleter = q.deleter.Where(sq.Eq{"status": status})
	return q
}

func (q *organizations) FilterByAccountID(accountID uuid.UUID) repository.OrganizationsQ {
	sub := sq.
		Select("organization_id").
		From(organizationMembersTable).
		Where(sq.Eq{"account_id": accountID})

	subSQL, subArgs, err := sub.ToSql()
	if err != nil {
		ex := sq.Expr("1=0")
		q.selector = q.selector.Where(ex)
		q.updater = q.updater.Where(ex)
		q.deleter = q.deleter.Where(ex)
		q.counter = q.counter.Where(ex)
		return q
	}

	expr := sq.Expr("id IN ("+subSQL+")", subArgs...)

	q.selector = q.selector.Where(expr)
	q.updater = q.updater.Where(expr)
	q.deleter = q.deleter.Where(expr)
	q.counter = q.counter.Where(expr)

	return q
}

func (q *organizations) FilterNameLike(name string) repository.OrganizationsQ {
	q.selector = q.selector.Where(sq.ILike{"name": "%" + name + "%"})
	q.counter = q.counter.Where(sq.ILike{"name": "%" + name + "%"})
	return q
}

func (q *organizations) Get(ctx context.Context) (repository.OrganizationRow, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return repository.OrganizationRow{}, fmt.Errorf("building select query for %s: %w", organizationsTable, err)
	}

	return scanOrganization(q.db.QueryRow(ctx, query, args...))
}

func (q *organizations) Select(ctx context.Context) ([]repository.OrganizationRow, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", organizationsTable, err)
	}

	rows, err := q.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing select query for %s: %w", organizationsTable, err)
	}
	defer rows.Close()

	out := make([]repository.OrganizationRow, 0)
	for rows.Next() {
		o, err := scanOrganization(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, o)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (q *organizations) UpdateOne(ctx context.Context) (repository.OrganizationRow, error) {
	q.updater = q.updater.Set("updated_at", time.Now().UTC())

	query, args, err := q.updater.
		Suffix("RETURNING " + organizationsColumns).
		ToSql()
	if err != nil {
		return repository.OrganizationRow{}, fmt.Errorf("building update query for %s: %w", organizationsTable, err)
	}

	return scanOrganization(q.db.QueryRow(ctx, query, args...))
}

func (q *organizations) UpdateMany(ctx context.Context) (int64, error) {
	q.updater = q.updater.Set("updated_at", time.Now().UTC())

	query, args, err := q.updater.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building update query for %s: %w", organizationsTable, err)
	}

	res, err := q.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("executing update query for %s: %w", organizationsTable, err)
	}

	return res.RowsAffected(), nil
}

func (q *organizations) UpdateName(v string) repository.OrganizationsQ {
	q.updater = q.updater.Set("name", v)
	return q
}

func (q *organizations) UpdateStatus(v string) repository.OrganizationsQ {
	q.updater = q.updater.Set("status", v)
	return q
}

func (q *organizations) UpdateIconKey(v *string) repository.OrganizationsQ {
	q.updater = q.updater.Set("icon_key", v)
	return q
}

func (q *organizations) UpdateBannerKey(v *string) repository.OrganizationsQ {
	q.updater = q.updater.Set("banner_key", v)
	return q
}

func (q *organizations) UpdateMaxRoles(maxRoles uint) repository.OrganizationsQ {
	q.updater = q.updater.Set("max_roles", int32(maxRoles))
	return q
}

func (q *organizations) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", organizationsTable, err)
	}

	if _, err = q.db.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("executing delete query for %s: %w", organizationsTable, err)
	}

	return nil
}

func (q *organizations) Page(limit, offset uint) repository.OrganizationsQ {
	q.selector = q.selector.Limit(uint64(limit)).Offset(uint64(offset))
	return q
}

func (q *organizations) Count(ctx context.Context) (uint, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for %s: %w", organizationsTable, err)
	}

	var count uint
	if err = q.db.QueryRow(ctx, query, args...).Scan(&count); err != nil {
		return 0, fmt.Errorf("scanning count for %s: %w", organizationsTable, err)
	}

	return count, nil
}
