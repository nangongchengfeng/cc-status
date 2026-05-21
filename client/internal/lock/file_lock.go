package lock

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type FileLock struct {
	path string
	file *os.File
}

func Acquire(path string) (*FileLock, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("create lock dir: %w", err)
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return nil, fmt.Errorf("another instance is already running")
		}
		return nil, fmt.Errorf("acquire lock: %w", err)
	}

	return &FileLock{
		path: path,
		file: file,
	}, nil
}

func (lock *FileLock) Release() error {
	if lock == nil {
		return nil
	}

	if lock.file != nil {
		if err := lock.file.Close(); err != nil {
			return err
		}
	}

	if err := os.Remove(lock.path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	return nil
}
