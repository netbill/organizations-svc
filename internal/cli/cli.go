package cli

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/alecthomas/kingpin"
	evcli "github.com/netbill/eventbox/pg/cli"
	"github.com/netbill/organizations-svc/internal/app"
	"github.com/netbill/organizations-svc/migrations"
)

func Run(args []string) bool {
	cfg := app.LoadConfig()
	log := app.NewLogger(cfg.Log)

	var (
		service = kingpin.New("chains-auth", "")
		runCmd  = service.Command("run", "run command")

		serviceCmd     = runCmd.Command("service", "run service")
		migrateCmd     = service.Command("migrate", "migrate command")
		migrateUpCmd   = migrateCmd.Command("up", "migrate db up")
		migrateDownCmd = migrateCmd.Command("down", "migrate db down")

		eventsCmd = service.Command("events", "events commands")

		eventsInbox        = eventsCmd.Command("inbox", "inbox events")
		eventsInboxCleanup = eventsInbox.Command("cleanup", "cleanup inbox events")
		eventsInboxFailed  = eventsInboxCleanup.Command(
			"failed", "cleanup inbox events in failed state",
		)
		eventsInboxProcessing = eventsInboxCleanup.Command(
			"processing", "cleanup inbox events stuck in processing state",
		)
		eventsInboxProcessingProcessIDs = eventsInboxProcessing.Flag(
			"process-id", "cleanup only events reserved by this process id (repeatable)",
		).Strings()

		eventsOutbox        = eventsCmd.Command("outbox", "outbox events")
		eventsOutboxCleanup = eventsOutbox.Command("cleanup", "cleanup outbox events")
		eventsOutboxFailed  = eventsOutboxCleanup.Command(
			"failed", "cleanup outbox events in failed state",
		)
		eventsOutboxProcessing = eventsOutboxCleanup.Command(
			"processing", "cleanup outbox events stuck in processing state",
		)
		eventsOutboxProcessingProcessIDs = eventsOutboxProcessing.Flag(
			"process-id", "cleanup only events reserved by this process id (repeatable)",
		).Strings()
	)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup

	command, err := service.Parse(args[1:])
	if err != nil {
		log.WithError(err).Error("failed to parse arguments")
		return false
	}

	switch command {
	case serviceCmd.FullCommand():
		app.StartServices(ctx, log, &wg, cfg)
	case migrateUpCmd.FullCommand():
		err = migrations.MigrateUp(ctx, log, cfg.Database.SQL.URL)
	case migrateDownCmd.FullCommand():
		err = migrations.MigrateDown(ctx, log, cfg.Database.SQL.URL)
	case eventsOutboxFailed.FullCommand():
		err = evcli.CleanupOutboxFailed(ctx, log, cfg.Database.SQL.URL)
	case eventsOutboxProcessing.FullCommand():
		err = evcli.CleanupOutboxProcessing(ctx, log, cfg.Database.SQL.URL, *eventsOutboxProcessingProcessIDs...)
	case eventsInboxFailed.FullCommand():
		err = evcli.CleanupInboxFailed(ctx, log, cfg.Database.SQL.URL)
	case eventsInboxProcessing.FullCommand():
		err = evcli.CleanupInboxProcessing(ctx, log, cfg.Database.SQL.URL, *eventsInboxProcessingProcessIDs...)
	default:
		log.Errorf("unknown command %s", command)
		return false
	}
	if err != nil {
		log.WithError(err).Error("failed to exec cmd")
		return false
	}

	wgch := make(chan struct{})
	go func() {
		wg.Wait()
		close(wgch)
	}()

	select {
	case <-ctx.Done():
		log.Warnf("Interrupt signal received: %v", ctx.Err())
		<-wgch
	case <-wgch:
		log.Warnf("All api stopped")
	}

	return true
}
