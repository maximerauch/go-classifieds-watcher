//go:build integration

package postgres

import (
	"context"
	"os"
	"testing"

	"github.com/maximerauch/go-classifieds-watcher/internal/core"
)

// TestRepository_Integration validates the interaction with a real PostgreSQL instance.
// It requires the TEST_DATABASE_URL environment variable to be set.
func TestRepository_Integration(t *testing.T) {
	// SETUP & PRECONDITION CHECKS
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		t.Skip("Skipping integration test: TEST_DATABASE_URL not set")
	}

	// Initialize the repository using the real constructor to test connection logic.
	// NOTE: If this fails locally with "server does not support SSL",
	// ensure your NewRepository logic allows "sslmode=disable" (see fix below).
	repo, err := NewRepository(dbURL)
	if err != nil {
		t.Fatalf("Failed to initialize repository: %v", err)
	}
	// Ensure the connection is closed after tests complete.
	defer func() {
		_ = repo.db.Close()
	}()

	// Clean up the table to ensure a pristine state before running assertions.
	// Since NewRepository runs migrations, the table is guaranteed to exist.
	if _, err := repo.db.Exec("TRUNCATE TABLE seen_items"); err != nil {
		t.Fatalf("Failed to truncate table: %v", err)
	}

	ctx := context.Background()
	itemID := "test-item-123"
	item := core.Item{
		ID:    itemID,
		Title: "Integration Test Item",
	}

	// TEST SCENARIOS

	// Scenario A: Check for an item that hasn't been saved yet.
	t.Run("Returns false when item does not exist", func(t *testing.T) {
		exists, err := repo.Exists(ctx, itemID)
		if err != nil {
			t.Fatalf("Exists() failed unexpectedly: %v", err)
		}
		if exists {
			t.Error("Expected Exists() to return false, got true")
		}
	})

	// Scenario B: Persist a new item.
	t.Run("Successfully saves a new item", func(t *testing.T) {
		err := repo.Save(ctx, item)
		if err != nil {
			t.Fatalf("Save() failed: %v", err)
		}
	})

	// Scenario C: Verify the item now exists.
	t.Run("Returns true when item exists", func(t *testing.T) {
		exists, err := repo.Exists(ctx, itemID)
		if err != nil {
			t.Fatalf("Exists() failed unexpectedly: %v", err)
		}
		if !exists {
			t.Error("Expected Exists() to return true, got false")
		}
	})

	// Scenario D: Idempotency Check.
	// Saving the same ID again should not return an error (ON CONFLICT DO NOTHING).
	t.Run("Handles duplicate inserts gracefully (Idempotency)", func(t *testing.T) {
		err := repo.Save(ctx, item)
		if err != nil {
			t.Errorf("Save() on duplicate returned error: %v", err)
		}

		// Verify that we still have only 1 row (no duplicates)
		var count int
		err = repo.db.QueryRow("SELECT COUNT(*) FROM seen_items WHERE id = $1", itemID).Scan(&count)
		if err != nil {
			t.Fatalf("Count query failed: %v", err)
		}
		if count != 1 {
			t.Errorf("Expected 1 record, found %d", count)
		}
	})
}
