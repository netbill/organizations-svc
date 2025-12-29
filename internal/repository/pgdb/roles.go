package pgdb

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/umisto/pgx"

	sq "github.com/Masterminds/squirrel"
)

const RoleTable = "roles"

const RoleColumns = "id, agglomeration_id, head, rank, name, description, color, created_at, updated_at"
const RoleColumnsR = "r.id, r.agglomeration_id, r.head, r.rank, r.name, r.description, r.color, r.created_at, r.updated_at"

type Role struct {
	ID              uuid.UUID `json:"id"`
	AgglomerationID uuid.UUID `json:"agglomeration_id"`
	Head            bool      `json:"head"`
	Rank            uint      `json:"rank"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Color           string    `json:"color"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (r *Role) scan(row sq.RowScanner) error {
	err := row.Scan(
		&r.ID,
		&r.AgglomerationID,
		&r.Head,
		&r.Rank,
		&r.Name,
		&r.Description,
		&r.Color,
		&r.CreatedAt,
		&r.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("scanning role: %w", err)
	}
	return nil
}

type RolesQ struct {
	db       pgx.DBTX
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewRolesQ(db pgx.DBTX) RolesQ {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return RolesQ{
		db:       db,
		selector: builder.Select(RoleColumnsR).From(RoleTable + " r"),
		inserter: builder.Insert(RoleTable),
		updater:  builder.Update(RoleTable + " r"),
		deleter:  builder.Delete(RoleTable + " r"),
		counter:  builder.Select("COUNT(*)").From(RoleTable + " r"),
	}
}

type InsertRoleParams struct {
	AgglomerationID uuid.UUID `json:"agglomeration_id"`
	Head            bool      `json:"head"`
	Rank            uint      `json:"rank"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Color           string    `json:"color"`
}

func (q RolesQ) Insert(ctx context.Context, data InsertRoleParams) (Role, error) {
	const sqlInsertAtRank = `
		WITH bumped AS (
			UPDATE roles
			SET
				rank = rank + 1,
				updated_at = now()
			WHERE agglomeration_id = $1
			  AND rank >= $2
			RETURNING 1
		),
		ins AS (
			INSERT INTO roles (agglomeration_id, head, rank, name, description, color)
			VALUES ($1, $3, $2, $4, $5, $6)
			RETURNING id, agglomeration_id, head, rank, name, description, color, created_at, updated_at
		)
		SELECT id, agglomeration_id, head, rank, name, description, color, created_at, updated_at
		FROM ins;
	`

	args := []any{
		data.AgglomerationID,
		data.Rank,
		data.Head,
		data.Name,
		data.Description,
		data.Color,
	}

	var inserted Role
	if err := inserted.scan(q.db.QueryRowContext(ctx, sqlInsertAtRank, args...)); err != nil {
		return Role{}, fmt.Errorf("insert role at rank: %w", err)
	}

	return inserted, nil
}

func (q RolesQ) Get(ctx context.Context) (Role, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return Role{}, fmt.Errorf("building select query for %s: %w", RoleTable, err)
	}

	var r Role
	if err = r.scan(q.db.QueryRowContext(ctx, query, args...)); err != nil {
		return Role{}, err
	}

	return r, nil
}

func (q RolesQ) Select(ctx context.Context) ([]Role, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", RoleTable, err)
	}

	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing select query for %s: %w", RoleTable, err)
	}
	defer rows.Close()

	var out []Role
	for rows.Next() {
		var r Role
		if err = r.scan(rows); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (q RolesQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", RoleTable, err)
	}

	_, err = q.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("executing delete query for %s: %w", RoleTable, err)
	}

	return nil
}

func (q RolesQ) Count(ctx context.Context) (int64, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for %s: %w", RoleTable, err)
	}

	var count int64
	if err = q.db.QueryRowContext(ctx, query, args...).Scan(&count); err != nil {
		return 0, fmt.Errorf("scanning count for %s: %w", RoleTable, err)
	}

	return count, nil
}

func (q RolesQ) UpdateOne(ctx context.Context) (Role, error) {
	q.updater = q.updater.Set("updated_at", time.Now().UTC())

	query, args, err := q.updater.Suffix("RETURNING " + RoleColumns).ToSql()
	if err != nil {
		return Role{}, fmt.Errorf("building update query for %s: %w", RoleTable, err)
	}

	var updated Role
	if err = updated.scan(q.db.QueryRowContext(ctx, query, args...)); err != nil {
		return Role{}, err
	}

	return updated, nil
}

func (q RolesQ) UpdateMany(ctx context.Context) (int64, error) {
	q.updater = q.updater.Set("updated_at", time.Now().UTC())

	query, args, err := q.updater.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building update query for %s: %w", RoleTable, err)
	}

	res, err := q.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("executing update query for %s: %w", RoleTable, err)
	}

	aff, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("rows affected for %s: %w", RoleTable, err)
	}

	return aff, nil
}

func (q RolesQ) UpdateName(name string) RolesQ {
	q.updater = q.updater.Set("name", name)
	return q
}

func (q RolesQ) UpdateDescription(description string) RolesQ {
	q.updater = q.updater.Set("description", description)
	return q
}

func (q RolesQ) UpdateColor(color string) RolesQ {
	q.updater = q.updater.Set("color", color)
	return q
}

func (q RolesQ) FilterByID(id ...uuid.UUID) RolesQ {
	q.selector = q.selector.Where(sq.Eq{"r.id": id})
	q.counter = q.counter.Where(sq.Eq{"r.id": id})
	q.updater = q.updater.Where(sq.Eq{"r.id": id})
	q.deleter = q.deleter.Where(sq.Eq{"r.id": id})
	return q
}

func (q RolesQ) FilterByAgglomerationID(id uuid.UUID) RolesQ {
	q.selector = q.selector.Where(sq.Eq{"r.agglomeration_id": id})
	q.counter = q.counter.Where(sq.Eq{"r.agglomeration_id": id})
	q.updater = q.updater.Where(sq.Eq{"r.agglomeration_id": id})
	q.deleter = q.deleter.Where(sq.Eq{"r.agglomeration_id": id})
	return q
}

func (q RolesQ) FilterByAccountID(accountID uuid.UUID) RolesQ {
	sub := sq.
		Select("DISTINCT mr.role_id").
		From("members m").
		Join("member_roles mr ON mr.member_id = m.id").
		Where(sq.Eq{"m.account_id": accountID})

	subSQL, subArgs, err := sub.ToSql()
	if err != nil {
		q.selector = q.selector.Where(sq.Expr("1=0"))
		q.counter = q.counter.Where(sq.Expr("1=0"))
		q.updater = q.updater.Where(sq.Expr("1=0"))
		q.deleter = q.deleter.Where(sq.Expr("1=0"))
		return q
	}

	expr := sq.Expr("r.id IN ("+subSQL+")", subArgs...)

	q.selector = q.selector.Where(expr)
	q.counter = q.counter.Where(expr)
	q.updater = q.updater.Where(expr)
	q.deleter = q.deleter.Where(expr)

	return q
}

func (q RolesQ) FilterByMemberID(memberID uuid.UUID) RolesQ {
	sub := sq.
		Select("mr.role_id").
		From("member_roles mr").
		Where(sq.Eq{"mr.member_id": memberID})

	subSQL, subArgs, err := sub.ToSql()
	if err != nil {
		q.selector = q.selector.Where(sq.Expr("1=0"))
		q.counter = q.counter.Where(sq.Expr("1=0"))
		q.updater = q.updater.Where(sq.Expr("1=0"))
		q.deleter = q.deleter.Where(sq.Expr("1=0"))
		return q
	}

	expr := sq.Expr("r.id IN ("+subSQL+")", subArgs...)

	q.selector = q.selector.Where(expr)
	q.counter = q.counter.Where(expr)
	q.updater = q.updater.Where(expr)
	q.deleter = q.deleter.Where(expr)

	return q
}

func (q RolesQ) FilterHead(head bool) RolesQ {
	q.selector = q.selector.Where(sq.Eq{"r.head": head})
	q.counter = q.counter.Where(sq.Eq{"r.head": head})
	q.updater = q.updater.Where(sq.Eq{"r.head": head})
	q.deleter = q.deleter.Where(sq.Eq{"r.head": head})
	return q
}

func (q RolesQ) FilterByRank(rank int) RolesQ {
	q.selector = q.selector.Where(sq.Eq{"r.rank": rank})
	q.counter = q.counter.Where(sq.Eq{"r.rank": rank})
	q.updater = q.updater.Where(sq.Eq{"r.rank": rank})
	q.deleter = q.deleter.Where(sq.Eq{"r.rank": rank})
	return q
}

func (q RolesQ) FilterLikeName(name string) RolesQ {
	q.selector = q.selector.Where(sq.ILike{"r.name": "%" + name + "%"})
	q.counter = q.counter.Where(sq.ILike{"r.name": "%" + name + "%"})
	q.updater = q.updater.Where(sq.ILike{"r.name": "%" + name + "%"})
	q.deleter = q.deleter.Where(sq.ILike{"r.name": "%" + name + "%"})
	return q
}

func (q RolesQ) OrderByRoleRank(asc bool) RolesQ {
	if asc {
		q.selector = q.selector.OrderBy("r.rank ASC", "r.id ASC")
	} else {
		q.selector = q.selector.OrderBy("r.rank DESC", "r.id DESC")
	}
	return q
}

func (q RolesQ) Page(limit, offset uint) RolesQ {
	q.selector = q.selector.Limit(uint64(limit)).Offset(uint64(offset))
	return q
}

//Special methods to interact with role ranks in agglomeration

func (q RolesQ) DeleteAndShiftRanks(ctx context.Context, roleID uuid.UUID) error {
	const sqlq = `
		WITH del AS (
			DELETE FROM roles
			WHERE id = $1
			RETURNING agglomeration_id, rank
		)
		UPDATE roles r
		SET rank = r.rank - 1,
		    updated_at = now()
		FROM del
		WHERE r.agglomeration_id = del.agglomeration_id
		  AND r.rank > del.rank
	`

	if _, err := q.db.ExecContext(ctx, sqlq, roleID); err != nil {
		return fmt.Errorf("executing delete+shift for %s: %w", RoleTable, err)
	}

	return nil
}

func (q RolesQ) UpdateRoleRank(ctx context.Context, roleID uuid.UUID, newRank uint) (Role, error) {
	var aggID uuid.UUID
	var oldRank int

	{
		const sqlGet = `SELECT agglomeration_id, rank FROM roles WHERE id = $1 LIMIT 1`
		if err := q.db.QueryRowContext(ctx, sqlGet, roleID).Scan(&aggID, &oldRank); err != nil {
			return Role{}, fmt.Errorf("scanning role rank: %w", err)
		}
	}

	if oldRank == int(newRank) {
		return NewRolesQ(q.db).FilterByID(roleID).Get(ctx)
	}

	const sqlMove = `
		WITH upd AS (
			UPDATE roles
			SET
				rank = CASE
					WHEN id = $1 THEN $2
					WHEN $2 < $3 AND rank >= $2 AND rank < $3 THEN rank + 1
					WHEN $2 > $3 AND rank <= $2 AND rank > $3 THEN rank - 1
					ELSE rank
				END,
				updated_at = now()
			WHERE agglomeration_id = $4
			RETURNING id, agglomeration_id, head, rank, name, description, color, created_at, updated_at
		)
		SELECT id, agglomeration_id, head, rank, name, description, color, created_at, updated_at
		FROM upd
		WHERE id = $1
	`

	args := []any{roleID, int(newRank), oldRank, aggID}

	var out Role
	if err := out.scan(q.db.QueryRowContext(ctx, sqlMove, args...)); err != nil {
		return Role{}, err
	}

	return out, nil
}

func (q RolesQ) UpdateRolesRanks(
	ctx context.Context,
	agglomerationID uuid.UUID,
	order map[uint]uuid.UUID,
) ([]Role, error) {
	roles, err := NewRolesQ(q.db).
		FilterByAgglomerationID(agglomerationID).
		OrderByRoleRank(true).
		Select(ctx)
	if err != nil {
		return nil, fmt.Errorf("select roles by agglomeration: %w", err)
	}
	if len(roles) == 0 {
		return nil, fmt.Errorf("no roles in agglomeration %s", agglomerationID)
	}

	n := len(roles)

	idToIndex := make(map[uuid.UUID]int, n)
	for i := range roles {
		idToIndex[roles[i].ID] = i
	}

	seenIDs := make(map[uuid.UUID]struct{}, len(order))
	for newRank, roleID := range order {
		if int(newRank) < 0 || int(newRank) >= n {
			return nil, fmt.Errorf("rank %d out of range [0..%d]", newRank, n-1)
		}
		if _, ok := idToIndex[roleID]; !ok {
			return nil, fmt.Errorf("role %s not in agglomeration %s", roleID, agglomerationID)
		}
		if _, ok := seenIDs[roleID]; ok {
			return nil, fmt.Errorf("duplicate role id %s in order map", roleID)
		}
		seenIDs[roleID] = struct{}{}
	}

	type mv struct {
		rank uint
		id   uuid.UUID
	}
	moves := make([]mv, 0, len(order))
	for r, id := range order {
		moves = append(moves, mv{rank: r, id: id})
	}
	for i := 0; i < len(moves); i++ {
		for j := i + 1; j < len(moves); j++ {
			if moves[j].rank < moves[i].rank {
				moves[i], moves[j] = moves[j], moves[i]
			}
		}
	}

	current := make([]uuid.UUID, 0, n)
	for i := range roles {
		current = append(current, roles[i].ID)
	}

	removeAt := func(s []uuid.UUID, idx int) []uuid.UUID {
		copy(s[idx:], s[idx+1:])
		return s[:len(s)-1]
	}
	insertAt := func(s []uuid.UUID, idx int, v uuid.UUID) []uuid.UUID {
		s = append(s, uuid.Nil)
		copy(s[idx+1:], s[idx:])
		s[idx] = v
		return s
	}

	for _, m := range moves {
		id := m.id
		newIdx := int(m.rank)

		var oldIdx int = -1
		for i := range current {
			if current[i] == id {
				oldIdx = i
				break
			}
		}
		if oldIdx == -1 {
			return nil, fmt.Errorf("role %s not found in current order", id)
		}

		if oldIdx == newIdx {
			continue
		}

		current = removeAt(current, oldIdx)
		if oldIdx < newIdx {
			newIdx--
		}
		current = insertAt(current, newIdx, id)
	}

	changed := make([]uuid.UUID, 0, n)
	newRanks := make([]int, 0, n)

	for idx, id := range current {
		oldIdx := idToIndex[id]
		if oldIdx != idx {
			changed = append(changed, id)
			newRanks = append(newRanks, idx)
		}
	}

	if len(changed) == 0 {
		return roles, nil
	}

	const sqlUpdate = `
		UPDATE roles r
		SET
			rank = v.rank,
			updated_at = now()
		FROM (
			SELECT UNNEST($1::uuid[]) AS id, UNNEST($2::int[]) AS rank
		) v
		WHERE r.id = v.id
		  AND r.agglomeration_id = $3
		RETURNING r.id, r.agglomeration_id, r.head, r.rank, r.name, r.description, r.color, r.created_at, r.updated_at
	`

	rows, err := q.db.QueryContext(ctx, sqlUpdate, changed, newRanks, agglomerationID)
	if err != nil {
		return nil, fmt.Errorf("updating roles ranks: %w", err)
	}
	defer rows.Close()

	out := make([]Role, 0, len(changed))
	for rows.Next() {
		var r Role
		if err := r.scan(rows); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}
