package cli

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	_ "modernc.org/sqlite"
)

func TestSyncReportsUsageAndAdvancesStateOnSuccess(t *testing.T) {
	t.Parallel()

	projectsDir := filepath.Join(t.TempDir(), "projects")
	projectDir := filepath.Join(projectsDir, "demo-project")
	if err := os.MkdirAll(projectDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() returned error: %v", err)
	}

	logPath := filepath.Join(projectDir, "session.jsonl")
	logContent := []byte("{\"type\":\"assistant\",\"message\":{\"id\":\"msg_sync_1\",\"model\":\"claude-opus-4-1\",\"usage\":{\"input_tokens\":15,\"output_tokens\":6,\"cache_read_input_tokens\":3,\"cache_creation_input_tokens\":1},\"stop_reason\":\"end_turn\"},\"timestamp\":\"2026-04-05T12:00:00Z\",\"sessionId\":\"session-sync-1\"}\n")
	if err := os.WriteFile(logPath, logContent, 0o644); err != nil {
		t.Fatalf("WriteFile() returned error: %v", err)
	}

	var captured struct {
		ClientID string `json:"client_id"`
		Reports  []struct {
			RequestID           string `json:"request_id"`
			AppType             string `json:"app_type"`
			Model               string `json:"model"`
			InputTokens         uint32 `json:"input_tokens"`
			OutputTokens        uint32 `json:"output_tokens"`
			CacheReadTokens     uint32 `json:"cache_read_tokens"`
			CacheCreationTokens uint32 `json:"cache_creation_tokens"`
			SessionID           string `json:"session_id"`
			DataSource          string `json:"data_source"`
		} `json:"reports"`
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/sync" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}

		if err := json.NewDecoder(r.Body).Decode(&captured); err != nil {
			t.Fatalf("Decode() returned error: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":0,"message":"success","accepted_count":1,"duplicate_count":0}`))
	}))
	defer server.Close()

	appDir := filepath.Join(t.TempDir(), ".cc-usage-client")
	var stdout bytes.Buffer
	app := NewApp(NewBootstrapRunner(BootstrapOptions{
		AppDir:            appDir,
		ClaudeProjectsDir: projectsDir,
		EnvLookup: func(key string) string {
			switch key {
			case "CC_USAGE_CLIENT_SERVER_URL":
				return server.URL
			case "CC_USAGE_CLIENT_AUTH_TOKEN":
				return "token-from-test"
			default:
				return ""
			}
		},
		Stdout: &stdout,
	}))

	if err := app.Run([]string{"sync"}); err != nil {
		t.Fatalf("Run(sync) returned error: %v", err)
	}

	if captured.ClientID == "" {
		t.Fatal("expected client_id to be sent")
	}

	if len(captured.Reports) != 1 {
		t.Fatalf("expected 1 report, got %d", len(captured.Reports))
	}

	report := captured.Reports[0]
	if report.RequestID != "session:msg_sync_1" {
		t.Fatalf("unexpected request_id: %q", report.RequestID)
	}
	if report.AppType != "claude" || report.DataSource != "session_log" {
		t.Fatalf("unexpected report routing fields: %#v", report)
	}
	if report.Model != "claude-opus-4-1" || report.SessionID != "session-sync-1" {
		t.Fatalf("unexpected report identity fields: %#v", report)
	}

	db, err := sql.Open("sqlite", filepath.Join(appDir, "client.db"))
	if err != nil {
		t.Fatalf("sql.Open() returned error: %v", err)
	}
	defer db.Close()

	var reportedCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM reported_ids WHERE request_id = ?`, "session:msg_sync_1").Scan(&reportedCount); err != nil {
		t.Fatalf("query reported_ids returned error: %v", err)
	}
	if reportedCount != 1 {
		t.Fatalf("expected reported_ids row, got %d", reportedCount)
	}

	var lastLineOffset int
	if err := db.QueryRow(`SELECT last_line_offset FROM sync_state WHERE file_path = ?`, logPath).Scan(&lastLineOffset); err != nil {
		t.Fatalf("query sync_state returned error: %v", err)
	}
	if lastLineOffset != 1 {
		t.Fatalf("expected last_line_offset=1, got %d", lastLineOffset)
	}

	if !strings.Contains(stdout.String(), "accepted=1") {
		t.Fatalf("expected sync summary, got %q", stdout.String())
	}
}

func TestSyncReturnsSuccessWhenNoNewRecordsExist(t *testing.T) {
	t.Parallel()

	projectsDir := filepath.Join(t.TempDir(), "projects")
	projectDir := filepath.Join(projectsDir, "demo-project")
	if err := os.MkdirAll(projectDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() returned error: %v", err)
	}

	logPath := filepath.Join(projectDir, "session.jsonl")
	logContent := []byte("{\"type\":\"user\",\"message\":{\"id\":\"msg_user\",\"model\":\"claude-opus-4-1\",\"usage\":{\"input_tokens\":15,\"output_tokens\":6},\"stop_reason\":\"end_turn\"},\"timestamp\":\"2026-04-05T12:00:00Z\",\"sessionId\":\"session-sync-1\"}\n")
	if err := os.WriteFile(logPath, logContent, 0o644); err != nil {
		t.Fatalf("WriteFile() returned error: %v", err)
	}

	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		t.Fatalf("unexpected upstream request when no new records exist")
	}))
	defer server.Close()

	var stdout bytes.Buffer
	app := NewApp(NewBootstrapRunner(BootstrapOptions{
		AppDir:            filepath.Join(t.TempDir(), ".cc-usage-client"),
		ClaudeProjectsDir: projectsDir,
		EnvLookup: func(key string) string {
			if key == "CC_USAGE_CLIENT_SERVER_URL" {
				return server.URL
			}
			return ""
		},
		Stdout: &stdout,
	}))

	if err := app.Run([]string{"sync"}); err != nil {
		t.Fatalf("Run(sync) returned error: %v", err)
	}

	if requestCount != 0 {
		t.Fatalf("expected no upstream request, got %d", requestCount)
	}

	if !strings.Contains(stdout.String(), "accepted=0") {
		t.Fatalf("expected zero summary, got %q", stdout.String())
	}
}

