package email

import (
	"context"
	"strings"
	"testing"

	"github.com/maximerauch/go-classifieds-watcher/internal/config"
	"github.com/maximerauch/go-classifieds-watcher/internal/core"
)

// TestEmailNotifier_buildBody tests the private method buildBody.
// By using 'package email' instead of 'package email_test', we gain access to private methods.
// This allows us to verify HTML generation logic without sending actual emails.
func TestEmailNotifier_buildBody(t *testing.T) {
	notifier := NewEmailNotifier(config.EmailConfig{})

	item := core.Item{
		ID:          "123",
		Title:       "Super Guitar",
		Description: "Mint condition",
		Price:       1500.00,
		Currency:    "EUR",
		Url:         "https://test.com/guitar",
	}

	body := notifier.buildBody(item)

	// We verify that critical information is present in the generated HTML.

	tests := []struct {
		name     string
		contains string
	}{
		{"Title presence", "Super Guitar"},
		{"Description presence", "Mint condition"},
		{"Price formatting", "1500.00 EUR"}, // Checks formatting %.2f
		{"Link presence", "https://test.com/guitar"},
		{"HTML Structure", "<h2>New Item Discovered!</h2>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(body, tt.contains) {
				t.Errorf("Email body missing expected content: '%s'", tt.contains)
			}
		})
	}
}

// TestEmailNotifier_Send_ContextCancelled ensures that the Send method
// respects the context cancellation and aborts immediately.
// This is a safe way to test part of the Send method without triggering a real network call.
func TestEmailNotifier_Send_ContextCancelled(t *testing.T) {
	// ARRANGE
	notifier := NewEmailNotifier(config.EmailConfig{
		SMTPHost: "smtp.example.com", // Dummy host, should not be reached
	})

	// Create a context that is ALREADY cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	item := core.Item{Title: "Test"}

	// ACT
	err := notifier.Send(ctx, item)

	// ASSERT
	if err == nil {
		t.Error("Expected error due to cancelled context, got nil")
	}

	if err != context.Canceled {
		t.Errorf("Expected context.Canceled error, got: %v", err)
	}
}
