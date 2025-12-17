package composite

import (
	"context"
	"errors"
	"strings"

	"github.com/maximerauch/go-classifieds-watcher/internal/core"
)

type CompositeNotifier struct {
	notifiers []core.Notifier
}

func NewCompositeNotifier(notifiers ...core.Notifier) *CompositeNotifier {
	return &CompositeNotifier{
		notifiers: notifiers,
	}
}

// Send dispatches the notification to all registered notifiers.
// It tries to notify everyone and collects errors if any occur.
func (m *CompositeNotifier) Send(ctx context.Context, item core.Item) error {
	var errs []string

	for _, n := range m.notifiers {
		if err := n.Send(ctx, item); err != nil {
			errs = append(errs, err.Error())
		}
	}

	if len(errs) > 0 {
		return errors.New("notification errors: " + strings.Join(errs, "; "))
	}
	return nil
}
