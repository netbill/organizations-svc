package outbound

import (
	"database/sql"

	"github.com/netbill/evebox/box/outbox"
	"github.com/netbill/logium"
)

type Outbound struct {
	log    logium.Logger
	outbox outbox.Box
}

func New(log logium.Logger, db *sql.DB) *Outbound {
	return &Outbound{
		log:    log,
		outbox: outbox.New(db),
	}
}
