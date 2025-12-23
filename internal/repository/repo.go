package repository

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/repository/pgdb"
)

type Service struct {
	sql *pgdb.Queries
}

func New(db *sql.DB) *Service {
	return &Service{
		sql: pgdb.New(db),
	}
}

func nullString(s *string) sql.NullString {
	if s != nil {
		return sql.NullString{String: *s, Valid: true}
	}
	return sql.NullString{Valid: false}
}

func nullUUID(id *uuid.UUID) uuid.NullUUID {
	if id != nil {
		return uuid.NullUUID{UUID: *id, Valid: true}
	}

	return uuid.NullUUID{Valid: false}
}

func nullInt32(i *int) sql.NullInt32 {
	if i != nil {
		return sql.NullInt32{Int32: int32(*i), Valid: true}
	}

	return sql.NullInt32{Valid: false}
}

func calculateLimit(limit, def, max int) int {
	if limit <= 0 {
		return def
	}
	if limit > max {
		return max
	}

	return limit
}
