package cli

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/alecthomas/kingpin"
	"github.com/netbill/organizations-svc/internal/app"
	"github.com/netbill/organizations-svc/internal/config"
)

func Run(args []string) bool {
	cfg := config.LoadConfig()
	log := cfg.Logger()

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

	application := app.New(log, cfg)
	switch command {
	case serviceCmd.FullCommand():
		application.Run(ctx)
	case migrateUpCmd.FullCommand():
		err = application.MigrateUp(ctx)
	case migrateDownCmd.FullCommand():
		err = application.MigrateDown(ctx)
	case eventsOutboxFailed.FullCommand():
		err = application.CleanupOutboxFailedEvents(ctx)
	case eventsOutboxProcessing.FullCommand():
		err = application.CleanupOutboxProcessingEvents(ctx, *eventsOutboxProcessingProcessIDs...)
	case eventsInboxFailed.FullCommand():
		err = application.CleanupInboxFailedEvents(ctx)
	case eventsInboxProcessing.FullCommand():
		err = application.CleanupInboxProcessingEvents(ctx, *eventsInboxProcessingProcessIDs...)
	default:
		log.Error("unknown command %s", command)
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
		log.Warn("Interrupt signal received: %v", ctx.Err())
		<-wgch
	case <-wgch:
		log.Warn("All api stopped")
	}

	return true
}
