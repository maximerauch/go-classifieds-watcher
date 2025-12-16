package core

import (
	"context"
	"fmt"
	"log/slog"
)

// WatcherService orchestrates the data flow between the Provider, Repository and Notifier.
type WatcherService struct {
	provider Provider
	repo     Repository
	notifier Notifier
	logger   *slog.Logger
}

// NewWatcherService creates a new service instance with injected dependencies.
func NewWatcherService(p Provider, r Repository, n Notifier, l *slog.Logger) *WatcherService {
	return &WatcherService{
		provider: p,
		repo:     r,
		notifier: n,
		logger:   l,
	}
}

// Run executes the main logic: fetch, filter, notify, and persist.
func (s *WatcherService) Run(ctx context.Context) error {
	s.logger.Info("starting watcher run", "provider", s.provider.Name())

	// 1. Fetch
	items, err := s.provider.FetchItems(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch items from %s: %w", s.provider.Name(), err)
	}

	s.logger.Info("items fetched", "count", len(items))

	newCount := 0
	for _, item := range items {
		// Defensive check
		if !item.IsValid() {
			s.logger.Warn("skipping invalid item", "item", item)
			continue
		}

		// 2. Check Dedup (Idempotency)
		exists, err := s.repo.Exists(ctx, item.ID)
		if err != nil {
			s.logger.Error("failed to check existence", "id", item.ID, "error", err)
			continue // Don't block the batch on single failure
		}

		if exists {
			continue
		}

		s.logger.Info("new item found", "id", item.ID, "title", item.Title)

		// 3. Notify
		if err := s.notifier.Send(ctx, item); err != nil {
			s.logger.Error("failed to notify", "id", item.ID, "error", err)
			// Strategy: If notification fails, do not save the ID.
			// We want to retry this item on the next run (At-Least-Once delivery).
			continue
		}

		// 4. Save
		if err := s.repo.Save(ctx, item); err != nil {
			s.logger.Error("failed to save id", "id", item.ID, "error", err)
		} else {
			newCount++
		}
	}

	s.logger.Info("watcher run finished", "new_items", newCount)
	return nil
}
