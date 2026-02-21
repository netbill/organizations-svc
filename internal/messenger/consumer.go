package messenger

import (
	"time"

	"github.com/netbill/eventbox"
	"github.com/netbill/organizations-svc/log"
	"github.com/netbill/organizations-svc/pkg/evtypes"
)

type ListenerConfig struct {
	GroupID string   `json:"group_id"`
	Brokers []string `json:"brokers"`

	MinBackoff time.Duration `json:"min_backoff"`
	MaxBackoff time.Duration `json:"max_backoff"`

	ProfilesV1 ListenKafkaConfig `json:"profiles_v1"`
}

type ListenKafkaConfig struct {
	Instances     int           `json:"instances"`
	MinBytes      int           `json:"min_bytes"`
	MaxBytes      int           `json:"max_bytes"`
	MaxWait       time.Duration `json:"max_wait"`
	StartOffset   int64         `json:"start_offset"`
	QueueCapacity int           `json:"queue_capacity"`
}

func NewConsumer(
	logger *log.Logger,
	inbox eventbox.Inbox,
	config ListenerConfig,
) *eventbox.Consumer {
	consumer := eventbox.NewConsumer(logger, inbox, eventbox.ConsumerConfig{
		MinBackoff: config.MinBackoff,
		MaxBackoff: config.MaxBackoff,
	})

	consumer.AddReader(eventbox.ReaderConfig{
		Brokers:  config.Brokers,
		GroupID:  config.GroupID,
		Topic:    evtypes.OrgMembersTopicV1,
		MinBytes: config.ProfilesV1.MinBytes,
		MaxBytes: config.ProfilesV1.MaxBytes,
	}, config.ProfilesV1.Instances)

	return consumer
}
