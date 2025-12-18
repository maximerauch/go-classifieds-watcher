package fs

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/maximerauch/go-classifieds-watcher/internal/core"
)

func TestJSONRepository(t *testing.T) {
	// Create a temporary directory that will be automatically cleaned up
	// when the test finishes. This ensures complete isolation.
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test-db.json")
	ctx := context.Background()

	item := core.Item{
		ID:    "xyz-123",
		Title: "Test Item",
		Price: 100,
	}

	// Scenario 1: Fresh start
	// Ensure a new repository works even if the file doesn't exist yet.
	t.Run("Initialize with missing file", func(t *testing.T) {
		repo := NewJSONRepository(dbPath)
		exists, err := repo.Exists(ctx, "any-id")
		if err != nil {
			t.Fatalf("Unexpected error checking existence: %v", err)
		}
		if exists {
			t.Error("Repository should be empty initially")
		}
	})

	// Scenario 2: Persistence (Write)
	// Save an item and verify it exists in memory.
	t.Run("Save item", func(t *testing.T) {
		repo := NewJSONRepository(dbPath)
		if err := repo.Save(ctx, item); err != nil {
			t.Fatalf("Failed to save item: %v", err)
		}

		exists, _ := repo.Exists(ctx, item.ID)
		if !exists {
			t.Error("Item should exist in memory after save")
		}
	})

	// Scenario 3: Persistence (Reload from Disk)
	// CRITICAL: We create a NEW instance pointing to the SAME file.
	// This proves that data was actually written to the disk, not just kept in RAM.
	t.Run("Reload data from disk", func(t *testing.T) {
		// Verify file physically exists
		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			t.Fatal("Database file was not created on disk")
		}

		// Create a fresh instance to simulate a restart
		repo := NewJSONRepository(dbPath)

		exists, err := repo.Exists(ctx, item.ID)
		if err != nil {
			t.Fatalf("Error checking existence after reload: %v", err)
		}
		if !exists {
			t.Error("Item was lost after reloading repository (persistence failed)")
		}
	})

	// Scenario 4: Corruption resilience
	// If the file contains garbage, the repo should handle it gracefully (start empty)
	// instead of crashing.
	t.Run("Handle corrupt file gracefully", func(t *testing.T) {
		// Overwrite the DB file with invalid JSON
		err := os.WriteFile(dbPath, []byte("{ this is not json }"), 0644)
		if err != nil {
			t.Fatal(err)
		}

		// Should not panic
		repo := NewJSONRepository(dbPath)

		// Should effectively be empty
		exists, _ := repo.Exists(ctx, item.ID)
		if exists {
			t.Error("Repository should ignore corrupt data and start empty")
		}
	})
}
