package cmd

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/netbill/awsx"
	"github.com/netbill/logium"
	"github.com/netbill/organizations-svc/internal/bucket"
	"github.com/netbill/organizations-svc/internal/core/modules/invite"
	"github.com/netbill/organizations-svc/internal/core/modules/member"
	"github.com/netbill/organizations-svc/internal/core/modules/organization"
	"github.com/netbill/organizations-svc/internal/core/modules/profile"
	"github.com/netbill/organizations-svc/internal/core/modules/role"
	"github.com/netbill/organizations-svc/internal/core/modules/rperm"
	"github.com/netbill/organizations-svc/internal/messenger"
	"github.com/netbill/organizations-svc/internal/messenger/inbound"
	"github.com/netbill/organizations-svc/internal/messenger/outbound"
	"github.com/netbill/organizations-svc/internal/repository"
	"github.com/netbill/organizations-svc/internal/repository/pg"
	"github.com/netbill/organizations-svc/internal/rest"
	"github.com/netbill/organizations-svc/internal/rest/controller"
	"github.com/netbill/organizations-svc/internal/rest/middlewares"
	"github.com/netbill/organizations-svc/internal/tokenmanager"
	"github.com/netbill/pgdbx"
	"github.com/netbill/restkit"
)

func StartServices(ctx context.Context, cfg Config, log *logium.Logger, wg *sync.WaitGroup) {
	run := func(f func()) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			f()
		}()
	}

	pool, err := pgxpool.New(ctx, cfg.Database.SQL.URL)
	if err != nil {
		log.Fatal("failed to connect to database", "error", err)
	}
	db := pgdbx.NewDB(pool)

	awsCfg := aws.Config{
		Region: cfg.S3.AWS.Region,
		Credentials: credentials.NewStaticCredentialsProvider(
			cfg.S3.AWS.AccessKeyID,
			cfg.S3.AWS.SecretAccessKey,
			"",
		),
	}

	s3Client := s3.NewFromConfig(awsCfg)
	presignClient := s3.NewPresignClient(s3Client)

	awsS3 := awsx.New(
		cfg.S3.AWS.BucketName,
		s3Client,
		presignClient,
	)

	orgIconValidator := &awsx.ImgObjectValidator{
		AllowedContentTypes: cfg.S3.Upload.Organization.Icon.AllowedContentTypes,
		AllowedFormats:      cfg.S3.Upload.Organization.Icon.AllowedFormats,
		MaxWidth:            cfg.S3.Upload.Organization.Icon.MaxWidth,
		MaxHeight:           cfg.S3.Upload.Organization.Icon.MaxHeight,
		ContentLengthMax:    cfg.S3.Upload.Organization.Icon.ContentLengthMax,
	}

	orgBannerValidator := &awsx.ImgObjectValidator{
		AllowedContentTypes: cfg.S3.Upload.Organization.Banner.AllowedContentTypes,
		AllowedFormats:      cfg.S3.Upload.Organization.Banner.AllowedFormats,
		MaxWidth:            cfg.S3.Upload.Organization.Banner.MaxWidth,
		MaxHeight:           cfg.S3.Upload.Organization.Banner.MaxHeight,
		ContentLengthMax:    cfg.S3.Upload.Organization.Banner.ContentLengthMax,
	}

	s3Bucket := bucket.New(bucket.Config{
		S3:                 awsS3,
		OrgIconValidator:   orgIconValidator,
		OrgBannerValidator: orgBannerValidator,
		UploadTokensTTL: bucket.UploadTokensTTL{
			Org: cfg.S3.Upload.Token.TTL.Organization,
		},
	})

	orgInvitesSql := pg.NewOrgInvitesQ(db)
	orgMemberRolesSql := pg.NewOrgMemberRolesQ(db)
	orgMembersSql := pg.NewOrgMembersQ(db)
	organizationsSql := pg.NewOrganizationsQ(db)
	orgRolePermSql := pg.NewOrgRolePermissionsQ(db)
	orgRolePermLinksSql := pg.NewOrgRolePermissionLinksQ(db)
	orgRolesSql := pg.NewOrgRolesQ(db)
	profilesSql := pg.NewProfilesQ(db)
	transactioner := pg.NewTransaction(db)

	repo := &repository.Repository{
		OrgInvitesSql:             orgInvitesSql,
		OrgMemberRolesSql:         orgMemberRolesSql,
		OrgMembersSql:             orgMembersSql,
		OrganizationsSql:          organizationsSql,
		OrgRolePermissionsSql:     orgRolePermSql,
		OrgRolePermissionLinksSql: orgRolePermLinksSql,
		OrgRolesSql:               orgRolesSql,
		ProfilesSql:               profilesSql,
		Transactioner:             transactioner,
	}

	kafkaOutbound := outbound.New(log, db)

	tokenManager := tokenmanager.New(cfg.Service.Name, cfg.S3.Upload.Token.TTL.Organization)

	orgSvc := organization.New(repo, kafkaOutbound, tokenManager, s3Bucket)
	memberSvc := member.New(repo, kafkaOutbound)
	roleSvc := role.New(repo, kafkaOutbound)
	permSvc := rperm.New(repo, kafkaOutbound)
	inviteSvc := invite.New(repo, kafkaOutbound)
	profileSvc := profile.New(repo)

	responser := restkit.NewResponser()
	ctrl := controller.New(log, responser, orgSvc, memberSvc, roleSvc, permSvc, inviteSvc)
	mdll := middlewares.New(log, middlewares.Config{
		AccountAccessSK: cfg.Auth.Account.Token.Access.SecretKey,
		UploadFilesSK:   cfg.S3.Upload.Token.SecretKey,
	}, responser)
	router := rest.New(log, mdll, ctrl)

	msgx := messenger.New(log, db, cfg.Kafka.Brokers...)

	log.Infof("starting kafka brokers %s", cfg.Kafka.Brokers)

	run(func() {
		router.Run(ctx, rest.Config{
			Port:              cfg.Rest.Port,
			TimeoutRead:       cfg.Rest.Timeouts.Read,
			TimeoutReadHeader: cfg.Rest.Timeouts.ReadHeader,
			TimeoutWrite:      cfg.Rest.Timeouts.Write,
			TimeoutIdle:       cfg.Rest.Timeouts.Idle,
		})
	})

	run(func() { msgx.RunConsumer(ctx, inbound.New(log, profileSvc)) })

	run(func() { msgx.RunProducer(ctx) })
}
