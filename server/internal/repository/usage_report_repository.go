package repository

import (
	"context"

	"cc-status/server/internal/model/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// InsertBatchResult 汇总批量写入后的接收与重复计数。
type InsertBatchResult struct {
	AcceptedCount  int
	DuplicateCount int
}

// UsageReportRepository 负责 usage_reports 的持久化写入。
type UsageReportRepository struct{}

// NewUsageReportRepository 创建使用记录 repository。
func NewUsageReportRepository() *UsageReportRepository {
	return &UsageReportRepository{}
}

// List 返回全部使用记录，供统计查询聚合使用。
func (repository *UsageReportRepository) List(
	ctx context.Context,
	db *gorm.DB,
) ([]entity.UsageReport, error) {
	var reports []entity.UsageReport
	if err := db.WithContext(ctx).Order("id ASC").Find(&reports).Error; err != nil {
		return nil, err
	}

	return reports, nil
}

// InsertBatch 在调用方提供的事务里逐条写入，并把唯一键冲突折算为重复数。
func (repository *UsageReportRepository) InsertBatch(
	_ context.Context,
	tx *gorm.DB,
	reports []entity.UsageReport,
) (InsertBatchResult, error) {
	result := InsertBatchResult{}

	for _, report := range reports {
		outcome := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "client_id"}, {Name: "request_id"}},
			DoNothing: true,
		}).Create(&report)
		if outcome.Error != nil {
			return InsertBatchResult{}, outcome.Error
		}
		if outcome.RowsAffected == 0 {
			result.DuplicateCount++
			continue
		}
		result.AcceptedCount++
	}

	return result, nil
}
