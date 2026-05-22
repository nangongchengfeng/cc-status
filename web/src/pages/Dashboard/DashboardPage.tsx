import { TIME_RANGE_OPTIONS } from '@/constants/timeRanges';
import { useDashboardQuery } from '@/hooks/useDashboardQuery';
import { useRecentLogsQuery } from '@/hooks/useRecentLogsQuery';
import { buildQueryTimeRange, getDashboardIntervalLabel, type TimeRangePreset } from '@/utils/timeRange';
import { useMemo, useState } from 'react';

function MetricCard(props: { title: string; value: string; note: string }) {
  return (
    <article className="rounded-[28px] border border-white/10 bg-black/20 p-5 shadow-[0_18px_60px_rgba(0,0,0,0.25)] backdrop-blur-sm">
      <p className="text-xs uppercase tracking-[0.28em] text-[#d9cdb8]/60">{props.title}</p>
      <p className="mt-3 text-3xl font-semibold text-[#f7f2e8]">{props.value}</p>
      <p className="mt-2 text-sm text-[#cbbda5]/72">{props.note}</p>
    </article>
  );
}

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

          <section className="grid gap-4 md:grid-cols-2 2xl:grid-cols-4">
            <MetricCard title="总 Token" value="--" note="等数据接进来" />
            <MetricCard title="总费用" value="--" note="别急，还没连后端" />
            <MetricCard title="总请求数" value="--" note="先把骨架立住" />
            <MetricCard title="活跃客户端" value="--" note="后面一起刷新" />
          </section>

          <section className="rounded-[32px] border border-white/10 bg-black/20 p-6 backdrop-blur-sm">
            <div className="flex items-center justify-between">
              <div>
                <h2 className="text-2xl font-semibold text-[#fff5e6]">主图区域</h2>
                <p className="mt-2 text-sm text-[#c9b89c]">这里先占坑，下一刀接趋势图。</p>
              </div>
              <div className="rounded-full border border-[#d8a978]/30 px-4 py-2 text-xs uppercase tracking-[0.25em] text-[#d8a978]">
                issue 02
              </div>
            </div>
            <div className="mt-6 h-[280px] rounded-[24px] border border-dashed border-white/15 bg-white/[0.03]" />
          </section>
        </section>

        <aside className="grid gap-6">
          <section className="rounded-[32px] border border-white/10 bg-black/20 p-6 backdrop-blur-sm">
            <h2 className="text-2xl font-semibold text-[#fff5e6]">查询状态</h2>
            <p className="mt-3 text-sm leading-7 text-[#cab99d]">
              dashboard: {dashboardQuery.status} / logs: {recentLogsQuery.status}
            </p>
          </section>
          <section className="rounded-[32px] border border-white/10 bg-black/20 p-6 backdrop-blur-sm">
            <h2 className="text-2xl font-semibold text-[#fff5e6]">卡片容器</h2>
            <p className="mt-3 text-sm leading-7 text-[#cab99d]">先把参数和请求走通，再慢慢塞图表。</p>
          </section>
        </aside>
      </div>
    </main>
  );
}
