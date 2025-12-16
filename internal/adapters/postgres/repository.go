package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/maximerauch/go-classifieds-watcher/internal/core"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(databaseURL string) (*Repository, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("database unreachable: %w", err)
	}

	query := `CREATE TABLE IF NOT EXISTS seen_items (
		id TEXT PRIMARY KEY,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(query); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	return &Repository{db: db}, nil
}

// Exists checks if an ID is already in the database
func (r *Repository) Exists(ctx context.Context, id string) (bool, error) {
	var exists int
	query := "SELECT 1 FROM seen_items WHERE id = $1"

	err := r.db.QueryRowContext(ctx, query, id).Scan(&exists)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// Save inserts a new ID into the database
func (r *Repository) Save(ctx context.Context, item core.Item) error {
	query := "INSERT INTO seen_items (id) VALUES ($1) ON CONFLICT (id) DO NOTHING"
	_, err := r.db.ExecContext(ctx, query, item.ID)
	return err
}
