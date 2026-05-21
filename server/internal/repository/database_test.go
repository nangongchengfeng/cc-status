package repository

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOpenDatabaseCreatesParentDirAndFile(t *testing.T) {
	t.Parallel()

	dbPath := filepath.Join(t.TempDir(), "server", "data", "server.db")
	db, err := OpenDatabase(dbPath)
	if err != nil {
		t.Fatalf("OpenDatabase() returned error: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("DB() returned error: %v", err)
	}
	t.Cleanup(func() {
		_ = sqlDB.Close()
	})

	if _, err := os.Stat(filepath.Dir(dbPath)); err != nil {
		t.Fatalf("expected parent dir to exist: %v", err)
	}
	if _, err := os.Stat(dbPath); err != nil {
		t.Fatalf("expected sqlite file to exist: %v", err)
	}
}
