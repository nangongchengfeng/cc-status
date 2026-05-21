package cli

import (
	"path/filepath"
	"strings"
	"testing"

	"cc-status/client/internal/lock"
)

func TestAppRunsKnownCommands(t *testing.T) {
	t.Parallel()

	var called []string
	app := NewApp(func(command string) error {
		called = append(called, command)
		return nil
	})

	if err := app.Run([]string{"sync"}); err != nil {
		t.Fatalf("Run(sync) returned error: %v", err)
	}

	if err := app.Run([]string{"dry-run"}); err != nil {
		t.Fatalf("Run(dry-run) returned error: %v", err)
	}

	if len(called) != 2 {
		t.Fatalf("expected 2 command invocations, got %d", len(called))
	}

	if called[0] != "sync" || called[1] != "dry-run" {
		t.Fatalf("unexpected command order: %#v", called)
	}
}

func TestAppReturnsHelpfulErrorWhenLocked(t *testing.T) {
	t.Parallel()

	appDir := filepath.Join(t.TempDir(), ".cc-usage-client")
	heldLock, err := lock.Acquire(filepath.Join(appDir, "client.lock"))
	if err != nil {
		t.Fatalf("Acquire() returned error: %v", err)
	}
	defer heldLock.Release()

	app := NewApp(NewBootstrapRunner(BootstrapOptions{
		AppDir: appDir,
		EnvLookup: func(string) string {
			return ""
		},
	}))

	err = app.Run([]string{"sync"})
	if err == nil {
		t.Fatal("expected locked app run to fail")
	}

	if !strings.Contains(err.Error(), "another instance is already running") {
		t.Fatalf("expected helpful lock error, got %v", err)
	}
}
