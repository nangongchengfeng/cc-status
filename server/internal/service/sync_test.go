package service

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"cc-status/server/internal/model/dto"
	"cc-status/server/internal/model/entity"
	"cc-status/server/internal/repository"

	"gorm.io/gorm"
)

func TestSyncServiceRollsBackBatchWhenInsertFails(t *testing.T) {
	t.Parallel()

	dbPath := filepath.Join(t.TempDir(), "server.db")
	db, err := repository.OpenDatabase(dbPath)
	if err != nil {
		t.Fatalf("OpenDatabase() returned error: %v", err)
	}
	if err := repository.InitializeSchema(db); err != nil {
		t.Fatalf("InitializeSchema() returned error: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("DB() returned error: %v", err)
	}
	t.Cleanup(func() {
		_ = sqlDB.Close()
	})

	syncService := NewSyncServiceWithWriter(db, rollbackWriter{})
	_, err = syncService.Ingest(context.Background(), dto.SyncRequest{
		ClientID: "client-1",
		Reports: []dto.SyncReport{
			{
				RequestID:   "session:msg-1",
				AppType:     "claude",
				Model:       "claude-sonnet-4-0",
				CreatedAt:   1743840000,
				SessionID:   "session-1",
				DataSource:  "session_log",
				InputTokens: 10,
			},
			{
				RequestID:  "session:msg-2",
				AppType:    "claude",
				Model:      "claude-sonnet-4-0",
				CreatedAt:  1743840001,
				SessionID:  "session-1",
				DataSource: "session_log",
			},
		},
	})
	if err == nil {
		t.Fatal("expected insert failure to bubble up")
	}

	var storedCount int64
	if err := db.Table("usage_reports").Count(&storedCount).Error; err != nil {
		t.Fatalf("Count() returned error: %v", err)
	}
	if storedCount != 0 {
		t.Fatalf("expected transaction rollback to leave zero rows, got %d", storedCount)
	}
}

type rollbackWriter struct{}

func (rollbackWriter) InsertBatch(
	_ context.Context,
	tx *gorm.DB,
	reports []entity.UsageReport,
) (repository.InsertBatchResult, error) {
	if len(reports) == 0 {
		return repository.InsertBatchResult{}, nil
	}

	if err := tx.Create(&reports[0]).Error; err != nil {
		return repository.InsertBatchResult{}, err
	}

	return repository.InsertBatchResult{}, errors.New("forced insert failure")
}
