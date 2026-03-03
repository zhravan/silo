package compress

import (
	"fmt"

	"github.com/klauspost/compress/zstd"
)

// Compressor compresses and decompresses data.
type Compressor interface {
	Compress(src []byte) ([]byte, error)
	Decompress(src []byte) ([]byte, error)
}

// None is a no-op compressor (pass-through).
type None struct{}

func (None) Compress(src []byte) ([]byte, error) {
	out := make([]byte, len(src))
	copy(out, src)
	return out, nil
}

func (None) Decompress(src []byte) ([]byte, error) {
	out := make([]byte, len(src))
	copy(out, src)
	return out, nil
}

// Zstd compresses with Zstandard.
type Zstd struct {
	enc *zstd.Encoder
	dec *zstd.Decoder
}

// NewZstd creates a Zstd compressor. Level 1-22; 3 is a good default.
func NewZstd(level int) (*Zstd, error) {
	if level <= 0 {
		level = 3
	}
	enc, err := zstd.NewWriter(nil, zstd.WithEncoderLevel(zstd.EncoderLevelFromZstd(level)))
	if err != nil {
		return nil, err
	}
	dec, err := zstd.NewReader(nil)
	if err != nil {
		enc.Close()
		return nil, err
	}
	return &Zstd{enc: enc, dec: dec}, nil
}

func (z *Zstd) Compress(src []byte) ([]byte, error) {
	return z.enc.EncodeAll(src, nil), nil
}

func (z *Zstd) Decompress(src []byte) ([]byte, error) {
	return z.dec.DecodeAll(src, nil)
}

// Close releases resources. Call when done with the compressor.
func (z *Zstd) Close() {
	z.enc.Close()
	z.dec.Close()
}

// New returns a Compressor for the given type (zstd, lz4, none or "").
func New(typ string, level int) (Compressor, error) {
	switch typ {
	case "", "none":
		return None{}, nil
	case "zstd":
		return NewZstd(level)
	default:
		return nil, fmt.Errorf("unknown compression type %q", typ)
	}
}
