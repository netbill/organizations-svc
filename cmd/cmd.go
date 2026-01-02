package cmd

import (
	"context"
	"database/sql"
	"sync"

	"github.com/netbill/kafkakit/box"
	"github.com/netbill/logium"
	"github.com/netbill/organizations-svc/internal"
	"github.com/netbill/organizations-svc/internal/domain/modules/invite"
	"github.com/netbill/organizations-svc/internal/domain/modules/member"
	"github.com/netbill/organizations-svc/internal/domain/modules/organization"
	"github.com/netbill/organizations-svc/internal/domain/modules/profile"
	"github.com/netbill/organizations-svc/internal/domain/modules/role"
	"github.com/netbill/organizations-svc/internal/messenger/consumer"
	"github.com/netbill/organizations-svc/internal/messenger/consumer/callbacker"
	"github.com/netbill/organizations-svc/internal/messenger/producer"
	"github.com/netbill/organizations-svc/internal/repository"
	"github.com/netbill/organizations-svc/internal/rest"
	"github.com/netbill/organizations-svc/internal/rest/controller"
	"github.com/netbill/restkit/mdlv"
)

func StartServices(ctx context.Context, cfg internal.Config, log logium.Logger, wg *sync.WaitGroup) {
	run := func(f func()) {
		wg.Add(1)
		go func() {
			f()
			wg.Done()
		}()
	}

	pg, err := sql.Open("postgres", cfg.Database.SQL.URL)
	if err != nil {
		log.Fatal("failed to connect to database", "error", err)
	}

	database := repository.New(pg)
	kafkaBox := box.New(pg)

	kafkaProducer := producer.New(log, cfg.Kafka.Brokers, kafkaBox)

	aggloSvc := organization.New(database, kafkaProducer)
	memberSvc := member.New(database, kafkaProducer)
	roleSvc := role.New(database, kafkaProducer)
	inviteSvc := invite.New(database, kafkaProducer)
	profileSvc := profile.New(database)

	kafkaConsumer := consumer.New(log, cfg.Kafka.Brokers, box.New(pg), callbacker.New(log, profileSvc))

	ctrl := controller.New(aggloSvc, memberSvc, roleSvc, inviteSvc, log)
	mdll := mdlv.New(cfg.JWT.User.AccessToken.SecretKey, rest.AccountDataCtxKey)
	router := rest.New(log, mdll, ctrl)

	run(func() { router.Run(ctx, cfg) })

	run(func() { kafkaConsumer.Run(ctx) })

}
