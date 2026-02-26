package pg

import (
	"context"

	"github.com/netbill/organizations-svc/internal/repository"
	"github.com/netbill/pgdbx"
)

type transaction struct {
	db *pgdbx.DB
}

func NewTransaction(db *pgdbx.DB) repository.Transaction {
	return &transaction{
		db: db,
	}
}

func (q *transaction) Begin(ctx context.Context, fn func(ctx context.Context) error) error {
	return q.db.Transaction(ctx, fn)
}
