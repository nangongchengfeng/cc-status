export interface RecentLogQuery {
  startAt: number;
  endAt: number;
  limit: number;
}

export interface RecentLogItem {
  id: number;
  clientId: string;
  requestId: string;
  model: string;
  inputTokens: number;
  outputTokens: number;
  cacheReadTokens: number;
  cacheCreationTokens: number;
  totalCostUsd: string;
  createdAt: number;
}

export interface RecentLogsResponse {
  data: RecentLogItem[];
  total: number;
  offset: number;
  limit: number;
  page: number;
}
