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
  it('在根路由渲染仪表盘页头和卡片骨架', () => {
    renderApp();

    expect(screen.getByRole('heading', { name: 'Claude 用量指挥舱' })).toBeInTheDocument();
    expect(screen.getByText('先看花费，再找来源。')).toBeInTheDocument();
    expect(screen.getByText('总 Token')).toBeInTheDocument();
  });
});
