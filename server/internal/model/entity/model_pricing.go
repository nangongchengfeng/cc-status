package entity

import "time"

// ModelPricing 表示模型的每百万 token 定价。
type ModelPricing struct {
	ID                          uint      `gorm:"primaryKey" json:"id"`
	ModelID                     string    `gorm:"size:128;not null;uniqueIndex" json:"model_id"`
	DisplayName                 string    `gorm:"size:255" json:"display_name"`
	InputCostPerMillion         string    `gorm:"type:decimal(20,10);not null" json:"input_cost_per_million"`
	OutputCostPerMillion        string    `gorm:"type:decimal(20,10);not null" json:"output_cost_per_million"`
	CacheReadCostPerMillion     string    `gorm:"type:decimal(20,10);not null;default:'0'" json:"cache_read_cost_per_million"`
	CacheCreationCostPerMillion string    `gorm:"type:decimal(20,10);not null;default:'0'" json:"cache_creation_cost_per_million"`
	IsPlaceholder               bool      `gorm:"not null;default:false;index" json:"is_placeholder"`
	CreatedAt                   time.Time `gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt                   time.Time `gorm:"not null;autoUpdateTime" json:"updated_at"`
}

func (ModelPricing) TableName() string {
	return "model_pricing"
}
