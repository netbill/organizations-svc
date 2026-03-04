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

const profilesTable = "profiles"

const profilesColumns = "account_id, username, pseudonym, avatar_key, version, source_created_at, source_updated_at, replica_created_at, replica_updated_at"
const profilesColumnsP = "p.account_id, p.username, p.pseudonym, p.avatar_key, p.version, p.source_created_at, p.source_updated_at, p.replica_created_at, p.replica_updated_at"

func scanProfile(row sq.RowScanner) (p repository.ProfileRow, err error) {
	var pseudonym pgtype.Text
	var avatarKey pgtype.Text

	err = row.Scan(
		&p.AccountID,
		&p.Username,
		&pseudonym,
		&avatarKey,
		&p.Version,
		&p.SourceCreatedAt,
		&p.SourceUpdatedAt,
		&p.ReplicaCreatedAt,
		&p.ReplicaUpdatedAt,
	)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return repository.ProfileRow{}, nil
	case err != nil:
		return repository.ProfileRow{}, fmt.Errorf("scanning profile: %w", err)
	}

	if pseudonym.Valid {
		p.Pseudonym = &pseudonym.String
	}
	if avatarKey.Valid {
		p.AvatarKey = &avatarKey.String
	}

	return p, nil
}

type profiles struct {
	db       *pgdbx.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewProfilesQ(db *pgdbx.DB) repository.ProfilesQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &profiles{
		db:       db,
		selector: b.Select(profilesColumnsP).From(profilesTable + " p"),
		inserter: b.Insert(profilesTable),
		updater:  b.Update(profilesTable + " p"),
		deleter:  b.Delete(profilesTable + " p"),
		counter:  b.Select("COUNT(*)").From(profilesTable + " p"),
	}
}

func (q *profiles) New() repository.ProfilesQ {
	return NewProfilesQ(q.db)
}

func (q *profiles) Insert(ctx context.Context, data repository.ProfileRow) (repository.ProfileRow, error) {
	query, args, err := q.inserter.SetMap(map[string]any{
		"account_id":        data.AccountID,
		"username":          data.Username,
		"pseudonym":         data.Pseudonym,
		"avatar_key":        data.AvatarKey,
		"source_created_at": data.SourceCreatedAt.UTC(),
		"source_updated_at": data.SourceUpdatedAt.UTC(),
	}).Suffix("RETURNING " + profilesColumns).ToSql()
	if err != nil {
		return repository.ProfileRow{}, fmt.Errorf("building insert query for %s: %w", profilesTable, err)
	}

	return scanProfile(q.db.QueryRow(ctx, query, args...))
}

func (q *profiles) Get(ctx context.Context) (repository.ProfileRow, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return repository.ProfileRow{}, fmt.Errorf("building select query for %s: %w", profilesTable, err)
	}

	return scanProfile(q.db.QueryRow(ctx, query, args...))
}

func (q *profiles) Select(ctx context.Context) ([]repository.ProfileRow, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", profilesTable, err)
	}

	rows, err := q.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing select query for %s: %w", profilesTable, err)
	}
	defer rows.Close()

	out := make([]repository.ProfileRow, 0)
	for rows.Next() {
		p, err := scanProfile(rows)
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

func (q *profiles) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", profilesTable, err)
	}

	if _, err = q.db.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("executing delete query for %s: %w", profilesTable, err)
	}

	return nil
}

func (q *profiles) FilterByAccountID(accountID ...uuid.UUID) repository.ProfilesQ {
	q.selector = q.selector.Where(sq.Eq{"p.account_id": accountID})
	q.counter = q.counter.Where(sq.Eq{"p.account_id": accountID})
	q.updater = q.updater.Where(sq.Eq{"p.account_id": accountID})
	q.deleter = q.deleter.Where(sq.Eq{"p.account_id": accountID})
	return q
}

func (q *profiles) FilterByUsername(username string) repository.ProfilesQ {
	q.selector = q.selector.Where(sq.Eq{"p.username": username})
	q.counter = q.counter.Where(sq.Eq{"p.username": username})
	q.updater = q.updater.Where(sq.Eq{"p.username": username})
	q.deleter = q.deleter.Where(sq.Eq{"p.username": username})
	return q
}

func (q *profiles) UpdateOne(ctx context.Context) (repository.ProfileRow, error) {
	q.updater = q.updater.Set("replica_updated_at", time.Now().UTC())

	query, args, err := q.updater.Suffix("RETURNING " + profilesColumns).ToSql()
	if err != nil {
		return repository.ProfileRow{}, fmt.Errorf("building update query for %s: %w", profilesTable, err)
	}

	return scanProfile(q.db.QueryRow(ctx, query, args...))
}

func (q *profiles) UpdateUsername(v string) repository.ProfilesQ {
	q.updater = q.updater.Set("username", v)
	return q
}

func (q *profiles) UpdatePseudonym(v *string) repository.ProfilesQ {
	q.updater = q.updater.Set("pseudonym", v)
	return q
}

func (q *profiles) UpdateAvatarKey(v *string) repository.ProfilesQ {
	q.updater = q.updater.Set("avatar_key", v)
	return q
}

func (q *profiles) UpdateVersion(v int32) repository.ProfilesQ {
	q.updater = q.updater.Set("version", v)
	return q
}

func (q *profiles) UpdateSourceUpdatedAt(updatedAt time.Time) repository.ProfilesQ {
	q.updater = q.updater.Set("source_updated_at", updatedAt.UTC())
	return q
}

func (q *profiles) Page(limit, offset uint) repository.ProfilesQ {
	q.selector = q.selector.Limit(uint64(limit)).Offset(uint64(offset))
	return q
}

func (q *profiles) Count(ctx context.Context) (uint, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for %s: %w", profilesTable, err)
	}

	var count uint
	if err = q.db.QueryRow(ctx, query, args...).Scan(&count); err != nil {
		return 0, fmt.Errorf("scanning count for %s: %w", profilesTable, err)
	}

	return count, nil
}
