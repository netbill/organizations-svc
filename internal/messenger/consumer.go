package messenger

import (
	"context"
	"sync"
	"time"

	eventpg "github.com/netbill/eventbox/pg"
	"github.com/netbill/organizations-svc/pkg/evtypes"
	"github.com/segmentio/kafka-go"
)

func (m *Manager) RunConsumer(ctx context.Context) {
	var wg sync.WaitGroup

	consumer := eventpg.NewConsumer(m.log, m.db, eventpg.ConsumerConfig{
		MinBackoff: time.Millisecond * 100,
		MaxBackoff: time.Second * 10,
	})

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        m.config.Brokers,
		Topic:          evtypes.ProfilesTopicV1,
		GroupID:        OrganizationsSvcGroup,
		QueueCapacity:  m.config.Reader.Topics.ProfilesV1.QueueCapacity,
		MaxBytes:       m.config.Reader.Topics.ProfilesV1.MaxBytes,
		MinBytes:       m.config.Reader.Topics.ProfilesV1.MinBytes,
		MaxWait:        m.config.Reader.Topics.ProfilesV1.MaxWait,
		CommitInterval: m.config.Reader.Topics.ProfilesV1.CommitInterval,
	})

	wg.Add(1)
	go func(r *kafka.Reader) {
		defer r.Close()
		defer wg.Done()

		consumer.Read(ctx, r)
	}(reader)

	wg.Wait()
}
