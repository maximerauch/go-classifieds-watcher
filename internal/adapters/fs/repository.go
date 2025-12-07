package fs

import (
	"context"
	"encoding/json"
	"os"
	"sync"
)

type JSONRepository struct {
	filePath string
	mu       sync.Mutex
	seenIDs  map[string]bool
}

func NewJSONRepository(filePath string) *JSONRepository {
	repo := &JSONRepository{
		filePath: filePath,
		seenIDs:  make(map[string]bool),
	}
	repo.load()
	return repo
}

func (r *JSONRepository) Exists(ctx context.Context, id string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.seenIDs[id], nil
}

func (r *JSONRepository) SaveID(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.seenIDs[id] = true
	return r.save()
}

func (r *JSONRepository) load() {
	file, err := os.ReadFile(r.filePath)
	if err != nil {
		// Si le fichier n'existe pas, on part de zéro
		return
	}
	_ = json.Unmarshal(file, &r.seenIDs)
}

func (r *JSONRepository) save() error {
	data, err := json.MarshalIndent(r.seenIDs, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.filePath, data, 0644)
}
