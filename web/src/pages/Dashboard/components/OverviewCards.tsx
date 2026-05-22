import { formatMetricValue } from '@/utils/format';
import type { DashboardTrendPoint } from '@/types/dashboard';

interface OverviewCardsProps {
  overview?: {
    totalTokens: number;
    totalCostUsd: string;
    totalRequests: number;
    activeClients: number;
  };
  trend?: DashboardTrendPoint[];
}

function PrimaryMetricCard(props: { title: string; value: string; note: string; accent: string; change?: string; isPositive?: boolean }) {
  return (
    <article
      className={[
        'overflow-hidden rounded-[32px] border border-white/80 p-6 shadow-[0_22px_60px_rgba(111,153,200,0.16)] backdrop-blur-xl',
        props.accent,
      ].join(' ')}
    >
      <p className="text-xs uppercase tracking-[0.3em] text-[#6c92b4]">{props.title}</p>
      <p className="mt-4 text-4xl font-semibold text-[#12304d] xl:text-[2.8rem]">{props.value}</p>
      <div className="mt-3 flex items-center gap-2">
        <p className="text-sm text-[#5f7f9e]">{props.note}</p>
        {props.change && (
          <span className={`text-sm font-medium ${props.isPositive ? 'text-[#16a34a]' : 'text-[#dc2626]'}`}>
            {props.isPositive ? '↑' : '↓'} {props.change}
          </span>
        )}
      </div>
    </article>
  );
}

function SecondaryMetricCard(props: { title: string; value: string; note: string; change?: string; isPositive?: boolean }) {
  return (
    <article className="rounded-[28px] border border-white/80 bg-white/72 p-5 shadow-[0_18px_48px_rgba(111,153,200,0.14)] backdrop-blur-xl">
      <p className="text-xs uppercase tracking-[0.28em] text-[#6c92b4]">{props.title}</p>
      <p className="mt-3 text-3xl font-semibold text-[#12304d]">{props.value}</p>
      <div className="mt-2 flex items-center gap-2">
        <p className="text-sm text-[#5f7f9e]">{props.note}</p>
        {props.change && (
          <span className={`text-sm font-medium ${props.isPositive ? 'text-[#16a34a]' : 'text-[#dc2626]'}`}>
            {props.isPositive ? '↑' : '↓'} {props.change}
          </span>
        )}
      </div>
    </article>
  );
}

function calculateChange(current: number, previous: number): { change: string; isPositive: boolean } | null {
  if (previous === 0) return null;
  const diff = ((current - previous) / previous) * 100;
  const change = `${Math.abs(diff).toFixed(1)}%`;
  return { change, isPositive: diff > 0 };
}

export function OverviewCards({ overview, trend }: OverviewCardsProps) {
  const totalCost = overview ? formatMetricValue(overview.totalCostUsd, 'currency') : '--';
  const totalTokens = overview ? formatMetricValue(overview.totalTokens, 'number') : '--';
  const totalRequests = overview ? formatMetricValue(overview.totalRequests, 'number') : '--';
  const activeClients = overview ? formatMetricValue(overview.activeClients, 'number') : '--';

  // 从趋势数据中获取最后两天的数据进行对比
  const sortedTrend = [...(trend ?? [])].sort((a, b) => a.bucket.localeCompare(b.bucket));
  const yesterday = sortedTrend[sortedTrend.length - 2];
  const today = sortedTrend[sortedTrend.length - 1];

  const costChange = yesterday && today
    ? calculateChange(Number(today.totalCostUsd), Number(yesterday.totalCostUsd))
    : null;

  const tokenChange = yesterday && today
    ? calculateChange(today.inputTokens + today.outputTokens, yesterday.inputTokens + yesterday.outputTokens)
    : null;

  const requestChange = yesterday && today
    ? calculateChange(today.totalRequests, yesterday.totalRequests)
    : null;

  return (
    <section className="grid gap-4 xl:grid-cols-[1.28fr_1fr_0.9fr]">
      <PrimaryMetricCard
        title="总费用"
        value={totalCost}
        note="费用为核心指标。"
        accent="bg-[linear-gradient(145deg,rgba(91,178,255,0.18),rgba(255,255,255,0.88))]"
        change={costChange?.change}
        isPositive={costChange?.isPositive}
      />
      <PrimaryMetricCard
        title="总 Token"
        value={totalTokens}
        note="优先展示用量数据。"
        accent="bg-[linear-gradient(145deg,rgba(181,224,255,0.68),rgba(255,255,255,0.92))]"
        change={tokenChange?.change}
        isPositive={tokenChange?.isPositive}
      />
      <div className="grid gap-4">
        <SecondaryMetricCard
          title="总请求数"
          value={totalRequests}
          note="请求总数统计。"
          change={requestChange?.change}
          isPositive={requestChange?.isPositive}
        />
        <SecondaryMetricCard title="活跃客户端" value={activeClients} note="活跃客户端数量。" />
      </div>
    </section>
  );
}
