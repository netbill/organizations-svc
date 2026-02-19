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
	"github.com/segmentio/kafka-go"
)

type Publisher interface {
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
	Brokers  []string            `mapstructure:"brokers"`
	Identity string              `mapstructure:"identity"`
	Topics   []WriterKafkaConfig `mapstructure:"topics"`
}

type WriterKafkaConfig struct {
	Topic        string        `mapstructure:"topic"`
	RequiredAcks string        `mapstructure:"required_acks"`
	Compression  string        `mapstructure:"compression"`
	Balancer     string        `mapstructure:"balancer"`
	BatchSize    int           `mapstructure:"batch_size"`
	BatchTimeout time.Duration `mapstructure:"batch_timeout"`
	DialTimeout  time.Duration `mapstructure:"dial_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

func NewProducer(cfg ProducerConfig) *eventbox.Producer {
	producer := eventbox.NewProducer()
	for _, tc := range cfg.Topics {
		w, err := buildWriter(cfg.Brokers, tc)
		if err != nil {
			panic(fmt.Errorf("topic %s: %w", tc.Topic, err))
		}
		if err := producer.AddTopic(tc.Topic, w); err != nil {
			panic(fmt.Errorf("topic %s: %w", tc.Topic, err))
		}
	}

	return producer
}
