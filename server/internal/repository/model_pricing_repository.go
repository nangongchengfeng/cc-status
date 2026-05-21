package repository

import (
	"context"

	"cc-status/server/internal/model/entity"

	"gorm.io/gorm"
)

// ModelPricingRepository 负责模型定价的查询。
type ModelPricingRepository struct{}

// NewModelPricingRepository 创建模型定价 repository。
func NewModelPricingRepository() *ModelPricingRepository {
	return &ModelPricingRepository{}
}

// FindMatch 按精确、最长前缀、默认价顺序查找匹配的定价规则。
func (repository *ModelPricingRepository) FindMatch(
	ctx context.Context,
	db *gorm.DB,
	model string,
) (entity.ModelPricing, string, error) {
	var pricing entity.ModelPricing

	err := db.WithContext(ctx).
		Where("model_id = ?", model).
		First(&pricing).Error
	if err == nil {
		return pricing, "exact", nil
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return entity.ModelPricing{}, "", err
	}

	err = db.WithContext(ctx).
		Where("is_placeholder = ?", false).
		Where("? LIKE model_id || '%'", model).
		Order("LENGTH(model_id) DESC").
		Order("id ASC").
		First(&pricing).Error
	if err == nil {
		return pricing, "prefix", nil
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return entity.ModelPricing{}, "", err
	}

	err = db.WithContext(ctx).
		Where("is_placeholder = ?", true).
		Order("id ASC").
		First(&pricing).Error
	if err != nil {
		return entity.ModelPricing{}, "", err
	}

	return pricing, "default", nil
}
