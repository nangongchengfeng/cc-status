import type { DashboardCacheAnalysis } from '@/types/dashboard';
import { formatMetricValue } from '@/utils/format';

interface CacheAnalysisProps {
  analysis?: DashboardCacheAnalysis;
}

const CACHE_METRICS = [
  {
    key: 'savedCostUsd',
    label: '节省金额',
    accent: 'text-[#1667a6]',
    tone: 'bg-[linear-gradient(145deg,rgba(89,180,255,0.16),rgba(255,255,255,0.92))]',
  },
  {
    key: 'cacheReadCostUsd',
    label: '缓存读取成本',
    accent: 'text-[#1a6e75]',
    tone: 'bg-white/72',
  },
  {
    key: 'cacheCreationCostUsd',
    label: '缓存建设成本',
    accent: 'text-[#335d9a]',
    tone: 'bg-white/72',
  },
] as const;

export function CacheAnalysis({ analysis }: CacheAnalysisProps) {
  if (!analysis) {
    return (
      <div className="grid h-[220px] place-items-center rounded-[28px] border border-dashed border-[#cfe0f1] bg-[linear-gradient(145deg,rgba(255,255,255,0.72),rgba(232,243,252,0.72))] text-sm text-[#6a86a3]">
        当前时间范围还没有缓存数据。
      </div>
    );
  }

  return (
    <div className="grid gap-4 md:grid-cols-[1.25fr_1fr_1fr]">
      {CACHE_METRICS.map((metric) => (
        <article
          key={metric.key}
          className={[
            'rounded-[28px] border border-white/80 p-5 shadow-[0_18px_48px_rgba(111,153,200,0.14)] backdrop-blur-xl',
            metric.tone,
          ].join(' ')}
        >
          <p className="text-xs uppercase tracking-[0.25em] text-[#6c92b4]">{metric.label}</p>
          <p className={`mt-4 text-3xl font-semibold ${metric.accent}`}>{formatMetricValue(analysis[metric.key], 'currency')}</p>
          <p className="mt-2 text-sm text-[#5f7f9e]">{metric.key === 'savedCostUsd' ? '先看省下多少。' : '这是换来的成本。'}</p>
        </article>
      ))}
    </div>
  );
}
