package messenger

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/eventbox"
	"github.com/netbill/organizations-svc/internal/messenger/inbound"
	"github.com/netbill/organizations-svc/log"
	"github.com/netbill/organizations-svc/pkg/evtypes"
)

type KafkaInboxConfig struct {
	Routines       int           `json:"routines"`
	Slots          int           `json:"slots"`
	BatchSize      int           `json:"batch_size"`
	Sleep          time.Duration `json:"sleep"`
	MinNextAttempt time.Duration `json:"min_next_attempt"`
	MaxNextAttempt time.Duration `json:"max_next_attempt"`
	MaxAttempts    int32         `json:"max_attempts"`
}

type Inbound struct {
	handlers inbound.Handlers
	worker   *eventbox.InboxWorker
}

func NewInbound(
	logger *log.Logger,
	inbox eventbox.Inbox,
	handlers inbound.Handlers,
	cfg eventbox.InboxWorkerConfig,
) *Inbound {
	id := uuid.New().String()

	return &Inbound{
		handlers: handlers,
		worker:   eventbox.NewInboxWorker(id, logger, inbox, cfg),
	}
}

func (i *Inbound) RunInbox(ctx context.Context) {
	defer func() {
		i.worker.Stop(context.Background())
	}()

	i.worker.Route(evtypes.ProfileCreatedEvent, i.handlers.ProfileCreated)
	i.worker.Route(evtypes.ProfileDeletedEvent, i.handlers.ProfileDeleted)
	i.worker.Route(evtypes.ProfileUpdatedEvent, i.handlers.ProfileUpdated)

	i.worker.Run(ctx)
}
