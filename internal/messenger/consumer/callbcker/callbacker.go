package callabcker

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/netbill/logium"
	"github.com/netbill/organizations-svc/internal/messenger/contracts"
	"github.com/segmentio/kafka-go"
)

type Inbox interface {
	CreateInboxEvent(
		ctx context.Context,
		message kafka.Message,
	) (box.InboxEvent, error)

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

	UpdateInboxEventStatus(
		ctx context.Context,
		id uuid.UUID,
		status string,
	) (box.InboxEvent, error)

	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type processor interface {
	AccountCreated(
		ctx context.Context,
		event box.InboxEvent,
	) string
	AccountDeleted(
		ctx context.Context,
		event box.InboxEvent,
	) string
	AccountUsernameChanged(
		ctx context.Context,
		event box.InboxEvent,
	) string
	ProfileUpdated(
		ctx context.Context,
		event box.InboxEvent,
	) string
}

type Callbacker struct {
	inbox     Inbox
	processor processor
	log       logium.Logger
}

func New(log logium.Logger, inbox Inbox, callbacks processor) *Callbacker {
	return &Callbacker{
		inbox:     inbox,
		processor: callbacks,
		log:       log,
	}
}

func (c Callbacker) RunInbox(ctx context.Context) {
	const (
		eventInboxRetryDelay = 1 * time.Minute
		batchSize            = int32(50)

		idleDelay = 3 * time.Second
		busyDelay = 100 * time.Millisecond
	)

	delay := time.Duration(0)

	for {
		if delay > 0 {
			select {
			case <-ctx.Done():
				return
			case <-time.After(delay):
			}
		} else {
			select {
			case <-ctx.Done():
				return
			default:
			}
		}

		events, err := c.inbox.GetPendingInboxEvents(ctx, batchSize)
		if err != nil {
			c.log.Errorf("failed to get pending inbox events: %v", err)
			delay = idleDelay
			continue
		}

		if len(events) == 0 {
			delay = idleDelay
			continue
		}

		processed := make([]uuid.UUID, 0, len(events))
		pending := make([]uuid.UUID, 0, len(events))
		failed := make([]uuid.UUID, 0, len(events))

		distribute := func(id uuid.UUID, status string) {
			switch status {
			case box.InboxStatusProcessed:
				processed = append(processed, id)
			case box.InboxStatusPending:
				pending = append(pending, id)
			case box.InboxStatusFailed:
				failed = append(failed, id)
			default:
				c.log.Errorf("unknown status for inbox event %s: %s", id, status)
				failed = append(failed, id)
			}
		}

		for _, event := range events {
			// per-event panic shield, so one bad event doesn't kill the whole worker
			func() {
				defer func() {
					if r := recover(); r != nil {
						c.log.Errorf("panic while handling inbox event id=%s type=%s: %v", event.ID, event.Type, r)
						failed = append(failed, event.ID)
					}
				}()

				c.log.Infof("processing inbox event: %s, type %s", event.ID, event.Type)

				var st string
				switch event.Type {
				case contracts.AccountCreatedEvent:
					st = c.processor.AccountCreated(ctx, event)
				case contracts.AccountDeletedEvent:
					st = c.processor.AccountDeleted(ctx, event)
				case contracts.AccountUsernameChangeEvent:
					st = c.processor.AccountUsernameChanged(ctx, event)
				case contracts.ProfileUpdatedEvent:
					st = c.processor.ProfileUpdated(ctx, event)
				default:
					c.log.Errorf("unknown inbox event type: %s, id: %s", event.Type, event.ID)
					st = box.InboxStatusFailed
				}

				distribute(event.ID, st)
			}()
		}

		if len(processed) > 0 {
			if _, err = c.inbox.MarkInboxEventsAsProcessed(ctx, processed); err != nil {
				c.log.Errorf("failed to mark processed: ids=%v err=%v", processed, err)
			}
		}

		if len(pending) > 0 {
			if _, err = c.inbox.MarkInboxEventsAsPending(ctx, pending, eventInboxRetryDelay); err != nil {
				c.log.Errorf("failed to mark pending: ids=%v err=%v", pending, err)
			}
		}

		if len(failed) > 0 {
			if _, err = c.inbox.MarkInboxEventsAsFailed(ctx, failed); err != nil {
				c.log.Errorf("failed to mark failed: ids=%v err=%v", failed, err)
			}
		}

		if int32(len(events)) < batchSize {
			delay = busyDelay
		} else {
			delay = 0
		}
	}
}
