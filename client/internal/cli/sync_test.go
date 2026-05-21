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
	if !strings.Contains(stdout.String(), "records=1") || !strings.Contains(stdout.String(), "skipped=0") {
		t.Fatalf("expected stable sync summary fields, got %q", stdout.String())
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
	if !strings.Contains(stdout.String(), "records=0") || !strings.Contains(stdout.String(), "skipped=0") {
		t.Fatalf("expected stable zero summary fields, got %q", stdout.String())
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

func TestSyncPreservesSuccessfulBatchWhenLaterBatchFails(t *testing.T) {
	t.Parallel()

	projectsDir := filepath.Join(t.TempDir(), "projects")
	projectDir := filepath.Join(projectsDir, "demo-project")
	if err := os.MkdirAll(projectDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() returned error: %v", err)
	}

	var content strings.Builder
	for index := 0; index < 501; index++ {
		content.WriteString("{\"type\":\"assistant\",\"message\":{\"id\":\"partial_")
		content.WriteString(strconv.Itoa(index))
		content.WriteString("\",\"model\":\"claude-opus-4-1\",\"usage\":{\"input_tokens\":1,\"output_tokens\":1},\"stop_reason\":\"end_turn\"},\"timestamp\":\"2026-04-05T12:00:00Z\",\"sessionId\":\"session-partial\"}\n")
	}
	logPath := filepath.Join(projectDir, "session.jsonl")
	if err := os.WriteFile(logPath, []byte(content.String()), 0o644); err != nil {
		t.Fatalf("WriteFile() returned error: %v", err)
	}

	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.Header().Set("Content-Type", "application/json")
		if requestCount == 1 {
			_, _ = w.Write([]byte(`{"code":0,"message":"success","accepted_count":500,"duplicate_count":0}`))
			return
		}

		http.Error(w, "upstream failed", http.StatusBadGateway)
	}))
	defer server.Close()

	appDir := filepath.Join(t.TempDir(), ".cc-usage-client")
	app := NewApp(NewBootstrapRunner(BootstrapOptions{
		AppDir:            appDir,
		ClaudeProjectsDir: projectsDir,
		EnvLookup: func(key string) string {
			if key == "CC_USAGE_CLIENT_SERVER_URL" {
				return server.URL
			}
			if key == "CC_USAGE_CLIENT_TIMEOUT_SECONDS" {
				return "1"
			}
			return ""
		},
		Stdout: &bytes.Buffer{},
	}))

	err := app.Run([]string{"sync"})
	if err == nil {
		t.Fatal("expected sync to fail on later batch")
	}

	db, openErr := sql.Open("sqlite", filepath.Join(appDir, "client.db"))
	if openErr != nil {
		t.Fatalf("sql.Open() returned error: %v", openErr)
	}
	defer db.Close()

	var reportedCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM reported_ids`).Scan(&reportedCount); err != nil {
		t.Fatalf("query reported_ids returned error: %v", err)
	}
	if reportedCount != 500 {
		t.Fatalf("expected first successful batch to persist 500 rows, got %d", reportedCount)
	}

	var syncStateCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM sync_state WHERE file_path = ?`, logPath).Scan(&syncStateCount); err != nil {
		t.Fatalf("query sync_state returned error: %v", err)
	}
	if syncStateCount != 0 {
		t.Fatalf("expected failed file not to advance sync_state, got %d rows", syncStateCount)
	}
}

