import { describe, expect, it, vi } from 'vitest';

import { getRecentLogs } from '@/api/logs';

const httpGetSpy = vi.fn();

vi.mock('@/api/http', () => ({
  http: {
    get: (...args: unknown[]) => httpGetSpy(...args),
  },
}));

describe('getRecentLogs', () => {
  it('把后端 snake_case 日志字段转换成前端 camelCase', async () => {
    httpGetSpy.mockResolvedValue({
      data: {
        data: [
          {
            id: 1,
            client_id: 'client-a',
            request_id: 'req-1',
            model: 'claude-sonnet-4-0',
            input_tokens: 12,
            output_tokens: 34,
            cache_read_tokens: 56,
            cache_creation_tokens: 78,
            total_cost_usd: '1.23',
            created_at: 1747879200,
          },
        ],
        total: 1,
        offset: 0,
        limit: 8,
        page: 1,
      },
    });

    const result = await getRecentLogs({ startAt: 1, endAt: 2, limit: 8 });

    expect(result.data[0]).toEqual({
      id: 1,
      clientId: 'client-a',
      requestId: 'req-1',
      model: 'claude-sonnet-4-0',
      inputTokens: 12,
      outputTokens: 34,
      cacheReadTokens: 56,
      cacheCreationTokens: 78,
      totalCostUsd: '1.23',
      createdAt: 1747879200,
    });
  });
});
