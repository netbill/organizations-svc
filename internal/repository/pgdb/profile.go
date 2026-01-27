package pgdb

import (
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/netbill/pgxtx"
)

const ProfileTable = "profiles"

const ProfileColumns = "account_id, username, official, pseudonym, avatar, source_created_at, source_updated_at, replica_created_at, replica_updated_at"
const ProfileColumnsP = "p.account_id, p.username, p.official, p.pseudonym, p.avatar, p.source_created_at, p.source_updated_at, p.replica_created_at, p.replica_updated_at"

type Profile struct {
	AccountID uuid.UUID `json:"account_id"`
	Username  string    `json:"username"`
	Official  bool      `json:"official"`

	Pseudonym pgtype.Text `json:"pseudonym,omitempty"`
	Avatar    pgtype.Text `json:"avatar,omitempty"`

	SourceCreatedAt  time.Time `json:"source_created_at"`
	SourceUpdatedAt  time.Time `json:"source_updated_at"`
	ReplicaCreatedAt time.Time `json:"replica_created_at"`
	ReplicaUpdatedAt time.Time `json:"replica_updated_at"`
}

func (p *Profile) scan(row sq.RowScanner) error {
	if err := row.Scan(
		&p.AccountID,
		&p.Username,
		&p.Official,
		&p.Pseudonym,
		&p.Avatar,
		&p.SourceCreatedAt,
		&p.SourceUpdatedAt,
		&p.ReplicaCreatedAt,
		&p.ReplicaUpdatedAt,
	); err != nil {
		return fmt.Errorf("scanning profile: %w", err)
	}
	return nil
}

type ProfilesQ struct {
	db       pgxtx.DBTX
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewProfilesQ(db pgxtx.DBTX) ProfilesQ {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return ProfilesQ{
		db:       db,
		selector: builder.Select(ProfileColumnsP).From(ProfileTable + " p"),
		inserter: builder.Insert(ProfileTable),
		updater:  builder.Update(ProfileTable + " p"),
		deleter:  builder.Delete(ProfileTable + " p"),
		counter:  builder.Select("COUNT(*)").From(ProfileTable + " p"),
	}
}

type ProfileInsertInput struct {
	AccountID uuid.UUID
	Username  string
	Official  bool

	Pseudonym *string
	Avatar    *string

	SourceCreatedAt time.Time
	SourceUpdatedAt time.Time
}

func (q ProfilesQ) Insert(ctx context.Context, data ProfileInsertInput) (Profile, error) {
	query, args, err := q.inserter.SetMap(map[string]any{
		"account_id":        data.AccountID,
		"username":          data.Username,
		"official":          data.Official,
		"pseudonym":         data.Pseudonym,
		"avatar":            data.Avatar,
		"source_created_at": data.SourceCreatedAt.UTC(),
		"source_updated_at": data.SourceUpdatedAt.UTC(),
		// replica_* defaults exist in schema, but keeping your previous behavior:
		"replica_created_at": time.Now().UTC(),
		"replica_updated_at": time.Now().UTC(),
	}).Suffix("RETURNING " + ProfileColumns).ToSql()
	if err != nil {
		return Profile{}, fmt.Errorf("building insert query for %s: %w", ProfileTable, err)
	}

	var inserted Profile
	if err = inserted.scan(q.db.QueryRow(ctx, query, args...)); err != nil {
		return Profile{}, err
	}
	return inserted, nil
}

func (q ProfilesQ) Get(ctx context.Context) (Profile, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return Profile{}, fmt.Errorf("building select query for %s: %w", ProfileTable, err)
	}

	var p Profile
	if err = p.scan(q.db.QueryRow(ctx, query, args...)); err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return Profile{}, nil
		default:
			return Profile{}, err
		}
	}
	return p, nil
}

func (q ProfilesQ) Select(ctx context.Context) ([]Profile, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", ProfileTable, err)
	}

	rows, err := q.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing select query for %s: %w", ProfileTable, err)
	}
	defer rows.Close()

	var out []Profile
	for rows.Next() {
		var p Profile
		if err = p.scan(rows); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (q ProfilesQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", ProfileTable, err)
	}

	if _, err = q.db.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("executing delete query for %s: %w", ProfileTable, err)
	}

	return nil
}

