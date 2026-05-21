package syncer

import (
	"fmt"

	"cc-status/client/internal/claude"
	"cc-status/client/internal/httpclient"
	"cc-status/client/internal/storage"
)

type Result struct {
	FilesScanned int
	Accepted     int
	Duplicates   int
	Errors       []string
}

func RunHappyPath(
	store *storage.Store,
	clientID string,
	syncClient *httpclient.SyncClient,
	fileResults []claude.FileScanResult,
	batchSize int,
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
		return result, nil
	}

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
			// 只有该文件当前轮次所有“未上报记录”都确认成功后，才推进它的 sync_state。
			if pendingByFile[fileResult.FilePath] == 0 {
				if err := store.UpdateSyncState(fileResult.FilePath, fileResult.LastModifiedNanos, fileResult.LastLineOffset); err != nil {
					return result, err
				}
			}
		}

		result.Accepted += response.AcceptedCount
		result.Duplicates += response.DuplicateCount
	}

	return result, nil
}
