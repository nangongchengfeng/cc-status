import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { describe, expect, it } from 'vitest';

import { DashboardPage } from '@/pages/Dashboard/DashboardPage';

function renderPage() {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
      },
    },
  });

  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter>
        <DashboardPage />
      </MemoryRouter>
    </QueryClientProvider>,
  );
}

describe('DashboardPage states', () => {
  it('渲染静态复刻页的主体区域和面板标题', () => {
    renderPage();

    expect(screen.getByText('核心指标')).toBeInTheDocument();
    expect(screen.getByText('综合效益')).toBeInTheDocument();
    expect(screen.getByText('使用趋势分析')).toBeInTheDocument();
    expect(screen.getByText('最近请求')).toBeInTheDocument();
  });

  it('渲染静态复刻页的关键指标、表格和底部信息', () => {
    renderPage();

    expect(screen.getByRole('heading', { name: 'Claude 用量看板' })).toBeInTheDocument();
    expect(screen.getByText(/数据更新于/)).toBeInTheDocument();
    expect(screen.getByText('总使用量（Tokens）')).toBeInTheDocument();
    expect(screen.getByText('总费用（USD）')).toBeInTheDocument();
    expect(screen.getByText('请求数')).toBeInTheDocument();
    expect(screen.getByText('活跃客户端')).toBeInTheDocument();
    expect(screen.getByText('节省费用')).toBeInTheDocument();
    expect(screen.getByText('节省时间成本')).toBeInTheDocument();
    expect(screen.getByText('节省资源成本')).toBeInTheDocument();
    expect(screen.getByText('费用趋势（USD）')).toBeInTheDocument();
    expect(screen.getByText('Token 趋势（万）')).toBeInTheDocument();
    expect(screen.getByText('模型排行')).toBeInTheDocument();
    expect(screen.getByText('客户端排行')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: '查看全部' })).toBeInTheDocument();
    expect(screen.getByText('时间')).toBeInTheDocument();
    expect(screen.getByText('模型')).toBeInTheDocument();
    expect(screen.getByText('输入')).toBeInTheDocument();
    expect(screen.getByText('输出')).toBeInTheDocument();
    expect(screen.getByText('费用（USD）')).toBeInTheDocument();
    expect(screen.getByText('客户端')).toBeInTheDocument();
    expect(screen.getByText(/数据来源：Claude API/)).toBeInTheDocument();
  });
});
