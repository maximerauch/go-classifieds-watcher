package composite

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/maximerauch/go-classifieds-watcher/internal/core"
)

type mockNotifier struct {
	shouldFail bool
	wasCalled  bool
}

func (m *mockNotifier) Send(ctx context.Context, item core.Item) error {
	m.wasCalled = true
	if m.shouldFail {
		return errors.New("mock failure")
	}
	return nil
}

func TestCompositeNotifier_Send(t *testing.T) {
	item := core.Item{ID: "test-item"}
	ctx := context.Background()

	t.Run("Success: Notify all without errors", func(t *testing.T) {
		// Arrange
		n1 := &mockNotifier{shouldFail: false}
		n2 := &mockNotifier{shouldFail: false}
		composite := NewCompositeNotifier(n1, n2)

		// Act
		err := composite.Send(ctx, item)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if !n1.wasCalled || !n2.wasCalled {
			t.Error("Expected all notifiers to be called")
		}
	})

	t.Run("Resilience: Continue notifying even if one fails", func(t *testing.T) {
		// The first one fails, the second one should still run.
		n1 := &mockNotifier{shouldFail: true}
		n2 := &mockNotifier{shouldFail: false}
		composite := NewCompositeNotifier(n1, n2)

		// Act
		err := composite.Send(ctx, item)

		// Assert
		if err == nil {
			t.Error("Expected an error because one notifier failed")
		}

		// Critical check: Ensure the loop didn't break after the first error
		if !n1.wasCalled {
			t.Error("Notifier 1 should have been called")
		}
		if !n2.wasCalled {
			t.Error("Notifier 2 should have been called despite Notifier 1 failure")
		}

		// Check error aggregation message
		if !strings.Contains(err.Error(), "notification errors") {
			t.Errorf("Error message format invalid: %v", err)
		}
	})

	t.Run("Aggregation: Return multiple errors", func(t *testing.T) {
		// Arrange
		n1 := &mockNotifier{shouldFail: true}
		n2 := &mockNotifier{shouldFail: true}
		composite := NewCompositeNotifier(n1, n2)

		// Act
		err := composite.Send(ctx, item)

		// Assert
		if err == nil {
			t.Fatal("Expected error")
		}

		// Count occurrences of "mock failure" in the error string
		// Since both failed, we expect the error message to contain the error twice.
		count := strings.Count(err.Error(), "mock failure")
		if count != 2 {
			t.Errorf("Expected 2 errors in message, found %d. Msg: %v", count, err)
		}
	})

	t.Run("Edge Case: No notifiers registered", func(t *testing.T) {
		// Arrange
		composite := NewCompositeNotifier() // Empty

		// Act
		err := composite.Send(ctx, item)

		// Assert
		if err != nil {
			t.Errorf("Expected success on empty list, got %v", err)
		}
	})
}
