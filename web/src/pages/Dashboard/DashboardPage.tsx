import { dashboardReplicaMock } from '@/pages/Dashboard/dashboardReplica.mock';
import { BenefitPanel } from '@/pages/Dashboard/components/BenefitPanel';
import { DashboardFooterMeta } from '@/pages/Dashboard/components/DashboardFooterMeta';
import { DashboardHero } from '@/pages/Dashboard/components/DashboardHero';
import { MetricsPanel } from '@/pages/Dashboard/components/MetricsPanel';
import { RankingPanel } from '@/pages/Dashboard/components/RankingPanel';
import { RequestsPanel } from '@/pages/Dashboard/components/RequestsPanel';
import { TrendPanel } from '@/pages/Dashboard/components/TrendPanel';

export function DashboardPage() {
  const {
    updatedAtLabel,
    coreMetrics,
    benefitSummary,
    benefitMetrics,
    costTrend,
    tokenTrend,
    modelRanking,
    clientRanking,
    recentRequests,
    footerMeta,
  } = dashboardReplicaMock;

  return (
    <main className="h-screen overflow-hidden px-[20px] py-[16px] text-[#f5f7ff]">
      <div className="mx-auto grid h-[calc(100vh-32px)] w-[calc(100vw-40px)] max-w-[1920px] grid-rows-[auto_auto_1fr_auto] gap-[16px]">
        <DashboardHero updatedAtLabel={updatedAtLabel} />

        <div className="grid min-h-0 gap-[16px] xl:grid-cols-[1.06fr_0.94fr]">
          <section className="space-y-[16px]">
            <MetricsPanel items={coreMetrics} />
          </section>

          <aside className="space-y-[16px]">
            <BenefitPanel description={benefitSummary} items={benefitMetrics} />
          </aside>
        </div>

        <div className="grid min-h-0 gap-[16px] xl:grid-cols-[1.06fr_0.94fr]">
          <section className="min-h-0 rounded-[18px] border border-[#25314d] bg-[linear-gradient(180deg,rgba(7,17,40,0.96),rgba(6,14,34,0.95))] px-[14px] py-[12px] shadow-[0_16px_42px_rgba(0,0,0,0.34)]">
            <div className="mb-4 flex items-center gap-2 text-[14px] font-semibold text-white">
              <span className="inline-block h-4 w-1 rounded-full bg-[linear-gradient(180deg,#c27cff,#7c8fff)]" />
              <h2>使用趋势分析</h2>
            </div>
            <div className="grid h-[calc(100%-32px)] min-h-0 grid-rows-[1.02fr_0.98fr] gap-4">
              <TrendPanel costTrend={costTrend} tokenTrend={tokenTrend} />
              <RankingPanel modelItems={modelRanking} clientItems={clientRanking} />
            </div>
          </section>

          <aside className="min-h-0">
            <RequestsPanel items={recentRequests} />
          </aside>
        </div>

        <DashboardFooterMeta left={footerMeta.left} right={footerMeta.right} />
      </div>
    </main>
  );
}
