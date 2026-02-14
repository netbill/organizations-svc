package messenger

import (
	"context"

	eventpg "github.com/netbill/eventbox/pg"
	"github.com/netbill/organizations-svc/internal/messenger/evtypes"
	"github.com/segmentio/kafka-go"
)

type handlers interface {
	AccountCreated(
		ctx context.Context,
		message kafka.Message,
	) error
	AccountDeleted(
		ctx context.Context,
		message kafka.Message,
	) error
	AccountUsernameUpdated(
		ctx context.Context,
		message kafka.Message,
	) error
}

func (m *Manager) RunInbox(ctx context.Context, handlers handlers) {
	id := BuildProcessID("inbox")
	worker := eventpg.NewInboxWorker(id, m.log, m.db, eventpg.InboxWorkerConfig{
		Routines:       m.config.Inbox.Routines,
		Slots:          m.config.Inbox.Slots,
		BatchSize:      m.config.Inbox.BatchSize,
		Sleep:          m.config.Inbox.Sleep,
		MinNextAttempt: m.config.Inbox.MinNextAttempt,
		MaxNextAttempt: m.config.Inbox.MaxNextAttempt,
		MaxAttempts:    m.config.Inbox.MaxAttempts,
	})

	defer func() {
		if err := worker.Stop(context.Background()); err != nil {
			m.log.WithError(err).Errorf("stop inbox worker %s failed", id)
		}
	}()

	worker.Route(evtypes.AccountCreatedEvent, handlers.AccountCreated)
	worker.Route(evtypes.AccountDeletedEvent, handlers.AccountDeleted)
	worker.Route(evtypes.AccountUsernameUpdatedEvent, handlers.AccountUsernameUpdated)

	worker.Run(ctx)
}
