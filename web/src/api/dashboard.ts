import { http } from '@/api/http';
import type { DashboardQuery, DashboardResponse } from '@/types/dashboard';

interface DashboardApiOverview {
  total_tokens?: number;
  total_cost_usd?: string;
  total_requests?: number;
  active_clients?: number;
  total_cache_tokens?: number;
  cache_read_tokens?: number;
  input_tokens?: number;
  output_tokens?: number;
}

interface DashboardApiTrendModelCost {
  model: string;
  display_name?: string;
  cost_usd: string;
}

interface DashboardApiTrendPoint {
  bucket: string;
  input_tokens?: number;
  output_tokens?: number;
  cache_read_tokens?: number;
  cache_creation_tokens?: number;
  total_requests?: number;
  total_cost_usd?: string;
  model_costs?: DashboardApiTrendModelCost[];
}

interface DashboardApiTopModel {
  model: string;
  display_name?: string;
  total_tokens?: number;
  total_cost_usd?: string;
}

interface DashboardApiTopClientModelCost {
  model: string;
  display_name?: string;
  cost_usd: string;
}

interface DashboardApiTopClient {
  client_id: string;
  total_cost_usd?: string;
  model_costs?: DashboardApiTopClientModelCost[];
}

interface DashboardApiCacheAnalysis {
  saved_cost_usd?: string;
  cache_read_cost_usd?: string;
  cache_creation_cost_usd?: string;
}

interface DashboardApiResponse {
  overview?: DashboardApiOverview;
  previous_overview?: DashboardApiOverview;
  trend?: DashboardApiTrendPoint[];
  top_models?: DashboardApiTopModel[];
  top_clients?: DashboardApiTopClient[];
  cache_analysis?: DashboardApiCacheAnalysis;
}

interface DashboardEnvelope {
  data: DashboardApiResponse;
}

export async function getDashboard(query: DashboardQuery): Promise<DashboardResponse> {
  const response = await http.get<DashboardEnvelope>('/stats/dashboard', {
    params: {
      start_at: query.startAt,
      end_at: query.endAt,
      interval: query.interval,
    },
  });

  const payload = response.data.data;

  return {
    overview: {
      totalTokens: payload.overview?.total_tokens ?? 0,
      totalCostUsd: payload.overview?.total_cost_usd ?? '0',
      totalRequests: payload.overview?.total_requests ?? 0,
      activeClients: payload.overview?.active_clients ?? 0,
      totalCacheTokens: payload.overview?.total_cache_tokens ?? 0,
      cacheReadTokens: payload.overview?.cache_read_tokens ?? 0,
      inputTokens: payload.overview?.input_tokens ?? 0,
      outputTokens: payload.overview?.output_tokens ?? 0,
    },
    previousOverview: {
      totalTokens: payload.previous_overview?.total_tokens ?? 0,
      totalCostUsd: payload.previous_overview?.total_cost_usd ?? '0',
      totalRequests: payload.previous_overview?.total_requests ?? 0,
      activeClients: payload.previous_overview?.active_clients ?? 0,
      totalCacheTokens: payload.previous_overview?.total_cache_tokens ?? 0,
      cacheReadTokens: payload.previous_overview?.cache_read_tokens ?? 0,
      inputTokens: payload.previous_overview?.input_tokens ?? 0,
      outputTokens: payload.previous_overview?.output_tokens ?? 0,
    },
    trend: (payload.trend ?? []).map((item) => ({
      bucket: item.bucket,
      inputTokens: item.input_tokens ?? 0,
      outputTokens: item.output_tokens ?? 0,
      cacheReadTokens: item.cache_read_tokens ?? 0,
      cacheCreationTokens: item.cache_creation_tokens ?? 0,
      totalRequests: item.total_requests ?? 0,
      totalCostUsd: item.total_cost_usd ?? '0',
      modelCosts: (item.model_costs ?? []).map((modelCost) => ({
        model: modelCost.model,
        displayName: modelCost.display_name ?? '',
        costUsd: modelCost.cost_usd ?? '0',
      })),
    })),
    topModels: (payload.top_models ?? []).map((item) => ({
      model: item.model,
      displayName: item.display_name ?? '',
      totalTokens: item.total_tokens ?? 0,
      totalCostUsd: item.total_cost_usd ?? '0.0000000000',
    })),
    topClients: (payload.top_clients ?? []).map((item) => ({
      clientId: item.client_id,
      totalCostUsd: item.total_cost_usd ?? '0',
      modelCosts: (item.model_costs ?? []).map((modelCost) => ({
        model: modelCost.model,
        displayName: modelCost.display_name ?? '',
        costUsd: modelCost.cost_usd ?? '0',
      })),
    })),
    cacheAnalysis: {
      savedCostUsd: payload.cache_analysis?.saved_cost_usd ?? '0',
      cacheReadCostUsd: payload.cache_analysis?.cache_read_cost_usd ?? '0',
      cacheCreationCostUsd: payload.cache_analysis?.cache_creation_cost_usd ?? '0',
    },
  };
}
