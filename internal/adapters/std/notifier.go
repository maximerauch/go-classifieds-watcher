package std

import (
	"context"
	"log/slog"

	"github.com/maximerauch/go-classifieds-watcher/internal/core"
)

// LoggerNotifier implements core.Notifier by logging details to stdout.
type LoggerNotifier struct {
	logger *slog.Logger
}

func NewLoggerNotifier(l *slog.Logger) *LoggerNotifier {
	return &LoggerNotifier{logger: l}
}

func (n *LoggerNotifier) Send(ctx context.Context, listing core.Listing) error {
	n.logger.Info("NOTIFICATION SENT",
		"title", listing.Title,
		"price", listing.Price,
		"url", listing.Url,
	)
	return nil
}
