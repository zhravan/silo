package backup

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/shravan20/silo/internal/backend"
	"github.com/shravan20/silo/internal/chunk"
	"github.com/shravan20/silo/internal/compress"
	"github.com/shravan20/silo/internal/config"
	"github.com/shravan20/silo/internal/crypto"
	"github.com/shravan20/silo/internal/index"
)

const saltKey = ".salt"

// Run runs a full backup: scan sources, chunk, compress, dedupe (index), encrypt, upload; then upload encrypted manifest.
func Run(ctx context.Context, cfg *config.Config, indexPath string) error {
	idx, err := index.Open(indexPath)
	if err != nil {
		return fmt.Errorf("index: %w", err)
	}
	defer idx.Close()

	var be backend.Backend
	if cfg.Backend.Type == "s3" {
		s3be, err := backend.NewS3(ctx, backend.S3Config{
			Bucket: cfg.Backend.Bucket,
			Prefix: cfg.Backend.Prefix,
			Region: cfg.Backend.Region,
		})
		if err != nil {
			return fmt.Errorf("s3: %w", err)
		}
		be = s3be
	} else {
		return fmt.Errorf("unsupported backend type %q", cfg.Backend.Type)
	}

	password := os.Getenv(cfg.Encryption.PasswordEnv)
	if password == "" {
		return fmt.Errorf("password env %q not set", cfg.Encryption.PasswordEnv)
	}

	salt, err := getOrCreateSalt(be)
	if err != nil {
		return err
	}
	key := crypto.DeriveKey([]byte(password), salt)

	comp, err := compress.New(cfg.Compression.Type, cfg.Compression.Level)
	if err != nil {
		return err
	}
	if z, ok := comp.(*compress.Zstd); ok {
		defer z.Close()
	}

	chunker := chunk.FixedChunker{Size: cfg.Chunking.TargetSize}
	if chunker.Size <= 0 {
		chunker.Size = 4 << 20
	}

	manifestBuilder := &chunk.Builder{}
	backendType := cfg.Backend.Type

	for _, root := range cfg.Sources {
		err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			if excluded(path, root, cfg.Exclude) {
				if d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			ids, payloads, err := ProcessFile(f, chunker, comp)
			f.Close()
			if err != nil {
				return err
			}
			manifestBuilder.Add(path, ids)
			for i, id := range ids {
				ok, err := idx.Has(id)
				if err != nil {
					return err
				}
				if ok {
					continue
				}
				encrypted, err := crypto.Encrypt(payloads[i], key)
				if err != nil {
					return err
				}
				objKey := chunkKey(id)
				if err := be.Put(objKey, bytes.NewReader(encrypted)); err != nil {
					return err
				}
				if err := idx.Put(id, objKey, backendType); err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	manifest := manifestBuilder.Build()
	manifestJSON, err := json.Marshal(manifest)
	if err != nil {
		return err
	}
	encryptedManifest, err := crypto.Encrypt(manifestJSON, key)
	if err != nil {
		return err
	}
	backupID := time.Now().UTC().Format("2006-01-02T15-04-05")
	manifestKey := "manifests/" + backupID + ".enc"
	if err := be.Put(manifestKey, bytes.NewReader(encryptedManifest)); err != nil {
		return err
	}
	_ = backupID
	return nil
}

func chunkKey(id chunk.ChunkID) string {
	s := id.String()
	if len(s) < 2 {
		return "chunks/xx/" + s + ".enc"
	}
	return "chunks/" + s[:2] + "/" + s + ".enc"
}

func getOrCreateSalt(be backend.Backend) ([]byte, error) {
	rc, err := be.Get(saltKey)
	if err == nil {
		defer rc.Close()
		return io.ReadAll(rc)
	}
	salt, err := crypto.NewSalt()
	if err != nil {
		return nil, err
	}
	if err := be.Put(saltKey, bytes.NewReader(salt)); err != nil {
		return nil, err
	}
	return salt, nil
}

func excluded(path, root string, patterns []string) bool {
	rel, _ := filepath.Rel(root, path)
	for _, p := range patterns {
		if p == "" {
			continue
		}
		if strings.HasSuffix(p, "/") {
			if strings.HasPrefix(rel, p) || strings.Contains(rel, p) {
				return true
			}
		}
		matched, _ := filepath.Match(p, filepath.Base(path))
		if matched {
			return true
		}
	}
	return false
}

// Restore restores a backup: get manifest, then for each file get chunks, decrypt, decompress, write.
func Restore(ctx context.Context, cfg *config.Config, indexPath, backupID, destDir string) error {
	idx, err := index.Open(indexPath)
	if err != nil {
		return fmt.Errorf("index: %w", err)
	}
	defer idx.Close()

	var be backend.Backend
	if cfg.Backend.Type == "s3" {
		s3be, err := backend.NewS3(ctx, backend.S3Config{
			Bucket: cfg.Backend.Bucket,
			Prefix: cfg.Backend.Prefix,
			Region: cfg.Backend.Region,
		})
		if err != nil {
			return fmt.Errorf("s3: %w", err)
		}
		be = s3be
	} else {
		return fmt.Errorf("unsupported backend type %q", cfg.Backend.Type)
	}

	password := os.Getenv(cfg.Encryption.PasswordEnv)
	if password == "" {
		return fmt.Errorf("password env %q not set", cfg.Encryption.PasswordEnv)
	}

	salt, err := getOrCreateSalt(be)
	if err != nil {
		return err
	}
	key := crypto.DeriveKey([]byte(password), salt)

	manifestKey := "manifests/" + backupID + ".enc"
	rc, err := be.Get(manifestKey)
	if err != nil {
		return fmt.Errorf("get manifest: %w", err)
	}
	encrypted, err := io.ReadAll(rc)
	rc.Close()
	if err != nil {
		return err
	}
	manifestJSON, err := crypto.Decrypt(encrypted, key)
	if err != nil {
		return err
	}
	var manifest chunk.Manifest
	if err := json.Unmarshal(manifestJSON, &manifest); err != nil {
		return err
	}

	comp, err := compress.New(cfg.Compression.Type, cfg.Compression.Level)
	if err != nil {
		return err
	}
	if z, ok := comp.(*compress.Zstd); ok {
		defer z.Close()
	}

	for _, e := range manifest.Entries {
		destPath := filepath.Join(destDir, e.Path)
		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return err
		}
		f, err := os.Create(destPath)
		if err != nil {
			return err
		}
		for _, id := range e.ChunkIDs {
			objKey, _, err := idx.Get(id)
			if err != nil {
				f.Close()
				return err
			}
			rc, err := be.Get(objKey)
			if err != nil {
				f.Close()
				return err
			}
			encrypted, err := io.ReadAll(rc)
			rc.Close()
			if err != nil {
				f.Close()
				return err
			}
			compressed, err := crypto.Decrypt(encrypted, key)
			if err != nil {
				f.Close()
				return err
			}
			data, err := comp.Decompress(compressed)
			if err != nil {
				f.Close()
				return err
			}
			if _, err := f.Write(data); err != nil {
				f.Close()
				return err
			}
		}
		f.Close()
	}
	return nil
}
