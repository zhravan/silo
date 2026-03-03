package index

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"

	"github.com/shravan20/silo/internal/chunk"
)

// Index maps chunk IDs to backend keys for deduplication.
type Index struct {
	db *sql.DB
}

// Open opens or creates the index database at path.
func Open(path string) (*Index, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open index: %w", err)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping index: %w", err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS chunks (
		chunk_id BLOB PRIMARY KEY,
		backend_key TEXT NOT NULL,
		backend_type TEXT NOT NULL
	)`)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("create table: %w", err)
	}
	return &Index{db: db}, nil
}

// Has returns true if the chunk ID is already in the index.
func (idx *Index) Has(id chunk.ChunkID) (bool, error) {
	var exists int
	err := idx.db.QueryRow("SELECT EXISTS(SELECT 1 FROM chunks WHERE chunk_id = ?)", id[:]).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}

// Put records a chunk ID and its backend key and type.
func (idx *Index) Put(id chunk.ChunkID, backendKey, backendType string) error {
	_, err := idx.db.Exec("INSERT OR REPLACE INTO chunks (chunk_id, backend_key, backend_type) VALUES (?, ?, ?)",
		id[:], backendKey, backendType)
	return err
}

// Get returns the backend key and type for a chunk ID, or error if not found.
func (idx *Index) Get(id chunk.ChunkID) (backendKey, backendType string, err error) {
	err = idx.db.QueryRow("SELECT backend_key, backend_type FROM chunks WHERE chunk_id = ?", id[:]).
		Scan(&backendKey, &backendType)
	return backendKey, backendType, err
}

// Close closes the index database.
func (idx *Index) Close() error {
	return idx.db.Close()
}
