package messenger

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/eventbox"
	"github.com/netbill/organizations-svc/log"
)

type OutboxConfig struct {
	Routines       int           `json:"routines"`
	Slots          int           `json:"slots"`
	BatchSize      int           `json:"batch_size"`
	Sleep          time.Duration `json:"sleep"`
	MinNextAttempt time.Duration `json:"min_next_attempt"`
	MaxNextAttempt time.Duration `json:"max_next_attempt"`
	MaxAttempts    int32         `json:"max_attempts"`
}

type Outbound struct {
	worker *eventbox.OutboxWorker
}

func NewOutbound(
	log *log.Logger,
	outbox eventbox.Outbox,
	producer *eventbox.Producer,
	cfg eventbox.OutboxWorkerConfig,
) *Outbound {
	id := uuid.New().String()

	return &Outbound{
		worker: eventbox.NewOutboxWorker(id, log, outbox, producer, cfg),
	}
}

func (o *Outbound) RunOutbox(ctx context.Context) {
	o.worker.Run(ctx)
}
