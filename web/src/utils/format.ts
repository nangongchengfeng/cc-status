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
