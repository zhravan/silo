package crypto

import (
	"crypto/rand"
	"io"

	"golang.org/x/crypto/chacha20poly1305"
)

const NonceSize = chacha20poly1305.NonceSize

// Encrypt encrypts plaintext with key (KeySize bytes). Prepends nonce to ciphertext.
func Encrypt(plaintext, key []byte) ([]byte, error) {
	aead, err := chacha20poly1305.New(key)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, NonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	ciphertext := aead.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// Decrypt decrypts ciphertext (nonce prepended) with key.
func Decrypt(ciphertext, key []byte) ([]byte, error) {
	aead, err := chacha20poly1305.New(key)
	if err != nil {
		return nil, err
	}
	if len(ciphertext) < NonceSize {
		return nil, errInvalidCiphertext
	}
	nonce, ct := ciphertext[:NonceSize], ciphertext[NonceSize:]
	return aead.Open(nil, nonce, ct, nil)
}

var errInvalidCiphertext = errMsg("ciphertext too short")

type errMsg string

func (e errMsg) Error() string { return string(e) }
