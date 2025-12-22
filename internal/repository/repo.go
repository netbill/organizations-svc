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
