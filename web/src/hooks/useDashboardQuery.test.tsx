import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { renderHook, waitFor } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

import { useDashboardQuery } from '@/hooks/useDashboardQuery';
import { useRecentLogsQuery } from '@/hooks/useRecentLogsQuery';
import type { DashboardQuery } from '@/types/dashboard';

const dashboardSpy = vi.fn();
const logsSpy = vi.fn();

vi.mock('@/api/dashboard', () => ({
  getDashboard: (params: unknown) => dashboardSpy(params),
}));

vi.mock('@/api/logs', () => ({
  getRecentLogs: (params: unknown) => logsSpy(params),
}));

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
      },
    },
  });

  return function Wrapper({ children }: { children: React.ReactNode }) {
    return <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>;
  };
}

describe('dashboard queries', () => {
  it('在参数变化时重新请求 dashboard 数据', async () => {
    dashboardSpy.mockResolvedValue({ overview: {}, trend: [], topModels: [], topClients: [], cacheAnalysis: {} });
    const initialQuery: DashboardQuery = { startAt: 1, endAt: 2, interval: 'hour' };
    const nextQuery: DashboardQuery = { startAt: 3, endAt: 4, interval: 'day' };

    const wrapper = createWrapper();
    const { rerender } = renderHook((query: DashboardQuery) => useDashboardQuery(query), {
      wrapper,
      initialProps: initialQuery,
    });

    await waitFor(() => {
      expect(dashboardSpy).toHaveBeenCalledWith(initialQuery);
    });

    rerender(nextQuery);

    await waitFor(() => {
      expect(dashboardSpy).toHaveBeenLastCalledWith(nextQuery);
    });
  });

  it('把最近请求接口错误稳定透传出来', async () => {
    logsSpy.mockRejectedValue(new Error('boom'));

    const wrapper = createWrapper();
    const { result } = renderHook(
      () => useRecentLogsQuery({ startAt: 1, endAt: 2, limit: 5 }),
      { wrapper },
    );

    await waitFor(() => {
      expect(result.current.isError).toBe(true);
    });
  });
});
