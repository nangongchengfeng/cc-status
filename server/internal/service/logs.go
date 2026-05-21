package service

import (
	"context"

	"cc-status/server/internal/model/dto"
	"cc-status/server/internal/model/entity"
	"cc-status/server/internal/repository"

	"gorm.io/gorm"
)

// LogsReader 定义日志查询所需的最小读取能力。
type LogsReader interface {
	QueryLogs(context.Context, *gorm.DB, string, string, string, int64, int64, int, int) ([]entity.UsageReport, int64, error)
}

// LogsService 承载原始日志查询能力。
type LogsService struct {
	db     *gorm.DB
	reader LogsReader
}

// NewLogsService 创建日志查询服务。
func NewLogsService(db *gorm.DB) *LogsService {
	return &LogsService{
		db:     db,
		reader: repository.NewUsageReportRepository(),
	}
}

// List 按过滤条件和分页规则返回原始日志记录。
func (service *LogsService) List(ctx context.Context, query dto.LogsQuery) (dto.LogsResponse, error) {
	normalizedQuery := normalizeLogsQuery(query)
	reports, total, err := service.reader.QueryLogs(
		ctx,
		service.db,
		normalizedQuery.ClientID,
		normalizedQuery.Model,
		normalizedQuery.RequestID,
		normalizedQuery.StartTime,
		normalizedQuery.EndTime,
		normalizedQuery.Offset,
		normalizedQuery.Limit,
	)
	if err != nil {
		return dto.LogsResponse{}, err
	}

	items := make([]dto.LogsItem, 0, len(reports))
	for _, report := range reports {
		items = append(items, dto.LogsItem{
			ID:                   report.ID,
			ClientID:             report.ClientID,
			RequestID:            report.RequestID,
			AppType:              report.AppType,
			Model:                report.Model,
			InputTokens:          report.InputTokens,
			OutputTokens:         report.OutputTokens,
			CacheReadTokens:      report.CacheReadTokens,
			CacheCreationTokens:  report.CacheCreationTokens,
			InputCostUSD:         report.InputCostUSD,
			OutputCostUSD:        report.OutputCostUSD,
			CacheReadCostUSD:     report.CacheReadCostUSD,
			CacheCreationCostUSD: report.CacheCreationCostUSD,
			TotalCostUSD:         report.TotalCostUSD,
			SessionID:            report.SessionID,
			PricingSource:        report.PricingSource,
			CreatedAt:            report.CreatedAtUnix,
			DataSource:           report.DataSource,
		})
	}

	return dto.LogsResponse{
		Data:   items,
		Total:  total,
		Offset: normalizedQuery.Offset,
		Limit:  normalizedQuery.Limit,
		Page:   normalizedQuery.Offset/normalizedQuery.Limit + 1,
	}, nil
}

func normalizeLogsQuery(query dto.LogsQuery) dto.LogsQuery {
	if query.Offset < 0 {
		query.Offset = 0
	}
	if query.Limit <= 0 {
		query.Limit = 20
	}
	if query.Limit > 100 {
		query.Limit = 100
	}
	return query
}
