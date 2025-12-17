package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/maximerauch/go-classifieds-watcher/internal/adapters/composite"
	"github.com/maximerauch/go-classifieds-watcher/internal/adapters/email"
	"github.com/maximerauch/go-classifieds-watcher/internal/adapters/postgres"
	"github.com/maximerauch/go-classifieds-watcher/internal/adapters/rememberme"
	"github.com/maximerauch/go-classifieds-watcher/internal/adapters/std"
	"github.com/maximerauch/go-classifieds-watcher/internal/config"
	"github.com/maximerauch/go-classifieds-watcher/internal/core"
)

func main() {
	// Initialize JSON logger for structured logging (Cloud-Native standard)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Create a context that cancels on SIGINT (Ctrl+C) or SIGTERM (Docker/PaaS stop signal)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := run(ctx, logger); err != nil {
		logger.Error("fatal application error", "error", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func run(ctx context.Context, logger *slog.Logger) error {
	// 1. Load Configuration
	cfg := config.Load()

	// Switched mode to "daemon-worker" to reflect the long-running nature
	logger.Info("starting go-classifieds-watcher", "mode", "daemon-worker")

	// 2. Wiring & Dependencies Initialization

	// Initialize PostgreSQL repository
	repo, err := postgres.NewRepository(cfg.Database.DSN)
	if err != nil {
		return err
	}

	// Setup Notifiers: We combine Logger (stdout) and Email
	notifier := composite.NewCompositeNotifier(
		std.NewLoggerNotifier(logger),
		email.NewEmailNotifier(cfg.Email),
	)

	// Setup Provider
	provider := rememberme.NewProvider(cfg.RememberMe.SearchURL)

	// Initialize the Domain Service
	svc := core.NewWatcherService(provider, repo, notifier, logger)

	// 3. Execution Loop (Daemon Mode)

	// Step A: Immediate execution on startup (Fail-safe check)
	logger.Info("executing initial scan")
	if err := svc.Run(ctx); err != nil {
		// In daemon mode, we log errors but do not crash the app unless it's critical.
		logger.Error("initial scan failed", "error", err)
	}

	// Step B: Scheduled execution (every 15 minutes)
	ticker := time.NewTicker(15 * time.Minute)
	defer ticker.Stop()

	logger.Info("worker started, waiting for next cycle", "interval", "15m")

	for {
		select {
		case <-ctx.Done():
			// Handle graceful shutdown signal from Docker
			logger.Info("shutdown signal received, stopping worker")
			return nil

		case <-ticker.C:
			// The timer triggered
			logger.Info("starting scheduled scan")

			// Create a timeout context for this specific job execution (2 minutes max)
			// This prevents a stuck network call from hanging the worker forever.
			jobCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)

			if err := svc.Run(jobCtx); err != nil {
				logger.Error("scheduled scan failed", "error", err)
			}

			cancel() // Always release context resources
		}
	}
}
