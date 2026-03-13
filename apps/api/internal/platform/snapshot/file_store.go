package snapshot

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

// FileStore persists snapshots on local filesystem.
type FileStore struct {
	dir string
}

func NewFileStore(dir string) *FileStore {
	return &FileStore{dir: dir}
}

func (s *FileStore) Load(_ context.Context, module string) ([]byte, error) {
	if s == nil {
		return nil, nil
	}
	path := filepath.Join(s.dir, module+".snapshot")
	payload, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("file snapshot load failed for %s: %w", module, err)
	}
	if len(payload) == 0 {
		return nil, nil
	}
	return payload, nil
}

func (s *FileStore) Save(_ context.Context, module string, payload []byte) error {
	if s == nil {
		return nil
	}
	if err := os.MkdirAll(s.dir, 0o700); err != nil {
		return fmt.Errorf("file snapshot ensure dir failed: %w", err)
	}
	if err := os.Chmod(s.dir, 0o700); err != nil {
		return fmt.Errorf("file snapshot chmod dir failed: %w", err)
	}

	path := filepath.Join(s.dir, module+".snapshot")
	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, payload, 0o600); err != nil {
		return fmt.Errorf("file snapshot temp write failed for %s: %w", module, err)
	}
	if err := os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("file snapshot rename failed for %s: %w", module, err)
	}
	if err := os.Chmod(path, 0o600); err != nil {
		return fmt.Errorf("file snapshot chmod file failed for %s: %w", module, err)
	}
	return nil
}
