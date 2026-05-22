import { useMemo, useState } from 'react';

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
import { formatUnixTimestamp } from '@/utils/format';
import { buildQueryTimeRange, getDashboardIntervalLabel, type TimeRangePreset } from '@/utils/timeRange';

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
  const isLoading = dashboardQuery.isPending || recentLogsQuery.isPending;
  const trend = dashboardQuery.data?.trend ?? [];
  const cacheAnalysis = dashboardQuery.data?.cacheAnalysis;
  const recentLogs = recentLogsQuery.data?.data ?? [];
  const intervalLabel = getDashboardIntervalLabel(queryRange.interval);
  const selectedRangeLabel = TIME_RANGE_OPTIONS.find((option) => option.value === preset)?.label ?? '最近 7 天';
  const statusTitle = hasError ? '数据异常' : isLoading ? '加载中' : '数据就绪';
  const statusNote = hasError ? '页面框架保持稳定。' : isLoading ? '数据正在加载中。' : '数据已准备就绪。';
  const overview = dashboardQuery.data?.overview;
  const previousOverview = dashboardQuery.data?.previousOverview;

  return (
    <main aria-busy={isLoading} className="min-h-screen px-6 py-6 text-[#18324a]">
      <div className="mx-auto grid max-w-[1720px] gap-6 xl:grid-cols-12">
        <header className="xl:col-span-8 overflow-hidden rounded-[40px] border border-white/75 bg-[linear-gradient(145deg,rgba(255,255,255,0.72),rgba(232,243,252,0.9))] p-8 shadow-[0_32px_120px_rgba(104,153,204,0.2)] backdrop-blur-xl">
          <div className="relative">
            <div className="pointer-events-none absolute inset-x-0 top-0 h-44 rounded-full bg-[radial-gradient(circle,rgba(102,187,255,0.22),transparent_68%)] blur-3xl" />
            <div className="relative flex flex-col gap-8">
              <div className="flex flex-wrap items-start justify-between gap-6">
                <div className="max-w-3xl">
                  <div className="inline-flex items-center rounded-full border border-[#cfe0f0] bg-white/55 px-4 py-2 text-xs uppercase tracking-[0.32em] text-[#4f86b7]">
                    Claude Usage Deck
                  </div>
                  <h1 className="mt-6 text-5xl font-semibold leading-[1.05] text-[#12304d] xl:text-6xl">Claude 用量监控</h1>
                  <p className="mt-4 max-w-2xl text-lg leading-8 text-[#52718f]">优先展示费用数据，其次为使用来源。</p>
                  <div className="mt-6 flex flex-wrap gap-3 text-sm text-[#3f6180]">
                    <span className="rounded-full border border-[#d6e6f4] bg-white/70 px-4 py-2 shadow-[0_12px_30px_rgba(124,164,208,0.12)]">
                      当前按{intervalLabel}粒度观察
                    </span>
                    <span className="rounded-full border border-[#d6e6f4] bg-white/70 px-4 py-2 shadow-[0_12px_30px_rgba(124,164,208,0.12)]">
                      范围：{selectedRangeLabel}
                    </span>
                    <span className="rounded-full border border-[#c9ddf4] bg-[linear-gradient(135deg,rgba(90,177,255,0.15),rgba(255,255,255,0.78))] px-4 py-2 text-[#1d5f92] shadow-[0_12px_30px_rgba(124,164,208,0.12)]">
                      费用数据为页面核心指标。
                    </span>
                  </div>
                </div>

                <div className="grid min-w-[280px] gap-4 sm:grid-cols-2 xl:max-w-[360px] xl:grid-cols-1">
                  <div className="rounded-[28px] border border-white/80 bg-[linear-gradient(145deg,rgba(72,150,255,0.14),rgba(255,255,255,0.86))] p-5 shadow-[0_16px_40px_rgba(112,151,194,0.16)]">
                    <p className="text-xs uppercase tracking-[0.28em] text-[#6b93b5]">状态快照</p>
                    <p className="mt-3 text-3xl font-semibold text-[#15324c]">{statusTitle}</p>
                    <p className="mt-2 text-sm text-[#60809f]">{statusNote}</p>
                  </div>
                  <div className="rounded-[28px] border border-white/80 bg-white/72 p-5 shadow-[0_16px_40px_rgba(112,151,194,0.16)]">
                    <p className="text-xs uppercase tracking-[0.28em] text-[#6b93b5]">当前视图</p>
                    <p className="mt-3 text-3xl font-semibold text-[#15324c]">{selectedRangeLabel}</p>
                    <p className="mt-2 text-sm text-[#60809f]">切换时间范围将刷新页面数据。</p>
                  </div>
                </div>
              </div>

              <div className="mt-4">
                <OverviewCards overview={overview} previousOverview={previousOverview} preset={preset} />
              </div>
            </div>
          </div>
        </header>

        <aside className="grid gap-4 xl:col-span-4">
          <section className="rounded-[36px] border border-white/75 bg-[linear-gradient(150deg,rgba(255,255,255,0.8),rgba(236,245,252,0.94))] p-5 shadow-[0_24px_90px_rgba(111,153,200,0.16)] backdrop-blur-xl">
            <div className="flex items-start justify-between gap-4">
              <div>
                <p className="text-xs uppercase tracking-[0.28em] text-[#6b93b5]">时间范围</p>
                <h2 className="mt-3 text-2xl font-semibold text-[#12304d]">选择时间范围，页面将同步更新</h2>
                <p className="mt-2 text-sm leading-7 text-[#5d7f9d]">起点 {formatUnixTimestamp(queryRange.startAt)}</p>
              </div>
              <div className="rounded-full border border-[#d6e5f4] bg-white/70 px-4 py-2 text-xs uppercase tracking-[0.28em] text-[#4e88bb]">
                {intervalLabel}
              </div>
            </div>

            <div className="mt-6 flex flex-wrap gap-3">
              {TIME_RANGE_OPTIONS.map((option) => (
                <button
                  key={option.value}
                  className={[
                    'rounded-full border px-4 py-2 text-sm shadow-[0_10px_24px_rgba(117,157,201,0.12)] transition-all duration-300 ease-[cubic-bezier(0.22,1,0.36,1)]',
                    option.value === preset
                      ? 'border-[#78b8ff] bg-[linear-gradient(135deg,#59b4ff,#d9f0ff)] text-[#10365b] -translate-y-0.5'
                      : 'border-[#d9e7f4] bg-white/72 text-[#4d6e8c] hover:-translate-y-0.5 hover:border-[#bddaf5] hover:bg-white',
                  ].join(' ')}
                  onClick={() => setPreset(option.value)}
                  type="button"
                >
                  {option.label}
                </button>
              ))}
            </div>
          </section>

          <section className="rounded-[36px] border border-white/75 bg-[linear-gradient(145deg,rgba(255,255,255,0.82),rgba(232,243,252,0.94))] p-5 shadow-[0_24px_90px_rgba(111,153,200,0.16)] backdrop-blur-xl">
            <p className="text-xs uppercase tracking-[0.28em] text-[#6b93b5]">系统状态</p>
            <div className="mt-4 rounded-[26px] border border-white/80 bg-white/68 p-5 shadow-[inset_0_1px_0_rgba(255,255,255,0.7)]">
              <p className="text-xs uppercase tracking-[0.24em] text-[#6b93b5]">页面状态</p>
              <p className="mt-3 text-3xl font-semibold text-[#15324c]">{statusTitle}</p>
              <p className="mt-2 text-sm text-[#60809f]">{statusNote}</p>
            </div>

            {hasError ? (
              <div className="mt-4 rounded-[24px] border border-[#f0c2b7] bg-[linear-gradient(145deg,rgba(255,255,255,0.8),rgba(255,235,230,0.9))] px-4 py-3 text-sm text-[#9b4c3b]">
                数据加载异常，页面状态保持稳定。
              </div>
            ) : null}

            {isLoading ? (
              <div className="mt-4 rounded-[24px] border border-[#cde2f6] bg-[linear-gradient(145deg,rgba(255,255,255,0.8),rgba(229,242,255,0.92))] px-4 py-3 text-sm text-[#2d638d]">
                数据加载中，正在刷新当前时间范围。
              </div>
            ) : null}

            <div className="mt-4 grid gap-4 sm:grid-cols-2">
              <div className="rounded-[24px] border border-white/80 bg-white/65 p-4">
                <p className="text-xs uppercase tracking-[0.24em] text-[#6b93b5]">数据口径</p>
                <p className="mt-2 text-lg font-semibold text-[#163553]">费用优先，用量次之</p>
              </div>
              <div className="rounded-[24px] border border-white/80 bg-white/65 p-4">
                <p className="text-xs uppercase tracking-[0.24em] text-[#6b93b5]">数据状态</p>
                <p className="mt-2 text-lg font-semibold text-[#163553]">
                  dashboard: {dashboardQuery.status} / logs: {recentLogsQuery.status}
                </p>
              </div>
            </div>
          </section>
        </aside>

        <section className="grid gap-6 xl:col-span-12 xl:grid-cols-12">
          <section className="min-w-0 rounded-[36px] border border-white/75 bg-[linear-gradient(145deg,rgba(255,255,255,0.82),rgba(232,243,252,0.94))] p-6 shadow-[0_24px_90px_rgba(111,153,200,0.16)] backdrop-blur-xl xl:col-span-7">
            <div className="flex items-start justify-between gap-4">
              <div>
                <p className="text-xs uppercase tracking-[0.28em] text-[#6b93b5]">Primary Trend</p>
                <h2 className="mt-3 text-3xl font-semibold text-[#12304d]">费用趋势</h2>
                <p className="mt-2 text-sm text-[#5d7f9d]">费用趋势一目了然。</p>
              </div>
              <div className="rounded-full border border-[#cfe2f4] bg-white/72 px-4 py-2 text-xs uppercase tracking-[0.28em] text-[#4e88bb]">
                cost first
              </div>
            </div>
            <div className="mt-6 min-w-0">
              <CostTrendChart trend={trend} interval={queryRange.interval} />
            </div>
          </section>

          <section className="min-w-0 rounded-[36px] border border-white/75 bg-[linear-gradient(145deg,rgba(255,255,255,0.82),rgba(232,243,252,0.94))] p-6 shadow-[0_24px_90px_rgba(111,153,200,0.16)] backdrop-blur-xl xl:col-span-5">
            <div className="flex items-start justify-between gap-4">
              <div>
                <p className="text-xs uppercase tracking-[0.28em] text-[#6b93b5]">Token Flow</p>
                <h2 className="mt-3 text-3xl font-semibold text-[#12304d]">Token 轨迹</h2>
                <p className="mt-2 text-sm text-[#5d7f9d]">Token 使用分布一目了然。</p>
              </div>
              <div className="rounded-full border border-[#cfe2f4] bg-white/72 px-4 py-2 text-xs uppercase tracking-[0.28em] text-[#4e88bb]">
                token flow
              </div>
            </div>
            <div className="mt-6 min-w-0">
              <TokenTrendChart trend={trend} interval={queryRange.interval} />
            </div>
          </section>
        </section>

        <section className="grid gap-6 xl:col-span-12 xl:grid-cols-12">
          <div className="grid min-w-0 gap-6 xl:col-span-5">
            <section className="min-w-0 rounded-[36px] border border-white/75 bg-[linear-gradient(145deg,rgba(255,255,255,0.82),rgba(232,243,252,0.94))] p-6 shadow-[0_24px_90px_rgba(111,153,200,0.16)] backdrop-blur-xl">
              <h2 className="text-3xl font-semibold text-[#12304d]">模型排行</h2>
              <p className="mt-2 text-sm text-[#5d7f9d]">模型使用排行一目了然。</p>
              <div className="mt-6 min-w-0">
                <ModelRanking items={dashboardQuery.data?.topModels ?? []} />
              </div>
            </section>

            <section className="min-w-0 rounded-[36px] border border-white/75 bg-[linear-gradient(145deg,rgba(255,255,255,0.82),rgba(232,243,252,0.94))] p-6 shadow-[0_24px_90px_rgba(111,153,200,0.16)] backdrop-blur-xl">
              <h2 className="text-3xl font-semibold text-[#12304d]">客户端排行</h2>
              <p className="mt-2 text-sm text-[#5d7f9d]">客户端费用排行一目了然。</p>
              <div className="mt-6 min-w-0">
                <ClientRanking items={dashboardQuery.data?.topClients ?? []} />
              </div>
            </section>
          </div>

          <div className="grid min-w-0 gap-6 xl:col-span-7">
            <section className="rounded-[36px] border border-white/75 bg-[linear-gradient(145deg,rgba(255,255,255,0.82),rgba(232,243,252,0.94))] p-5 shadow-[0_24px_90px_rgba(111,153,200,0.16)] backdrop-blur-xl">
              <div className="flex items-start justify-between gap-4">
                <div>
                  <h2 className="text-3xl font-semibold text-[#12304d]">缓存效益</h2>
                  <p className="mt-2 text-sm leading-7 text-[#5d7f9d]">优先展示缓存效益，其次为成本数据。</p>
                </div>
                <div className="rounded-full border border-[#cfe2f4] bg-white/72 px-4 py-2 text-xs uppercase tracking-[0.28em] text-[#4e88bb]">
                  cache
                </div>
              </div>
              <div className="mt-6">
                <CacheAnalysis analysis={cacheAnalysis} />
              </div>
            </section>

            <section className="rounded-[36px] border border-white/75 bg-[linear-gradient(145deg,rgba(255,255,255,0.82),rgba(232,243,252,0.94))] p-5 shadow-[0_24px_90px_rgba(111,153,200,0.16)] backdrop-blur-xl">
              <div className="flex items-start justify-between gap-4">
                <div>
                  <h2 className="text-3xl font-semibold text-[#12304d]">最近请求</h2>
                  <p className="mt-2 text-sm leading-7 text-[#5d7f9d]">展示最近八条请求记录。</p>
                </div>
                <div className="rounded-full border border-[#cfe2f4] bg-white/72 px-4 py-2 text-xs uppercase tracking-[0.28em] text-[#4e88bb]">
                  live feed
                </div>
              </div>
              <div className="mt-6">
                <RecentRequestsTable items={recentLogs} />
              </div>
            </section>
          </div>
        </section>
      </div>
    </main>
  );
}
