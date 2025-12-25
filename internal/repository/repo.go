package repository

import (
	"context"
	"database/sql"

	"github.com/umisto/cities-svc/internal/repository/pgdb"
	"github.com/umisto/pgx"
)

type Service struct {
	db *sql.DB
}

func New(db *sql.DB) Service {
	return Service{db: db}
}

func (s Service) exec(ctx context.Context) pgdb.DBTX {
	return pgx.Exec(s.db, ctx)
}

func (s Service) sql(ctx context.Context) *pgdb.Queries {
	return pgdb.New(s.exec(ctx))
}

func (s Service) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return pgx.Transaction(s.db, ctx, fn)
}
