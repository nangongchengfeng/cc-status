package repository

import (
	"context"
	"strings"

	"cc-status/server/internal/model/entity"

	"gorm.io/gorm"
)

// ModelPricingRepository 负责模型定价的查询。
type ModelPricingRepository struct{}

// NewModelPricingRepository 创建模型定价 repository。
func NewModelPricingRepository() *ModelPricingRepository {
	return &ModelPricingRepository{}
}

// List 返回全部模型定价，供管理接口直接展示。
func (repository *ModelPricingRepository) List(
	ctx context.Context,
	db *gorm.DB,
) ([]entity.ModelPricing, error) {
	var pricings []entity.ModelPricing
	if err := db.WithContext(ctx).
		Order("is_placeholder DESC").
		Order("id ASC").
		Find(&pricings).Error; err != nil {
		return nil, err
	}

	return pricings, nil
}

// Create 写入一条新的模型定价记录。
func (repository *ModelPricingRepository) Create(
	ctx context.Context,
	db *gorm.DB,
	pricing *entity.ModelPricing,
) error {
	return db.WithContext(ctx).Create(pricing).Error
}

// GetByID 按主键读取一条模型定价记录。
func (repository *ModelPricingRepository) GetByID(
	ctx context.Context,
	db *gorm.DB,
	id uint,
) (entity.ModelPricing, error) {
	var pricing entity.ModelPricing
	if err := db.WithContext(ctx).First(&pricing, id).Error; err != nil {
		return entity.ModelPricing{}, err
	}
	return pricing, nil
}

// Update 保存一条已存在的模型定价记录。
func (repository *ModelPricingRepository) Update(
	ctx context.Context,
	db *gorm.DB,
	pricing *entity.ModelPricing,
) error {
	return db.WithContext(ctx).Save(pricing).Error
}

// HasPlaceholder 判断除指定 ID 外是否还存在其他 placeholder 默认定价。
func (repository *ModelPricingRepository) HasPlaceholder(
	ctx context.Context,
	db *gorm.DB,
	excludeID uint,
) (bool, error) {
	var count int64
	query := db.WithContext(ctx).Model(&entity.ModelPricing{}).Where("is_placeholder = ?", true)
	if excludeID > 0 {
		query = query.Where("id <> ?", excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
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
		return normalizeModelPricing(pricing), "exact", nil
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
		return normalizeModelPricing(pricing), "prefix", nil
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

	return normalizeModelPricing(pricing), "default", nil
}

func normalizeModelPricing(pricing entity.ModelPricing) entity.ModelPricing {
	pricing.ModelID = strings.ToLower(strings.TrimSpace(pricing.ModelID))
	return pricing
}
