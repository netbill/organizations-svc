package inboxer

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/messenger/consumer/inboxer/handler"
	"github.com/umisto/kafkakit/box"
	"github.com/umisto/logium"
)

type Service struct {
	log    logium.Logger
	router router
	inbox  inbox
}

type router interface {
	Handle(ctx context.Context, event box.InboxEvent) handler.Status
}

type inbox interface {
	GetInboxEventByID(
		ctx context.Context,
		id uuid.UUID,
	) (box.InboxEvent, error)

	GetPendingInboxEvents(
		ctx context.Context,
		limit int32,
	) ([]box.InboxEvent, error)

	MarkInboxEventsAsProcessed(
		ctx context.Context,
		ids []uuid.UUID,
	) ([]box.InboxEvent, error)

	MarkInboxEventsAsFailed(
		ctx context.Context,
		ids []uuid.UUID,
	) ([]box.InboxEvent, error)

	MarkInboxEventsAsPending(
		ctx context.Context,
		ids []uuid.UUID,
		delay time.Duration,
	) ([]box.InboxEvent, error)
}

func New(
	log logium.Logger,
	router router,
	inbox inbox,

) Service {
	return Service{
		log:    log,
		router: router,
		inbox:  inbox,
	}
}

const eventInboxRetryDelay = 1 * time.Minute

func (s Service) Run(ctx context.Context) {
	const (
		batchSize = int32(50)
		idleDelay = 3 * time.Second
		busyDelay = 100 * time.Millisecond
	)

	timer := time.NewTimer(0)
	defer timer.Stop()

	delay := time.Duration(0)

	for {
		if delay > 0 {
			timer.Reset(delay)
		} else {
			timer.Reset(0)
		}

		select {
		case <-ctx.Done():
			return
		case <-timer.C:
		}

		events, err := s.inbox.GetPendingInboxEvents(ctx, batchSize)
		if err != nil {
			s.log.Errorf("failed to get pending inbox events: %v", err)
			delay = idleDelay
			continue
		}

		if len(events) == 0 {
			delay = idleDelay
			continue
		}

		processed := make([]uuid.UUID, 0, len(events))
		delayed := make([]uuid.UUID, 0, len(events))
		failed := make([]uuid.UUID, 0, len(events))

		for _, event := range events {
			s.log.Infof("processing inbox event: %s, type %s", event.ID, event.Type)

			st := s.router.Handle(ctx, event)

			switch st {
			case handler.StatusProcessed:
				processed = append(processed, event.ID)

			case handler.StatusDelayed:
				delayed = append(delayed, event.ID)

			case handler.StatusFailed:
				failed = append(failed, event.ID)

			default:
				s.log.Errorf("unknown handler status: %d, id=%s, type=%s", st, event.ID, event.Type)
				failed = append(failed, event.ID)
			}
		}

		if len(processed) > 0 {
			if _, err = s.inbox.MarkInboxEventsAsProcessed(ctx, processed); err != nil {
				s.log.Errorf("failed to mark processed: ids=%v err=%v", processed, err)
			}
		}

		if len(delayed) > 0 {
			if _, err = s.inbox.MarkInboxEventsAsPending(ctx, delayed, eventInboxRetryDelay); err != nil {
				s.log.Errorf("failed to mark pending: ids=%v err=%v", delayed, err)
			}
		}

		if len(failed) > 0 {
			if _, err = s.inbox.MarkInboxEventsAsFailed(ctx, failed); err != nil {
				s.log.Errorf("failed to mark failed: ids=%v err=%v", failed, err)
			}
		}

		if len(events) < int(batchSize) {
			delay = busyDelay
		} else {
			delay = 0
		}
	}
}
