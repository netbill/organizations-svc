package pgdbsq

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/repository/pgdb"

	sq "github.com/Masterminds/squirrel"
)

const RoleTable = "roles"

const RoleColumns = "id, agglomeration_id, head, editable, rank, name, created_at, updated_at"
const RoleColumnsR = "r.id, r.agglomeration_id, r.head, r.editable, r.rank, r.name, r.created_at, r.updated_at"

type Role struct {
	ID              uuid.UUID `json:"id"`
	AgglomerationID uuid.UUID `json:"agglomeration_id"`
	Head            bool      `json:"head"`
	Editable        bool      `json:"editable"`
	Rank            int       `json:"rank"`
	Name            string    `json:"name"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (r *Role) scan(row sq.RowScanner) error {
	err := row.Scan(
		&r.ID,
		&r.AgglomerationID,
		&r.Head,
		&r.Editable,
		&r.Rank,
		&r.Name,
		&r.CreatedAt,
		&r.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("scanning role: %w", err)
	}
	return nil
}

type RoleQ struct {
	db       pgdb.DBTX
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewRoleQ(db pgdb.DBTX) RoleQ {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return RoleQ{
		db:       db,
		selector: builder.Select(RoleColumnsR).From(RoleTable + " r"),
		inserter: builder.Insert(RoleTable),
		updater:  builder.Update(RoleTable + " r"),
		deleter:  builder.Delete(RoleTable + " r"),
		counter:  builder.Select("COUNT(*)").From(RoleTable + " r"),
	}
}

func (q RoleQ) New() RoleQ {
	return NewRoleQ(q.db)
}

func (q RoleQ) Insert(ctx context.Context, data Role) (Role, error) {
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
			INSERT INTO roles (agglomeration_id, head, editable, rank, name)
			VALUES ($1, $3, $4, $2, $5)
			RETURNING id, agglomeration_id, head, editable, rank, name, created_at, updated_at
		)
		SELECT id, agglomeration_id, head, editable, rank, name, created_at, updated_at
		FROM ins;
	`

	args := []any{
		data.AgglomerationID,
		data.Rank,
		data.Head,
		data.Editable,
		data.Name,
	}

	var inserted Role
	if err := inserted.scan(q.db.QueryRowContext(ctx, sqlInsertAtRank, args...)); err != nil {
		return Role{}, fmt.Errorf("insert role at rank: %w", err)
	}

	return inserted, nil
}

