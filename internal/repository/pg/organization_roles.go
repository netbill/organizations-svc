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

const organizationRoleTable = "organization_roles"

const organizationRoleColumns = "id, organization_id, rank, name, description, color, created_at, updated_at"
const organizationRoleColumnsR = "r.id, r.organization_id, r.rank, r.name, r.description, r.color, r.created_at, r.updated_at"

func scanOrganizationRole(row sq.RowScanner) (r repository.OrganizationRoleRow, err error) {
	err = row.Scan(
		&r.ID,
		&r.OrganizationID,
		&r.Rank,
		&r.Name,
		&r.Description,
		&r.Color,
		&r.CreatedAt,
		&r.UpdatedAt,
	)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return repository.OrganizationRoleRow{}, nil
	case err != nil:
		return repository.OrganizationRoleRow{}, fmt.Errorf("scanning role: %w", err)
	}
	return r, nil
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
		selector: b.Select(organizationRoleColumnsR).From(organizationRoleTable + " r"),
		inserter: b.Insert(organizationRoleTable),
		updater:  b.Update(organizationRoleTable + " r"),
		deleter:  b.Delete(organizationRoleTable + " r"),
		counter:  b.Select("COUNT(*)").From(organizationRoleTable + " r"),
	}
}

func (q *orgRoles) New() repository.OrgRolesQ {
	return NewOrgRolesQ(q.db)
}

func (q *orgRoles) Insert(ctx context.Context, data repository.OrganizationRoleRow) (repository.OrganizationRoleRow, error) {
	const sqlInsertAtRank = `
		WITH bumped AS (
			UPDATE organization_roles
			SET
				rank = rank + 1,
				updated_at = now()
			WHERE organization_id = $1
			  AND rank >= $2
			RETURNING 1
		),
		ins AS (
			INSERT INTO organization_roles (organization_id, rank, name, description, color)
			VALUES ($1, $3, $2, $4, $5, $6)
			RETURNING id, organization_id, rank, name, description, color, created_at, updated_at
		)
		SELECT id, organization_id, rank, name, description, color, created_at, updated_at
		FROM ins;
	`

	args := []any{
		data.OrganizationID,
		data.Rank,
		data.Name,
		data.Description,
		data.Color,
	}

	return scanOrganizationRole(q.db.QueryRow(ctx, sqlInsertAtRank, args...))
}

func (q *orgRoles) Get(ctx context.Context) (repository.OrganizationRoleRow, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return repository.OrganizationRoleRow{}, fmt.Errorf("building select query for %s: %w", organizationRoleTable, err)
	}

	return scanOrganizationRole(q.db.QueryRow(ctx, query, args...))
}

func (q *orgRoles) Select(ctx context.Context) ([]repository.OrganizationRoleRow, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", organizationRoleTable, err)
	}

	rows, err := q.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing select query for %s: %w", organizationRoleTable, err)
	}
	defer rows.Close()

	out := make([]repository.OrganizationRoleRow, 0)
	for rows.Next() {
		r, err := scanOrganizationRole(rows)
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

func (q *orgRoles) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", organizationRoleTable, err)
	}

	if _, err = q.db.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("executing delete query for %s: %w", organizationRoleTable, err)
	}

	return nil
}

func (q *orgRoles) Count(ctx context.Context) (uint, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for %s: %w", organizationRoleTable, err)
	}

	var count uint
	if err = q.db.QueryRow(ctx, query, args...).Scan(&count); err != nil {
		return 0, fmt.Errorf("scanning count for %s: %w", organizationRoleTable, err)
	}

	return count, nil
}

func (q *orgRoles) UpdateOne(ctx context.Context) (repository.OrganizationRoleRow, error) {
	q.updater = q.updater.Set("updated_at", time.Now().UTC())

	query, args, err := q.updater.Suffix("RETURNING " + organizationRoleColumns).ToSql()
	if err != nil {
		return repository.OrganizationRoleRow{}, fmt.Errorf("building update query for %s: %w", organizationRoleTable, err)
	}

	return scanOrganizationRole(q.db.QueryRow(ctx, query, args...))
}

func (q *orgRoles) UpdateMany(ctx context.Context) (int64, error) {
	q.updater = q.updater.Set("updated_at", time.Now().UTC())

	query, args, err := q.updater.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building update query for %s: %w", organizationRoleTable, err)
	}

	res, err := q.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("executing update query for %s: %w", organizationRoleTable, err)
	}

	return res.RowsAffected(), nil
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

