package messenger

import (
	"fmt"
	"strings"
	"time"

	"github.com/netbill/eventbox"
	"github.com/netbill/organizations-svc/pkg/evtypes"
	"github.com/segmentio/kafka-go"
)

type ProducerConfig struct {
	Producer string   `json:"producer"`
	Brokers  []string `json:"brokers"`

	OrganizationV1 WriterKafkaConfig `json:"organizations_v1"`
	OrgMembersV1   WriterKafkaConfig `json:"organization_members_v1"`
}

type WriterKafkaConfig struct {
	RequiredAcks string        `json:"required_acks"`
	Compression  string        `json:"compression"`
	Balancer     string        `json:"balancer"`
	BatchSize    int           `json:"batch_size"`
	BatchTimeout time.Duration `json:"batch_timeout"`
	DialTimeout  time.Duration `json:"dial_timeout"`
	IdleTimeout  time.Duration `json:"idle_timeout"`
}

func (c WriterKafkaConfig) ToWriter(topic string, brokers ...string) (*kafka.Writer, error) {
	ra, err := parseRequiredAcks(c.RequiredAcks)
	if err != nil {
		return nil, fmt.Errorf("topic %s: %w", topic, err)
	}

	bal, err := parseBalancer(c.Balancer)
	if err != nil {
		return nil, fmt.Errorf("topic %s: %w", topic, err)
	}

	comp, err := parseCompression(c.Compression)
	if err != nil {
		return nil, fmt.Errorf("topic %s: %w", topic, err)
	}

	return &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		RequiredAcks: ra,
		Balancer:     bal,
		BatchSize:    c.BatchSize,
		BatchTimeout: c.BatchTimeout,
		Compression:  comp,
		Transport: &kafka.Transport{
			DialTimeout: c.DialTimeout,
			IdleTimeout: c.IdleTimeout,
		},
	}, nil
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

func NewProducer(config ProducerConfig) *eventbox.Producer {
	prod := eventbox.NewProducer()

	wc, err := config.OrganizationV1.ToWriter(evtypes.OrganizationsTopicV1, config.Brokers...)
	if err != nil {
		panic(fmt.Errorf("topic %s: %w", evtypes.OrganizationsTopicV1, err))
	}
	if err = prod.AddTopic(wc); err != nil {
		panic(fmt.Errorf("topic %s: %w", evtypes.OrganizationsTopicV1, err))
	}

	wc, err = config.OrgMembersV1.ToWriter(evtypes.OrgMembersTopicV1, config.Brokers...)
	if err != nil {
		panic(fmt.Errorf("topic %s: %w", evtypes.OrgMembersTopicV1, err))
	}
	if err = prod.AddTopic(wc); err != nil {
		panic(fmt.Errorf("topic %s: %w", evtypes.OrgMembersTopicV1, err))
	}

	return prod
}
