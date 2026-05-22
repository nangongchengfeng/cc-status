import { describe, expect, it, vi } from 'vitest';

import { getDashboard } from '@/api/dashboard';

const httpGetSpy = vi.fn();

vi.mock('@/api/http', () => ({
  http: {
    get: (...args: unknown[]) => httpGetSpy(...args),
  },
}));

describe('getDashboard', () => {
  it('把后端 snake_case 仪表盘字段转换成前端 camelCase', async () => {
    httpGetSpy.mockResolvedValue({
      data: {
        data: {
          overview: {
            total_tokens: 123,
            total_cost_usd: '9.99',
            total_requests: 4,
            active_clients: 2,
          },
          trend: [
            {
              bucket: '2026-05-22T10:00:00+08:00',
              input_tokens: 10,
              output_tokens: 20,
              cache_read_tokens: 30,
              cache_creation_tokens: 40,
              total_requests: 1,
              total_cost_usd: '1.23',
            },
          ],
          top_models: [
            {
              model: 'claude-sonnet-4-0',
              display_name: 'Claude Sonnet 4',
              total_tokens: 60,
            },
          ],
          top_clients: [
            {
              client_id: 'client-a',
              total_cost_usd: '7.89',
            },
          ],
          cache_analysis: {
            saved_cost_usd: '3.21',
            cache_read_cost_usd: '1.11',
            cache_creation_cost_usd: '0.22',
          },
        },
      },
    });

    const result = await getDashboard({ startAt: 1, endAt: 2, interval: 'hour' });

    expect(result).toEqual({
      overview: {
        totalTokens: 123,
        totalCostUsd: '9.99',
        totalRequests: 4,
        activeClients: 2,
      },
      trend: [
        {
          bucket: '2026-05-22T10:00:00+08:00',
          inputTokens: 10,
          outputTokens: 20,
          cacheReadTokens: 30,
          cacheCreationTokens: 40,
          totalRequests: 1,
          totalCostUsd: '1.23',
        },
      ],
      topModels: [
        {
          model: 'claude-sonnet-4-0',
          displayName: 'Claude Sonnet 4',
          totalTokens: 60,
        },
      ],
      topClients: [
        {
          clientId: 'client-a',
          totalCostUsd: '7.89',
        },
      ],
      cacheAnalysis: {
        savedCostUsd: '3.21',
        cacheReadCostUsd: '1.11',
        cacheCreationCostUsd: '0.22',
      },
    });
  });
});