func TestSyncRetriesOnServerErrorsAndEventuallySucceeds(t *testing.T) {
	t.Parallel()

	projectsDir := filepath.Join(t.TempDir(), "projects")
	projectDir := filepath.Join(projectsDir, "demo-project")
	if err := os.MkdirAll(projectDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() returned error: %v", err)
	}

	logPath := filepath.Join(projectDir, "session.jsonl")
	logContent := []byte("{\"type\":\"assistant\",\"message\":{\"id\":\"retry_msg\",\"model\":\"claude-opus-4-1\",\"usage\":{\"input_tokens\":2,\"output_tokens\":1},\"stop_reason\":\"end_turn\"},\"timestamp\":\"2026-04-05T12:00:00Z\",\"sessionId\":\"session-retry\"}\n")
	if err := os.WriteFile(logPath, logContent, 0o644); err != nil {
		t.Fatalf("WriteFile() returned error: %v", err)
	}

	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		if requestCount < 3 {
			http.Error(w, "temporary failure", http.StatusBadGateway)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":0,"message":"success","accepted_count":1,"duplicate_count":0}`))
	}))
	defer server.Close()

	app := NewApp(NewBootstrapRunner(BootstrapOptions{
		AppDir:            filepath.Join(t.TempDir(), ".cc-usage-client"),
		ClaudeProjectsDir: projectsDir,
		EnvLookup: func(key string) string {
			if key == "CC_USAGE_CLIENT_SERVER_URL" {
				return server.URL
			}
			if key == "CC_USAGE_CLIENT_TIMEOUT_SECONDS" {
				return "1"
			}
			return ""
		},
		Stdout: &bytes.Buffer{},
	}))

	if err := app.Run([]string{"sync"}); err != nil {
		t.Fatalf("Run(sync) returned error: %v", err)
	}

	if requestCount != 3 {
		t.Fatalf("expected 3 attempts, got %d", requestCount)
	}
}

func TestSyncTreatsBusinessCodeFailureAsWholeBatchFailure(t *testing.T) {
	t.Parallel()

	projectsDir := filepath.Join(t.TempDir(), "projects")
	projectDir := filepath.Join(projectsDir, "demo-project")
	if err := os.MkdirAll(projectDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() returned error: %v", err)
	}

	logPath := filepath.Join(projectDir, "session.jsonl")
	logContent := []byte("{\"type\":\"assistant\",\"message\":{\"id\":\"biz_fail_msg\",\"model\":\"claude-opus-4-1\",\"usage\":{\"input_tokens\":2,\"output_tokens\":1},\"stop_reason\":\"end_turn\"},\"timestamp\":\"2026-04-05T12:00:00Z\",\"sessionId\":\"session-biz\"}\n")
	if err := os.WriteFile(logPath, logContent, 0o644); err != nil {
		t.Fatalf("WriteFile() returned error: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":1,"message":"rejected","accepted_count":0,"duplicate_count":0}`))
	}))
	defer server.Close()

	appDir := filepath.Join(t.TempDir(), ".cc-usage-client")
	app := NewApp(NewBootstrapRunner(BootstrapOptions{
		AppDir:            appDir,
		ClaudeProjectsDir: projectsDir,
		EnvLookup: func(key string) string {
			if key == "CC_USAGE_CLIENT_SERVER_URL" {
				return server.URL
			}
			return ""
		},
		Stdout: &bytes.Buffer{},
	}))

	err := app.Run([]string{"sync"})
	if err == nil {
		t.Fatal("expected sync to fail on business code rejection")
	}

	db, openErr := sql.Open("sqlite", filepath.Join(appDir, "client.db"))
	if openErr != nil {
		t.Fatalf("sql.Open() returned error: %v", openErr)
	}
	defer db.Close()

	var reportedCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM reported_ids`).Scan(&reportedCount); err != nil {
		t.Fatalf("query reported_ids returned error: %v", err)
	}
	if reportedCount != 0 {
		t.Fatalf("expected whole batch failure to keep reported_ids empty, got %d", reportedCount)
	}
}

