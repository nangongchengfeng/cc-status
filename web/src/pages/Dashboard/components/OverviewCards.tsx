import { formatMetricValue } from '@/utils/format';

interface OverviewCardsProps {
  overview?: {
    totalTokens: number;
    totalCostUsd: string;
    totalRequests: number;
    activeClients: number;
  };
}

function OverviewCard(props: { title: string; value: string; note: string }) {
  return (
    <article className="rounded-[28px] border border-white/10 bg-[linear-gradient(155deg,rgba(255,255,255,0.08),rgba(255,255,255,0.02))] p-5 shadow-[0_20px_60px_rgba(0,0,0,0.28)] backdrop-blur-sm">
      <p className="text-xs uppercase tracking-[0.28em] text-[#d9cdb8]/60">{props.title}</p>
      <p className="mt-3 text-3xl font-semibold text-[#f7f2e8]">{props.value}</p>
      <p className="mt-2 text-sm text-[#cbbda5]/72">{props.note}</p>
    </article>
  );
}

export function OverviewCards({ overview }: OverviewCardsProps) {
  return (
    <section className="grid gap-4 md:grid-cols-2 2xl:grid-cols-4">
      <OverviewCard
        title="总 Token"
        value={overview ? formatMetricValue(overview.totalTokens, 'number') : '--'}
        note="把消耗先看透。"
      />
      <OverviewCard
        title="总费用"
        value={overview ? formatMetricValue(overview.totalCostUsd, 'currency') : '--'}
        note="钱花在哪，一眼看见。"
      />
      <OverviewCard
        title="总请求数"
        value={overview ? formatMetricValue(overview.totalRequests, 'number') : '--'}
        note="请求热度先立住。"
      />
      <OverviewCard
        title="活跃客户端"
        value={overview ? formatMetricValue(overview.activeClients, 'number') : '--'}
        note="谁在跑，一目了然。"
      />
    </section>
  );
}