func (q *orgRoles) FilterByID(id ...uuid.UUID) repository.OrgRolesQ {
	q.selector = q.selector.Where(sq.Eq{"r.id": id})
	q.counter = q.counter.Where(sq.Eq{"r.id": id})
	q.updater = q.updater.Where(sq.Eq{"r.id": id})
	q.deleter = q.deleter.Where(sq.Eq{"r.id": id})
	return q
}

func (q *orgRoles) FilterByOrganizationID(id uuid.UUID) repository.OrgRolesQ {
	q.selector = q.selector.Where(sq.Eq{"r.organization_id": id})
	q.counter = q.counter.Where(sq.Eq{"r.organization_id": id})
	q.updater = q.updater.Where(sq.Eq{"r.organization_id": id})
	q.deleter = q.deleter.Where(sq.Eq{"r.organization_id": id})
	return q
}

func (q *orgRoles) FilterByAccountID(accountID uuid.UUID) repository.OrgRolesQ {
	sub := sq.
		Select("DISTINCT mr.role_id").
		From("organization_members m").
		Join("organization_member_roles mr ON mr.member_id = m.id").
		Where(sq.Eq{"m.account_id": accountID})

	subSQL, subArgs, err := sub.ToSql()
	if err != nil {
		ex := sq.Expr("1=0")
		q.selector = q.selector.Where(ex)
		q.counter = q.counter.Where(ex)
		q.updater = q.updater.Where(ex)
		q.deleter = q.deleter.Where(ex)
		return q
	}

	expr := sq.Expr("r.id IN ("+subSQL+")", subArgs...)
	q.selector = q.selector.Where(expr)
	q.counter = q.counter.Where(expr)
	q.updater = q.updater.Where(expr)
	q.deleter = q.deleter.Where(expr)
	return q
}

func (q *orgRoles) FilterByMemberID(memberID uuid.UUID) repository.OrgRolesQ {
	sub := sq.
		Select("mr.role_id").
		From("organization_member_roles mr").
		Where(sq.Eq{"mr.member_id": memberID})

	subSQL, subArgs, err := sub.ToSql()
	if err != nil {
		ex := sq.Expr("1=0")
		q.selector = q.selector.Where(ex)
		q.counter = q.counter.Where(ex)
		q.updater = q.updater.Where(ex)
		q.deleter = q.deleter.Where(ex)
		return q
	}

	expr := sq.Expr("r.id IN ("+subSQL+")", subArgs...)
	q.selector = q.selector.Where(expr)
	q.counter = q.counter.Where(expr)
	q.updater = q.updater.Where(expr)
	q.deleter = q.deleter.Where(expr)
	return q
}

func (q *orgRoles) FilterByRank(rank int) repository.OrgRolesQ {
	q.selector = q.selector.Where(sq.Eq{"r.rank": rank})
	q.counter = q.counter.Where(sq.Eq{"r.rank": rank})
	q.updater = q.updater.Where(sq.Eq{"r.rank": rank})
	q.deleter = q.deleter.Where(sq.Eq{"r.rank": rank})
	return q
}

func (q *orgRoles) FilterLikeName(name string) repository.OrgRolesQ {
	q.selector = q.selector.Where(sq.ILike{"r.name": "%" + name + "%"})
	q.counter = q.counter.Where(sq.ILike{"r.name": "%" + name + "%"})
	q.updater = q.updater.Where(sq.ILike{"r.name": "%" + name + "%"})
	q.deleter = q.deleter.Where(sq.ILike{"r.name": "%" + name + "%"})
	return q
}

func (q *orgRoles) OrderByRoleRank(asc bool) repository.OrgRolesQ {
	if asc {
		q.selector = q.selector.OrderBy("r.rank ASC", "r.id ASC")
	} else {
		q.selector = q.selector.OrderBy("r.rank DESC", "r.id DESC")
	}
	return q
}

func (q *orgRoles) Page(limit, offset uint) repository.OrgRolesQ {
	q.selector = q.selector.Limit(uint64(limit)).Offset(uint64(offset))
	return q
}

// ===== special rank methods =====

func (q *orgRoles) DeleteAndShiftRanks(ctx context.Context, roleID uuid.UUID) error {
	const query = `
		WITH del AS (
			DELETE FROM organization_roles
			WHERE id = $1
			RETURNING organization_id, rank
		)
		UPDATE organization_roles r
		SET rank = r.rank - 1,
		    updated_at = now()
		FROM del
		WHERE r.organization_id = del.organization_id
		  AND r.rank > del.rank
	`

	if _, err := q.db.Exec(ctx, query, roleID); err != nil {
		return fmt.Errorf("executing delete+shift for %s: %w", organizationRoleTable, err)
	}
	return nil
}

