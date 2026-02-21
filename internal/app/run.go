package app

import (
	"context"
	"fmt"
	"os"
	"sync"

	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/netbill/awsx"
	"github.com/netbill/eventbox"
	eventpg "github.com/netbill/eventbox/pg"
	"github.com/netbill/organizations-svc/internal/bucket"
	"github.com/netbill/organizations-svc/internal/core/modules/invite"
	"github.com/netbill/organizations-svc/internal/core/modules/member"
	"github.com/netbill/organizations-svc/internal/core/modules/organization"
	"github.com/netbill/organizations-svc/internal/core/modules/profile"
	"github.com/netbill/organizations-svc/internal/core/modules/role"
	"github.com/netbill/organizations-svc/internal/messenger"
	"github.com/netbill/organizations-svc/internal/messenger/handler"
	"github.com/netbill/organizations-svc/internal/messenger/publisher"
	"github.com/netbill/organizations-svc/internal/repository"
	"github.com/netbill/organizations-svc/internal/repository/pg"
	"github.com/netbill/organizations-svc/internal/rest"
	"github.com/netbill/organizations-svc/internal/rest/controller"
	"github.com/netbill/organizations-svc/internal/rest/middlewares"
	"github.com/netbill/organizations-svc/internal/tokenmanager"
	"github.com/netbill/pgdbx"
	"github.com/netbill/restkit"
)

