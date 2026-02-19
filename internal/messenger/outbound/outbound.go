package outbound

import (
	"context"

	"github.com/netbill/eventbox"
)

type Outbound struct {
	worker *eventbox.OutboxWorker
}

func NewOutboxWorker(
	log eventbox.Logger,
	outbox eventbox.Outbox,
	producer *eventbox.Producer,
	config eventbox.OutboxWorkerConfig,
) *Outbound {
	id := "TODO"

	return &Outbound{
		worker: eventbox.NewOutboxWorker(id, log, outbox, producer, config),
	}
}

func (m *Outbound) Run(ctx context.Context) {
	m.worker.Run(ctx)
}