func (q *orgRoles) UpdateRoleRank(ctx context.Context, roleID uuid.UUID, newRank uint) (repository.OrganizationRoleRow, error) {
	var orgID uuid.UUID
	var oldRank int

	const sqlGet = `
		SELECT organization_id, rank
		FROM organization_roles
		WHERE id = $1
		LIMIT 1
	`
	if err := q.db.QueryRow(ctx, sqlGet, roleID).Scan(&orgID, &oldRank); err != nil {
		return repository.OrganizationRoleRow{}, fmt.Errorf("scanning role rank: %w", err)
	}

	if oldRank == int(newRank) {
		return NewOrgRolesQ(q.db).FilterByID(roleID).Get(ctx)
	}

	const sqlMove = `
		WITH upd AS (
			UPDATE organization_roles
			SET
				rank = CASE
					WHEN id = $1 THEN $2
					WHEN $2 < $3 AND rank >= $2 AND rank < $3 THEN rank + 1
					WHEN $2 > $3 AND rank <= $2 AND rank > $3 THEN rank - 1
					ELSE rank
				END,
				updated_at = now()
			WHERE organization_id = $4
			RETURNING id, organization_id, rank, name, description, color, created_at, updated_at
		)
		SELECT id, organization_id, rank, name, description, color, created_at, updated_at
		FROM upd
		WHERE id = $1
	`

	args := []any{roleID, int(newRank), oldRank, orgID}
	return scanOrganizationRole(q.db.QueryRow(ctx, sqlMove, args...))
}

func (q *orgRoles) UpdateRolesRanks(
	ctx context.Context,
	organizationID uuid.UUID,
	order map[uuid.UUID]uint,
) ([]repository.OrganizationRoleRow, error) {
	roles, err := NewOrgRolesQ(q.db).
		FilterByOrganizationID(organizationID).
		OrderByRoleRank(true).
		Select(ctx)
	if err != nil {
		return nil, fmt.Errorf("select roles by organization: %w", err)
	}
	if len(roles) == 0 {
		return nil, fmt.Errorf("no roles in organization %s", organizationID)
	}

	n := uint(len(roles))

	idToRole := make(map[uuid.UUID]repository.OrganizationRoleRow, n)
	for i := range roles {
		idToRole[roles[i].ID] = roles[i]
	}

	usedRank := make(map[uint]uuid.UUID, len(order))
	for roleID, newRank := range order {
		if newRank >= n {
			return nil, fmt.Errorf("rank %d out of range [0..%d]", newRank, n-1)
		}
		if _, ok := idToRole[roleID]; !ok {
			return nil, fmt.Errorf("role %s not in organization %s", roleID, organizationID)
		}
		if prev, ok := usedRank[newRank]; ok && prev != roleID {
			return nil, fmt.Errorf("duplicate rank %d for roles %s and %s", newRank, prev, roleID)
		}
		usedRank[newRank] = roleID
	}

	target := make([]uuid.UUID, n)
	filled := make([]bool, n)

	for rnk, id := range usedRank {
		target[rnk] = id
		filled[rnk] = true
	}

	rest := make([]uuid.UUID, 0, n-uint(len(order)))
	for i := range roles {
		id := roles[i].ID
		if _, ok := order[id]; ok {
			continue
		}
		rest = append(rest, id)
	}

	j := 0
	for i := 0; uint(i) < n; i++ {
		if filled[i] {
			continue
		}
		target[i] = rest[j]
		j++
	}

	changed := make([]uuid.UUID, 0, n)
	newRanks := make([]int32, 0, n)

	for newRank, id := range target {
		if roles[newRank].ID != id {
			changed = append(changed, id)
			newRanks = append(newRanks, int32(newRank))
		}
	}

	if len(changed) == 0 {
		return roles, nil
	}

	const sqlUpdate = `
		UPDATE organization_roles r
		SET
			rank = v.rank,
			updated_at = now()
		FROM (
			SELECT UNNEST($1::uuid[]) AS id, UNNEST($2::int4[]) AS rank
		) v
		WHERE r.id = v.id
		  AND r.organization_id = $3
		RETURNING r.id, r.organization_id, r.rank, r.name, r.description, r.color, r.created_at, r.updated_at
	`

	rows, err := q.db.Query(
		ctx,
		sqlUpdate,
		pgtype.FlatArray[uuid.UUID](changed),
		pgtype.FlatArray[int32](newRanks),
		organizationID,
	)

	if err != nil {
		return nil, fmt.Errorf("updating roles ranks: %w", err)
	}
	defer rows.Close()

	out := make([]repository.OrganizationRoleRow, 0, len(changed))
	for rows.Next() {
		r, err := scanOrganizationRole(rows)
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
