import { describe, expect, it } from 'vitest';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';

import App from '@/App';

function renderApp() {
  const queryClient = new QueryClient();

  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter initialEntries={['/']}>
        <App />
      </MemoryRouter>
    </QueryClientProvider>,
  );
}

describe('DashboardPage', () => {
  it('在根路由渲染静态复刻页面的关键模块', () => {
    renderApp();

    expect(screen.getByRole('heading', { name: 'Claude 用量看板' })).toBeInTheDocument();
    expect(screen.getByText('核心指标')).toBeInTheDocument();
    expect(screen.getByText('综合效益')).toBeInTheDocument();
    expect(screen.getByText('最近请求')).toBeInTheDocument();
  });
});
