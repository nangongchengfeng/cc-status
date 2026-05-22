import { http } from '@/api/http';
import type { DashboardQuery, DashboardResponse } from '@/types/dashboard';

interface DashboardEnvelope {
  data: DashboardResponse;
}

export async function getDashboard(query: DashboardQuery): Promise<DashboardResponse> {
  const response = await http.get<DashboardEnvelope>('/stats/dashboard', {
    params: {
      start_at: query.startAt,
      end_at: query.endAt,
      interval: query.interval,
    },
  });

  return response.data.data;
}
