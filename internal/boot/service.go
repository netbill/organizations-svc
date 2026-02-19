package boot

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/netbill/logium"
	"github.com/netbill/organizations-svc/internal/bucket"
	"github.com/netbill/organizations-svc/internal/core/modules/invite"
	"github.com/netbill/organizations-svc/internal/core/modules/member"
	"github.com/netbill/organizations-svc/internal/core/modules/organization"
	"github.com/netbill/organizations-svc/internal/core/modules/profile"
	"github.com/netbill/organizations-svc/internal/core/modules/role"
	"github.com/netbill/organizations-svc/internal/messenger/inbound"
	"github.com/netbill/organizations-svc/internal/messenger/inbound/handlers"
	"github.com/netbill/organizations-svc/internal/repository"
	"github.com/netbill/organizations-svc/internal/repository/pg"
	"github.com/netbill/organizations-svc/internal/rest"
	"github.com/netbill/organizations-svc/internal/rest/controller"
	"github.com/netbill/organizations-svc/internal/rest/middlewares"
	"github.com/netbill/organizations-svc/internal/tokenmanager"
	"github.com/netbill/pgdbx"
	"github.com/netbill/restkit"
)

func StartServices(ctx context.Context, log *logium.Entry, wg *sync.WaitGroup, cfg Config) {
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

	s3Bucket, err := bucket.New(cfg.S3)
	if err != nil {
		log.Fatal("failed to create s3 bucket", "error", err)
	}

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

	msg := messenger.NewManager(log, db, cfg.Kafka)

	kafkaProducer := messenger.NewWriter()

	tokenManager := tokenmanager.New(cfg.Auth.Tokens)

	orgSvc := organization.New(repo, kafkaProducer, s3Bucket)
	memberSvc := member.New(repo, kafkaProducer)
	roleSvc := role.New(repo, kafkaProducer)
	inviteSvc := invite.New(repo, kafkaProducer)
	profileSvc := profile.New(repo)

	responser := restkit.NewResponser()
	ctrl := controller.New(&controller.Modules{
		Organization: orgSvc,
		Role:         roleSvc,
		Invite:       inviteSvc,
		Member:       memberSvc,
	}, responser)
	mdll := middlewares.New(responser, tokenManager)
	router := rest.New(mdll, ctrl)

	run(func() {
		router.Run(ctx, log, cfg.Rest)
	})

	kafkaConsumer := messenger.NewConsumer()

	kafkaInbox := inbound.NewInbox()

	kafkaOutbox := messenger.NewOutbox()

	run(func() {
		msg.RunInbox(ctx, handlers.New(handlers.Modules{Profile: profileSvc}))
	})

	run(func() { msg.RunConsumer(ctx) })

	run(func() { msg.RunOutbox(ctx) })
}
