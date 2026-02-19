package app

import (
	"context"
	"time"

	"github.com/netbill/eventbox"
	eventpg "github.com/netbill/eventbox/pg"
	"github.com/netbill/organizations-svc/internal/messenger/outbound"
)

type OutboxConfig struct {
	Routines       int           `mapstructure:"routines"`
	Slots          int           `mapstructure:"slots"`
	BatchSize      int           `mapstructure:"batch_size"`
	Sleep          time.Duration `mapstructure:"sleep"`
	MinNextAttempt time.Duration `mapstructure:"min_next_attempt"`
	MaxNextAttempt time.Duration `mapstructure:"max_next_attempt"`
	MaxAttempts    int32         `mapstructure:"max_attempts"`
}

func (a *App) RunOutbox(ctx context.Context, producer *eventbox.Producer) {
	outbox := eventpg.NewOutbox(a.db)

	worker := outbound.NewOutboxWorker(nil, outbox, producer, eventbox.OutboxWorkerConfig{
		Routines:       a.config.Kafka.Outbox.Routines,
		Slots:          a.config.Kafka.Outbox.Slots,
		BatchSize:      a.config.Kafka.Outbox.BatchSize,
		Sleep:          a.config.Kafka.Outbox.Sleep,
		MinNextAttempt: a.config.Kafka.Outbox.MinNextAttempt,
		MaxNextAttempt: a.config.Kafka.Outbox.MaxNextAttempt,
		MaxAttempts:    a.config.Kafka.Outbox.MaxAttempts,
	})

	worker.Run(ctx)
}
