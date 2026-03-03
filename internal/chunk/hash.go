package chunk

import (
	"encoding/hex"

	"github.com/zeebo/blake3"
)

// ChunkID is the BLAKE3 hash of a chunk (32 bytes).
type ChunkID [32]byte

// Hash returns the BLAKE3 hash of data as a ChunkID.
func Hash(data []byte) ChunkID {
	return ChunkID(blake3.Sum256(data))
}

// String returns the hex-encoded chunk ID.
func (c ChunkID) String() string {
	return hex.EncodeToString(c[:])
}
