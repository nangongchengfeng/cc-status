package claude

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanProjectsDirAppliesClaudeFilteringAndDedup(t *testing.T) {
	t.Parallel()

	projectsDir := filepath.Join(t.TempDir(), "projects")
	projectDir := filepath.Join(projectsDir, "demo-project")
	if err := os.MkdirAll(projectDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() returned error: %v", err)
	}

	logPath := filepath.Join(projectDir, "session.jsonl")
	logContent := "" +
		"{\"type\":\"user\",\"message\":{\"id\":\"ignored-user\",\"model\":\"claude-sonnet-4\",\"usage\":{\"input_tokens\":1,\"output_tokens\":1},\"stop_reason\":\"end_turn\"},\"timestamp\":\"2026-04-05T12:00:00Z\",\"sessionId\":\"session-1\"}\n" +
		"{\"type\":\"assistant\",\"message\":{\"id\":\"ignored-no-stop\",\"model\":\"claude-sonnet-4\",\"usage\":{\"input_tokens\":2,\"output_tokens\":5},\"stop_reason\":\"\"},\"timestamp\":\"2026-04-05T12:00:01Z\",\"sessionId\":\"session-1\"}\n" +
		"{\"type\":\"assistant\",\"message\":{\"id\":\"ignored-no-output\",\"model\":\"claude-sonnet-4\",\"usage\":{\"input_tokens\":2,\"output_tokens\":0},\"stop_reason\":\"end_turn\"},\"timestamp\":\"2026-04-05T12:00:02Z\",\"sessionId\":\"session-1\"}\n" +
		"{\"type\":\"assistant\",\"message\":{\"id\":\"dup-msg\",\"model\":\"claude-sonnet-4\",\"usage\":{\"input_tokens\":3,\"output_tokens\":4,\"cache_read_input_tokens\":1,\"cache_creation_input_tokens\":0},\"stop_reason\":\"end_turn\"},\"timestamp\":\"2026-04-05T12:00:03Z\",\"sessionId\":\"session-1\"}\n" +
		"{\"type\":\"assistant\",\"message\":{\"id\":\"dup-msg\",\"model\":\"claude-sonnet-4\",\"usage\":{\"input_tokens\":3,\"output_tokens\":9,\"cache_read_input_tokens\":2,\"cache_creation_input_tokens\":1},\"stop_reason\":\"end_turn\"},\"timestamp\":\"2026-04-05T12:00:04Z\",\"sessionId\":\"session-1\"}\n"
	if err := os.WriteFile(logPath, []byte(logContent), 0o644); err != nil {
		t.Fatalf("WriteFile() returned error: %v", err)
	}

	result, err := ScanProjectsDir(projectsDir)
	if err != nil {
		t.Fatalf("ScanProjectsDir() returned error: %v", err)
	}

	if len(result.Records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(result.Records))
	}

	record := result.Records[0]
	if record.MessageID != "dup-msg" {
		t.Fatalf("unexpected message id: %q", record.MessageID)
	}

	if record.OutputTokens != 9 {
		t.Fatalf("expected deduped output_tokens=9, got %d", record.OutputTokens)
	}

	if record.CacheReadTokens != 2 || record.CacheCreationTokens != 1 {
		t.Fatalf("unexpected cache tokens: %#v", record)
	}
}

func TestScanProjectsDirReportsBrokenTailLineWithoutDroppingPriorRecords(t *testing.T) {
	t.Parallel()

	projectsDir := filepath.Join(t.TempDir(), "projects")
	projectDir := filepath.Join(projectsDir, "demo-project")
	if err := os.MkdirAll(projectDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() returned error: %v", err)
	}

	logPath := filepath.Join(projectDir, "session.jsonl")
	logContent := "" +
		"{\"type\":\"assistant\",\"message\":{\"id\":\"msg-1\",\"model\":\"claude-sonnet-4\",\"usage\":{\"input_tokens\":5,\"output_tokens\":2},\"stop_reason\":\"end_turn\"},\"timestamp\":\"2026-04-05T12:00:00Z\",\"sessionId\":\"session-1\"}\n" +
		"{\"type\":\"assistant\""
	if err := os.WriteFile(logPath, []byte(logContent), 0o644); err != nil {
		t.Fatalf("WriteFile() returned error: %v", err)
	}

	result, err := ScanProjectsDir(projectsDir)
	if err != nil {
		t.Fatalf("ScanProjectsDir() returned error: %v", err)
	}

	if len(result.Records) != 1 {
		t.Fatalf("expected prior valid record to be preserved, got %d", len(result.Records))
	}

	if len(result.Errors) != 1 {
		t.Fatalf("expected 1 file error, got %d", len(result.Errors))
	}
}

func TestScanProjectsDirIncludesSubagentLogs(t *testing.T) {
	t.Parallel()

	projectsDir := filepath.Join(t.TempDir(), "projects")
	subagentDir := filepath.Join(projectsDir, "demo-project", "session-1", "subagents")
	if err := os.MkdirAll(subagentDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() returned error: %v", err)
	}

	logPath := filepath.Join(subagentDir, "agent.jsonl")
	logContent := []byte("{\"type\":\"assistant\",\"message\":{\"id\":\"subagent-msg\",\"model\":\"claude-sonnet-4\",\"usage\":{\"input_tokens\":7,\"output_tokens\":3},\"stop_reason\":\"end_turn\"},\"timestamp\":\"2026-04-05T12:00:00Z\",\"sessionId\":\"session-1\"}\n")
	if err := os.WriteFile(logPath, logContent, 0o644); err != nil {
		t.Fatalf("WriteFile() returned error: %v", err)
	}

	result, err := ScanProjectsDir(projectsDir)
	if err != nil {
		t.Fatalf("ScanProjectsDir() returned error: %v", err)
	}

	if result.FilesScanned != 1 {
		t.Fatalf("expected 1 scanned file, got %d", result.FilesScanned)
	}

	if len(result.Records) != 1 || result.Records[0].MessageID != "subagent-msg" {
		t.Fatalf("expected subagent record to be included, got %#v", result.Records)
	}
}
