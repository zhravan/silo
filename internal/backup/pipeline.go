package backup

import (
	"io"

	"github.com/shravan20/silo/internal/chunk"
	"github.com/shravan20/silo/internal/compress"
)

// ProcessFile reads from r, chunks with ch, compresses each chunk with comp, hashes, and returns chunk IDs and compressed payloads in order.
// Pipeline: read → chunk → compress → hash. Payloads are the compressed bytes (for encrypt and upload).
func ProcessFile(r io.Reader, ch chunk.Chunker, comp compress.Compressor) (ids []chunk.ChunkID, payloads [][]byte, err error) {
	rawChunks, err := ch.Chunk(r)
	if err != nil {
		return nil, nil, err
	}
	ids = make([]chunk.ChunkID, 0, len(rawChunks))
	payloads = make([][]byte, 0, len(rawChunks))
	for _, raw := range rawChunks {
		compressed, err := comp.Compress(raw)
		if err != nil {
			return nil, nil, err
		}
		ids = append(ids, chunk.Hash(compressed))
		payloads = append(payloads, compressed)
	}
	return ids, payloads, nil
}
