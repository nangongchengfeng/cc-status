import { formatMetricValue } from '@/utils/format';

interface OverviewCardsProps {
  overview?: {
    totalTokens: number;
    totalCostUsd: string;
    totalRequests: number;
    activeClients: number;
  };
}

function PrimaryMetricCard(props: { title: string; value: string; note: string; accent: string }) {
  return (
    <article
      className={[
        'overflow-hidden rounded-[32px] border border-white/80 p-6 shadow-[0_22px_60px_rgba(111,153,200,0.16)] backdrop-blur-xl',
        props.accent,
      ].join(' ')}
    >
      <p className="text-xs uppercase tracking-[0.3em] text-[#6c92b4]">{props.title}</p>
      <p className="mt-4 text-4xl font-semibold text-[#12304d] xl:text-[2.8rem]">{props.value}</p>
      <p className="mt-3 text-sm text-[#5f7f9e]">{props.note}</p>
    </article>
  );
}

function SecondaryMetricCard(props: { title: string; value: string; note: string }) {
  return (
    <article className="rounded-[28px] border border-white/80 bg-white/72 p-5 shadow-[0_18px_48px_rgba(111,153,200,0.14)] backdrop-blur-xl">
      <p className="text-xs uppercase tracking-[0.28em] text-[#6c92b4]">{props.title}</p>
      <p className="mt-3 text-3xl font-semibold text-[#12304d]">{props.value}</p>
      <p className="mt-2 text-sm text-[#5f7f9e]">{props.note}</p>
    </article>
  );
}

export function OverviewCards({ overview }: OverviewCardsProps) {
  const totalCost = overview ? formatMetricValue(overview.totalCostUsd, 'currency') : '--';
  const totalTokens = overview ? formatMetricValue(overview.totalTokens, 'number') : '--';
  const totalRequests = overview ? formatMetricValue(overview.totalRequests, 'number') : '--';
  const activeClients = overview ? formatMetricValue(overview.activeClients, 'number') : '--';

  return (
    <section className="grid gap-4 xl:grid-cols-[1.28fr_1fr_0.9fr]">
      <PrimaryMetricCard
        title="总费用"
        value={totalCost}
        note="费用为核心指标。"
        accent="bg-[linear-gradient(145deg,rgba(91,178,255,0.18),rgba(255,255,255,0.88))]"
      />
      <PrimaryMetricCard
        title="总 Token"
        value={totalTokens}
        note="优先展示用量数据。"
        accent="bg-[linear-gradient(145deg,rgba(181,224,255,0.68),rgba(255,255,255,0.92))]"
      />
      <div className="grid gap-4">
        <SecondaryMetricCard title="总请求数" value={totalRequests} note="请求总数统计。" />
        <SecondaryMetricCard title="活跃客户端" value={activeClients} note="活跃客户端数量。" />
      </div>
    </section>
  );
}
