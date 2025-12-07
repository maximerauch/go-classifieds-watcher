package fs

import (
	"context"
	"encoding/json"
	"os"
	"sync"

	"github.com/maximerauch/go-classifieds-watcher/internal/core"
)

type JSONRepository struct {
	filePath string
	mu       sync.Mutex
	storage  map[string]core.Listing
}

func NewJSONRepository(filePath string) *JSONRepository {
	repo := &JSONRepository{
		filePath: filePath,
		storage:  make(map[string]core.Listing),
	}
	repo.load()
	return repo
}

func (r *JSONRepository) Exists(ctx context.Context, id string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, exists := r.storage[id]
	return exists, nil
}

// Save stores the full listing in the map and persists to disk
func (r *JSONRepository) Save(ctx context.Context, listing core.Listing) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.storage[listing.ID] = listing
	return r.save()
}

func (r *JSONRepository) load() {
	file, err := os.ReadFile(r.filePath)
	if err != nil {
		return // Start with empty map if file doesn't exist
	}
	_ = json.Unmarshal(file, &r.storage)
}

func (r *JSONRepository) save() error {
	data, err := json.MarshalIndent(r.storage, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.filePath, data, 0644)
}
