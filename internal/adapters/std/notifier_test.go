package std

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"

	"github.com/maximerauch/go-classifieds-watcher/internal/core"
)

func TestLoggerNotifier_Send(t *testing.T) {
	// 1. ARRANGE
	// On crée un buffer en mémoire pour capturer les logs au lieu de les afficher
	var buf bytes.Buffer

	// On configure un logger qui écrit dans ce buffer.
	// On retire l'heure (TimeKey: nil) pour éviter que le test échoue à chaque seconde différente.
	opts := &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{} // Supprime le timestamp du log
			}
			return a
		},
	}
	logger := slog.New(slog.NewTextHandler(&buf, opts))

	notifier := NewLoggerNotifier(logger)

	item := core.Item{
		ID:    "123",
		Title: "Gibson Les Paul",
		Price: 2500.50,
		Url:   "https://example.com/guitar",
	}

	// 2. ACT
	err := notifier.Send(context.Background(), item)

	// 3. ASSERT
	if err != nil {
		t.Errorf("Send() returned an error: %v", err)
	}

	// On récupère ce qui a été écrit dans le buffer
	logOutput := buf.String()

	// On vérifie que les informations clés sont présentes dans le log
	expectedSubstrings := []string{
		"level=INFO",
		"msg=\"NOTIFICATION SENT\"",
		"title=\"Gibson Les Paul\"",
		"price=2500.5",
		"url=https://example.com/guitar",
	}

	for _, s := range expectedSubstrings {
		if !strings.Contains(logOutput, s) {
			t.Errorf("Log output missing expected string: '%s'.\nFull Output: %s", s, logOutput)
		}
	}
}
