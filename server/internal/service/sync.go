package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"cc-status/server/internal/model/dto"
	"cc-status/server/internal/model/entity"
	"cc-status/server/internal/repository"

	"gorm.io/gorm"
)

var minCreatedAt = time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC).Unix()

// UsageReportWriter 定义同步服务所需的最小写入能力。
type UsageReportWriter interface {
	InsertBatch(context.Context, *gorm.DB, []entity.UsageReport) (repository.InsertBatchResult, error)
}

// ValidationError 用于把请求校验错误映射成 400。
type ValidationError struct {
	message string
}

func (validationError ValidationError) Error() string {
	return validationError.message
}

// SyncService 承载批量同步的校验、规范化和事务控制。
type SyncService struct {
	db     *gorm.DB
	writer UsageReportWriter
	now    func() time.Time
}

// NewSyncService 创建默认同步服务。
func NewSyncService(db *gorm.DB) *SyncService {
	return NewSyncServiceWithWriter(db, repository.NewUsageReportRepository())
}

// NewSyncServiceWithWriter 为测试或特殊场景注入自定义 writer。
func NewSyncServiceWithWriter(db *gorm.DB, writer UsageReportWriter) *SyncService {
	return &SyncService{
		db:     db,
		writer: writer,
		now:    time.Now,
	}
}

// Ingest 执行单批次同步，并确保整个批次在单事务内处理。
func (service *SyncService) Ingest(ctx context.Context, request dto.SyncRequest) (dto.SyncResult, error) {
	reports, err := service.normalizeReports(request)
	if err != nil {
		return dto.SyncResult{}, err
	}

	var result dto.SyncResult
	err = service.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		insertResult, insertErr := service.writer.InsertBatch(ctx, tx, reports)
		if insertErr != nil {
			return insertErr
		}
		result.AcceptedCount = insertResult.AcceptedCount
		result.DuplicateCount = insertResult.DuplicateCount
		return nil
	})
	if err != nil {
		return dto.SyncResult{}, fmt.Errorf("insert sync batch: %w", err)
	}

	return result, nil
}

func (service *SyncService) normalizeReports(request dto.SyncRequest) ([]entity.UsageReport, error) {
	clientID := strings.TrimSpace(request.ClientID)
	if clientID == "" {
		return nil, ValidationError{message: "client_id 不能为空"}
	}
	if len(clientID) > 64 {
		return nil, ValidationError{message: "client_id 长度不能超过 64"}
	}
	if len(request.Reports) == 0 {
		return nil, ValidationError{message: "reports 至少需要一条记录"}
	}

	reports := make([]entity.UsageReport, 0, len(request.Reports))
	maxCreatedAt := service.now().Add(24 * time.Hour).Unix()

	for _, report := range request.Reports {
		normalized, err := service.normalizeReport(clientID, report, maxCreatedAt)
		if err != nil {
			return nil, err
		}
		reports = append(reports, normalized)
	}

	return reports, nil
}

func (service *SyncService) normalizeReport(
	clientID string,
	report dto.SyncReport,
	maxCreatedAt int64,
) (entity.UsageReport, error) {
	requestID := strings.TrimSpace(report.RequestID)
	if requestID == "" {
		return entity.UsageReport{}, ValidationError{message: "request_id 不能为空"}
	}
	if len(requestID) > 255 {
		return entity.UsageReport{}, ValidationError{message: "request_id 长度不能超过 255"}
	}

	appType := strings.ToLower(strings.TrimSpace(report.AppType))
	if appType != "claude" {
		return entity.UsageReport{}, ValidationError{message: "app_type 仅支持 claude"}
	}

	dataSource := strings.ToLower(strings.TrimSpace(report.DataSource))
	if dataSource != "session_log" {
		return entity.UsageReport{}, ValidationError{message: "data_source 仅支持 session_log"}
	}

	model := strings.ToLower(strings.TrimSpace(report.Model))
	if model == "" {
		return entity.UsageReport{}, ValidationError{message: "model 不能为空"}
	}
	if len(model) > 128 {
		return entity.UsageReport{}, ValidationError{message: "model 长度不能超过 128"}
	}

	sessionID := strings.TrimSpace(report.SessionID)
	if len(sessionID) > 255 {
		return entity.UsageReport{}, ValidationError{message: "session_id 长度不能超过 255"}
	}

	if report.CreatedAt < minCreatedAt || report.CreatedAt > maxCreatedAt {
		return entity.UsageReport{}, ValidationError{message: "created_at 超出允许范围"}
	}

	for _, tokenCount := range []int64{
		report.InputTokens,
		report.OutputTokens,
		report.CacheReadTokens,
		report.CacheCreationTokens,
	} {
		if tokenCount < 0 {
			return entity.UsageReport{}, ValidationError{message: "token 数不能为负数"}
		}
	}

	return entity.UsageReport{
		ClientID:             clientID,
		RequestID:            requestID,
		AppType:              appType,
		Model:                model,
		InputTokens:          report.InputTokens,
		OutputTokens:         report.OutputTokens,
		CacheReadTokens:      report.CacheReadTokens,
		CacheCreationTokens:  report.CacheCreationTokens,
		InputCostUSD:         "0",
		OutputCostUSD:        "0",
		CacheReadCostUSD:     "0",
		CacheCreationCostUSD: "0",
		TotalCostUSD:         "0",
		SessionID:            sessionID,
		PricingSource:        "placeholder",
		CreatedAtUnix:        report.CreatedAt,
		DataSource:           dataSource,
	}, nil
}

// IsValidationError 判断错误是否属于请求参数错误。
func IsValidationError(err error) bool {
	var validationError ValidationError
	return errors.As(err, &validationError)
}
