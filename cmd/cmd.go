package cmd

import (
	"context"
	"database/sql"
	"sync"

	"github.com/umisto/cities-svc/internal"
	"github.com/umisto/cities-svc/internal/domain/modules/agglomeration"
	"github.com/umisto/cities-svc/internal/domain/modules/city"
	"github.com/umisto/cities-svc/internal/domain/modules/invite"
	"github.com/umisto/cities-svc/internal/domain/modules/member"
	"github.com/umisto/cities-svc/internal/domain/modules/role"
	"github.com/umisto/cities-svc/internal/repository"
	"github.com/umisto/cities-svc/internal/rest/controller"
	"github.com/umisto/logium"
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

	aggloSvc := agglomeration.New(database, nil)
	citySvc := city.New(database, nil)
	memberSvc := member.New(database, nil)
	roleSvc := role.New(database, nil)
	inviteSvc := invite.New(database, nil)

	_ = controller.New(aggloSvc, citySvc, memberSvc, roleSvc, inviteSvc)
}
