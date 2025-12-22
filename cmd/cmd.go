package cmd

import (
	"context"
	"sync"

	"github.com/umisto/cities-svc/internal"
	"github.com/umisto/logium"
)

func StartServices(ctx context.Context, cfg internal.Config, log logium.Logger, wg *sync.WaitGroup) {
	//run := func(f func()) {
	//	wg.Add(1)
	//	go func() {
	//		f()
	//		wg.Done()
	//	}()
	//}
	//
	//pg, err := sql.Open("postgres", cfg.Database.SQL.URL)
	//if err != nil {
	//	log.Fatal("failed to connect to database", "error", err)
	//}
	//
	//database := repo.NewDatabase(pg)
	//
	//eventPublish := publisher.New(cfg.Kafka.Broker)
	//
	//citySvc := city.NewService(database, eventPublish)
	//
	//ctrl := controller.New(log, citySvc)
	//mdlv := middlewares.New(log)
	//
	//run(func() { rest.Run(ctx, cfg, log, mdlv, ctrl) })
}
