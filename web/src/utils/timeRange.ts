import type { DashboardInterval, DashboardQuery } from '@/types/dashboard';

export type TimeRangePreset = 'today' | 'last7Days' | 'last30Days' | 'thisMonth' | 'all';

export const ALL_RANGE_START_AT = 1704038400;
const SHANGHAI_OFFSET_MS = 8 * 60 * 60 * 1000;

function toShanghaiDateParts(date: Date) {
  const shanghaiDate = new Date(date.getTime() + SHANGHAI_OFFSET_MS);

  return {
    year: shanghaiDate.getUTCFullYear(),
    month: shanghaiDate.getUTCMonth(),
    day: shanghaiDate.getUTCDate(),
  };
}

function shanghaiDateToUnix(year: number, month: number, day: number) {
  return Math.floor(Date.UTC(year, month, day, 0, 0, 0) / 1000) - 8 * 60 * 60;
}

export function buildQueryTimeRange(preset: TimeRangePreset, now = new Date()): DashboardQuery {
  const endAt = Math.floor(now.getTime() / 1000);
  const { year, month, day } = toShanghaiDateParts(now);

  if (preset === 'today') {
    return {
      startAt: shanghaiDateToUnix(year, month, day),
      endAt,
      interval: 'hour',
    };
  }

  if (preset === 'last7Days') {
    return {
      startAt: shanghaiDateToUnix(year, month, day - 6),
      endAt,
      interval: 'day',
    };
  }

  if (preset === 'last30Days') {
    return {
      startAt: shanghaiDateToUnix(year, month, day - 29),
      endAt,
      interval: 'day',
    };
  }

  if (preset === 'thisMonth') {
    return {
      startAt: shanghaiDateToUnix(year, month, 1),
      endAt,
      interval: 'day',
    };
  }

  return {
    startAt: ALL_RANGE_START_AT,
    endAt,
    interval: 'day',
  };
}

export function getDashboardQueryKey(query: DashboardQuery) {
  return ['dashboard', query.startAt, query.endAt, query.interval] as const;
}

export function getRecentLogsQueryKey(query: { startAt: number; endAt: number; limit: number }) {
  return ['recent-logs', query.startAt, query.endAt, query.limit] as const;
}

export function getDashboardIntervalLabel(interval: DashboardInterval) {
  return interval === 'hour' ? '小时' : '天';
}
