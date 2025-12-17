package std

import (
	"context"
	"log/slog"

	"github.com/maximerauch/go-classifieds-watcher/internal/core"
)

type LoggerNotifier struct {
	logger *slog.Logger
}

func NewLoggerNotifier(l *slog.Logger) *LoggerNotifier {
	return &LoggerNotifier{logger: l}
}

func (n *LoggerNotifier) Send(ctx context.Context, item core.Item) error {
	n.logger.Info("NOTIFICATION SENT",
		"title", item.Title,
		"price", item.Price,
		"url", item.Url,
	)
	return nil
}
