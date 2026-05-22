import { http } from '@/api/http';
import type { RecentLogQuery, RecentLogsResponse } from '@/types/logs';

export async function getRecentLogs(query: RecentLogQuery): Promise<RecentLogsResponse> {
  const response = await http.get<RecentLogsResponse>('/logs', {
    params: {
      start_time: query.startAt,
      end_time: query.endAt,
      limit: query.limit,
      offset: 0,
    },
  });

  return response.data;
}
