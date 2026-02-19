package inbound

import (
	"context"
	"fmt"
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
	log    eventbox.Logger
	inbox  eventbox.Inbox
	config ReceiverConfig
}

func NewReceiver(
	log eventbox.Logger,
	inbox eventbox.Inbox,
	config ReceiverConfig,
) *Receiver {
	return &Receiver{
		log:    log,
		inbox:  inbox,
		config: config,
	}
}

func (r *Receiver) Run(ctx context.Context) {
	consumer := eventbox.NewConsumer(r.log, r.inbox, eventbox.ConsumerConfig{
		MinBackoff: r.config.MinBackoff,
		MaxBackoff: r.config.MaxBackoff,
	})

	for _, topicCfg := range r.config.Topics {
		err := consumer.AddTopic(eventbox.TopicReaderConfig{
			Instances: topicCfg.Instances,
			Reader: kafka.ReaderConfig{
				Brokers:     r.config.Brokers,
				GroupID:     r.config.GroupID,
				Topic:       topicCfg.Topic,
				MinBytes:    topicCfg.MinBytes,
				MaxBytes:    topicCfg.MaxBytes,
				MaxWait:     topicCfg.MaxWait,
				StartOffset: topicCfg.StartOffset,
			},
		})
		if err != nil {
			panic(fmt.Errorf("failed to add topic reader config: %v", err))
		}
	}

	consumer.Run(ctx)
}