func (q ProfilesQ) FilterByAccountID(accountID uuid.UUID) ProfilesQ {
	q.selector = q.selector.Where(sq.Eq{"p.account_id": accountID})
	q.counter = q.counter.Where(sq.Eq{"p.account_id": accountID})
	q.updater = q.updater.Where(sq.Eq{"p.account_id": accountID})
	q.deleter = q.deleter.Where(sq.Eq{"p.account_id": accountID})
	return q
}

func (q ProfilesQ) FilterByUsername(username string) ProfilesQ {
	q.selector = q.selector.Where(sq.Eq{"p.username": username})
	q.counter = q.counter.Where(sq.Eq{"p.username": username})
	q.updater = q.updater.Where(sq.Eq{"p.username": username})
	q.deleter = q.deleter.Where(sq.Eq{"p.username": username})
	return q
}

func (q ProfilesQ) FilterOfficial(official bool) ProfilesQ {
	q.selector = q.selector.Where(sq.Eq{"p.official": official})
	q.counter = q.counter.Where(sq.Eq{"p.official": official})
	q.updater = q.updater.Where(sq.Eq{"p.official": official})
	q.deleter = q.deleter.Where(sq.Eq{"p.official": official})
	return q
}

func (q ProfilesQ) FilterLikeUsername(username string) ProfilesQ {
	q.selector = q.selector.Where(sq.ILike{"p.username": "%" + username + "%"})
	q.counter = q.counter.Where(sq.ILike{"p.username": "%" + username + "%"})
	q.updater = q.updater.Where(sq.ILike{"p.username": "%" + username + "%"})
	q.deleter = q.deleter.Where(sq.ILike{"p.username": "%" + username + "%"})
	return q
}

func (q ProfilesQ) FilterLikePseudonym(pseudonym string) ProfilesQ {
	q.selector = q.selector.Where(sq.ILike{"p.pseudonym": "%" + pseudonym + "%"})
	q.counter = q.counter.Where(sq.ILike{"p.pseudonym": "%" + pseudonym + "%"})
	q.updater = q.updater.Where(sq.ILike{"p.pseudonym": "%" + pseudonym + "%"})
	q.deleter = q.deleter.Where(sq.ILike{"p.pseudonym": "%" + pseudonym + "%"})
	return q
}

func (q ProfilesQ) UpdateOne(ctx context.Context) (Profile, error) {
	q.updater = q.updater.Set("replica_updated_at", time.Now().UTC())

	query, args, err := q.updater.Suffix("RETURNING " + ProfileColumns).ToSql()
	if err != nil {
		return Profile{}, fmt.Errorf("building update query for %s: %w", ProfileTable, err)
	}

	var updated Profile
	if err = updated.scan(q.db.QueryRow(ctx, query, args...)); err != nil {
		return Profile{}, err
	}
	return updated, nil
}

func (q ProfilesQ) UpdateMany(ctx context.Context) (int64, error) {
	q.updater = q.updater.Set("replica_updated_at", time.Now().UTC())

	query, args, err := q.updater.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building update query for %s: %w", ProfileTable, err)
	}

	res, err := q.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("executing update query for %s: %w", ProfileTable, err)
	}

	return res.RowsAffected(), nil
}

func (q ProfilesQ) UpdateUsername(username string) ProfilesQ {
	q.updater = q.updater.Set("username", username)
	return q
}

func (q ProfilesQ) UpdateOfficial(official bool) ProfilesQ {
	q.updater = q.updater.Set("official", official)
	return q
}

func (q ProfilesQ) UpdatePseudonym(pseudonym *string) ProfilesQ {
	q.updater = q.updater.Set("pseudonym", pseudonym)
	return q
}

func (q ProfilesQ) UpdateAvatar(avatar *string) ProfilesQ {
	q.updater = q.updater.Set("avatar", avatar)
	return q
}

func (q ProfilesQ) UpdateSourceUpdatedAt(updatedAt time.Time) ProfilesQ {
	q.updater = q.updater.Set("source_updated_at", updatedAt.UTC())
	return q
}

func (q ProfilesQ) Page(limit, offset uint) ProfilesQ {
	q.selector = q.selector.Limit(uint64(limit)).Offset(uint64(offset))
	return q
}

func (q ProfilesQ) Count(ctx context.Context) (uint, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for %s: %w", ProfileTable, err)
	}

	var count uint
	if err = q.db.QueryRow(ctx, query, args...).Scan(&count); err != nil {
		return 0, fmt.Errorf("scanning count for %s: %w", ProfileTable, err)
	}

	return count, nil
}
