import { getRecentLogs } from '@/api/logs';
import type { RecentLogQuery } from '@/types/logs';
import { getRecentLogsQueryKey } from '@/utils/timeRange';
import { useQuery } from '@tanstack/react-query';

export function useRecentLogsQuery(query: RecentLogQuery) {
  return useQuery({
    queryKey: getRecentLogsQueryKey(query),
    queryFn: () => getRecentLogs(query),
  });
}
