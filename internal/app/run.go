package app

import (
	"context"
	"fmt"
	"sync"

	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/netbill/awsx"
	"github.com/netbill/eventbox"
	eventpg "github.com/netbill/eventbox/pg"
	"github.com/netbill/organizations-svc/internal/core"
	"github.com/netbill/organizations-svc/internal/media"
	"github.com/netbill/organizations-svc/internal/messenger"
	"github.com/netbill/organizations-svc/internal/messenger/evcontroller"
	"github.com/netbill/organizations-svc/internal/messenger/publisher"
	"github.com/netbill/organizations-svc/internal/repository"
	"github.com/netbill/organizations-svc/internal/repository/pg"
	"github.com/netbill/organizations-svc/internal/rest"
	"github.com/netbill/organizations-svc/internal/rest/controller"
	"github.com/netbill/organizations-svc/internal/rest/middlewares"
	"github.com/netbill/organizations-svc/internal/tokenmanager"
	"github.com/netbill/pgdbx"
)

func (a *App) Run(ctx context.Context) error {
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
		return fmt.Errorf("connect to database: %w", err)
	}
	defer pool.Close()

	a.log.Info("starting application")

	db := pgdbx.NewDB(pool)

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
		return fmt.Errorf("load aws config: %w", err)
	}

	uploader := media.NewUploader(awsx.New(a.config.S3.Aws.BucketName, cfg), media.Config{
		LinkTTL:   a.config.S3.Media.Link.TTL,
		OrgIcon:   a.config.S3.Media.Resources.Organization.Icon,
		OrgBanner: a.config.S3.Media.Resources.Organization.Banner,
	})

	outbox := eventpg.NewOutbox(db)
	inbox := eventpg.NewInbox(db)

	producer, err := messenger.NewProducer(a.log, messenger.ProducerConfig{
		Producer: a.config.Kafka.Identity,
		Brokers:  a.config.Kafka.Brokers,
		OrganizationV1: messenger.ProduceKafkaConfig{
			RequiredAcks: a.config.Kafka.Produce.Topics.OrganizationV1.RequiredAcks,
			Compression:  a.config.Kafka.Produce.Topics.OrganizationV1.Compression,
			Balancer:     a.config.Kafka.Produce.Topics.OrganizationV1.Balancer,
			BatchSize:    a.config.Kafka.Produce.Topics.OrganizationV1.BatchSize,
			BatchTimeout: a.config.Kafka.Produce.Topics.OrganizationV1.BatchTimeout,
		},
		OrgMembersV1: messenger.ProduceKafkaConfig{
			RequiredAcks: a.config.Kafka.Produce.Topics.OrgMemberV1.RequiredAcks,
			Compression:  a.config.Kafka.Produce.Topics.OrgMemberV1.Compression,
			Balancer:     a.config.Kafka.Produce.Topics.OrgMemberV1.Balancer,
			BatchSize:    a.config.Kafka.Produce.Topics.OrgMemberV1.BatchSize,
			BatchTimeout: a.config.Kafka.Produce.Topics.OrgMemberV1.BatchTimeout,
		},
	})
	if err != nil {
		return fmt.Errorf("create kafka producer: %w", err)
	}
	defer producer.Close()

	outbound := publisher.New(a.config.Kafka.Identity, outbox)

	profileRepo := repository.NewProfileRepo(pg.NewProfilesQ(db))
	organizationRepo := repository.NewOrganizationRepo(pg.NewOrganizationsQ(db))
	orgMemberRepo := repository.NewMemberRepo(pg.NewOrgMembersQ(db))
	orgInviteRepo := repository.NewInviteRepo(pg.NewOrgInvitesQ(db))
	placeRepo := repository.NewPlaceRepo(pg.NewPlacesQ(db))
	tombstoneRepo := pg.NewTombstonesQ(db)

	profileCore := core.NewProfileModule(core.ProfileDeps{
		Repo:      profileRepo,
		Tombstone: tombstoneRepo,
		Tx:        db,
	})

	placeCore := core.NewPlaceModule(core.PlaceDeps{
		Repo: placeRepo,
		Tx:   db,
	})

	orgCore := core.NewOrganizationModule(core.OrganizationDeps{
		Repo:      organizationRepo,
		Member:    orgMemberRepo,
		Place:     placeRepo,
		Tombstone: tombstoneRepo,
		Tx:        db,
		Messenger: outbound,
		Media:     uploader,
	})

	memberCore := core.NewMemberModule(core.MemberDeps{
		Auth:      orgCore,
		Repo:      orgMemberRepo,
		Tombstone: tombstoneRepo,
		Tx:        db,
		Messenger: outbound,
	})

	inviteCore := core.NewInviteModule(core.InviteDeps{
		Auth:      orgCore,
		Repo:      orgInviteRepo,
		Profile:   profileRepo,
		Member:    orgMemberRepo,
		Tx:        db,
		Messenger: outbound,
	})

	tokenManager := tokenmanager.New(tokenmanager.Config{
		Issuer:   a.config.Auth.Tokens.Issuer,
		AccessSK: a.config.Auth.Tokens.AccountAccess.SecretKey,
	})

	orgController := controller.NewOrganizationController(controller.OrganizationControllerDeps{
		Organizations: orgCore,
		Members:       memberCore,
		Profiles:      profileCore,
	})

	memberController := controller.NewMemberController(controller.MemberControllerDeps{
		Members:       memberCore,
		Profiles:      profileCore,
		Organizations: orgCore,
	})

	inviteController := controller.NewInviteController(controller.InviteControllerDeps{
		Invites:       inviteCore,
		Organizations: orgCore,
		Profiles:      profileCore,
	})

	router := rest.NewServer(rest.ServerDeps{
		Middlewares: middlewares.New(tokenManager),
		Org:         orgController,
		Member:      memberController,
		Invite:      inviteController,

		Log:           a.log,
		MediaResolver: media.NewResolver(a.config.S3.Aws.BaseURL),
		Config: rest.Config{
			Port:              a.config.Rest.Port,
			ReadTimeout:       a.config.Rest.Timeouts.Read,
			ReadHeaderTimeout: a.config.Rest.Timeouts.ReadHeader,
			WriteTimeout:      a.config.Rest.Timeouts.Write,
			IdleTimeout:       a.config.Rest.Timeouts.Idle,
		},
	})

	run(func() { router.Run(ctx) })

	outboxWorker := messenger.NewOutboxWorker(a.log, outbox, producer, eventbox.OutboxWorkerConfig{
		Routines:       a.config.Kafka.Outbox.Routines,
		Slots:          a.config.Kafka.Outbox.Slots,
		BatchSize:      a.config.Kafka.Outbox.BatchSize,
		Sleep:          a.config.Kafka.Outbox.Sleep,
		MinNextAttempt: a.config.Kafka.Outbox.MinNextAttempt,
		MaxNextAttempt: a.config.Kafka.Outbox.MaxNextAttempt,
		MaxAttempts:    a.config.Kafka.Outbox.MaxAttempts,
	})
	defer outboxWorker.Clean()

	run(func() { outboxWorker.Run(ctx) })

	evProfileController := evcontroller.NewProfileController(a.log, profileCore)
	evPlaceController := evcontroller.NewPlaceController(a.log, placeCore)

	inboxWorker := messenger.NewInboxWorker(messenger.InboxWorkerDeps{
		Logger:            a.log,
		Inbox:             inbox,
		ProfileController: evProfileController,
		PlaceController:   evPlaceController,
		Config: eventbox.InboxWorkerConfig{
			Routines:       a.config.Kafka.Inbox.Routines,
			Slots:          a.config.Kafka.Inbox.Slots,
			BatchSize:      a.config.Kafka.Inbox.BatchSize,
			Sleep:          a.config.Kafka.Inbox.Sleep,
			MinNextAttempt: a.config.Kafka.Inbox.MinNextAttempt,
			MaxNextAttempt: a.config.Kafka.Inbox.MaxNextAttempt,
			MaxAttempts:    a.config.Kafka.Inbox.MaxAttempts,
		},
	})
	defer inboxWorker.Clean()

	run(func() { inboxWorker.Run(ctx) })

	consumer := messenger.NewConsumer(a.log, inbox, messenger.ConsumerConfig{
		GroupID:    a.config.Kafka.Identity,
		Brokers:    a.config.Kafka.Brokers,
		MinBackoff: a.config.Kafka.Consume.Backoff.Min,
		MaxBackoff: a.config.Kafka.Consume.Backoff.Max,
		ProfilesV1: messenger.ConsumeKafkaConfig{
			Instances:     a.config.Kafka.Consume.Topics.ProfilesV1.Instances,
			MinBytes:      a.config.Kafka.Consume.Topics.ProfilesV1.MinBytes,
			MaxBytes:      a.config.Kafka.Consume.Topics.ProfilesV1.MaxBytes,
			MaxWait:       a.config.Kafka.Consume.Topics.ProfilesV1.MaxWait,
			QueueCapacity: a.config.Kafka.Consume.Topics.ProfilesV1.QueueCapacity,
		},
		PlacesV1: messenger.ConsumeKafkaConfig{
			Instances:     a.config.Kafka.Consume.Topics.PlacesV1.Instances,
			MinBytes:      a.config.Kafka.Consume.Topics.PlacesV1.MinBytes,
			MaxBytes:      a.config.Kafka.Consume.Topics.PlacesV1.MaxBytes,
			MaxWait:       a.config.Kafka.Consume.Topics.PlacesV1.MaxWait,
			QueueCapacity: a.config.Kafka.Consume.Topics.PlacesV1.QueueCapacity,
		},
	})
	defer consumer.Close()

	run(func() { consumer.Run(ctx) })

	wg.Wait()
	return nil
}
