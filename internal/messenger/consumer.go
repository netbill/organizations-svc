package messenger

import (
	"context"
	"sync"

	eventpg "github.com/netbill/eventbox/pg"
	"github.com/netbill/organizations-svc/internal/messenger/evtypes"
	"github.com/segmentio/kafka-go"
)

func (m *Manager) RunConsumer(ctx context.Context) {
	var wg sync.WaitGroup

	consumer := eventpg.NewConsumer(m.log, m.db, eventpg.ConsumerConfig{})

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        m.config.Brokers,
		Topic:          evtypes.AccountsTopicV1,
		GroupID:        evtypes.ProfilesSvcGroup,
		QueueCapacity:  m.config.Reader.Topics.AccountsV1.QueueCapacity,
		MaxBytes:       m.config.Reader.Topics.AccountsV1.MaxBytes,
		MinBytes:       m.config.Reader.Topics.AccountsV1.MinBytes,
		MaxWait:        m.config.Reader.Topics.AccountsV1.MaxWait,
		CommitInterval: m.config.Reader.Topics.AccountsV1.CommitInterval,
	})

	wg.Add(1)
	go func(r *kafka.Reader) {
		defer r.Close()
		defer wg.Done()

		consumer.Read(ctx, r)
	}(reader)

	wg.Wait()
}
