package chunk

import "io"

// Chunker splits a stream into chunks.
type Chunker interface {
	// Chunk reads from r and returns chunk payloads (each is a complete chunk).
	// Caller may reuse or discard the returned byte slices.
	Chunk(r io.Reader) ([][]byte, error)
}

// FixedChunker splits data into fixed-size chunks.
type FixedChunker struct {
	Size int
}

// Chunk implements Chunker by reading Size-sized blocks.
func (f FixedChunker) Chunk(r io.Reader) ([][]byte, error) {
	if f.Size <= 0 {
		f.Size = 4 << 20 // 4 MiB default
	}
	var chunks [][]byte
	buf := make([]byte, f.Size)
	for {
		n, err := io.ReadFull(r, buf)
		if n > 0 {
			chunk := make([]byte, n)
			copy(chunk, buf[:n])
			chunks = append(chunks, chunk)
		}
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		}
		if err != nil {
			return chunks, err
		}
	}
	return chunks, nil
}
