import type { DashboardInterval } from '@/types/dashboard';

function toFiniteNumber(value: number | string | null | undefined) {
  if (value === null || value === undefined) {
    return null;
  }

  if (typeof value === 'string' && value.trim() === '') {
    return null;
  }

  const numericValue = typeof value === 'number' ? value : Number(value);
  if (!Number.isFinite(numericValue)) {
    return null;
  }

  return numericValue;
}

export function formatMetricValue(value: number | string | null | undefined, type: 'number' | 'currency') {
  const numericValue = toFiniteNumber(value);
  if (numericValue === null) {
    return '--';
  }

  if (type === 'currency') {
    return `$${numericValue.toFixed(2)}`;
  }

  return new Intl.NumberFormat('en-US').format(numericValue);
}

export function formatLargeNumber(value: number) {
  if (value >= 1_000_000_000) {
    return `${(value / 1_000_000_000).toFixed(1)}B`;
  }
  if (value >= 1_000_000) {
    return `${(value / 1_000_000).toFixed(1)}M`;
  }
  if (value >= 1_000) {
    return `${(value / 1_000).toFixed(1)}K`;
  }
  return String(value);
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
  const date = new Date(createdAt * 1000);
  if (Number.isNaN(date.getTime())) {
    return '--';
  }

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
  const parts = formatter.formatToParts(date);
  const valueByType = Object.fromEntries(parts.filter((part) => part.type !== 'literal').map((part) => [part.type, part.value]));

  return `${valueByType.year}-${valueByType.month}-${valueByType.day} ${valueByType.hour}:${valueByType.minute}`;
}

export function truncateLabel(value: string, maxLength = 16) {
  if (value.length <= maxLength) {
    return value;
  }

  return `${value.slice(0, Math.max(1, maxLength - 1))}…`;
}

export function formatUnixTimestamp(timestamp: number) {
  const date = new Date(timestamp * 1000);
  if (Number.isNaN(date.getTime())) {
    return '--';
  }

  const formatter = new Intl.DateTimeFormat('zh-CN', {
    timeZone: 'Asia/Shanghai',
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    hour12: false,
  });

  const parts = formatter.formatToParts(date);
  const valueByType = Object.fromEntries(parts.filter((part) => part.type !== 'literal').map((part) => [part.type, part.value]));

  return `${valueByType.year}-${valueByType.month}-${valueByType.day} ${valueByType.hour}:${valueByType.minute}`;
}

export function formatNumberInWanYi(value: number): string | null {
  if (value >= 100_000_000) {
    return `≈ ${(value / 100_000_000).toFixed(2)} 亿`;
  }
  if (value >= 10_000) {
    return `≈ ${(value / 10_000).toFixed(2)} 万`;
  }
  return null;
}
