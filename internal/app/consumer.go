package app

import (
	"context"
	"sync"
	"time"

	"github.com/netbill/eventbox"
	eventpg "github.com/netbill/eventbox/pg"
	"github.com/netbill/organizations-svc/pkg/evtypes"
	"github.com/segmentio/kafka-go"
)

type ReaderKafkaConfig struct {
	Instances      int           `mapstructure:"instances"`
	MinBytes       int           `mapstructure:"min_bytes"`
	MaxBytes       int           `mapstructure:"max_bytes"`
	MaxWait        time.Duration `mapstructure:"max_wait"`
	CommitInterval time.Duration `mapstructure:"commit_interval"`
	StartOffset    int64         `mapstructure:"start_offset"`
	QueueCapacity  int           `mapstructure:"queue_capacity"`
}

func (a *App) RunConsumer(ctx context.Context) {
	var wg sync.WaitGroup

	inbox := eventpg.NewInbox(a.db)

	consumer := eventbox.NewConsumer(a.log, inbox, eventbox.ConsumerConfig{
		MinBackoff: 100 * time.Millisecond,
		MaxBackoff: 5 * time.Second,
	})

	subscribe := func(instances int, cfg kafka.ReaderConfig) {
		for i := 0; i < instances; i++ {
			wg.Add(1)
			go func(cfg kafka.ReaderConfig) {
				defer wg.Done()

				r := kafka.NewReader(cfg)
				defer r.Close()

				consumer.Subscribe(ctx, r)
			}(cfg)
		}
	}

	subscribe(a.config.Kafka.Reader.Topics.ProfilesV1.Instances, kafka.ReaderConfig{
		Brokers:        a.config.Kafka.Brokers,
		GroupID:        a.config.Kafka.Identity,
		Topic:          evtypes.ProfilesTopicV1,
		MinBytes:       a.config.Kafka.Reader.Topics.ProfilesV1.MinBytes,
		MaxBytes:       a.config.Kafka.Reader.Topics.ProfilesV1.MaxBytes,
		MaxWait:        a.config.Kafka.Reader.Topics.ProfilesV1.MaxWait,
		StartOffset:    a.config.Kafka.Reader.Topics.ProfilesV1.StartOffset,
		CommitInterval: 0,
	})

	wg.Wait()
}
