package dto

// SyncRequest 表示 client 发来的整批同步请求。
type SyncRequest struct {
	ClientID string       `json:"client_id" binding:"required"`
	Reports  []SyncReport `json:"reports" binding:"required"`
}

// SyncReport 表示单条上报记录的最小字段集。
type SyncReport struct {
	RequestID           string `json:"request_id" binding:"required"`
	AppType             string `json:"app_type" binding:"required"`
	Model               string `json:"model" binding:"required"`
	InputTokens         int64  `json:"input_tokens"`
	OutputTokens        int64  `json:"output_tokens"`
	CacheReadTokens     int64  `json:"cache_read_tokens"`
	CacheCreationTokens int64  `json:"cache_creation_tokens"`
	CreatedAt           int64  `json:"created_at" binding:"required"`
	SessionID           string `json:"session_id"`
	DataSource          string `json:"data_source" binding:"required"`
}

// SyncResult 表示同步成功后的计数结果。
type SyncResult struct {
	AcceptedCount  int `json:"accepted_count"`
	DuplicateCount int `json:"duplicate_count"`
}
