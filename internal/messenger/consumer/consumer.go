package consumer

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/logium"
	"github.com/netbill/msgx/subscriber"
	"github.com/netbill/organizations-svc/internal/messenger/contracts"
	"github.com/segmentio/kafka-go"
	"golang.org/x/sync/errgroup"
)

type Inbox interface {
	CreateInboxEvent(
		ctx context.Context,
		message kafka.Message,
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

type callbacks interface {
	AccountCreated(ctx context.Context, event kafka.Message) error
	AccountDeleted(ctx context.Context, event kafka.Message) error
	AccountUsernameChanged(ctx context.Context, event kafka.Message) error
	ProfileUpdated(ctx context.Context, event kafka.Message) error

	RunInbox(ctx context.Context)
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
		accountSub := subscriber.New(c.addr, contracts.AccountsTopicV1, contracts.OrganizationsSvcGroup)
		err := accountSub.Consume(ctx, func(m kafka.Message) (subscriber.HandlerFunc, bool) {
			et, ok := subscriber.Header(m, "event_type")
			if !ok {
				return nil, false
			}
			switch et {
			case contracts.AccountCreatedEvent:
				return c.callbacks.AccountCreated, true
			case contracts.AccountDeletedEvent:
				return c.callbacks.AccountDeleted, true
			case contracts.AccountUsernameChangeEvent:
				return c.callbacks.AccountUsernameChanged, true
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
		profileSub := subscriber.New(c.addr, contracts.ProfilesTopicV1, contracts.OrganizationsSvcGroup)
		err := profileSub.Consume(ctx, func(m kafka.Message) (subscriber.HandlerFunc, bool) {
			et, ok := subscriber.Header(m, "event_type")
			if !ok {
				return nil, false
			}
			switch et {
			case contracts.ProfileUpdatedEvent:
				return c.callbacks.ProfileUpdated, true
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
		c.callbacks.RunInbox(ctx)
		return nil
	})

	_ = g.Wait()
}
