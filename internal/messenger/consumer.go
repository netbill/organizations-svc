package messenger

import (
	"context"
	"sync"
	"time"

	"github.com/netbill/evebox/box/inbox"
	"github.com/netbill/evebox/consumer"
	"github.com/netbill/organizations-svc/internal/messenger/contracts"
)

type handlers interface {
	ProfileCreated(
		ctx context.Context,
		event inbox.Event,
	) inbox.EventStatus
	ProfileDeleted(
		ctx context.Context,
		event inbox.Event,
	) inbox.EventStatus
	ProfileUpdated(
		ctx context.Context,
		event inbox.Event,
	) inbox.EventStatus
}

func (m Messenger) RunConsumer(ctx context.Context, handlers handlers) {
	wg := &sync.WaitGroup{}
	run := func(f func()) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			f()
		}()
	}

	profileConsumer := consumer.New(m.log, m.db, "profiles-svc-profile-consumer", consumer.OnUnknownDoNothing, m.addr...)

	profileConsumer.Handle(contracts.ProfileCreatedEvent, handlers.ProfileCreated)
	profileConsumer.Handle(contracts.ProfileDeletedEvent, handlers.ProfileDeleted)
	profileConsumer.Handle(contracts.ProfileUpdatedEvent, handlers.ProfileUpdated)

	inboxer1 := consumer.NewInboxer(m.log, m.db, consumer.ConfigInboxer{
		Name:       "profiles-svc-inbox-worker-1",
		BatchSize:  10,
		RetryDelay: 1 * time.Minute,
		MinSleep:   100 * time.Millisecond,
		MaxSleep:   1 * time.Second,
	})
	inboxer1.Handle(contracts.ProfileCreatedEvent, handlers.ProfileCreated)
	inboxer1.Handle(contracts.ProfileDeletedEvent, handlers.ProfileDeleted)
	inboxer1.Handle(contracts.ProfileUpdatedEvent, handlers.ProfileUpdated)

	inboxer2 := consumer.NewInboxer(m.log, m.db, consumer.ConfigInboxer{
		Name:       "profiles-svc-inbox-worker-2",
		BatchSize:  10,
		RetryDelay: 1 * time.Minute,
		MinSleep:   100 * time.Millisecond,
		MaxSleep:   1 * time.Second,
	})
	inboxer2.Handle(contracts.ProfileCreatedEvent, handlers.ProfileCreated)
	inboxer2.Handle(contracts.ProfileDeletedEvent, handlers.ProfileDeleted)
	inboxer2.Handle(contracts.ProfileUpdatedEvent, handlers.ProfileUpdated)

	run(func() {
		profileConsumer.Run(ctx, contracts.OrganizationsSvcGroup, contracts.ProfilesTopicV1, m.addr...)
	})

	run(func() {
		inboxer1.Run(ctx)
	})

	run(func() {
		inboxer2.Run(ctx)
	})

	wg.Wait()
}
