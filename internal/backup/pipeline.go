package backup

import (
	"io"

	"github.com/shravan20/silo/internal/chunk"
	"github.com/shravan20/silo/internal/compress"
)

// ProcessFile reads from r, chunks with ch, compresses each chunk with comp, hashes, and returns chunk IDs in order.
// Pipeline: read → chunk → compress → hash.
func ProcessFile(r io.Reader, ch chunk.Chunker, comp compress.Compressor) ([]chunk.ChunkID, error) {
	rawChunks, err := ch.Chunk(r)
	if err != nil {
		return nil, err
	}
	ids := make([]chunk.ChunkID, 0, len(rawChunks))
	for _, raw := range rawChunks {
		compressed, err := comp.Compress(raw)
		if err != nil {
			return nil, err
		}
		ids = append(ids, chunk.Hash(compressed))
	}
	return ids, nil
}
