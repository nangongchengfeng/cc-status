import { describe, expect, it } from 'vitest';

import { formatBucketLabel, formatMetricValue, formatRecentRequestTime, getModelDisplayName, truncateLabel } from '@/utils/format';

describe('format helpers', () => {
  it('格式化金额与大数字卡片值', () => {
    expect(formatMetricValue('12.34', 'currency')).toBe('$12.34');
    expect(formatMetricValue(1234567, 'number')).toBe('1,234,567');
  });

  it('把趋势桶格式化成适合 tooltip 的标签', () => {
    expect(formatBucketLabel('2026-05-22T10:00:00+08:00', 'hour')).toBe('05-22 10:00');
    expect(formatBucketLabel('2026-05-22T00:00:00+08:00', 'day')).toBe('05-22');
  });

  it('模型展示名缺失时回退到原始模型名', () => {
    expect(getModelDisplayName({ displayName: 'Claude Sonnet 4', model: 'claude-sonnet-4-0' })).toBe('Claude Sonnet 4');
    expect(getModelDisplayName({ displayName: '', model: 'claude-sonnet-4-0' })).toBe('claude-sonnet-4-0');
  });

  it('最近请求时间统一格式化为上海时区', () => {
    expect(formatRecentRequestTime(1747879200)).toBe('2025-05-22 10:00');
  });

  it('长标签会被统一截断，避免撑坏布局', () => {
    expect(truncateLabel('short-label')).toBe('short-label');
    expect(truncateLabel('claude-super-long-model-name', 12)).toBe('claude-supe…');
  });
});
