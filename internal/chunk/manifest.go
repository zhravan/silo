package chunk

// ManifestEntry is a single file's metadata and its chunk IDs in order.
type ManifestEntry struct {
	Path     string
	ChunkIDs []ChunkID
}

// Manifest holds the list of entries (path → chunk IDs) for a backup.
type Manifest struct {
	Entries []ManifestEntry
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
