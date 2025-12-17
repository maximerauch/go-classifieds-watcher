package core

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
)

type mockProvider struct {
	items []Item
	err   error
}

func (m *mockProvider) FetchItems(ctx context.Context) ([]Item, error) {
	return m.items, m.err
}
func (m *mockProvider) Name() string { return "MockProvider" }

type mockRepository struct {
	exists map[string]bool // Simulate database state
	saved  []Item          // Store saved items to verify assertions
	err    error           // Simulate DB error
}

func (m *mockRepository) Exists(ctx context.Context, id string) (bool, error) {
	if m.err != nil {
		return false, m.err
	}
	return m.exists[id], nil
}

func (m *mockRepository) Save(ctx context.Context, item Item) error {
	if m.err != nil {
		return m.err
	}
	m.saved = append(m.saved, item)
	return nil
}

type mockNotifier struct {
	sent []Item // Store sent notifications
	err  error  // Simulate SMTP error
}

func (m *mockNotifier) Send(ctx context.Context, item Item) error {
	if m.err != nil {
		return m.err
	}
	m.sent = append(m.sent, item)
	return nil
}

func TestWatcherService_Run(t *testing.T) {
	// We use a discarded logger to avoid polluting test output
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	validItem := Item{ID: "1", Title: "Valid Item", Url: "https://test.com"}
	invalidItem := Item{ID: "", Title: "Invalid", Url: ""} // ID missing

	tests := []struct {
		name          string
		providerItems []Item
		providerErr   error
		repoExisting  map[string]bool
		repoErr       error
		notifierErr   error
		expectError   bool // Do we expect Run() to return an error?
		expectedSaved int  // How many items should be persisted?
		expectedSent  int  // How many emails should be sent?
	}{
		{
			name:          "Nominal Case: New item found",
			providerItems: []Item{validItem},
			repoExisting:  map[string]bool{}, // Empty DB
			expectError:   false,
			expectedSaved: 1,
			expectedSent:  1,
		},
		{
			name:          "Idempotency: Item already exists",
			providerItems: []Item{validItem},
			repoExisting:  map[string]bool{"1": true}, // Item "1" already seen
			expectError:   false,
			expectedSaved: 0, // Should NOT save again
			expectedSent:  0, // Should NOT notify again
		},
		{
			name:          "Resilience: Provider failure",
			providerItems: nil,
			providerErr:   errors.New("network timeout"),
			expectError:   true,
			expectedSaved: 0,
			expectedSent:  0,
		},
		{
			name:          "Logic: Invalid items are skipped",
			providerItems: []Item{invalidItem},
			expectError:   false,
			expectedSaved: 0,
			expectedSent:  0,
		},
		{
			name:          "At-Least-Once Delivery: If Notifier fails, do NOT save",
			providerItems: []Item{validItem},
			notifierErr:   errors.New("smtp down"),
			expectError:   false,
			// CRITICAL: We expect 0 saved.
			// Because if we save it, we will never retry notifying it.
			expectedSaved: 0,
			expectedSent:  0,
		},
		{
			name:          "Resilience: Repository failure on check",
			providerItems: []Item{validItem},
			repoErr:       errors.New("db connection lost"),
			expectError:   false,
			expectedSaved: 0, // Should skip this item safely
			expectedSent:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup Mocks
			mockProv := &mockProvider{items: tt.providerItems, err: tt.providerErr}
			mockRepo := &mockRepository{exists: tt.repoExisting, err: tt.repoErr}
			mockNotif := &mockNotifier{err: tt.notifierErr}

			// Instantiate Service
			svc := NewWatcherService(mockProv, mockRepo, mockNotif, logger)

			// Execute
			err := svc.Run(context.Background())

			// Assertions
			if (err != nil) != tt.expectError {
				t.Errorf("Run() error = %v, expectError %v", err, tt.expectError)
			}

			if len(mockRepo.saved) != tt.expectedSaved {
				t.Errorf("Repo.Save() called %d times, want %d", len(mockRepo.saved), tt.expectedSaved)
			}

			if len(mockNotif.sent) != tt.expectedSent {
				t.Errorf("Notifier.Send() called %d times, want %d", len(mockNotif.sent), tt.expectedSent)
			}
		})
	}
}
