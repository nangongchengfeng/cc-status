package cli

import (
	"bytes"
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"cc-status/client/internal/lock"

	_ "modernc.org/sqlite"
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

func TestDryRunScansClaudeLogsAndPrintsSummary(t *testing.T) {
	t.Parallel()

	projectsDir := filepath.Join(t.TempDir(), "projects")
	projectDir := filepath.Join(projectsDir, "demo-project")
	if err := os.MkdirAll(projectDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() returned error: %v", err)
	}

	logPath := filepath.Join(projectDir, "session.jsonl")
	logContent := []byte("{\"type\":\"assistant\",\"message\":{\"id\":\"msg_1\",\"model\":\"claude-opus-4-1\",\"usage\":{\"input_tokens\":12,\"output_tokens\":3,\"cache_read_input_tokens\":2,\"cache_creation_input_tokens\":1},\"stop_reason\":\"end_turn\"},\"timestamp\":\"2026-04-05T12:00:00Z\",\"sessionId\":\"session-1\"}\n")
	if err := os.WriteFile(logPath, logContent, 0o644); err != nil {
		t.Fatalf("WriteFile() returned error: %v", err)
	}

	var stdout bytes.Buffer
	app := NewApp(NewBootstrapRunner(BootstrapOptions{
		AppDir:            filepath.Join(t.TempDir(), ".cc-usage-client"),
		ClaudeProjectsDir: projectsDir,
		EnvLookup: func(string) string {
			return ""
		},
		Stdout: &stdout,
	}))

	if err := app.Run([]string{"dry-run"}); err != nil {
		t.Fatalf("Run(dry-run) returned error: %v", err)
	}

	output := stdout.String()
	if !strings.Contains(output, "files_scanned=1") {
		t.Fatalf("expected files_scanned summary, got %q", output)
	}

	if !strings.Contains(output, "records=1") {
		t.Fatalf("expected record summary, got %q", output)
	}
}

func TestDryRunDoesNotWriteBusinessState(t *testing.T) {
	t.Parallel()

	projectsDir := filepath.Join(t.TempDir(), "projects")
	projectDir := filepath.Join(projectsDir, "demo-project")
	if err := os.MkdirAll(projectDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() returned error: %v", err)
	}

	logPath := filepath.Join(projectDir, "session.jsonl")
	logContent := []byte("{\"type\":\"assistant\",\"message\":{\"id\":\"msg_2\",\"model\":\"claude-opus-4-1\",\"usage\":{\"input_tokens\":9,\"output_tokens\":4},\"stop_reason\":\"end_turn\"},\"timestamp\":\"2026-04-05T12:00:00Z\",\"sessionId\":\"session-2\"}\n")
	if err := os.WriteFile(logPath, logContent, 0o644); err != nil {
		t.Fatalf("WriteFile() returned error: %v", err)
	}

	appDir := filepath.Join(t.TempDir(), ".cc-usage-client")
	app := NewApp(NewBootstrapRunner(BootstrapOptions{
		AppDir:            appDir,
		ClaudeProjectsDir: projectsDir,
		EnvLookup: func(string) string {
			return ""
		},
		Stdout: &bytes.Buffer{},
	}))

	if err := app.Run([]string{"dry-run"}); err != nil {
		t.Fatalf("Run(dry-run) returned error: %v", err)
	}

	db, err := sql.Open("sqlite", filepath.Join(appDir, "client.db"))
	if err != nil {
		t.Fatalf("sql.Open() returned error: %v", err)
	}
	defer db.Close()

	for _, tableName := range []string{"sync_state", "reported_ids"} {
		var count int
		if err := db.QueryRow("SELECT COUNT(*) FROM " + tableName).Scan(&count); err != nil {
			t.Fatalf("query %s count returned error: %v", tableName, err)
		}
		if count != 0 {
			t.Fatalf("expected %s to stay empty after dry-run, got %d rows", tableName, count)
		}
	}
}

func TestDryRunPrintsFileErrorsForBrokenTailLine(t *testing.T) {
	t.Parallel()

	projectsDir := filepath.Join(t.TempDir(), "projects")
	projectDir := filepath.Join(projectsDir, "demo-project")
	if err := os.MkdirAll(projectDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() returned error: %v", err)
	}

	logPath := filepath.Join(projectDir, "session.jsonl")
	logContent := []byte("{\"type\":\"assistant\",\"message\":{\"id\":\"msg_3\",\"model\":\"claude-opus-4-1\",\"usage\":{\"input_tokens\":9,\"output_tokens\":4},\"stop_reason\":\"end_turn\"},\"timestamp\":\"2026-04-05T12:00:00Z\",\"sessionId\":\"session-2\"}\n{\"type\":\"assistant\"")
	if err := os.WriteFile(logPath, logContent, 0o644); err != nil {
		t.Fatalf("WriteFile() returned error: %v", err)
	}

	var stdout bytes.Buffer
	app := NewApp(NewBootstrapRunner(BootstrapOptions{
		AppDir:            filepath.Join(t.TempDir(), ".cc-usage-client"),
		ClaudeProjectsDir: projectsDir,
		EnvLookup: func(string) string {
			return ""
		},
		Stdout: &stdout,
	}))

	if err := app.Run([]string{"dry-run"}); err != nil {
		t.Fatalf("Run(dry-run) returned error: %v", err)
	}

	output := stdout.String()
	if !strings.Contains(output, "errors=1") {
		t.Fatalf("expected error count in output, got %q", output)
	}

	if !strings.Contains(output, "session.jsonl") {
		t.Fatalf("expected file error details in output, got %q", output)
	}
}
