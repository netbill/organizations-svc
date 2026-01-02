package consumer

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/umisto/agglomerations-svc/internal/messenger/contracts"
	"github.com/umisto/kafkakit/box"
	"github.com/umisto/kafkakit/subscriber"
	"github.com/umisto/logium"
	"golang.org/x/sync/errgroup"
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

type callbacks interface {
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
	ProfileUpdate(
		ctx context.Context,
		event box.InboxEvent,
	) string
}

type Consumer struct {
	addr      []string
	inbox     Inbox
	callbacks callbacks
	log       logium.Logger
}

func New(log logium.Logger, addr []string, inbox Inbox, callbacks callbacks) *Consumer {
	return &Consumer{
		addr:      addr,
		inbox:     inbox,
		callbacks: callbacks,
		log:       log,
	}
}

func (c Consumer) Run(ctx context.Context) {
	c.log.Info("starting events consumer", "addr", c.addr)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		accountSub := subscriber.New(c.addr, contracts.AccountsTopicV1, contracts.AgglomerationsSvcGroup)
		err := accountSub.Consume(ctx, func(m kafka.Message) (subscriber.HandlerFunc, bool) {
			et, ok := subscriber.Header(m, "event_type")
			if !ok {
				return nil, false
			}
			switch et {
			case contracts.AccountCreatedEvent:
				return c.AccountCreated, true
			case contracts.AccountDeletedEvent:
				return c.AccountDeleted, true
			case contracts.AccountUsernameChangeEvent:
				return c.AccountUsernameChanged, true
			default:
				return nil, false
			}
		})
		if err != nil {
			c.log.Warnf("accounts consumer stopped: %v", err)
		}
		return err
	})

	g.Go(func() error {
		profileSub := subscriber.New(c.addr, contracts.ProfilesTopicV1, contracts.AgglomerationsSvcGroup)
		err := profileSub.Consume(ctx, func(m kafka.Message) (subscriber.HandlerFunc, bool) {
			et, ok := subscriber.Header(m, "event_type")
			if !ok {
				return nil, false
			}
			switch et {
			case contracts.ProfileUpdatedEvent:
				return c.ProfileUpdated, true
			default:
				return nil, false
			}
		})
		if err != nil {
			c.log.Warnf("profiles consumer stopped: %v", err)
		}
		return err
	})

	g.Go(func() error {
		c.runInbox(ctx)
		return nil
	})

	_ = g.Wait()
}

func (c Consumer) runInbox(ctx context.Context) {
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
					st = c.callbacks.AccountCreated(ctx, event)
				case contracts.AccountDeletedEvent:
					st = c.callbacks.AccountDeleted(ctx, event)
				case contracts.AccountUsernameChangeEvent:
					st = c.callbacks.AccountUsernameChanged(ctx, event)
				case contracts.ProfileUpdatedEvent:
					st = c.callbacks.ProfileUpdate(ctx, event)
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
