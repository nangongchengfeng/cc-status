package repository

import (
	"os"
	"path/filepath"

	"cc-status/server/internal/model/entity"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// OpenDatabase 打开 SQLite 数据库，并确保父目录存在。
func OpenDatabase(path string) (*gorm.DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}

	return gorm.Open(sqlite.Open(path), &gorm.Config{})
}

// InitializeSchema 创建首版持久化结构，并初始化默认模型定价。
func InitializeSchema(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&entity.UsageReport{},
		&entity.ModelPricing{},
	); err != nil {
		return err
	}

	return seedModelPricings(db)
}

func seedModelPricings(db *gorm.DB) error {
	seedRows := []entity.ModelPricing{
		{
			ModelID:                     "claude-opus-4-1",
			DisplayName:                 "Claude Opus 4.1",
			InputCostPerMillion:         "15",
			OutputCostPerMillion:        "75",
			CacheReadCostPerMillion:     "1.5",
			CacheCreationCostPerMillion: "18.75",
		},
		{
			ModelID:                     "claude-sonnet-4-0",
			DisplayName:                 "Claude Sonnet 4.0",
			InputCostPerMillion:         "3",
			OutputCostPerMillion:        "15",
			CacheReadCostPerMillion:     "0.3",
			CacheCreationCostPerMillion: "3.75",
		},
		{
			ModelID:                     "claude-3-7-sonnet",
			DisplayName:                 "Claude 3.7 Sonnet",
			InputCostPerMillion:         "3",
			OutputCostPerMillion:        "15",
			CacheReadCostPerMillion:     "0.3",
			CacheCreationCostPerMillion: "3.75",
		},
		{
			ModelID:                     "claude-3-5-sonnet",
			DisplayName:                 "Claude 3.5 Sonnet",
			InputCostPerMillion:         "3",
			OutputCostPerMillion:        "15",
			CacheReadCostPerMillion:     "0.3",
			CacheCreationCostPerMillion: "3.75",
		},
		{
			ModelID:                     "__default__",
			DisplayName:                 "Global Placeholder Pricing",
			InputCostPerMillion:         "0.657",
			OutputCostPerMillion:        "3.429",
			CacheReadCostPerMillion:     "0",
			CacheCreationCostPerMillion: "0",
			IsPlaceholder:               true,
		},
	}

	for _, row := range seedRows {
		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "model_id"}},
			DoNothing: true,
		}).Create(&row).Error; err != nil {
			return err
		}
	}

	return nil
}
