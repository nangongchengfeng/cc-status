package syncer

import (
	"fmt"
	"io"

	"cc-status/client/internal/claude"
	"cc-status/client/internal/httpclient"
	"cc-status/client/internal/storage"
)

type Result struct {
	FilesScanned int
	Records      int
	Accepted     int
	Skipped      int
	Errors       []string
}

func RunHappyPath(
	store *storage.Store,
	clientID string,
	syncClient *httpclient.SyncClient,
	fileResults []claude.FileScanResult,
	batchSize int,
	progress io.Writer,
) (Result, error) {
	result := Result{
		FilesScanned: len(fileResults),
	}

	if batchSize <= 0 {
		batchSize = 500
	}

	reports := make([]httpclient.Report, 0)
	filesByRequestID := make(map[string]claude.FileScanResult)
	pendingByFile := make(map[string]int)
	for _, fileResult := range fileResults {
		if fileResult.Error != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", fileResult.FilePath, fileResult.Error))
			continue
		}
		for _, record := range fileResult.Records {
			requestID := "session:" + record.MessageID
			reported, err := store.HasReported(requestID)
			if err != nil {
				return result, err
			}
			if reported {
				result.Skipped++
				continue
			}

			reports = append(reports, httpclient.Report{
				RequestID:           requestID,
				AppType:             "claude",
				Model:               record.Model,
				InputTokens:         record.InputTokens,
				OutputTokens:        record.OutputTokens,
				CacheReadTokens:     record.CacheReadTokens,
				CacheCreationTokens: record.CacheCreationTokens,
				CreatedAt:           record.Timestamp,
				SessionID:           record.SessionID,
				DataSource:          "session_log",
			})
			filesByRequestID[requestID] = fileResult
			pendingByFile[fileResult.FilePath]++
		}
	}

	if len(reports) == 0 {
		if progress != nil {
			_, _ = fmt.Fprintf(progress, "[同步] 无新记录, 已跳过 %d 条已上报记录\n", result.Skipped)
		}
		return result, nil
	}

	result.Records = len(reports)

	if progress != nil {
		_, _ = fmt.Fprintf(progress, "[同步] %d 条待上报 (已跳过 %d 条已同步)\n", len(reports), result.Skipped)
	}

	totalBatches := (len(reports) + batchSize - 1) / batchSize
	for start := 0; start < len(reports); start += batchSize {
		end := start + batchSize
		if end > len(reports) {
			end = len(reports)
		}

		batch := reports[start:end]
		response, err := syncClient.Sync(clientID, batch)
		if err != nil {
			return result, err
		}

		if response.Code != 0 {
			return result, fmt.Errorf("sync rejected: %s", response.Message)
		}

		if response.AcceptedCount+response.DuplicateCount < len(batch) {
			return result, fmt.Errorf("sync acknowledged only %d of %d reports", response.AcceptedCount+response.DuplicateCount, len(batch))
		}

		for _, report := range batch {
			if err := store.MarkReported(report.RequestID); err != nil {
				return result, err
			}

			fileResult := filesByRequestID[report.RequestID]
			pendingByFile[fileResult.FilePath]--
			if pendingByFile[fileResult.FilePath] == 0 {
				if err := store.UpdateSyncState(fileResult.FilePath, fileResult.LastModifiedNanos, fileResult.LastLineOffset); err != nil {
					return result, err
				}
			}
		}

		result.Accepted += response.AcceptedCount
		result.Skipped += response.DuplicateCount

		if progress != nil {
			batchNum := start/batchSize + 1
			_, _ = fmt.Fprintf(progress, "[同步] 批次 %d/%d: 上报 %d 条 -> 成功 %d 条\n",
				batchNum, totalBatches, len(batch), response.AcceptedCount)
		}
	}

	return result, nil
}