package cmd

import (
	"context"
	"database/sql"
	"sync"

	"github.com/netbill/logium"
	"github.com/netbill/organizations-svc/internal"
	"github.com/netbill/organizations-svc/internal/core/modules/invite"
	"github.com/netbill/organizations-svc/internal/core/modules/member"
	"github.com/netbill/organizations-svc/internal/core/modules/organization"
	"github.com/netbill/organizations-svc/internal/core/modules/profile"
	"github.com/netbill/organizations-svc/internal/core/modules/role"
	"github.com/netbill/organizations-svc/internal/messenger"
	"github.com/netbill/organizations-svc/internal/messenger/inbound"
	"github.com/netbill/organizations-svc/internal/messenger/outbound"
	"github.com/netbill/organizations-svc/internal/repository"
	"github.com/netbill/organizations-svc/internal/rest"
	"github.com/netbill/organizations-svc/internal/rest/controller"
	"github.com/netbill/restkit/mdlv"
)

func StartServices(ctx context.Context, cfg internal.Config, log logium.Logger, wg *sync.WaitGroup) {
	run := func(f func()) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			f()
		}()
	}

	pg, err := sql.Open("postgres", cfg.Database.SQL.URL)
	if err != nil {
		log.Fatal("failed to connect to database", "error", err)
	}

	database := repository.New(pg)

	kafkaOutbound := outbound.New(log, pg)

	orgSvc := organization.New(database, kafkaOutbound)
	memberSvc := member.New(database, kafkaOutbound)
	roleSvc := role.New(database, kafkaOutbound)
	inviteSvc := invite.New(database, kafkaOutbound)
	profileSvc := profile.New(database)

	ctrl := controller.New(orgSvc, memberSvc, roleSvc, inviteSvc, log)
	mdll := mdlv.New(cfg.JWT.User.AccessToken.SecretKey, rest.AccountDataCtxKey)
	router := rest.New(log, mdll, ctrl)

	msgx := messenger.New(log, pg, cfg.Kafka.Brokers...)

	log.Infof("starting kafka brokers %s", cfg.Kafka.Brokers)

	run(func() { router.Run(ctx, cfg) })

	run(func() { msgx.RunConsumer(ctx, inbound.New(log, profileSvc)) })

	run(func() { msgx.RunProducer(ctx) })
}
