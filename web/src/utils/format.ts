import type { DashboardInterval } from '@/types/dashboard';

export function formatMetricValue(value: number | string, type: 'number' | 'currency') {
  if (type === 'currency') {
    const numericValue = typeof value === 'string' ? Number(value) : value;
    return `$${numericValue.toFixed(2)}`;
  }

  const numericValue = typeof value === 'number' ? value : Number(value);
  return new Intl.NumberFormat('en-US').format(numericValue);
}

export function formatBucketLabel(bucket: string, interval: DashboardInterval) {
  const date = new Date(bucket);
  const month = `${date.getMonth() + 1}`.padStart(2, '0');
  const day = `${date.getDate()}`.padStart(2, '0');

  if (interval === 'day') {
    return `${month}-${day}`;
  }

  const hour = `${date.getHours()}`.padStart(2, '0');
  return `${month}-${day} ${hour}:00`;
}

export function getModelDisplayName(model: { displayName?: string; model: string }) {
  return model.displayName?.trim() ? model.displayName : model.model;
}

export function formatRecentRequestTime(createdAt: number) {
  const formatter = new Intl.DateTimeFormat('zh-CN', {
    timeZone: 'Asia/Shanghai',
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    hour12: false,
  });

  // 大屏时间口径统一固定到上海时区，避免浏览器所在时区影响展示。
  const parts = formatter.formatToParts(new Date(createdAt * 1000));
  const valueByType = Object.fromEntries(parts.filter((part) => part.type !== 'literal').map((part) => [part.type, part.value]));

  return `${valueByType.year}-${valueByType.month}-${valueByType.day} ${valueByType.hour}:${valueByType.minute}`;
}