func (a *App) Run(ctx context.Context) {
	var wg = &sync.WaitGroup{}

	run := func(f func()) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			f()
		}()
	}

	pool, err := a.config.PoolDB(ctx)
	if err != nil {
		a.log.WithError(err).Error("failed to connect to database")
		os.Exit(1)
	}

	defer pool.Close()
	db := pgdbx.NewDB(pool)

	a.log.Info("starting application")

	repo := &repository.Repository{
		OrganizationsSql:          pg.NewOrganizationsQ(db),
		OrgMembersSql:             pg.NewOrgMembersQ(db),
		OrgMemberRolesSql:         pg.NewOrgMemberRolesQ(db),
		OrgRolesSql:               pg.NewOrgRolesQ(db),
		OrgRolePermissionLinksSql: pg.NewOrgRolePermissionLinksQ(db),
		OrgRolePermissionsSql:     pg.NewOrgRolePermissionsQ(db),
		OrgInvitesSql:             pg.NewOrgInvitesQ(db),
		ProfilesSql:               pg.NewProfilesQ(db),
	}

	cfg, err := awscfg.LoadDefaultConfig(
		context.Background(),
		awscfg.WithRegion(a.config.S3.Aws.Region),
		awscfg.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				a.config.S3.Aws.AccessKeyID,
				a.config.S3.Aws.SecretAccessKey,
				a.config.S3.Aws.SessionToken,
			),
		),
	)
	if err != nil {
		panic(fmt.Sprintf("unable to load S3 config: %v", err))
	}

	s3 := bucket.NewStorage(awsx.New(a.config.S3.Aws.BucketName, cfg), bucket.Config{
		LinkTTL:   a.config.S3.Media.Link.TTL,
		OrgIcon:   a.config.S3.Media.Resources.Organization.Icon,
		OrgBanner: a.config.S3.Media.Resources.Organization.Banner,
	})

	outbox := eventpg.NewOutbox(db)
	inbox := eventpg.NewInbox(db)

	producer := messenger.NewProducer(messenger.ProducerConfig{
		Producer: a.config.Kafka.Identity,
		Brokers:  a.config.Kafka.Brokers,
		OrganizationV1: messenger.WriterKafkaConfig{
			RequiredAcks: a.config.Kafka.Writer.Topics.OrganizationV1.RequiredAcks,
			Compression:  a.config.Kafka.Writer.Topics.OrganizationV1.Compression,
			Balancer:     a.config.Kafka.Writer.Topics.OrganizationV1.Balancer,
			BatchSize:    a.config.Kafka.Writer.Topics.OrganizationV1.BatchSize,
			BatchTimeout: a.config.Kafka.Writer.Topics.OrganizationV1.BatchTimeout,
		},
		OrgMembersV1: messenger.WriterKafkaConfig{
			RequiredAcks: a.config.Kafka.Writer.Topics.OrgMemberV1.RequiredAcks,
			Compression:  a.config.Kafka.Writer.Topics.OrgMemberV1.Compression,
			Balancer:     a.config.Kafka.Writer.Topics.OrgMemberV1.Balancer,
			BatchSize:    a.config.Kafka.Writer.Topics.OrgMemberV1.BatchSize,
			BatchTimeout: a.config.Kafka.Writer.Topics.OrgMemberV1.BatchTimeout,
		},
	})
	defer producer.Close()

	outbound := publisher.New(a.config.Kafka.Identity, outbox, producer)

	profileCore := profile.New(repo)
	orgCore := organization.New(repo, outbound, s3)
	orgMemberCore := member.New(repo, outbound)
	orgRoleCore := role.New(repo, outbound)
	orgInviteCore := invite.New(repo, outbound)

	tokenManager := tokenmanager.New(tokenmanager.Config{
		Issuer:   a.config.Auth.Tokens.Issuer,
		AccessSK: a.config.Auth.Tokens.AccountAccess.SecretKey,
	})

	responser := restkit.NewResponser()
	ctrl := controller.New(&controller.Modules{
		Organization: orgCore,
		Member:       orgMemberCore,
		Role:         orgRoleCore,
		Invite:       orgInviteCore,
	}, responser)
	mdll := middlewares.New(responser, tokenManager)
	router := rest.New(mdll, ctrl)

	run(func() {
		router.Run(ctx, a.log, rest.Config{
			Port:              a.config.Rest.Port,
			ReadTimeout:       a.config.Rest.Timeouts.Read,
			ReadHeaderTimeout: a.config.Rest.Timeouts.ReadHeader,
			WriteTimeout:      a.config.Rest.Timeouts.Write,
			IdleTimeout:       a.config.Rest.Timeouts.Idle,
		})
	})

	outboxWorker := messenger.NewOutboxWorker(a.log, outbox, producer, eventbox.OutboxWorkerConfig{
		Routines:       a.config.Kafka.Outbox.Routines,
		Slots:          a.config.Kafka.Outbox.Slots,
		BatchSize:      a.config.Kafka.Outbox.BatchSize,
		Sleep:          a.config.Kafka.Outbox.Sleep,
		MinNextAttempt: a.config.Kafka.Outbox.MinNextAttempt,
		MaxNextAttempt: a.config.Kafka.Outbox.MaxNextAttempt,
		MaxAttempts:    a.config.Kafka.Outbox.MaxAttempts,
	})

	run(func() {
		outboxWorker.Run(ctx)
	})

	inbound := handler.New(handler.Modules{
		Profile: profileCore,
	})

	inboxWorker := messenger.NewInboxWorker(a.log, inbox, eventbox.InboxWorkerConfig{
		Routines:       a.config.Kafka.Inbox.Routines,
		Slots:          a.config.Kafka.Inbox.Slots,
		BatchSize:      a.config.Kafka.Inbox.BatchSize,
		Sleep:          a.config.Kafka.Inbox.Sleep,
		MinNextAttempt: a.config.Kafka.Inbox.MinNextAttempt,
		MaxNextAttempt: a.config.Kafka.Inbox.MaxNextAttempt,
		MaxAttempts:    a.config.Kafka.Inbox.MaxAttempts,
	}, *inbound)

	run(func() {
		inboxWorker.Run(ctx)
	})

	consumer := messenger.NewConsumer(a.log, inbox, messenger.ListenerConfig{
		GroupID:    a.config.Kafka.Identity,
		Brokers:    a.config.Kafka.Brokers,
		MinBackoff: a.config.Kafka.Reader.Backoff.Min,
		MaxBackoff: a.config.Kafka.Reader.Backoff.Max,
		ProfilesV1: messenger.ListenKafkaConfig{
			Instances:     a.config.Kafka.Reader.Topics.ProfilesV1.Instances,
			MinBytes:      a.config.Kafka.Reader.Topics.ProfilesV1.MinBytes,
			MaxBytes:      a.config.Kafka.Reader.Topics.ProfilesV1.MaxBytes,
			MaxWait:       a.config.Kafka.Reader.Topics.ProfilesV1.MaxWait,
			StartOffset:   a.config.Kafka.Reader.Topics.ProfilesV1.StartOffset,
			QueueCapacity: a.config.Kafka.Reader.Topics.ProfilesV1.QueueCapacity,
		},
	})
	defer consumer.Close()

	run(func() {
		consumer.Run(ctx)
	})

	wg.Wait()
}
