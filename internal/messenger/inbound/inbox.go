package inbound

import (
	"context"
	"fmt"
	"os"

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
	worker   *eventbox.InboxWorker
	handlers Handlers
}

func NewInbox(
	log eventbox.Logger,
	inbox eventbox.Inbox,
	config eventbox.InboxWorkerConfig,
	handlers Handlers,
) *Inbound {
	id := buildProcessID("inbox")

	return &Inbound{
		worker:   eventbox.NewInboxWorker(id, log, inbox, config),
		handlers: handlers,
	}
}

func (b *Inbound) Run(ctx context.Context) {
	defer func() {
		b.worker.Stop(context.Background())
	}()

	b.worker.Route(evtypes.ProfileCreatedEvent, b.handlers.ProfileCreated)
	b.worker.Route(evtypes.ProfileDeletedEvent, b.handlers.ProfileDeleted)
	b.worker.Route(evtypes.ProfileUpdatedEvent, b.handlers.ProfileUpdated)

	b.worker.Run(ctx)
}

func buildProcessID(service string) string {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	return fmt.Sprintf("%s-%s-%d", service, hostname, os.Getpid())
}
