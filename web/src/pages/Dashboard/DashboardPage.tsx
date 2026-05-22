import { TIME_RANGE_OPTIONS } from '@/constants/timeRanges';
import { useDashboardQuery } from '@/hooks/useDashboardQuery';
import { useRecentLogsQuery } from '@/hooks/useRecentLogsQuery';
import { CacheAnalysis } from '@/pages/Dashboard/components/CacheAnalysis';
import { ClientRanking } from '@/pages/Dashboard/components/ClientRanking';
import { CostTrendChart } from '@/pages/Dashboard/components/CostTrendChart';
import { ModelRanking } from '@/pages/Dashboard/components/ModelRanking';
import { OverviewCards } from '@/pages/Dashboard/components/OverviewCards';
import { RecentRequestsTable } from '@/pages/Dashboard/components/RecentRequestsTable';
import { TokenTrendChart } from '@/pages/Dashboard/components/TokenTrendChart';
import { buildQueryTimeRange, getDashboardIntervalLabel, type TimeRangePreset } from '@/utils/timeRange';
import { useMemo, useState } from 'react';

export function DashboardPage() {
  const [preset, setPreset] = useState<TimeRangePreset>('last7Days');
  const queryRange = useMemo(() => buildQueryTimeRange(preset), [preset]);
  const dashboardQuery = useDashboardQuery(queryRange);
  const recentLogsQuery = useRecentLogsQuery({
    startAt: queryRange.startAt,
    endAt: queryRange.endAt,
    limit: 8,
  });

  const hasError = dashboardQuery.isError || recentLogsQuery.isError;
  const trend = dashboardQuery.data?.trend ?? [];
  const cacheAnalysis = dashboardQuery.data?.cacheAnalysis;
  const recentLogs = recentLogsQuery.data?.data ?? [];

  return (
    <main className="min-h-screen px-6 py-8 text-[#f7f2e8]">
      <div className="mx-auto grid max-w-[1680px] gap-6 xl:grid-cols-[1.45fr_0.85fr]">
        <section className="space-y-6">
          <header className="rounded-[36px] border border-white/10 bg-[linear-gradient(135deg,rgba(255,255,255,0.08),rgba(255,255,255,0.02))] p-8 shadow-[0_25px_90px_rgba(0,0,0,0.32)] backdrop-blur-md">
            <div className="flex flex-col gap-5 xl:flex-row xl:items-start xl:justify-between">
              <div>
                <p className="text-sm uppercase tracking-[0.35em] text-[#d8a978]">Claude Usage Dashboard</p>
                <h1 className="mt-4 text-5xl font-semibold leading-tight text-[#fff6e8]">Claude 用量看板</h1>
                <p className="mt-4 max-w-xl text-base leading-7 text-[#d7c8ae]">先把今天花在哪看明白</p>
              </div>

              <div className="max-w-[420px] rounded-[24px] border border-white/10 bg-black/20 p-4">
                <p className="text-xs uppercase tracking-[0.28em] text-[#d9cdb8]/60">时间范围</p>
                <div className="mt-4 flex flex-wrap gap-2">
                  {TIME_RANGE_OPTIONS.map((option) => (
                    <button
                      key={option.value}
                      className={[
                        'rounded-full border px-4 py-2 text-sm transition',
                        option.value === preset
                          ? 'border-[#d8a978] bg-[#d8a978] text-[#1b1612]'
                          : 'border-white/10 bg-white/5 text-[#f4e6cf]',
                      ].join(' ')}
                      onClick={() => setPreset(option.value)}
                      type="button"
                    >
                      {option.label}
                    </button>
                  ))}
                </div>
                <p className="mt-4 text-sm text-[#cab99d]">
                  当前按 {getDashboardIntervalLabel(queryRange.interval)} 粒度查询，起点 {queryRange.startAt}。
                </p>
              </div>
            </div>

            {hasError ? (
              <div className="mt-5 rounded-2xl border border-[#d36b4b]/40 bg-[#41261f]/60 px-4 py-3 text-sm text-[#ffd5c7]">
                数据暂时没接上，但页面状态是稳的。
              </div>
            ) : null}
          </header>

          <OverviewCards overview={dashboardQuery.data?.overview} />

          <section className="grid gap-6 2xl:grid-cols-[1.1fr_1fr]">
            <section className="rounded-[32px] border border-white/10 bg-black/20 p-6 backdrop-blur-sm">
              <div className="flex items-center justify-between">
                <div>
                  <h2 className="text-2xl font-semibold text-[#fff5e6]">费用趋势</h2>
                  <p className="mt-2 text-sm text-[#c9b89c]">共享时间轴，先看钱线怎么走。</p>
                </div>
                <div className="rounded-full border border-[#d8a978]/30 px-4 py-2 text-xs uppercase tracking-[0.25em] text-[#d8a978]">
                  cost trend
                </div>
              </div>
              <div className="mt-6">
                <CostTrendChart trend={trend} interval={queryRange.interval} />
              </div>
            </section>

            <section className="rounded-[32px] border border-white/10 bg-black/20 p-6 backdrop-blur-sm">
              <div className="flex items-center justify-between">
                <div>
                  <h2 className="text-2xl font-semibold text-[#fff5e6]">Token 细分</h2>
                  <p className="mt-2 text-sm text-[#c9b89c]">输入、输出、缓存，一根轴看完。</p>
                </div>
                <div className="rounded-full border border-[#63b59c]/30 px-4 py-2 text-xs uppercase tracking-[0.25em] text-[#63b59c]">
                  token trend
                </div>
              </div>
              <div className="mt-6">
                <TokenTrendChart trend={trend} interval={queryRange.interval} />
              </div>
            </section>
          </section>

          <section className="grid gap-6 2xl:grid-cols-[1fr_1fr]">
            <section className="rounded-[32px] border border-white/10 bg-black/20 p-6 backdrop-blur-sm">
              <div className="flex items-center justify-between">
                <div>
                  <h2 className="text-2xl font-semibold text-[#fff5e6]">模型排行</h2>
                  <p className="mt-2 text-sm text-[#c9b89c]">TOP 10 模型，谁最忙一眼看见。</p>
                </div>
              </div>
              <div className="mt-6">
                <ModelRanking items={dashboardQuery.data?.topModels ?? []} />
              </div>
            </section>

            <section className="rounded-[32px] border border-white/10 bg-black/20 p-6 backdrop-blur-sm">
              <div className="flex items-center justify-between">
                <div>
                  <h2 className="text-2xl font-semibold text-[#fff5e6]">客户端排行</h2>
                  <p className="mt-2 text-sm text-[#c9b89c]">TOP 10 客户端，谁最烧钱马上知道。</p>
                </div>
              </div>
              <div className="mt-6">
                <ClientRanking items={dashboardQuery.data?.topClients ?? []} />
              </div>
            </section>
          </section>
        </section>

        <aside className="grid gap-6">
          <section className="rounded-[32px] border border-white/10 bg-black/20 p-6 backdrop-blur-sm">
            <div className="flex items-center justify-between">
              <div>
                <h2 className="text-2xl font-semibold text-[#fff5e6]">缓存效益</h2>
                <p className="mt-2 text-sm leading-7 text-[#cab99d]">
                  节省、读取和建设三项成本，跟随当前时间范围一起刷新。
                </p>
              </div>
            </div>
            <div className="mt-6">
              <CacheAnalysis analysis={cacheAnalysis} />
            </div>
          </section>

          <section className="rounded-[32px] border border-white/10 bg-black/20 p-6 backdrop-blur-sm">
            <div className="flex items-center justify-between">
              <div>
                <h2 className="text-2xl font-semibold text-[#fff5e6]">最近请求</h2>
                <p className="mt-2 text-sm leading-7 text-[#cab99d]">
                  直接复用日志接口的最新排序结果。dashboard: {dashboardQuery.status} / logs: {recentLogsQuery.status}
                </p>
              </div>
            </div>
            <div className="mt-6">
              <RecentRequestsTable items={recentLogs} />
            </div>
          </section>
        </aside>
      </div>
    </main>
  );
}
