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
