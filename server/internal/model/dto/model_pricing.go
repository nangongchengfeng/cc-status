package dto

import "time"

// ModelPricingResponse 表示对外暴露的模型定价信息。
type ModelPricingResponse struct {
	ID                          uint      `json:"id"`
	ModelID                     string    `json:"model_id"`
	DisplayName                 string    `json:"display_name"`
	InputCostPerMillion         string    `json:"input_cost_per_million"`
	OutputCostPerMillion        string    `json:"output_cost_per_million"`
	CacheReadCostPerMillion     string    `json:"cache_read_cost_per_million"`
	CacheCreationCostPerMillion string    `json:"cache_creation_cost_per_million"`
	IsPlaceholder               bool      `json:"is_placeholder"`
	CreatedAt                   time.Time `json:"created_at"`
	UpdatedAt                   time.Time `json:"updated_at"`
}

// ModelPricingUpsertRequest 表示创建或全量更新模型定价的请求体。
type ModelPricingUpsertRequest struct {
	ModelID                     string `json:"model_id" binding:"required"`
	DisplayName                 string `json:"display_name"`
	InputCostPerMillion         string `json:"input_cost_per_million" binding:"required"`
	OutputCostPerMillion        string `json:"output_cost_per_million" binding:"required"`
	CacheReadCostPerMillion     string `json:"cache_read_cost_per_million" binding:"required"`
	CacheCreationCostPerMillion string `json:"cache_creation_cost_per_million" binding:"required"`
	IsPlaceholder               *bool  `json:"is_placeholder" binding:"required"`
}

// ModelPricingURI 表示模型定价资源的路径参数。
type ModelPricingURI struct {
	ID uint `uri:"id" binding:"required"`
}
