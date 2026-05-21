package lock

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestAcquirePreventsSecondLock(t *testing.T) {
	t.Parallel()

	lockPath := filepath.Join(t.TempDir(), "client.lock")

	firstLock, err := Acquire(lockPath)
	if err != nil {
		t.Fatalf("first Acquire() returned error: %v", err)
	}
	defer firstLock.Release()

	secondLock, err := Acquire(lockPath)
	if err == nil {
		_ = secondLock.Release()
		t.Fatal("expected second Acquire() to fail")
	}

	if !strings.Contains(err.Error(), "another instance is already running") {
		t.Fatalf("unexpected error: %v", err)
	}
}
