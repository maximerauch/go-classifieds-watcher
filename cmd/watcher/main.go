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
	"github.com/maximerauch/go-classifieds-watcher/internal/config" // Import Config
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
	os.Exit(0)
}

func run(ctx context.Context, logger *slog.Logger) error {
	// 1. Load Configuration
	cfg := config.Load()

	logger.Info("starting go-classifieds-watcher",
		"mode", "one-shot-job",
		"items_per_page", cfg.ItemsPerPage,
	)

	// 2. Wiring
	repo := fs.NewJSONRepository(cfg.DataFilePath)
	notifier := std.NewLoggerNotifier(logger)
	provider := asi67.NewProvider(cfg.APIURL, cfg.ItemsPerPage)

	svc := core.NewWatcherService(provider, repo, notifier, logger)

	// 3. Execution
	start := time.Now()
	if err := svc.Run(ctx); err != nil {
		return err
	}

	logger.Info("run execution time", "duration", time.Since(start).String())
	return nil
}
