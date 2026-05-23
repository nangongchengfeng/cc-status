export type DashboardInterval = 'hour' | 'day';

export interface DashboardQuery {
  startAt: number;
  endAt: number;
  interval: DashboardInterval;
}

export interface DashboardOverview {
  totalTokens: number;
  totalCostUsd: string;
  totalRequests: number;
  activeClients: number;
  totalCacheTokens: number;
  cacheReadTokens: number;
  inputTokens: number;
  outputTokens: number;
}

export interface DashboardTrendModelCost {
  model: string;
  displayName: string;
  costUsd: string;
}

export interface DashboardTrendPoint {
  bucket: string;
  inputTokens: number;
  outputTokens: number;
  cacheReadTokens: number;
  cacheCreationTokens: number;
  totalRequests: number;
  totalCostUsd: string;
  modelCosts?: DashboardTrendModelCost[];
}

export interface DashboardTopModel {
  model: string;
  displayName: string;
  totalTokens: number;
  totalCostUsd: string;
}

export interface DashboardTopClientModelCost {
  model: string;
  displayName: string;
  costUsd: string;
}

export interface DashboardTopClient {
  clientId: string;
  totalCostUsd: string;
  modelCosts?: DashboardTopClientModelCost[];
}

export interface DashboardCacheAnalysis {
  savedCostUsd: string;
  cacheReadCostUsd: string;
  cacheCreationCostUsd: string;
}

export interface DashboardResponse {
  overview: DashboardOverview;
  previousOverview: DashboardOverview;
  trend: DashboardTrendPoint[];
  topModels: DashboardTopModel[];
  topClients: DashboardTopClient[];
  cacheAnalysis: DashboardCacheAnalysis;
}
