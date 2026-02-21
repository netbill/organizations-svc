package messenger

import (
	"github.com/google/uuid"
	"github.com/netbill/eventbox"
	"github.com/netbill/organizations-svc/internal/messenger/handler"
	"github.com/netbill/organizations-svc/log"
	"github.com/netbill/organizations-svc/pkg/evtypes"
)

func NewInboxWorker(
	logger *log.Logger,
	inbox eventbox.Inbox,
	cfg eventbox.InboxWorkerConfig,
	handlers handler.Handler,
) *eventbox.InboxWorker {
	id := uuid.New().String()

	worker := eventbox.NewInboxWorker(id, logger, inbox, cfg)

	worker.Route(evtypes.ProfileCreatedEvent, handlers.ProfileCreated)
	worker.Route(evtypes.ProfileDeletedEvent, handlers.ProfileDeleted)
	worker.Route(evtypes.ProfileUpdatedEvent, handlers.ProfileUpdated)

	return worker
}
