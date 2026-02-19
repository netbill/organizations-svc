package app

import (
	"context"
	"time"

	"github.com/netbill/eventbox"
	eventpg "github.com/netbill/eventbox/pg"
	"github.com/netbill/organizations-svc/internal/messenger/inbound"
	"github.com/netbill/organizations-svc/pkg/evtypes"
)

type KafkaInboxConfig struct {
	Routines       int           `mapstructure:"routines"`
	Slots          int           `mapstructure:"slots"`
	BatchSize      int           `mapstructure:"batch_size"`
	Sleep          time.Duration `mapstructure:"sleep"`
	MinNextAttempt time.Duration `mapstructure:"min_next_attempt"`
	MaxNextAttempt time.Duration `mapstructure:"max_next_attempt"`
	MaxAttempts    int32         `mapstructure:"max_attempts"`
}

func (a *App) RunInbox(ctx context.Context, handlers inbound.Handlers) {
	inbox := eventpg.NewInbox(a.db)

	worker := eventbox.NewInboxWorker("TODO", nil, inbox, eventbox.InboxWorkerConfig{
		Routines:       a.config.Kafka.Inbox.Routines,
		Slots:          a.config.Kafka.Inbox.Slots,
		BatchSize:      a.config.Kafka.Inbox.BatchSize,
		Sleep:          a.config.Kafka.Inbox.Sleep,
		MinNextAttempt: a.config.Kafka.Inbox.MinNextAttempt,
		MaxNextAttempt: a.config.Kafka.Inbox.MaxNextAttempt,
		MaxAttempts:    a.config.Kafka.Inbox.MaxAttempts,
	})

	defer func() {
		worker.Stop(context.Background())
	}()

	worker.Route(evtypes.ProfileCreatedEvent, handlers.ProfileCreated)
	worker.Route(evtypes.ProfileDeletedEvent, handlers.ProfileDeleted)
	worker.Route(evtypes.ProfileUpdatedEvent, handlers.ProfileUpdated)

	worker.Run(ctx)
}
