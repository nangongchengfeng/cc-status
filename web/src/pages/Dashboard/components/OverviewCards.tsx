import { formatMetricValue, formatNumberInWanYi } from '@/utils/format';
import type { DashboardOverview } from '@/types/dashboard';
import type { TimeRangePreset } from '@/utils/timeRange';

interface OverviewCardsProps {
  overview?: DashboardOverview;
  previousOverview?: DashboardOverview;
  preset: TimeRangePreset;
  selectedRangeLabel: string;
  statusTitle: string;
  statusNote: string;
}

function PrimaryMetricCard(props: { title: string; value: string; note: string; accent: string; change?: string; isPositive?: boolean; compareLabel?: string; hideArrow?: boolean; secondaryValue?: string; unit?: string }) {
  return (
    <article
      className={[
        'overflow-hidden rounded-[32px] border border-white/80 px-5 pt-4 pb-5 shadow-[0_22px_60px_rgba(111,153,200,0.16)] backdrop-blur-xl flex flex-col justify-between',
        props.accent,
      ].join(' ')}
    >
      <div>
        <p className="text-xs uppercase tracking-[0.3em] text-[#6c92b4]">{props.title}</p>
        <div className="mt-3 flex items-baseline gap-1">
          <p className="text-4xl font-semibold text-[#12304d] xl:text-[2.8rem]">{props.value}</p>
          {props.unit && (
            <p className="text-2xl font-semibold text-[#12304d] xl:text-[1.8rem]">{props.unit}</p>
          )}
        </div>
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

function SmallMetricCard(props: { title: string; value: string; note: string; change?: string; isPositive?: boolean; compareLabel?: string; secondaryValue?: string; unit?: string; accent?: string; compact?: boolean }) {
  return (
    <article className={[
        props.compact
          ? 'rounded-[20px] border border-white/80 bg-white/72 px-3 pt-2 pb-3 shadow-[0_12px_32px_rgba(111,153,200,0.14)] backdrop-blur-xl flex flex-col justify-between'
          : 'rounded-[28px] border border-white/80 bg-white/72 px-4 pt-3 pb-4 shadow-[0_18px_48px_rgba(111,153,200,0.14)] backdrop-blur-xl flex flex-col justify-between',
        props.accent,
      ].join(' ')}
    >
      <div>
        <p className={props.compact ? 'text-[10px] uppercase tracking-[0.24em] text-[#6c92b4]' : 'text-xs uppercase tracking-[0.28em] text-[#6c92b4]'}>{props.title}</p>
        <div className="mt-1 flex items-baseline gap-1">
          <p className={props.compact ? 'text-2xl font-semibold text-[#12304d]' : 'text-3xl font-semibold text-[#12304d]'}>{props.value}</p>
          {props.unit && (
            <p className={props.compact ? 'text-lg font-semibold text-[#12304d]' : 'text-xl font-semibold text-[#12304d]'}>{props.unit}</p>
          )}
        </div>
        {props.secondaryValue && (
          <p className="mt-1 text-xs text-[#5f7f9e]">{props.secondaryValue}</p>
        )}
      </div>
      {(props.note || (props.change && props.compareLabel)) && (
        <div className="mt-2">
          {props.note && <p className={props.compact ? 'text-xs text-[#5f7f9e]' : 'text-sm text-[#5f7f9e]'}>{props.note}</p>}
          {props.change && props.compareLabel && (
            <div className="mt-1">
              <p className="text-xs text-[#8ba5bf]">{props.compareLabel}</p>
              <span className={`text-sm font-medium ${props.isPositive ? 'text-[#16a34a]' : 'text-[#dc2626]'}`}>
                {props.isPositive ? '↓' : '↑'} {props.change}
              </span>
            </div>
          )}
        </div>
      )}
    </article>
  );
}

function InfoCard(props: { title: string; value: string; note: string; accent?: string; compact?: boolean; extraCompact?: boolean }) {
  return (
    <article className={[
        props.extraCompact
          ? 'rounded-[20px] border border-white/80 bg-white/72 px-3 py-8 shadow-[0_12px_32px_rgba(112,151,194,0.16)] backdrop-blur-xl'
          : props.compact
          ? 'rounded-[20px] border border-white/80 bg-white/72 px-3 pt-2 pb-3 shadow-[0_12px_32px_rgba(112,151,194,0.16)] backdrop-blur-xl flex flex-col justify-between'
          : 'rounded-[28px] border border-white/80 bg-white/72 px-4 pt-3 pb-4 shadow-[0_16px_40px_rgba(112,151,194,0.16)] backdrop-blur-xl flex flex-col justify-between',
        props.accent,
      ].join(' ')}
    >
      <p className={props.extraCompact ? 'text-[10px] uppercase tracking-[0.24em] text-[#6b93b5]' : props.compact ? 'text-[10px] uppercase tracking-[0.24em] text-[#6b93b5]' : 'text-xs uppercase tracking-[0.28em] text-[#6b93b5]'}>{props.title}</p>
      <p className={props.extraCompact ? 'mt-3 text-[2.6rem] font-semibold text-[#15324c]' : props.compact ? 'mt-1 text-2xl font-semibold text-[#15324c]' : 'mt-2 text-3xl font-semibold text-[#15324c]'}>{props.value}</p>
      <p className={props.extraCompact ? 'mt-3 text-xs text-[#60809f]' : props.compact ? 'mt-1 text-xs text-[#60809f]' : 'mt-2 text-sm text-[#60809f]'}>{props.note}</p>
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

// 优先用万/亿显示主数字，完整版作为辅助
function getDisplayValue(value: number): { main: string; unit?: string; secondary?: string } {
  const wanYi = formatNumberInWanYi(value);
  const full = formatMetricValue(value, 'number');
  if (wanYi) {
    return { main: wanYi.number, unit: wanYi.unit, secondary: full };
  }
  return { main: full };
}

// 计算缓存命中率: cache_read_tokens / (input_tokens + cache_read_tokens)
function calculateCacheHitRate(cacheRead: number, input: number): string | null {
  const denominator = input + cacheRead;
  if (denominator === 0) return null;
  const rate = (cacheRead / denominator) * 100;
  return `${rate.toFixed(1)}%`;
}

export function OverviewTopCards({ overview, previousOverview, preset, selectedRangeLabel, statusTitle, statusNote }: OverviewCardsProps) {
  return (
    <div className="grid gap-4 xl:grid-cols-[1fr_1fr] items-end">
      {/* 左侧标题区域 */}
      <div className="flex flex-col justify-center">
        <div className="inline-flex items-center rounded-full border border-[#cfe0f0] bg-white/55 px-4 py-2 text-xs uppercase tracking-[0.32em] text-[#4f86b7] w-fit">
          Claude Usage Deck
        </div>
        <h1 className="mt-6 text-5xl font-semibold leading-[1.05] text-[#12304d] xl:text-6xl">Claude 用量监控</h1>
        <p className="mt-4 max-w-2xl text-lg leading-8 text-[#52718f]">优先展示费用数据，其次为使用来源。</p>
        <div className="mt-6 flex flex-wrap gap-3 text-sm text-[#3f6180]">
          <span className="rounded-full border border-[#d6e6f4] bg-white/70 px-4 py-2 shadow-[0_12px_30px_rgba(124,164,208,0.12)]">
            当前按天粒度观察
          </span>
          <span className="rounded-full border border-[#d6e6f4] bg-white/70 px-4 py-2 shadow-[0_12px_30px_rgba(124,164,208,0.12)]">
            范围：{selectedRangeLabel}
          </span>
          <span className="rounded-full border border-[#c9ddf4] bg-[linear-gradient(135deg,rgba(90,177,255,0.15),rgba(255,255,255,0.78))] px-4 py-2 text-[#1d5f92] shadow-[0_12px_30px_rgba(124,164,208,0.12)]">
            费用数据为页面核心指标。
          </span>
        </div>
      </div>

      {/* 右侧 2x1 卡片网格 */}
      <div className="grid grid-cols-2 gap-4">
        <InfoCard
          extraCompact
          title="当前视图"
          value={selectedRangeLabel}
          note="切换时间范围将刷新页面数据。"
        />
        <InfoCard
          extraCompact
          title="状态快照"
          value={statusTitle}
          note={statusNote}
          accent="bg-[linear-gradient(145deg,rgba(72,150,255,0.14),rgba(255,255,255,0.86))]"
        />
      </div>
    </div>
  );
}

// 导出用于系统状态区域的请求和客户端卡片
export function RequestAndClientCards({ overview, previousOverview, preset }: { overview?: any; previousOverview?: any; preset: TimeRangePreset }) {
  const compareLabel = getCompareLabel(preset);

  const requestChange = overview && previousOverview
    ? calculateChange(overview.totalRequests, previousOverview.totalRequests)
    : null;

  const activeClientsChange = overview && previousOverview
    ? calculateChange(overview.activeClients, previousOverview.activeClients)
    : null;

  const activeClientsValue = overview ? formatMetricValue(overview.activeClients, 'number') : '--';
  const totalRequestsValue = overview ? formatMetricValue(overview.totalRequests, 'number') : '--';

  return (
    <div className="grid gap-4 sm:grid-cols-2">
      <SmallMetricCard
        title="总请求数"
        value={totalRequestsValue}
        note="请求总数统计。"
        change={requestChange?.change}
        isPositive={requestChange?.isPositive}
        compareLabel={requestChange ? compareLabel : undefined}
      />
      <SmallMetricCard
        title="活跃客户端"
        value={activeClientsValue}
        note="活跃客户端数量。"
        change={activeClientsChange?.change}
        isPositive={activeClientsChange?.isPositive}
        compareLabel={activeClientsChange ? compareLabel : undefined}
      />
    </div>
  );
}

export function OverviewMainCards({ overview, previousOverview, preset }: OverviewCardsProps) {
  const totalCost = overview ? formatMetricValue(overview.totalCostUsd, 'currency') : '--';

  const totalTokensDisplay = overview ? getDisplayValue(overview.totalTokens) : { main: '--' };
  const totalCacheTokensDisplay = overview ? getDisplayValue(overview.totalCacheTokens) : { main: '--' };
  const inputTokensDisplay = overview ? getDisplayValue(overview.inputTokens) : { main: '--' };
  const outputTokensDisplay = overview ? getDisplayValue(overview.outputTokens) : { main: '--' };

  const costChange = overview && previousOverview
    ? calculateChange(parseFloat(overview.totalCostUsd), parseFloat(previousOverview.totalCostUsd))
    : null;

  const tokenChange = overview && previousOverview
    ? calculateChange(overview.totalTokens, previousOverview.totalTokens)
    : null;

  const inputChange = overview && previousOverview
    ? calculateChange(overview.inputTokens, previousOverview.inputTokens)
    : null;

  const outputChange = overview && previousOverview
    ? calculateChange(overview.outputTokens, previousOverview.outputTokens)
    : null;

  const cacheHitRate = overview
    ? calculateCacheHitRate(overview.cacheReadTokens, overview.inputTokens)
    : null;

  const compareLabel = getCompareLabel(preset);

  return (
    <div className="grid gap-4 xl:grid-cols-[1fr_1fr_1fr_0.9fr]">
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
        value={totalTokensDisplay.main}
        unit={totalTokensDisplay.unit}
        secondaryValue={totalTokensDisplay.secondary}
        note="优先展示用量数据。"
        accent="bg-[linear-gradient(145deg,rgba(181,224,255,0.68),rgba(255,255,255,0.92))]"
        change={tokenChange?.change}
        isPositive={tokenChange?.isPositive}
        compareLabel={tokenChange ? compareLabel : undefined}
      />
      <PrimaryMetricCard
        title="总缓存"
        value={totalCacheTokensDisplay.main}
        unit={totalCacheTokensDisplay.unit}
        secondaryValue={totalCacheTokensDisplay.secondary}
        note="缓存 Token 统计。"
        accent="bg-[linear-gradient(145deg,rgba(134,239,172,0.35),rgba(255,255,255,0.9))]"
        change={cacheHitRate ?? undefined}
        isPositive={true}
        compareLabel={cacheHitRate ? '缓存命中率' : undefined}
        hideArrow={true}
      />
      <div className="grid gap-2">
        <SmallMetricCard
          compact
          title="输入 Token"
          value={inputTokensDisplay.main}
          unit={inputTokensDisplay.unit}
          secondaryValue={inputTokensDisplay.secondary}
          note=""
          change={inputChange?.change}
          isPositive={inputChange?.isPositive}
          compareLabel={inputChange ? compareLabel : undefined}
        />
        <SmallMetricCard
          compact
          title="输出 Token"
          value={outputTokensDisplay.main}
          unit={outputTokensDisplay.unit}
          secondaryValue={outputTokensDisplay.secondary}
          note=""
          change={outputChange?.change}
          isPositive={outputChange?.isPositive}
          compareLabel={outputChange ? compareLabel : undefined}
        />
      </div>
    </div>
  );
}
