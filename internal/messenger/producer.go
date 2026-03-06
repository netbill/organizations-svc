package messenger

import (
	"time"

	"github.com/netbill/eventbox"
	"github.com/netbill/evtypes"
	"github.com/netbill/organizations-svc/pkg/log"
)

type ProducerConfig struct {
	Producer string   `json:"producer"`
	Brokers  []string `json:"brokers"`

	OrganizationV1 ProduceKafkaConfig `json:"organizations_v1"`
	OrgMembersV1   ProduceKafkaConfig `json:"organization_members_v1"`
}

type ProduceKafkaConfig struct {
	RequiredAcks string        `json:"required_acks"`
	Compression  string        `json:"compression"`
	Balancer     string        `json:"balancer"`
	BatchSize    int           `json:"batch_size"`
	BatchTimeout time.Duration `json:"batch_timeout"`
}

func NewProducer(log *log.Logger, cfg ProducerConfig) (*eventbox.Producer, error) {
	producer := eventbox.NewProducer(log, cfg.Brokers...)

	err := producer.AddWriter(evtypes.OrganizationsTopicV1, eventbox.WriterTopicConfig{
		RequiredAcks: cfg.OrganizationV1.RequiredAcks,
		Compression:  cfg.OrganizationV1.Compression,
		Balancer:     cfg.OrganizationV1.Balancer,
		BatchSize:    cfg.OrganizationV1.BatchSize,
		BatchTimeout: cfg.OrganizationV1.BatchTimeout,
	})
	if err != nil {
		return nil, err
	}

	err = producer.AddWriter(evtypes.OrgMembersTopicV1, eventbox.WriterTopicConfig{
		RequiredAcks: cfg.OrgMembersV1.RequiredAcks,
		Compression:  cfg.OrgMembersV1.Compression,
		Balancer:     cfg.OrgMembersV1.Balancer,
		BatchSize:    cfg.OrgMembersV1.BatchSize,
		BatchTimeout: cfg.OrgMembersV1.BatchTimeout,
	})
	if err != nil {
		return nil, err
	}

	return producer, nil
}
