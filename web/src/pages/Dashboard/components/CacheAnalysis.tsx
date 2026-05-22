import type { DashboardCacheAnalysis } from '@/types/dashboard';
import { formatMetricValue } from '@/utils/format';

interface CacheAnalysisProps {
  analysis?: DashboardCacheAnalysis;
}

const CACHE_METRICS = [
  {
    key: 'savedCostUsd',
    label: '节省金额',
    accent: 'text-[#d8a978]',
  },
  {
    key: 'cacheReadCostUsd',
    label: '缓存读取成本',
    accent: 'text-[#63b59c]',
  },
  {
    key: 'cacheCreationCostUsd',
    label: '缓存建设成本',
    accent: 'text-[#8cb8ff]',
  },
] as const;

export function CacheAnalysis({ analysis }: CacheAnalysisProps) {
  if (!analysis) {
    return <div className="grid h-[220px] place-items-center rounded-[24px] border border-dashed border-white/15 bg-white/[0.03] text-sm text-[#cab99d]">当前时间范围还没有缓存数据。</div>;
  }

  return (
    <div className="grid gap-4 md:grid-cols-3">
      {CACHE_METRICS.map((metric) => (
        <article key={metric.key} className="rounded-[24px] border border-white/10 bg-black/10 p-5">
          <p className="text-xs uppercase tracking-[0.25em] text-[#d4c5a8]/70">{metric.label}</p>
          <p className={`mt-4 text-3xl font-semibold ${metric.accent}`}>{formatMetricValue(analysis[metric.key], 'currency')}</p>
        </article>
      ))}
    </div>
  );
}
