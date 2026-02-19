package app

import (
	"fmt"
	"strings"
	"time"

	"github.com/netbill/eventbox"
	"github.com/netbill/organizations-svc/pkg/evtypes"
	"github.com/segmentio/kafka-go"
)

type WriterKafkaConfig struct {
	RequiredAcks string        `mapstructure:"required_acks"`
	Compression  string        `mapstructure:"compression"`
	Balancer     string        `mapstructure:"balancer"`
	BatchSize    int           `mapstructure:"batch_size"`
	BatchTimeout time.Duration `mapstructure:"batch_timeout"`
	DialTimeout  time.Duration `mapstructure:"dial_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

func buildWriter(topic string, addrs []string, tc WriterKafkaConfig) *kafka.Writer {
	ra, err := parseRequiredAcks(tc.RequiredAcks)
	if err != nil {
		panic(fmt.Errorf("topic %s: %w", topic, err))
	}

	bal, err := parseBalancer(tc.Balancer)
	if err != nil {
		panic(fmt.Errorf("topic %s: %w", topic, err))
	}

	comp, err := parseCompression(tc.Compression)
	if err != nil {
		panic(fmt.Errorf("topic %s: %w", topic, err))
	}

	return &kafka.Writer{
		Addr:         kafka.TCP(addrs...),
		Topic:        topic,
		RequiredAcks: ra,
		Balancer:     bal,
		BatchSize:    tc.BatchSize,
		BatchTimeout: tc.BatchTimeout,
		Compression:  comp,
		Transport: &kafka.Transport{
			DialTimeout: tc.DialTimeout,
			IdleTimeout: tc.IdleTimeout,
		},
	}
}

func (a *App) BuildProducer() *eventbox.Producer {
	prod := eventbox.NewProducer()

	err := prod.AddTopic(buildWriter(evtypes.OrganizationsTopicV1, a.config.Kafka.Brokers, WriterKafkaConfig{
		RequiredAcks: a.config.Kafka.Writer.Topics.OrganizationV1.RequiredAcks,
		Compression:  a.config.Kafka.Writer.Topics.OrganizationV1.Compression,
		Balancer:     a.config.Kafka.Writer.Topics.OrganizationV1.Balancer,
		BatchSize:    a.config.Kafka.Writer.Topics.OrganizationV1.BatchSize,
		BatchTimeout: a.config.Kafka.Writer.Topics.OrganizationV1.BatchTimeout,
		DialTimeout:  a.config.Kafka.Writer.Topics.OrganizationV1.DialTimeout,
		IdleTimeout:  a.config.Kafka.Writer.Topics.OrganizationV1.IdleTimeout,
	}))
	if err != nil {
		panic(fmt.Errorf("profiles topic: %w", err))
	}

	err = prod.AddTopic(buildWriter(evtypes.OrgMembersTopicV1, a.config.Kafka.Brokers, WriterKafkaConfig{
		RequiredAcks: a.config.Kafka.Writer.Topics.OrgMemberV1.RequiredAcks,
		Compression:  a.config.Kafka.Writer.Topics.OrgMemberV1.Compression,
		Balancer:     a.config.Kafka.Writer.Topics.OrgMemberV1.Balancer,
		BatchSize:    a.config.Kafka.Writer.Topics.OrgMemberV1.BatchSize,
		BatchTimeout: a.config.Kafka.Writer.Topics.OrgMemberV1.BatchTimeout,
		DialTimeout:  a.config.Kafka.Writer.Topics.OrgMemberV1.DialTimeout,
		IdleTimeout:  a.config.Kafka.Writer.Topics.OrgMemberV1.IdleTimeout,
	}))
	if err != nil {
		panic(fmt.Errorf("profiles topic: %w", err))
	}

	return prod
}

func parseRequiredAcks(v string) (kafka.RequiredAcks, error) {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "", "all", "-1":
		return kafka.RequireAll, nil
	case "none", "0":
		return kafka.RequireNone, nil
	case "one", "1":
		return kafka.RequireOne, nil
	default:
		return 0, fmt.Errorf("invalid required_acks: %q", v)
	}
}

func parseBalancer(v string) (kafka.Balancer, error) {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "", "hash":
		return &kafka.Hash{}, nil
	case "leastbytes", "least_bytes":
		return &kafka.LeastBytes{}, nil
	case "roundrobin", "round_robin":
		return &kafka.RoundRobin{}, nil
	default:
		return nil, fmt.Errorf("invalid balancer: %q", v)
	}
}

func parseCompression(v string) (kafka.Compression, error) {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "", "none":
		return kafka.Snappy, nil
	case "gzip":
		return kafka.Gzip, nil
	case "snappy":
		return kafka.Snappy, nil
	case "lz4":
		return kafka.Lz4, nil
	case "zstd":
		return kafka.Zstd, nil
	default:
		return 0, fmt.Errorf("invalid compression: %q", v)
	}
}
