package boot

import (
	"github.com/netbill/logium"
	"github.com/netbill/organizations-svc/internal/messenger/inbound"
	"github.com/netbill/organizations-svc/internal/messenger/outbound"
	"github.com/netbill/organizations-svc/internal/messenger/outbound/sender"
	"github.com/netbill/pgdbx"
	"github.com/segmentio/kafka-go"
)

func (c *Config) NewSender(db *pgdbx.DB) outbound.Sender {
	writer := &kafka.Writer{} //TODO
	return sender.New(c.Kafka.Identity, writer, db)
}

func (c *Config) NewOutboxWorker(
	log *logium.Entry,
	db *pgdbx.DB,
) *outbound.OutboxWorker {
	return outbound.NewOutboxWorker(log, db, outbound.OutboxWorkerConfig{
		Brokers: c.Kafka.Brokers,
		Writer: outbound.WriterConfig{
			RequiredAcks: c.Kafka.Writer.RequiredAcks,
			Compression:  c.Kafka.Writer.Compression,
			Balancer:     c.Kafka.Writer.Balancer,
			BatchSize:    c.Kafka.Writer.BatchSize,
			BatchTimeout: c.Kafka.Writer.BatchTimeout,
			DialTimeout:  c.Kafka.Writer.DialTimeout,
			IdleTimeout:  c.Kafka.Writer.IdleTimeout,
		},
		Outbox: outbound.OutboxConfig{
			Routines:       c.Kafka.Outbox.Routines,
			Slots:          c.Kafka.Outbox.Slots,
			BatchSize:      c.Kafka.Outbox.BatchSize,
			Sleep:          c.Kafka.Outbox.Sleep,
			MinNextAttempt: c.Kafka.Outbox.MinNextAttempt,
			MaxNextAttempt: c.Kafka.Outbox.MaxNextAttempt,
			MaxAttempts:    c.Kafka.Outbox.MaxAttempts,
		},
	})
}

func (c *Config) NewReceiver(
	log *logium.Entry,
	db *pgdbx.DB,
) *inbound.Receiver {
	return inbound.NewReceiver(log, db, inbound.ReceiverConfig{
		GroupID: c.Kafka.Identity,
		Brokers: c.Kafka.Brokers,
		Topics: inbound.TopicsConfig{
			Profiles: inbound.TopicReaderConfig{
				Instances:      c.Kafka.Reader.Topics.Profiles.Instances,
				MinBytes:       c.Kafka.Reader.Topics.Profiles.MinBytes,
				MaxBytes:       c.Kafka.Reader.Topics.Profiles.MaxBytes,
				MaxWait:        c.Kafka.Reader.Topics.Profiles.MaxWait,
				CommitInterval: c.Kafka.Reader.Topics.Profiles.CommitInterval,
				StartOffset:    c.Kafka.Reader.Topics.Profiles.StartOffset,
				QueueCapacity:  c.Kafka.Reader.Topics.Profiles.QueueCapacity,
			},
		},
	})
}

func (c *Config) NewInboxWorker(
	log *logium.Entry,
	db *pgdbx.DB,
	handlers inbound.Handlers,
) *inbound.Inbound {
	return inbound.NewInbox(log, db, handlers, inbound.InboxConfig{
		Routines:       c.Kafka.InboxConfig.Routines,
		Slots:          c.Kafka.InboxConfig.Slots,
		BatchSize:      c.Kafka.InboxConfig.BatchSize,
		Sleep:          c.Kafka.InboxConfig.Sleep,
		MinNextAttempt: c.Kafka.InboxConfig.MinNextAttempt,
		MaxNextAttempt: c.Kafka.InboxConfig.MaxNextAttempt,
		MaxAttempts:    c.Kafka.InboxConfig.MaxAttempts,
	})
}
