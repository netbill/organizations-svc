package outbound

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/eventbox"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/internal/core/modules/role"
	"github.com/netbill/organizations-svc/internal/messenger/outbound/sender"
	"github.com/segmentio/kafka-go"
)

type Producer interface {
	WriteOrganizationCreated(ctx context.Context, organization models.Organization) error
	WriteOrganizationActivated(ctx context.Context, organization models.Organization) error
	WriteOrganizationDeactivated(ctx context.Context, organization models.Organization) error
	WriteOrganizationUpdated(ctx context.Context, organization models.Organization) error
	WriteOrganizationDeleted(ctx context.Context, organization models.Organization) error

	WriteOrgMemberCreated(ctx context.Context, member models.Member) error
	WriteOrgMemberUpdated(ctx context.Context, member models.Member) error
	WriteOrgMemberDeleted(ctx context.Context, memberID uuid.UUID) error

	WriteOrgInviteCreated(ctx context.Context, invite models.Invite) error
	WriteOrgInviteAccepted(ctx context.Context, invite models.Invite) error
	WriteOrgInviteDeclined(ctx context.Context, invite models.Invite) error
	WriteOrgInviteDeleted(ctx context.Context, invite models.Invite) error

	WriteOrgRoleCreated(ctx context.Context, role models.Role) error
	WriteOrgRoleUpdated(ctx context.Context, role models.Role) error
	WriteOrgRoleDeleted(ctx context.Context, role models.Role) error

	WriteOrgRolesRanksUpdated(
		ctx context.Context,
		organizationID uuid.UUID,
		order map[uuid.UUID]uint,
	) error
	WriteOrgRolePermissionsUpdated(
		ctx context.Context,
		role models.Role,
		params role.SetPermissions,
	) error

	WriteOrgMemberRoleAdd(
		ctx context.Context,
		link models.OrgMemberRolesLink,
	) error
	WriteOrgMemberRoleRemove(
		ctx context.Context,
		memberID uuid.UUID,
		roleID uuid.UUID,
	) error

	Close() error
}

type ProducerConfig struct {
	Addrs  []string                       `mapstructure:"addrs"`
	Topics map[string]TopicProducerConfig `mapstructure:"topics"`
}

type TopicProducerConfig struct {
	RequiredAcks string        `mapstructure:"required_acks"`
	Compression  string        `mapstructure:"compression"`
	Balancer     string        `mapstructure:"balancer"`
	BatchSize    int           `mapstructure:"batch_size"`
	BatchTimeout time.Duration `mapstructure:"batch_timeout"`
	DialTimeout  time.Duration `mapstructure:"dial_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

func NewProducer(outbox eventbox.Outbox, cfg ProducerConfig) (Producer, error) {
	id := fmt.Sprintf("producer-%d", time.Now().UnixNano())

	prod := eventbox.NewProducer()

	addrs := make([]string, 0, len(cfg.Addrs))
	for _, a := range cfg.Addrs {
		a = strings.TrimSpace(a)
		if a != "" {
			addrs = append(addrs, a)
		}
	}
	if len(addrs) == 0 {
		return nil, fmt.Errorf("kafka addrs are empty")
	}

	for topic, c := range cfg.Topics {
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

		w := &kafka.Writer{
			Addr:         kafka.TCP(addrs...),
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
		}

		prod.AddTopic(topic, w)
	}

	return sender.New(id, outbox, prod), nil
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
