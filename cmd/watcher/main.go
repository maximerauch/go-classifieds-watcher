package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/maximerauch/go-classifieds-watcher/internal/adapters/asi67"
	"github.com/maximerauch/go-classifieds-watcher/internal/adapters/fs"
	"github.com/maximerauch/go-classifieds-watcher/internal/adapters/std"
	"github.com/maximerauch/go-classifieds-watcher/internal/core"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := run(ctx, logger); err != nil {
		logger.Error("job failed", "error", err)
		os.Exit(1)
	}

	logger.Info("job finished successfully")
	os.Exit(0)
}

func run(ctx context.Context, logger *slog.Logger) error {
	logger.Info("starting go-classifieds-watcher", "mode", "one-shot-job")

	// 1. Configuration
	targetAPI := "https://www.asi67.com/webapi/getJson/Templates/ProductsList"

	// 2. Wiring
	repo := fs.NewJSONRepository("data/seen.json")
	notifier := std.NewLoggerNotifier(logger)
	provider := asi67.NewProvider(targetAPI)

	svc := core.NewWatcherService(provider, repo, notifier, logger)

	// 3. Execution
	start := time.Now()
	if err := svc.Run(ctx); err != nil {
		return err
	}

	logger.Info("run execution time", "duration", time.Since(start).String())
	return nil
}
