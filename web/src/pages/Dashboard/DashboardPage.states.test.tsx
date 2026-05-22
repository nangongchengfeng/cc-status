import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { beforeEach, describe, expect, it, vi } from 'vitest';

import { DashboardPage } from '@/pages/Dashboard/DashboardPage';

const mockUseDashboardQuery = vi.fn();
const mockUseRecentLogsQuery = vi.fn();

vi.mock('@/hooks/useDashboardQuery', () => ({
  useDashboardQuery: (...args: unknown[]) => mockUseDashboardQuery(...args),
}));

vi.mock('@/hooks/useRecentLogsQuery', () => ({
  useRecentLogsQuery: (...args: unknown[]) => mockUseRecentLogsQuery(...args),
}));

function renderPage() {
  return render(
    <MemoryRouter>
      <DashboardPage />
    </MemoryRouter>,
  );
}

const successDashboardData = {
  overview: {
    totalTokens: 123456,
    totalCostUsd: '99.01',
    totalRequests: 18,
    activeClients: 3,
  },
  trend: [
    {
      bucket: '2025-05-22T10:00:00+08:00',
      inputTokens: 1000,
      outputTokens: 800,
      cacheReadTokens: 200,
      cacheCreationTokens: 50,
      totalRequests: 2,
      totalCostUsd: '12.34',
    },
  ],
  topModels: [
    {
      model: 'claude-3-7-sonnet-very-long-model-name',
      displayName: 'Claude 3.7 Sonnet Very Long Name',
      totalTokens: 1800,
    },
  ],
  topClients: [
    {
      clientId: 'client-with-a-very-long-identifier',
      totalCostUsd: '45.67',
    },
  ],
  cacheAnalysis: {
    savedCostUsd: '45.67',
    cacheReadCostUsd: '12.00',
    cacheCreationCostUsd: '4.00',
  },
};

beforeEach(() => {
  vi.clearAllMocks();
});

describe('DashboardPage states', () => {
  it('在加载中时展示稳定的加载提示', () => {
    mockUseDashboardQuery.mockReturnValue({
      status: 'pending',
      isPending: true,
      isError: false,
      data: undefined,
    });
    mockUseRecentLogsQuery.mockReturnValue({
      status: 'success',
      isPending: false,
      isError: false,
      data: { data: [] },
    });

    renderPage();

    expect(screen.getByText('数据加载中，正在刷新当前时间范围。')).toBeInTheDocument();
  });

  it('在错误态时展示统一错误提示', () => {
    mockUseDashboardQuery.mockReturnValue({
      status: 'error',
      isPending: false,
      isError: true,
      data: undefined,
    });
    mockUseRecentLogsQuery.mockReturnValue({
      status: 'success',
      isPending: false,
      isError: false,
      data: { data: [] },
    });

    renderPage();

    expect(screen.getByText('数据暂时没接上，但页面状态是稳的。')).toBeInTheDocument();
  });

  it('在空态时展示稳定的模块空提示', () => {
    mockUseDashboardQuery.mockReturnValue({
      status: 'success',
      isPending: false,
      isError: false,
      data: {
        overview: {
          totalTokens: 0,
          totalCostUsd: '0.00',
          totalRequests: 0,
          activeClients: 0,
        },
        trend: [],
        topModels: [],
        topClients: [],
        cacheAnalysis: undefined,
      },
    });
    mockUseRecentLogsQuery.mockReturnValue({
      status: 'success',
      isPending: false,
      isError: false,
      data: { data: [] },
    });

    renderPage();

    expect(screen.getByText('当前时间范围还没有缓存数据。')).toBeInTheDocument();
    expect(screen.getByText('当前时间范围还没有最近请求。')).toBeInTheDocument();
  });

  it('在成功态时展示格式化后的核心数据', () => {
    mockUseDashboardQuery.mockReturnValue({
      status: 'success',
      isPending: false,
      isError: false,
      data: successDashboardData,
    });
    mockUseRecentLogsQuery.mockReturnValue({
      status: 'success',
      isPending: false,
      isError: false,
      data: {
        data: [
          {
            id: 1,
            clientId: 'client-with-a-very-long-identifier',
            requestId: 'req-1',
            model: 'claude-3-7-sonnet-very-long-model-name',
            inputTokens: 1000,
            outputTokens: 800,
            cacheReadTokens: 200,
            cacheCreationTokens: 50,
            totalCostUsd: '12.34',
            createdAt: 1747879200,
          },
        ],
      },
    });

    renderPage();

    expect(screen.getByText('2025-05-22 10:00')).toBeInTheDocument();
    expect(screen.getByText('$45.67')).toBeInTheDocument();
    expect(screen.getByTitle('claude-3-7-sonnet-very-long-model-name')).toHaveTextContent('claude-3-7-sonnet-very-…');
    expect(screen.getByTitle('client-with-a-very-long-identifier')).toHaveTextContent('client-with-a-very-…');
  });
});
