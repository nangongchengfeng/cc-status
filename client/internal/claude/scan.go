package claude

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type TokenUsage struct {
	MessageID           string
	Model               string
	InputTokens         uint32
	OutputTokens        uint32
	CacheReadTokens     uint32
	CacheCreationTokens uint32
	Timestamp           int64
	SessionID           string
}

type ScanResult struct {
	FilesScanned int
	Records      []TokenUsage
	Errors       []string
}

type FileScanResult struct {
	FilePath          string
	LastLineOffset    int
	LastModifiedNanos int64
	Records           []TokenUsage
	Error             error
}

type rawEnvelope struct {
	Type      string `json:"type"`
	Timestamp string `json:"timestamp"`
	SessionID string `json:"sessionId"`
	Message   struct {
		ID         string `json:"id"`
		Model      string `json:"model"`
		StopReason string `json:"stop_reason"`
		Usage      struct {
			InputTokens         uint32 `json:"input_tokens"`
			OutputTokens        uint32 `json:"output_tokens"`
			CacheReadTokens     uint32 `json:"cache_read_input_tokens"`
			CacheCreationTokens uint32 `json:"cache_creation_input_tokens"`
		} `json:"usage"`
	} `json:"message"`
}

const maxJSONLLineSize = 4 * 1024 * 1024

func ScanProjectsDir(projectsDir string) (ScanResult, error) {
	fileResults, err := ScanProjectFiles(projectsDir)
	if err != nil {
		return ScanResult{}, err
	}

	result := ScanResult{
		FilesScanned: len(fileResults),
	}
	for _, fileResult := range fileResults {
		if fileResult.Error != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", fileResult.FilePath, fileResult.Error))
		}
		result.Records = append(result.Records, fileResult.Records...)
	}

	return result, nil
}

func ScanProjectFiles(projectsDir string) ([]FileScanResult, error) {
	files, err := collectJSONLFiles(projectsDir)
	if err != nil {
		return nil, err
	}

	results := make([]FileScanResult, 0, len(files))
	for _, filePath := range files {
		fileResult, scanErr := scanFile(filePath)
		fileResult.FilePath = filePath
		fileResult.Error = scanErr
		results = append(results, fileResult)
	}

	return results, nil
}

func collectJSONLFiles(projectsDir string) ([]string, error) {
	var files []string

	err := filepath.WalkDir(projectsDir, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		if strings.EqualFold(filepath.Ext(path), ".jsonl") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk projects dir: %w", err)
	}

	sort.Strings(files)
	return files, nil
}

func scanFile(filePath string) (FileScanResult, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return FileScanResult{}, fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	bestByMessageID := make(map[string]TokenUsage)
	seenStopReason := make(map[string]bool)
	lineCount := 0

	scanner := bufio.NewScanner(file)
	// Claude 会话里可能出现超长单行 JSONL，放大 scanner 缓冲区避免 "token too long"。
	scanner.Buffer(make([]byte, 0, 64*1024), maxJSONLLineSize)
	for scanner.Scan() {
		lineCount++
		line := scanner.Bytes()

		usage, hasStopReason, keep, parseErr := parseLine(line)
		if parseErr != nil {
			// 尾部半写入时保留当前文件之前已经解析出的有效记录，并把错误交给上层汇总。
			return buildFileScanResult(filePath, lineCount, mapValues(bestByMessageID)), parseErr
		}
		if !keep {
			continue
		}

		existing, exists := bestByMessageID[usage.MessageID]
		// 单文件内按 Rust 逻辑去重：优先 stop_reason 完整的消息，其次取输出 token 更大的版本。
		if !exists ||
			(hasStopReason && !seenStopReason[usage.MessageID]) ||
			(hasStopReason == seenStopReason[usage.MessageID] && usage.OutputTokens > existing.OutputTokens) {
			bestByMessageID[usage.MessageID] = usage
			seenStopReason[usage.MessageID] = hasStopReason
		}
	}

	if err := scanner.Err(); err != nil {
		return buildFileScanResult(filePath, lineCount, mapValues(bestByMessageID)), fmt.Errorf("scan file: %w", err)
	}

	return buildFileScanResult(filePath, lineCount, mapValues(bestByMessageID)), nil
}

func parseLine(line []byte) (TokenUsage, bool, bool, error) {
	var envelope rawEnvelope
	if err := json.Unmarshal(line, &envelope); err != nil {
		return TokenUsage{}, false, false, fmt.Errorf("parse json line: %w", err)
	}

	if envelope.Type != "assistant" {
		return TokenUsage{}, false, false, nil
	}

	if envelope.Message.StopReason == "" || envelope.Message.Usage.OutputTokens == 0 {
		return TokenUsage{}, false, false, nil
	}

	createdAt, err := parseTimestamp(envelope.Timestamp)
	if err != nil {
		return TokenUsage{}, false, false, err
	}

	return TokenUsage{
		MessageID:           envelope.Message.ID,
		Model:               envelope.Message.Model,
		InputTokens:         envelope.Message.Usage.InputTokens,
		OutputTokens:        envelope.Message.Usage.OutputTokens,
		CacheReadTokens:     envelope.Message.Usage.CacheReadTokens,
		CacheCreationTokens: envelope.Message.Usage.CacheCreationTokens,
		Timestamp:           createdAt,
		SessionID:           envelope.SessionID,
	}, true, true, nil
}

func parseTimestamp(raw string) (int64, error) {
	parsed, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return 0, fmt.Errorf("parse timestamp: %w", err)
	}
	return parsed.Unix(), nil
}

func mapValues(bestByMessageID map[string]TokenUsage) []TokenUsage {
	values := make([]TokenUsage, 0, len(bestByMessageID))
	for _, usage := range bestByMessageID {
		values = append(values, usage)
	}
	sort.Slice(values, func(i, j int) bool {
		return values[i].MessageID < values[j].MessageID
	})
	return values
}

func buildFileScanResult(filePath string, lineCount int, records []TokenUsage) FileScanResult {
	lastModifiedNanos := int64(0)
	if info, err := os.Stat(filePath); err == nil {
		lastModifiedNanos = info.ModTime().UnixNano()
	}

	return FileScanResult{
		FilePath:          filePath,
		LastLineOffset:    lineCount,
		LastModifiedNanos: lastModifiedNanos,
		Records:           records,
	}
}