func TestSyncUsesDefaultBatchSizeOfFiveHundred(t *testing.T) {
	t.Parallel()

	projectsDir := filepath.Join(t.TempDir(), "projects")
	projectDir := filepath.Join(projectsDir, "demo-project")
	if err := os.MkdirAll(projectDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() returned error: %v", err)
	}

	var content strings.Builder
	for index := 0; index < 501; index++ {
		content.WriteString("{\"type\":\"assistant\",\"message\":{\"id\":\"msg_")
		content.WriteString(strconv.Itoa(index))
		content.WriteString("\",\"model\":\"claude-opus-4-1\",\"usage\":{\"input_tokens\":1,\"output_tokens\":1},\"stop_reason\":\"end_turn\"},\"timestamp\":\"2026-04-05T12:00:00Z\",\"sessionId\":\"session-batch\"}\n")
	}
	if err := os.WriteFile(filepath.Join(projectDir, "session.jsonl"), []byte(content.String()), 0o644); err != nil {
		t.Fatalf("WriteFile() returned error: %v", err)
	}

	batchSizes := make([]int, 0)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			Reports []json.RawMessage `json:"reports"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("Decode() returned error: %v", err)
		}
		batchSizes = append(batchSizes, len(payload.Reports))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":0,"message":"success","accepted_count":500,"duplicate_count":0}`))
		if len(payload.Reports) == 1 {
			_, _ = w.Write([]byte{})
		}
	}))
	defer server.Close()

	app := NewApp(NewBootstrapRunner(BootstrapOptions{
		AppDir:            filepath.Join(t.TempDir(), ".cc-usage-client"),
		ClaudeProjectsDir: projectsDir,
		EnvLookup: func(key string) string {
			if key == "CC_USAGE_CLIENT_SERVER_URL" {
				return server.URL
			}
			return ""
		},
		Stdout: &bytes.Buffer{},
	}))

	if err := app.Run([]string{"sync"}); err != nil {
		t.Fatalf("Run(sync) returned error: %v", err)
	}

	if len(batchSizes) != 2 || batchSizes[0] != 500 || batchSizes[1] != 1 {
		t.Fatalf("expected batch sizes [500 1], got %#v", batchSizes)
	}
}
