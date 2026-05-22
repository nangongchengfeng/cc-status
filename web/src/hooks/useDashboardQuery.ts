import { getDashboard } from '@/api/dashboard';
import type { DashboardQuery } from '@/types/dashboard';
import { getDashboardQueryKey } from '@/utils/timeRange';
import { useQuery } from '@tanstack/react-query';

export function useDashboardQuery(query: DashboardQuery) {
  return useQuery({
    queryKey: getDashboardQueryKey(query),
    queryFn: () => getDashboard(query),
  });
}
