import { describe, expect, it } from 'vitest';

import { formatBucketLabel, formatMetricValue } from '@/utils/format';

describe('format helpers', () => {
  it('格式化金额与大数字卡片值', () => {
    expect(formatMetricValue('12.34', 'currency')).toBe('$12.34');
    expect(formatMetricValue(1234567, 'number')).toBe('1,234,567');
  });

  it('把趋势桶格式化成适合 tooltip 的标签', () => {
    expect(formatBucketLabel('2026-05-22T10:00:00+08:00', 'hour')).toBe('05-22 10:00');
    expect(formatBucketLabel('2026-05-22T00:00:00+08:00', 'day')).toBe('05-22');
  });
});
