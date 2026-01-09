package messenger

import (
	"context"
	"sync"
	"time"

	"github.com/netbill/evebox/box/inbox"
	"github.com/netbill/evebox/consumer"
	"github.com/netbill/logium"
	"github.com/netbill/organizations-svc/internal/messenger/contracts"
)

type Outbound interface {
	AccountCreated(
		ctx context.Context,
		event inbox.Event,
	) inbox.EventStatus
	AccountDeleted(
		ctx context.Context,
		event inbox.Event,
	) inbox.EventStatus
	AccountUsernameChanged(
		ctx context.Context,
		event inbox.Event,
	) inbox.EventStatus
	AccountProfileUpdated(
		ctx context.Context,
		event inbox.Event,
	) inbox.EventStatus
}

type Consumer struct {
	addr     []string
	log      logium.Logger
	inbox    inbox.Box
	handlers Outbound
}

func NewConsumer(
	log logium.Logger,
	inbox inbox.Box,
	handlers Outbound,
	addr ...string,
) Consumer {
	return Consumer{
		addr:     addr,
		log:      log,
		inbox:    inbox,
		handlers: handlers,
	}
}

func (c Consumer) Run(ctx context.Context) {
	wg := &sync.WaitGroup{}
	run := func(f func()) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			f()
		}()
	}

	accountConsumer := consumer.New(c.log, "profiles-svc-account-consumer", c.inbox)

	accountConsumer.Handle(contracts.AccountCreatedEvent, c.handlers.AccountCreated)
	accountConsumer.Handle(contracts.AccountDeletedEvent, c.handlers.AccountDeleted)
	accountConsumer.Handle(contracts.AccountUsernameChangedEvent, c.handlers.AccountUsernameChanged)
	accountConsumer.Handle(contracts.AccountProfileUpdatedEvent, c.handlers.AccountUsernameChanged)

	inboxer1 := inbox.NewWorker(c.log, c.inbox, inbox.ConfigWorker{
		Name:       "profiles-svc-inbox-worker-1",
		BatchSize:  10,
		RetryDelay: 1 * time.Minute,
		MinSleep:   100 * time.Millisecond,
		MaxSleep:   1 * time.Second,
	})
	inboxer1.Handle(contracts.AccountCreatedEvent, c.handlers.AccountCreated)
	inboxer1.Handle(contracts.AccountDeletedEvent, c.handlers.AccountDeleted)
	inboxer1.Handle(contracts.AccountUsernameChangedEvent, c.handlers.AccountUsernameChanged)
	inboxer1.Handle(contracts.AccountProfileUpdatedEvent, c.handlers.AccountProfileUpdated)

	inboxer2 := inbox.NewWorker(c.log, c.inbox, inbox.ConfigWorker{
		Name:       "profiles-svc-inbox-worker-2",
		BatchSize:  10,
		RetryDelay: 1 * time.Minute,
		MinSleep:   100 * time.Millisecond,
		MaxSleep:   1 * time.Second,
	})
	inboxer2.Handle(contracts.AccountCreatedEvent, c.handlers.AccountCreated)
	inboxer2.Handle(contracts.AccountDeletedEvent, c.handlers.AccountDeleted)
	inboxer2.Handle(contracts.AccountUsernameChangedEvent, c.handlers.AccountUsernameChanged)
	inboxer2.Handle(contracts.AccountProfileUpdatedEvent, c.handlers.AccountProfileUpdated)

	run(func() {
		accountConsumer.Run(ctx, contracts.OrganizationsSvcGroup, contracts.AccountsTopicV1, c.addr...)
	})

	run(func() {
		inboxer1.Run(ctx)
	})

	run(func() {
		inboxer2.Run(ctx)
	})

	wg.Wait()
}
