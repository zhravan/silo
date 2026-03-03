package crypto

import (
	"crypto/rand"
	"crypto/subtle"
	"io"

	"golang.org/x/crypto/argon2"
)

const (
	KeySize   = 32
	SaltSize  = 16
	Time      = 1
	Memory    = 64 * 1024 // 64 MiB
	Threads   = 4
)

// DeriveKey derives a 32-byte key from password and salt using Argon2id.
func DeriveKey(password, salt []byte) []byte {
	return argon2.IDKey(password, salt, Time, Memory, Threads, KeySize)
}

// NewSalt returns a random salt of SaltSize bytes.
func NewSalt() ([]byte, error) {
	salt := make([]byte, SaltSize)
	_, err := io.ReadFull(rand.Reader, salt)
	return salt, err
}

// ConstantTimeCompare compares two byte slices in constant time.
func ConstantTimeCompare(a, b []byte) bool {
	return subtle.ConstantTimeCompare(a, b) == 1
}
