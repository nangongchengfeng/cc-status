package runtime

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
)

func TestBootstrapCreatesAppDataDirAndDatabase(t *testing.T) {
	t.Parallel()

	appDir := filepath.Join(t.TempDir(), ".cc-usage-client")
	state, err := Bootstrap(Options{
		AppDir:    appDir,
		EnvLookup: func(string) string { return "" },
	})
	if err != nil {
		t.Fatalf("Bootstrap() returned error: %v", err)
	}
	t.Cleanup(func() {
		_ = state.Close()
	})

	if state.AppDir != appDir {
		t.Fatalf("expected AppDir %q, got %q", appDir, state.AppDir)
	}

	if _, err := os.Stat(appDir); err != nil {
		t.Fatalf("expected app dir to exist: %v", err)
	}

	dbPath := filepath.Join(appDir, "client.db")
	if _, err := os.Stat(dbPath); err != nil {
		t.Fatalf("expected database file to exist: %v", err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("sql.Open() returned error: %v", err)
	}
	defer db.Close()

	for _, tableName := range []string{"sync_state", "reported_ids", "metadata"} {
		var exists string
		err = db.QueryRow(
			`SELECT name FROM sqlite_master WHERE type = 'table' AND name = ?`,
			tableName,
		).Scan(&exists)
		if err != nil {
			t.Fatalf("expected table %q to exist: %v", tableName, err)
		}
	}
}

func TestBootstrapPersistsStableClientID(t *testing.T) {
	t.Parallel()

	appDir := filepath.Join(t.TempDir(), ".cc-usage-client")

	firstState, err := Bootstrap(Options{
		AppDir:    appDir,
		EnvLookup: func(string) string { return "" },
	})
	if err != nil {
		t.Fatalf("first Bootstrap() returned error: %v", err)
	}

	firstClientID := firstState.ClientID
	if firstClientID == "" {
		t.Fatal("expected client ID to be generated")
	}

	if err := firstState.Close(); err != nil {
		t.Fatalf("first Close() returned error: %v", err)
	}

	secondState, err := Bootstrap(Options{
		AppDir:    appDir,
		EnvLookup: func(string) string { return "" },
	})
	if err != nil {
		t.Fatalf("second Bootstrap() returned error: %v", err)
	}
	t.Cleanup(func() {
		_ = secondState.Close()
	})

	if secondState.ClientID != firstClientID {
		t.Fatalf("expected stable client ID, got %q then %q", firstClientID, secondState.ClientID)
	}
}
