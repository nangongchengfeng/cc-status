import App from '@/App';
import { createRoot } from 'react-dom/client';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { BrowserRouter } from 'react-router-dom';

import '@/styles/index.css';

const queryClient = new QueryClient();
const container = document.getElementById('root');

if (!container) {
  throw new Error('找不到根节点');
}

createRoot(container).render(
  <QueryClientProvider client={queryClient}>
    <BrowserRouter>
      <App />
    </BrowserRouter>
  </QueryClientProvider>,
);
