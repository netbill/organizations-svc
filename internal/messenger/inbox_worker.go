package messenger

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/eventbox"
	"github.com/netbill/evtypes"
	"github.com/netbill/organizations-svc/pkg/log"
)

type handlers interface {
	ProfileCreated(
		ctx context.Context,
		event eventbox.InboxEvent,
	) error
	ProfileUpdated(
		ctx context.Context,
		event eventbox.InboxEvent,
	) error
	ProfileDeleted(
		ctx context.Context,
		event eventbox.InboxEvent,
	) error

	PlaceCreated(
		ctx context.Context,
		event eventbox.InboxEvent,
	) error
	PlaceDeleted(
		ctx context.Context,
		event eventbox.InboxEvent,
	) error
}

func NewInboxWorker(
	logger *log.Logger,
	inbox eventbox.Inbox,
	cfg eventbox.InboxWorkerConfig,
	handlers handlers,
) *eventbox.InboxWorker {
	id := uuid.New().String()

	worker := eventbox.NewInboxWorker(id, logger, inbox, cfg)

	worker.Route(evtypes.ProfileCreatedEvent, handlers.ProfileCreated)
	worker.Route(evtypes.ProfileDeletedEvent, handlers.ProfileDeleted)
	worker.Route(evtypes.ProfileUpdatedEvent, handlers.ProfileUpdated)

	worker.Route(evtypes.PlaceCreatedEvent, handlers.PlaceCreated)
	worker.Route(evtypes.PlaceDeletedEvent, handlers.PlaceDeleted)

	return worker
}
