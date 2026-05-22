package dto

// StatsOverviewResponse 表示总览统计接口的响应体。
type StatsOverviewResponse struct {
	TotalTokens   int64             `json:"total_tokens"`
	TotalCostUSD  string            `json:"total_cost_usd"`
	TotalRequests int64             `json:"total_requests"`
	ActiveClients int64             `json:"active_clients"`
	TopModels     []StatsModelRank  `json:"top_models"`
	TopClients    []StatsClientRank `json:"top_clients"`
}

// StatsModelRank 表示模型排行项。
type StatsModelRank struct {
	Model  string `json:"model"`
	Tokens int64  `json:"tokens"`
}

// StatsClientRank 表示客户端排行项。
type StatsClientRank struct {
	ClientID     string `json:"client_id"`
	TotalCostUSD string `json:"total_cost_usd"`
}

// StatsTrendQuery 表示趋势统计接口的查询参数。
type StatsTrendQuery struct {
	Interval string `form:"interval" binding:"required"`
	StartAt  int64  `form:"start_at"`
	EndAt    int64  `form:"end_at"`
}

// StatsDashboardQuery 表示仪表盘统计接口的查询参数。
type StatsDashboardQuery struct {
	Interval string `form:"interval" binding:"required"`
	StartAt  int64  `form:"start_at" binding:"required"`
	EndAt    int64  `form:"end_at" binding:"required"`
}

// StatsTrendPoint 表示单个趋势桶的聚合结果。
type StatsTrendPoint struct {
	Bucket        string `json:"bucket"`
	TotalTokens   int64  `json:"total_tokens"`
	TotalRequests int64  `json:"total_requests"`
	TotalCostUSD  string `json:"total_cost_usd"`
}

// StatsDashboardResponse 表示仪表盘统计接口的响应体。
type StatsDashboardResponse struct {
	Overview         StatsDashboardOverview      `json:"overview"`
	PreviousOverview StatsDashboardOverview      `json:"previous_overview"`
	Trend            []StatsDashboardTrendPoint  `json:"trend"`
	TopModels        []StatsDashboardModelRank   `json:"top_models"`
	TopClients       []StatsClientRank           `json:"top_clients"`
	CacheAnalysis    StatsDashboardCacheAnalysis `json:"cache_analysis"`
}

// StatsDashboardOverview 表示仪表盘总览卡片数据。
type StatsDashboardOverview struct {
	TotalTokens      int64  `json:"total_tokens"`
	TotalCostUSD     string `json:"total_cost_usd"`
	TotalRequests    int64  `json:"total_requests"`
	ActiveClients    int64  `json:"active_clients"`
	TotalCacheTokens int64  `json:"total_cache_tokens"`
	CacheReadTokens  int64  `json:"cache_read_tokens"`
	InputTokens      int64  `json:"input_tokens"`
}

// StatsDashboardTrendPoint 表示仪表盘统一时间桶数据。
type StatsDashboardTrendPoint struct {
	Bucket              string `json:"bucket"`
	InputTokens         int64  `json:"input_tokens"`
	OutputTokens        int64  `json:"output_tokens"`
	CacheReadTokens     int64  `json:"cache_read_tokens"`
	CacheCreationTokens int64  `json:"cache_creation_tokens"`
	TotalRequests       int64  `json:"total_requests"`
	TotalCostUSD        string `json:"total_cost_usd"`
}

// StatsDashboardModelRank 表示仪表盘模型排行项。
type StatsDashboardModelRank struct {
	Model       string `json:"model"`
	DisplayName string `json:"display_name"`
	TotalTokens int64  `json:"total_tokens"`
}

// StatsDashboardCacheAnalysis 表示仪表盘缓存效益分析数据。
type StatsDashboardCacheAnalysis struct {
	SavedCostUSD         string `json:"saved_cost_usd"`
	CacheReadCostUSD     string `json:"cache_read_cost_usd"`
	CacheCreationCostUSD string `json:"cache_creation_cost_usd"`
}

// LogsQuery 表示原始日志查询接口的查询参数。
type LogsQuery struct {
	ClientID  string `form:"client_id"`
	Model     string `form:"model"`
	RequestID string `form:"request_id"`
	StartTime int64  `form:"start_time"`
	EndTime   int64  `form:"end_time"`
	Offset    int    `form:"offset"`
	Limit     int    `form:"limit"`
}

// LogsItem 表示单条原始日志记录。
type LogsItem struct {
	ID                   uint   `json:"id"`
	ClientID             string `json:"client_id"`
	RequestID            string `json:"request_id"`
	AppType              string `json:"app_type"`
	Model                string `json:"model"`
	InputTokens          int64  `json:"input_tokens"`
	OutputTokens         int64  `json:"output_tokens"`
	CacheReadTokens      int64  `json:"cache_read_tokens"`
	CacheCreationTokens  int64  `json:"cache_creation_tokens"`
	InputCostUSD         string `json:"input_cost_usd"`
	OutputCostUSD        string `json:"output_cost_usd"`
	CacheReadCostUSD     string `json:"cache_read_cost_usd"`
	CacheCreationCostUSD string `json:"cache_creation_cost_usd"`
	TotalCostUSD         string `json:"total_cost_usd"`
	SessionID            string `json:"session_id"`
	PricingSource        string `json:"pricing_source"`
	CreatedAt            int64  `json:"created_at"`
	DataSource           string `json:"data_source"`
}

// LogsResponse 表示原始日志查询结果。
type LogsResponse struct {
	Data   []LogsItem `json:"data"`
	Total  int64      `json:"total"`
	Offset int        `json:"offset"`
	Limit  int        `json:"limit"`
	Page   int        `json:"page"`
}