func TestSyncRerunsFailedFileByRescanningAndSendingOnlyRemainingRecords(t *testing.T) {
	t.Parallel()

	projectsDir := filepath.Join(t.TempDir(), "projects")
	projectDir := filepath.Join(projectsDir, "demo-project")
	if err := os.MkdirAll(projectDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() returned error: %v", err)
	}

	var content strings.Builder
	for index := 0; index < 501; index++ {
		content.WriteString("{\"type\":\"assistant\",\"message\":{\"id\":\"resume_")
		content.WriteString(strconv.Itoa(index))
		content.WriteString("\",\"model\":\"claude-opus-4-1\",\"usage\":{\"input_tokens\":1,\"output_tokens\":1},\"stop_reason\":\"end_turn\"},\"timestamp\":\"2026-04-05T12:00:00Z\",\"sessionId\":\"session-resume\"}\n")
	}
	logPath := filepath.Join(projectDir, "session.jsonl")
	if err := os.WriteFile(logPath, []byte(content.String()), 0o644); err != nil {
		t.Fatalf("WriteFile() returned error: %v", err)
	}

	appDir := filepath.Join(t.TempDir(), ".cc-usage-client")
	failingCount := 0
	failingServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		failingCount++
		w.Header().Set("Content-Type", "application/json")
		if failingCount == 1 {
			_, _ = w.Write([]byte(`{"code":0,"message":"success","accepted_count":500,"duplicate_count":0}`))
			return
		}
		http.Error(w, "temporary failure", http.StatusBadGateway)
	}))
	defer failingServer.Close()

	firstApp := NewApp(NewBootstrapRunner(BootstrapOptions{
		AppDir:            appDir,
		ClaudeProjectsDir: projectsDir,
		EnvLookup: func(key string) string {
			if key == "CC_USAGE_CLIENT_SERVER_URL" {
				return failingServer.URL
			}
			return ""
		},
		Stdout: &bytes.Buffer{},
	}))
	if err := firstApp.Run([]string{"sync"}); err == nil {
		t.Fatal("expected first sync to fail")
	}

	sentOnRetry := 0
	retryServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			Reports []json.RawMessage `json:"reports"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("Decode() returned error: %v", err)
		}
		sentOnRetry = len(payload.Reports)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":0,"message":"success","accepted_count":1,"duplicate_count":0}`))
	}))
	defer retryServer.Close()

	secondApp := NewApp(NewBootstrapRunner(BootstrapOptions{
		AppDir:            appDir,
		ClaudeProjectsDir: projectsDir,
		EnvLookup: func(key string) string {
			if key == "CC_USAGE_CLIENT_SERVER_URL" {
				return retryServer.URL
			}
			return ""
		},
		Stdout: &bytes.Buffer{},
	}))
	if err := secondApp.Run([]string{"sync"}); err != nil {
		t.Fatalf("second Run(sync) returned error: %v", err)
	}

	if sentOnRetry != 1 {
		t.Fatalf("expected strict rescan to skip already reported rows and send 1 remaining record, got %d", sentOnRetry)
	}

	db, openErr := sql.Open("sqlite", filepath.Join(appDir, "client.db"))
	if openErr != nil {
		t.Fatalf("sql.Open() returned error: %v", openErr)
	}
	defer db.Close()

	var lastLineOffset int
	if err := db.QueryRow(`SELECT last_line_offset FROM sync_state WHERE file_path = ?`, logPath).Scan(&lastLineOffset); err != nil {
		t.Fatalf("query sync_state returned error: %v", err)
	}
	if lastLineOffset != 501 {
		t.Fatalf("expected rerun to advance sync_state to 501, got %d", lastLineOffset)
	}
}

func TestSyncTreatsInsufficientAcknowledgementAsWholeBatchFailure(t *testing.T) {
	t.Parallel()

	projectsDir := filepath.Join(t.TempDir(), "projects")
	projectDir := filepath.Join(projectsDir, "demo-project")
	if err := os.MkdirAll(projectDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() returned error: %v", err)
	}

	logPath := filepath.Join(projectDir, "session.jsonl")
	logContent := []byte("{\"type\":\"assistant\",\"message\":{\"id\":\"ack_fail_msg\",\"model\":\"claude-opus-4-1\",\"usage\":{\"input_tokens\":2,\"output_tokens\":1},\"stop_reason\":\"end_turn\"},\"timestamp\":\"2026-04-05T12:00:00Z\",\"sessionId\":\"session-ack\"}\n")
	if err := os.WriteFile(logPath, logContent, 0o644); err != nil {
		t.Fatalf("WriteFile() returned error: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":0,"message":"partial","accepted_count":0,"duplicate_count":0}`))
	}))
	defer server.Close()

	appDir := filepath.Join(t.TempDir(), ".cc-usage-client")
	app := NewApp(NewBootstrapRunner(BootstrapOptions{
		AppDir:            appDir,
		ClaudeProjectsDir: projectsDir,
		EnvLookup: func(key string) string {
			if key == "CC_USAGE_CLIENT_SERVER_URL" {
				return server.URL
			}
			return ""
		},
		Stdout: &bytes.Buffer{},
	}))

	err := app.Run([]string{"sync"})
	if err == nil {
		t.Fatal("expected sync to fail on insufficient acknowledgement")
	}

	db, openErr := sql.Open("sqlite", filepath.Join(appDir, "client.db"))
	if openErr != nil {
		t.Fatalf("sql.Open() returned error: %v", openErr)
	}
	defer db.Close()

	var reportedCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM reported_ids`).Scan(&reportedCount); err != nil {
		t.Fatalf("query reported_ids returned error: %v", err)
	}
	if reportedCount != 0 {
		t.Fatalf("expected whole batch failure to keep reported_ids empty, got %d", reportedCount)
	}
}

