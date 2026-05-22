import type { TimeRangePreset } from '@/utils/timeRange';

export interface TimeRangeOption {
  value: TimeRangePreset;
  label: string;
  shortLabel: string;
}

export const TIME_RANGE_OPTIONS: TimeRangeOption[] = [
  { value: 'today', label: '今天', shortLabel: '今天' },
  { value: 'last7Days', label: '最近 7 天', shortLabel: '7 天' },
  { value: 'last30Days', label: '最近 30 天', shortLabel: '30 天' },
  { value: 'thisMonth', label: '本月', shortLabel: '本月' },
  { value: 'all', label: '全部', shortLabel: '全部' },
];
