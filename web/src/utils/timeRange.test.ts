import { describe, expect, it, vi } from 'vitest';

import { ALL_RANGE_START_AT, buildQueryTimeRange, type TimeRangePreset } from '@/utils/timeRange';

describe('buildQueryTimeRange', () => {
  it('把今天转换成小时粒度的显式范围', () => {
    const now = new Date('2026-05-22T12:30:00+08:00');

    const result = buildQueryTimeRange('today', now);

    expect(result.interval).toBe('hour');
    expect(result.startAt).toBe(Math.floor(new Date('2026-05-22T00:00:00+08:00').getTime() / 1000));
    expect(result.endAt).toBe(Math.floor(now.getTime() / 1000));
  });

  it.each<[TimeRangePreset, 'hour' | 'day']>([
    ['last7Days', 'day'],
    ['last30Days', 'day'],
    ['thisMonth', 'day'],
    ['all', 'day'],
  ])('把 %s 转成稳定参数', (preset, expectedInterval) => {
    const now = new Date('2026-05-22T12:30:00+08:00');

    const result = buildQueryTimeRange(preset, now);

    expect(result.interval).toBe(expectedInterval);
    expect(result.endAt).toBe(Math.floor(now.getTime() / 1000));
    expect(result.startAt).toBeGreaterThan(0);
  });

  it('把全部范围固定到项目允许的最早时间', () => {
    const now = new Date('2026-05-22T12:30:00+08:00');

    const result = buildQueryTimeRange('all', now);

    expect(result.startAt).toBe(ALL_RANGE_START_AT);
    expect(result.interval).toBe('day');
  });
});
