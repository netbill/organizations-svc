package messenger

import (
	"context"
	"sync"
	"time"

	"github.com/netbill/eventbox"
	"github.com/netbill/organizations-svc/log"
	"github.com/netbill/organizations-svc/pkg/evtypes"
	"github.com/segmentio/kafka-go"
)

type ConsumerConfig struct {
	GroupID string   `json:"group_id"`
	Brokers []string `json:"brokers"`

	MinBackoff time.Duration `json:"min_backoff"`
	MaxBackoff time.Duration `json:"max_backoff"`

	ProfilesV1 ReaderKafkaConfig `json:"profiles_v1"`
}

type ReaderKafkaConfig struct {
	Instances      int           `json:"instances"`
	MinBytes       int           `json:"min_bytes"`
	MaxBytes       int           `json:"max_bytes"`
	MaxWait        time.Duration `json:"max_wait"`
	CommitInterval time.Duration `json:"commit_interval"`
	StartOffset    int64         `json:"start_offset"`
	QueueCapacity  int           `json:"queue_capacity"`
}

type ConsumerWorker struct {
	consumer *eventbox.Consumer
	config   ConsumerConfig
}

func NewConsumerWorker(
	logger *log.Logger,
	inbox eventbox.Inbox,
	config ConsumerConfig,
) *ConsumerWorker {
	consumer := eventbox.NewConsumer(logger, inbox, eventbox.ConsumerConfig{
		MinBackoff: config.MinBackoff,
		MaxBackoff: config.MaxBackoff,
	})

	return &ConsumerWorker{
		consumer: consumer,
		config:   config,
	}
}

func (w *ConsumerWorker) Run(ctx context.Context) {
	var wg sync.WaitGroup

	for i := 0; i < w.config.ProfilesV1.Instances; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			w.reading(ctx, kafka.ReaderConfig{
				Brokers:        w.config.Brokers,
				GroupID:        w.config.GroupID,
				Topic:          evtypes.ProfilesTopicV1,
				MinBytes:       w.config.ProfilesV1.MinBytes,
				MaxBytes:       w.config.ProfilesV1.MaxBytes,
				MaxWait:        w.config.ProfilesV1.MaxWait,
				StartOffset:    w.config.ProfilesV1.StartOffset,
				CommitInterval: w.config.ProfilesV1.CommitInterval,
				QueueCapacity:  w.config.ProfilesV1.QueueCapacity,
			})
		}()
	}

	wg.Wait()
}

func (w *ConsumerWorker) reading(ctx context.Context, config kafka.ReaderConfig) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        config.Brokers,
		GroupID:        config.GroupID,
		Topic:          evtypes.ProfilesTopicV1,
		MinBytes:       config.MinBytes,
		MaxBytes:       config.MaxBytes,
		MaxWait:        config.MaxWait,
		StartOffset:    config.StartOffset,
		CommitInterval: 0,
	})
	defer reader.Close()

	w.consumer.Subscribe(ctx, reader)
}
