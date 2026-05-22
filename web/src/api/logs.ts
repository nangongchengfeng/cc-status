import { http } from '@/api/http';
import type { RecentLogQuery, RecentLogsResponse } from '@/types/logs';

interface RecentLogsApiItem {
  id: number;
  client_id: string;
  request_id: string;
  model: string;
  input_tokens: number;
  output_tokens: number;
  cache_read_tokens: number;
  cache_creation_tokens: number;
  total_cost_usd: string;
  created_at?: number;
}

interface RecentLogsApiResponse {
  data: RecentLogsApiItem[];
  total: number;
  offset: number;
  limit: number;
  page: number;
}

export async function getRecentLogs(query: RecentLogQuery): Promise<RecentLogsResponse> {
  const response = await http.get<RecentLogsApiResponse>('/logs', {
    params: {
      start_time: query.startAt,
      end_time: query.endAt,
      limit: query.limit,
      offset: 0,
    },
  });

  return {
    ...response.data,
    data: response.data.data.map((item) => ({
      id: item.id,
      clientId: item.client_id,
      requestId: item.request_id,
      model: item.model,
      inputTokens: item.input_tokens,
      outputTokens: item.output_tokens,
      cacheReadTokens: item.cache_read_tokens,
      cacheCreationTokens: item.cache_creation_tokens,
      totalCostUsd: item.total_cost_usd,
      createdAt: item.created_at ?? 0,
    })),
  };
}
