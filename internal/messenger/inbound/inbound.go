package inbound

import (
	"context"

	"github.com/netbill/eventbox"
	"github.com/netbill/organizations-svc/pkg/evtypes"
	"github.com/segmentio/kafka-go"
)

type Handlers interface {
	ProfileCreated(
		ctx context.Context,
		message kafka.Message,
	) error
	ProfileDeleted(
		ctx context.Context,
		message kafka.Message,
	) error
	ProfileUpdated(
		ctx context.Context,
		message kafka.Message,
	) error
}

type Inbound struct {
	log      eventbox.Logger
	inbox    eventbox.Inbox
	config   eventbox.InboxWorkerConfig
	handlers Handlers
}

func NewInbound(
	log eventbox.Logger,
	inbox eventbox.Inbox,
	config eventbox.InboxWorkerConfig,
	handlers Handlers,
) *Inbound {
	return &Inbound{
		log:      log,
		inbox:    inbox,
		config:   config,
		handlers: handlers,
	}
}

func (b *Inbound) Run(ctx context.Context) {
	worker := eventbox.NewInboxWorker(
		"TODO",
		b.log,
		b.inbox,
		b.config,
	)

	defer func() {
		worker.Stop(context.Background())
	}()

	worker.Route(evtypes.ProfileCreatedEvent, b.handlers.ProfileCreated)
	worker.Route(evtypes.ProfileDeletedEvent, b.handlers.ProfileDeleted)
	worker.Route(evtypes.ProfileUpdatedEvent, b.handlers.ProfileUpdated)

	worker.Run(ctx)
}
