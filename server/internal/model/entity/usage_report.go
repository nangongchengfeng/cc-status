package entity

import "time"

// UsageReport 表示一条已入库的使用记录。
type UsageReport struct {
	ID                   uint      `gorm:"primaryKey" json:"id"`
	ClientID             string    `gorm:"size:64;not null;index:idx_usage_reports_client_request,unique;index:idx_usage_reports_client_id" json:"client_id"`
	RequestID            string    `gorm:"size:255;not null;index:idx_usage_reports_client_request,unique" json:"request_id"`
	AppType              string    `gorm:"size:32;not null" json:"app_type"`
	Model                string    `gorm:"size:128;not null;index:idx_usage_reports_model" json:"model"`
	InputTokens          int64     `gorm:"not null;default:0" json:"input_tokens"`
	OutputTokens         int64     `gorm:"not null;default:0" json:"output_tokens"`
	CacheReadTokens      int64     `gorm:"not null;default:0" json:"cache_read_tokens"`
	CacheCreationTokens  int64     `gorm:"not null;default:0" json:"cache_creation_tokens"`
	InputCostUSD         string    `gorm:"type:decimal(20,10);not null;default:'0'" json:"input_cost_usd"`
	OutputCostUSD        string    `gorm:"type:decimal(20,10);not null;default:'0'" json:"output_cost_usd"`
	CacheReadCostUSD     string    `gorm:"type:decimal(20,10);not null;default:'0'" json:"cache_read_cost_usd"`
	CacheCreationCostUSD string    `gorm:"type:decimal(20,10);not null;default:'0'" json:"cache_creation_cost_usd"`
	TotalCostUSD         string    `gorm:"type:decimal(20,10);not null;default:'0'" json:"total_cost_usd"`
	SessionID            string    `gorm:"size:255" json:"session_id"`
	PricingSource        string    `gorm:"size:32;not null" json:"pricing_source"`
	CreatedAtUnix        int64     `gorm:"column:created_at;not null;index:idx_usage_reports_created_at" json:"created_at"`
	DataSource           string    `gorm:"size:32;not null" json:"data_source"`
	InsertedAt           time.Time `gorm:"not null;autoCreateTime" json:"inserted_at"`
}

func (UsageReport) TableName() string {
	return "usage_reports"
}