func TestSyncRescansFromBeginningAfterFileTruncation(t *testing.T) {
	t.Parallel()

	projectsDir := filepath.Join(t.TempDir(), "projects")
	projectDir := filepath.Join(projectsDir, "demo-project")
	if err := os.MkdirAll(projectDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() returned error: %v", err)
	}

	logPath := filepath.Join(projectDir, "session.jsonl")
	initialContent := "" +
		"{\"type\":\"assistant\",\"message\":{\"id\":\"old_1\",\"model\":\"claude-opus-4-1\",\"usage\":{\"input_tokens\":1,\"output_tokens\":1},\"stop_reason\":\"end_turn\"},\"timestamp\":\"2026-04-05T12:00:00Z\",\"sessionId\":\"session-reset\"}\n" +
		"{\"type\":\"assistant\",\"message\":{\"id\":\"old_2\",\"model\":\"claude-opus-4-1\",\"usage\":{\"input_tokens\":1,\"output_tokens\":1},\"stop_reason\":\"end_turn\"},\"timestamp\":\"2026-04-05T12:00:01Z\",\"sessionId\":\"session-reset\"}\n"
	if err := os.WriteFile(logPath, []byte(initialContent), 0o644); err != nil {
		t.Fatalf("WriteFile() returned error: %v", err)
	}

	firstServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":0,"message":"success","accepted_count":2,"duplicate_count":0}`))
	}))
	defer firstServer.Close()

	appDir := filepath.Join(t.TempDir(), ".cc-usage-client")
	firstApp := NewApp(NewBootstrapRunner(BootstrapOptions{
		AppDir:            appDir,
		ClaudeProjectsDir: projectsDir,
		EnvLookup: func(key string) string {
			if key == "CC_USAGE_CLIENT_SERVER_URL" {
				return firstServer.URL
			}
			return ""
		},
		Stdout: &bytes.Buffer{},
	}))
	if err := firstApp.Run([]string{"sync"}); err != nil {
		t.Fatalf("first Run(sync) returned error: %v", err)
	}

	resetContent := "{\"type\":\"assistant\",\"message\":{\"id\":\"new_after_reset\",\"model\":\"claude-opus-4-1\",\"usage\":{\"input_tokens\":1,\"output_tokens\":1},\"stop_reason\":\"end_turn\"},\"timestamp\":\"2026-04-05T12:00:02Z\",\"sessionId\":\"session-reset\"}\n"
	if err := os.WriteFile(logPath, []byte(resetContent), 0o644); err != nil {
		t.Fatalf("WriteFile() returned error: %v", err)
	}

	retrySent := 0
	secondServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			Reports []struct {
				RequestID string `json:"request_id"`
			} `json:"reports"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("Decode() returned error: %v", err)
		}
		retrySent = len(payload.Reports)
		if retrySent != 1 || payload.Reports[0].RequestID != "session:new_after_reset" {
			t.Fatalf("unexpected reset payload: %#v", payload.Reports)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":0,"message":"success","accepted_count":1,"duplicate_count":0}`))
	}))
	defer secondServer.Close()

	secondApp := NewApp(NewBootstrapRunner(BootstrapOptions{
		AppDir:            appDir,
		ClaudeProjectsDir: projectsDir,
		EnvLookup: func(key string) string {
			if key == "CC_USAGE_CLIENT_SERVER_URL" {
				return secondServer.URL
			}
			return ""
		},
		Stdout: &bytes.Buffer{},
	}))
	if err := secondApp.Run([]string{"sync"}); err != nil {
		t.Fatalf("second Run(sync) returned error: %v", err)
	}

	if retrySent != 1 {
		t.Fatalf("expected truncated file rerun to send 1 new record, got %d", retrySent)
	}
}
