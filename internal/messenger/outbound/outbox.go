package outbound

import (
	"context"
	"fmt"
	"os"

	"github.com/netbill/eventbox"
)

type OutboxWorker struct {
	worker *eventbox.OutboxWorker
}

func NewOutboxWorker(
	log eventbox.Logger,
	outbox eventbox.Outbox,
	producer eventbox.Producer,
	config eventbox.OutboxWorkerConfig,
) *OutboxWorker {
	id := buildProcessID("outbox")

	return &OutboxWorker{
		worker: eventbox.NewOutboxWorker(id, log, outbox, producer, config),
	}
}

func (m *OutboxWorker) RunOutbox(ctx context.Context) {
	m.worker.Run(ctx)
}

func buildProcessID(service string) string {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	return fmt.Sprintf("%s-%s-%d", service, hostname, os.Getpid())
}
