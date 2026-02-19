package inbound

import (
	"context"
	"sync"
	"time"

	"github.com/netbill/eventbox"
	"github.com/segmentio/kafka-go"
)

type ReceiverConfig struct {
	GroupID    string              `json:"group_id"`
	Brokers    []string            `json:"brokers"`
	Topics     []TopicReaderConfig `json:"topics"`
	MinBackoff time.Duration       `json:"min_backoff"`
	MaxBackoff time.Duration       `json:"max_backoff"`
}

type TopicReaderConfig struct {
	Topic          string        `json:"topic"`
	Instances      int           `json:"instances"`
	MinBytes       int           `json:"min_bytes"`
	MaxBytes       int           `json:"max_bytes"`
	MaxWait        time.Duration `json:"max_wait"`
	CommitInterval time.Duration `json:"commit_interval"`
	StartOffset    int64         `json:"start_offset"`
	QueueCapacity  int           `json:"queue_capacity"`
}

type Receiver struct {
	consumer *eventbox.Consumer
	config   ReceiverConfig
}

func NewReceiver(log eventbox.Logger, inbox eventbox.Inbox, config ReceiverConfig) *Receiver {
	return &Receiver{
		consumer: eventbox.NewConsumer(log, inbox, eventbox.ConsumerConfig{
			MinBackoff: config.MinBackoff,
			MaxBackoff: config.MaxBackoff,
		}),
		config: config,
	}
}

func (r *Receiver) Run(ctx context.Context) {
	var wg sync.WaitGroup

	for _, cfg := range r.config.Topics {
		cfg := cfg

		wg.Add(1)
		go func() {
			defer wg.Done()

			r.consumer.Subscribe(ctx, kafka.ReaderConfig{
				Brokers:        r.config.Brokers,
				GroupID:        r.config.GroupID,
				Topic:          cfg.Topic,
				MinBytes:       cfg.MinBytes,
				MaxBytes:       cfg.MaxBytes,
				MaxWait:        cfg.MaxWait,
				CommitInterval: cfg.CommitInterval,
				StartOffset:    cfg.StartOffset,
				QueueCapacity:  cfg.QueueCapacity,
			}, cfg.Instances)
		}()
	}

	wg.Wait()
}
