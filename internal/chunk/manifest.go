package chunk

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// ManifestEntry is a single file's metadata and its chunk IDs in order.
type ManifestEntry struct {
	Path     string
	ChunkIDs []ChunkID
}

// Manifest holds the list of entries (path → chunk IDs) for a backup.
type Manifest struct {
	Entries []ManifestEntry
}

// MarshalJSON serializes the manifest for encryption and storage.
func (m *Manifest) MarshalJSON() ([]byte, error) {
	entries := make([]struct {
		Path     string   `json:"path"`
		ChunkIDs []string `json:"chunk_ids"`
	}, len(m.Entries))
	for i, e := range m.Entries {
		entries[i].Path = e.Path
		entries[i].ChunkIDs = make([]string, len(e.ChunkIDs))
		for j, id := range e.ChunkIDs {
			entries[i].ChunkIDs[j] = id.String()
		}
	}
	return json.Marshal(entries)
}

// UnmarshalJSON deserializes the manifest after decryption.
func (m *Manifest) UnmarshalJSON(data []byte) error {
	var entries []struct {
		Path     string   `json:"path"`
		ChunkIDs []string `json:"chunk_ids"`
	}
	if err := json.Unmarshal(data, &entries); err != nil {
		return err
	}
	m.Entries = make([]ManifestEntry, len(entries))
	for i, e := range entries {
		m.Entries[i].Path = e.Path
		m.Entries[i].ChunkIDs = make([]ChunkID, len(e.ChunkIDs))
		for j, s := range e.ChunkIDs {
			b, err := hex.DecodeString(s)
			if err != nil {
				return err
			}
			if len(b) != 32 {
				return fmt.Errorf("chunk_id must be 32 bytes")
			}
			copy(m.Entries[i].ChunkIDs[j][:], b)
		}
	}
	return nil
}

// Builder builds a manifest by adding path → chunk IDs.
type Builder struct {
	entries []ManifestEntry
}

// Add appends an entry for the given path and ordered chunk IDs.
func (b *Builder) Add(path string, ids []ChunkID) {
	b.entries = append(b.entries, ManifestEntry{Path: path, ChunkIDs: ids})
}

// Build returns the manifest.
func (b *Builder) Build() *Manifest {
	entries := make([]ManifestEntry, len(b.entries))
	copy(entries, b.entries)
	return &Manifest{Entries: entries}
}