func (q RoleQ) Get(ctx context.Context) (Role, error) {
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

func (q RoleQ) Select(ctx context.Context) ([]Role, error) {
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

func (q RoleQ) Delete(ctx context.Context) error {
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

func (q RoleQ) Count(ctx context.Context) (int64, error) {
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

func (q RoleQ) UpdateOne(ctx context.Context) (Role, error) {
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

func (q RoleQ) UpdateMany(ctx context.Context) (int64, error) {
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

func (q RoleQ) UpdateRoleRank(ctx context.Context, roleID uuid.UUID, newRank uint) (Role, error) {
	var aggID uuid.UUID
	var oldRank int

	{
		const sqlGet = `SELECT agglomeration_id, rank FROM roles WHERE id = $1 LIMIT 1`
		if err := q.db.QueryRowContext(ctx, sqlGet, roleID).Scan(&aggID, &oldRank); err != nil {
			return Role{}, fmt.Errorf("scanning role rank: %w", err)
		}
	}

	if oldRank == int(newRank) {
		return q.New().FilterByID(roleID).Get(ctx)
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
			RETURNING id, agglomeration_id, head, editable, rank, name, created_at, updated_at
		)
		SELECT id, agglomeration_id, head, editable, rank, name, created_at, updated_at
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

func (q RoleQ) FilterByID(id uuid.UUID) RoleQ {
	q.selector = q.selector.Where(sq.Eq{"r.id": id})
	q.counter = q.counter.Where(sq.Eq{"r.id": id})
	q.updater = q.updater.Where(sq.Eq{"r.id": id})
	q.deleter = q.deleter.Where(sq.Eq{"r.id": id})
	return q
}

func (q RoleQ) FilterByAgglomerationID(id uuid.UUID) RoleQ {
	q.selector = q.selector.Where(sq.Eq{"r.agglomeration_id": id})
	q.counter = q.counter.Where(sq.Eq{"r.agglomeration_id": id})
	q.updater = q.updater.Where(sq.Eq{"r.agglomeration_id": id})
	q.deleter = q.deleter.Where(sq.Eq{"r.agglomeration_id": id})
	return q
}

func (q RoleQ) FilterHead(head bool) RoleQ {
	q.selector = q.selector.Where(sq.Eq{"r.head": head})
	q.counter = q.counter.Where(sq.Eq{"r.head": head})
	q.updater = q.updater.Where(sq.Eq{"r.head": head})
	q.deleter = q.deleter.Where(sq.Eq{"r.head": head})
	return q
}

func (q RoleQ) FilterEditable(editable bool) RoleQ {
	q.selector = q.selector.Where(sq.Eq{"r.editable": editable})
	q.counter = q.counter.Where(sq.Eq{"r.editable": editable})
	q.updater = q.updater.Where(sq.Eq{"r.editable": editable})
	q.deleter = q.deleter.Where(sq.Eq{"r.editable": editable})
	return q
}

func (q RoleQ) FilterByRank(rank int) RoleQ {
	q.selector = q.selector.Where(sq.Eq{"r.rank": rank})
	q.counter = q.counter.Where(sq.Eq{"r.rank": rank})
	q.updater = q.updater.Where(sq.Eq{"r.rank": rank})
	q.deleter = q.deleter.Where(sq.Eq{"r.rank": rank})
	return q
}

func (q RoleQ) FilterByName(name string) RoleQ {
	q.selector = q.selector.Where(sq.Eq{"r.name": name})
	q.counter = q.counter.Where(sq.Eq{"r.name": name})
	q.updater = q.updater.Where(sq.Eq{"r.name": name})
	q.deleter = q.deleter.Where(sq.Eq{"r.name": name})
	return q
}

func (q RoleQ) FilterLikeName(name string) RoleQ {
	q.selector = q.selector.Where(sq.ILike{"r.name": "%" + name + "%"})
	q.counter = q.counter.Where(sq.ILike{"r.name": "%" + name + "%"})
	q.updater = q.updater.Where(sq.ILike{"r.name": "%" + name + "%"})
	q.deleter = q.deleter.Where(sq.ILike{"r.name": "%" + name + "%"})
	return q
}

func (q RoleQ) UpdateHead(head bool) RoleQ {
	q.updater = q.updater.Set("head", head)
	return q
}

func (q RoleQ) UpdateEditable(editable bool) RoleQ {
	q.updater = q.updater.Set("editable", editable)
	return q
}

func (q RoleQ) UpdateRank(rank int) RoleQ {
	q.updater = q.updater.Set("rank", rank)
	return q
}

func (q RoleQ) UpdateName(name string) RoleQ {
	q.updater = q.updater.Set("name", name)
	return q
}

func (q RoleQ) CursorCreatedAt(limit uint, asc bool, createdAt time.Time, id uuid.UUID) RoleQ {
	if asc {
		q.selector = q.selector.OrderBy("r.created_at ASC", "r.id ASC")
	} else {
		q.selector = q.selector.OrderBy("r.created_at DESC", "r.id DESC")
	}

	q.selector = q.selector.Limit(uint64(limit))

	if asc {
		q.selector = q.selector.Where(sq.Expr("(r.created_at, r.id) > (?, ?)", createdAt, id))
	} else {
		q.selector = q.selector.Where(sq.Expr("(r.created_at, r.id) < (?, ?)", createdAt, id))
	}

	return q
}

func (q RoleQ) DeleteAndShiftRanks(ctx context.Context, roleID uuid.UUID) (Role, error) {
	const sqlq = `
		WITH del AS (
			DELETE FROM roles
			WHERE id = $1
			RETURNING id, agglomeration_id, head, editable, rank, name, created_at, updated_at
		),
		bump AS (
			UPDATE roles r
			SET
				rank = r.rank - 1,
				updated_at = now()
			FROM del
			WHERE r.agglomeration_id = del.agglomeration_id
			  AND r.rank > del.rank
			RETURNING 1
		)
		SELECT id, agglomeration_id, head, editable, rank, name, created_at, updated_at
		FROM del;
	`

	var deleted Role
	if err := deleted.scan(q.db.QueryRowContext(ctx, sqlq, roleID)); err != nil {
		return Role{}, fmt.Errorf("delete role and shift ranks: %w", err)
	}

	return deleted, nil
}
