import { formatMetricValue, formatNumberInWanYi } from '@/utils/format';
import type { DashboardOverview } from '@/types/dashboard';
import type { TimeRangePreset } from '@/utils/timeRange';

interface OverviewCardsProps {
  overview?: DashboardOverview;
  previousOverview?: DashboardOverview;
  preset: TimeRangePreset;
}

function PrimaryMetricCard(props: { title: string; value: string; note: string; accent: string; change?: string; isPositive?: boolean; compareLabel?: string; hideArrow?: boolean; secondaryValue?: string }) {
  return (
    <article
      className={[
        'overflow-hidden rounded-[32px] border border-white/80 px-5 pt-4 pb-5 shadow-[0_22px_60px_rgba(111,153,200,0.16)] backdrop-blur-xl flex flex-col justify-between',
        props.accent,
      ].join(' ')}
    >
      <div>
        <p className="text-xs uppercase tracking-[0.3em] text-[#6c92b4]">{props.title}</p>
        <p className="mt-3 text-4xl font-semibold text-[#12304d] xl:text-[2.8rem]">{props.value}</p>
        {props.secondaryValue && (
          <p className="mt-1 text-sm text-[#5f7f9e]">{props.secondaryValue}</p>
        )}
      </div>
      <div className="mt-4">
        <p className="text-sm text-[#5f7f9e]">{props.note}</p>
        {props.change && props.compareLabel && (
          <div className="mt-1">
            <p className="text-xs text-[#8ba5bf]">{props.compareLabel}</p>
            <span className={`text-sm font-medium ${props.isPositive ? 'text-[#16a34a]' : 'text-[#dc2626]'}`}>
              {props.hideArrow ? '' : (props.isPositive ? '↓' : '↑')} {props.change}
            </span>
          </div>
        )}
      </div>
    </article>
  );
}

function SecondaryMetricCard(props: { title: string; value: string; note: string; change?: string; isPositive?: boolean; compareLabel?: string }) {
  return (
    <article className="rounded-[28px] border border-white/80 bg-white/72 px-4 pt-3 pb-4 shadow-[0_18px_48px_rgba(111,153,200,0.14)] backdrop-blur-xl flex flex-col justify-between">
      <div>
        <p className="text-xs uppercase tracking-[0.28em] text-[#6c92b4]">{props.title}</p>
        <p className="mt-2 text-3xl font-semibold text-[#12304d]">{props.value}</p>
      </div>
      <div className="mt-3">
        <p className="text-sm text-[#5f7f9e]">{props.note}</p>
        {props.change && props.compareLabel && (
          <div className="mt-1">
            <p className="text-xs text-[#8ba5bf]">{props.compareLabel}</p>
            <span className={`text-sm font-medium ${props.isPositive ? 'text-[#16a34a]' : 'text-[#dc2626]'}`}>
              {props.isPositive ? '↓' : '↑'} {props.change}
            </span>
          </div>
        )}
      </div>
    </article>
  );
}

function calculateChange(current: number | undefined, previous: number | undefined): { change: string; isPositive: boolean } | null {
  if (current == null || previous == null || previous === 0) return null;
  const diff = ((current - previous) / previous) * 100;
  if (isNaN(diff) || !isFinite(diff)) return null;
  const change = `${Math.abs(diff).toFixed(1)}%`;
  return { change, isPositive: diff < 0 };
}

function getCompareLabel(preset: TimeRangePreset): string {
  switch (preset) {
    case 'today':
      return '相对于昨天';
    case 'last7Days':
      return '相对于上七天';
    case 'last30Days':
      return '相对于上一个月';
    case 'thisMonth':
      return '相对于上月';
    case 'all':
      return '相对于上一周期';
    default:
      return '相对于上一周期';
  }
}

export function OverviewCards({ overview, previousOverview, preset }: OverviewCardsProps) {
  const totalCost = overview ? formatMetricValue(overview.totalCostUsd, 'currency') : '--';
  const totalTokens = overview ? formatMetricValue(overview.totalTokens, 'number') : '--';
  const totalCacheTokens = overview ? formatMetricValue(overview.totalCacheTokens, 'number') : '--';
  const totalRequests = overview ? formatMetricValue(overview.totalRequests, 'number') : '--';
  const activeClients = overview ? formatMetricValue(overview.activeClients, 'number') : '--';

  const totalTokensWanYi = overview ? formatNumberInWanYi(overview.totalTokens) : null;
  const totalCacheTokensWanYi = overview ? formatNumberInWanYi(overview.totalCacheTokens) : null;

  const costChange = overview && previousOverview
    ? calculateChange(parseFloat(overview.totalCostUsd), parseFloat(previousOverview.totalCostUsd))
    : null;

  const tokenChange = overview && previousOverview
    ? calculateChange(overview.totalTokens, previousOverview.totalTokens)
    : null;

  const requestChange = overview && previousOverview
    ? calculateChange(overview.totalRequests, previousOverview.totalRequests)
    : null;

  // 计算缓存命中率: cache_read_tokens / (input_tokens + cache_read_tokens)
  function calculateCacheHitRate(cacheRead: number, input: number): string | null {
    const denominator = input + cacheRead;
    if (denominator === 0) return null;
    const rate = (cacheRead / denominator) * 100;
    return `${rate.toFixed(1)}%`;
  }

  const cacheHitRate = overview
    ? calculateCacheHitRate(overview.cacheReadTokens, overview.inputTokens)
    : null;

  const compareLabel = getCompareLabel(preset);

  return (
    <section className="grid gap-4 xl:grid-cols-[1fr_1fr_1fr_0.9fr]">
      <PrimaryMetricCard
        title="总费用"
        value={totalCost}
        note="费用为核心指标。"
        accent="bg-[linear-gradient(145deg,rgba(91,178,255,0.18),rgba(255,255,255,0.88))]"
        change={costChange?.change}
        isPositive={costChange?.isPositive}
        compareLabel={costChange ? compareLabel : undefined}
      />
      <PrimaryMetricCard
        title="总 Token"
        value={totalTokens}
        secondaryValue={totalTokensWanYi ?? undefined}
        note="优先展示用量数据。"
        accent="bg-[linear-gradient(145deg,rgba(181,224,255,0.68),rgba(255,255,255,0.92))]"
        change={tokenChange?.change}
        isPositive={tokenChange?.isPositive}
        compareLabel={tokenChange ? compareLabel : undefined}
      />
      <PrimaryMetricCard
        title="总缓存"
        value={totalCacheTokens}
        secondaryValue={totalCacheTokensWanYi ?? undefined}
        note="缓存 Token 统计。"
        accent="bg-[linear-gradient(145deg,rgba(134,239,172,0.35),rgba(255,255,255,0.9))]"
        change={cacheHitRate ?? undefined}
        isPositive={true}
        compareLabel={cacheHitRate ? '缓存命中率' : undefined}
        hideArrow={true}
      />
      <div className="grid gap-4">
        <SecondaryMetricCard
          title="总请求数"
          value={totalRequests}
          note="请求总数统计。"
          change={requestChange?.change}
          isPositive={requestChange?.isPositive}
          compareLabel={requestChange ? compareLabel : undefined}
        />
        <SecondaryMetricCard title="活跃客户端" value={activeClients} note="活跃客户端数量。" />
      </div>
    </section>
  );
}
