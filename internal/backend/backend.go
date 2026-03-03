package backend

import "io"

// Backend stores and retrieves opaque blobs by key.
type Backend interface {
	// Put uploads data under the given key.
	Put(key string, r io.Reader) error
	// Get downloads the blob for the given key. Returns ErrNotFound if missing.
	Get(key string) (io.ReadCloser, error)
	// List returns keys with the given prefix (e.g. "chunks/ab/").
	List(prefix string) ([]string, error)
}

// ErrNotFound is returned by Get when the key does not exist.
var ErrNotFound = errNotFound{}

type errNotFound struct{}

func (errNotFound) Error() string { return "key not found" }
