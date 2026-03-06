package messenger

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/eventbox"
	"github.com/netbill/evtypes"
	"github.com/netbill/organizations-svc/pkg/log"
)

type profileController interface {
	Created(ctx context.Context, event eventbox.InboxEvent) error
	Updated(ctx context.Context, event eventbox.InboxEvent) error
	Deleted(ctx context.Context, event eventbox.InboxEvent) error
}

type placeController interface {
	Created(ctx context.Context, event eventbox.InboxEvent) error
	Deleted(ctx context.Context, event eventbox.InboxEvent) error
}

type InboxWorkerDeps struct {
	ProfileController profileController
	PlaceController   placeController

	Logger *log.Logger
	Inbox  eventbox.Inbox
	Config eventbox.InboxWorkerConfig
}

func NewInboxWorker(deps InboxWorkerDeps) *eventbox.InboxWorker {
	worker := eventbox.NewInboxWorker(uuid.New().String(), deps.Logger, deps.Inbox, deps.Config)

	worker.Route(evtypes.ProfileCreatedEvent, deps.ProfileController.Created)
	worker.Route(evtypes.ProfileDeletedEvent, deps.ProfileController.Deleted)
	worker.Route(evtypes.ProfileUpdatedEvent, deps.ProfileController.Updated)

	worker.Route(evtypes.PlaceCreatedEvent, deps.PlaceController.Created)
	worker.Route(evtypes.PlaceDeletedEvent, deps.PlaceController.Deleted)

	return worker
}
