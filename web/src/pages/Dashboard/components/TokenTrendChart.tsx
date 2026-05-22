import type { DashboardInterval, DashboardTrendPoint } from '@/types/dashboard';
import { ChartViewport } from '@/pages/Dashboard/components/ChartViewport';
import { formatBucketLabel } from '@/utils/format';
import { Area, AreaChart, CartesianGrid, Tooltip, XAxis, YAxis } from 'recharts';

interface TokenTrendChartProps {
  trend: DashboardTrendPoint[];
  interval: DashboardInterval;
}

export function TokenTrendChart({ trend, interval }: TokenTrendChartProps) {
  if (trend.length === 0) {
    return (
      <div className="grid h-[300px] place-items-center rounded-[28px] border border-dashed border-[#cfe0f1] bg-[linear-gradient(145deg,rgba(255,255,255,0.72),rgba(232,243,252,0.72))] text-sm text-[#6a86a3]">
        还没有 token 趋势。
      </div>
    );
  }

  return (
    <div className="h-[300px] min-w-0 rounded-[28px] border border-white/80 bg-[linear-gradient(145deg,rgba(255,255,255,0.8),rgba(237,246,252,0.96))] p-4 shadow-[inset_0_1px_0_rgba(255,255,255,0.72)]">
      <ChartViewport>
        {({ width, height }) => (
        <AreaChart width={width} height={height} data={trend}>
          <CartesianGrid stroke="rgba(120,155,193,0.18)" vertical={false} />
          <XAxis
            dataKey="bucket"
            tickFormatter={(value) => formatBucketLabel(value, interval)}
            stroke="#7191b0"
            tickLine={false}
            axisLine={false}
          />
          <YAxis stroke="#7191b0" tickLine={false} axisLine={false} width={56} />
          <Tooltip
            contentStyle={{
              background: 'rgba(245, 250, 255, 0.96)',
              border: '1px solid rgba(138, 176, 214, 0.36)',
              borderRadius: '18px',
              color: '#17324b',
              boxShadow: '0 18px 42px rgba(104, 153, 204, 0.18)',
            }}
            labelFormatter={(label) => formatBucketLabel(label, interval)}
          />
          <Area type="monotone" dataKey="inputTokens" stackId="tokens" stroke="#3f8cff" fill="#7cc7ff" fillOpacity={0.65} />
          <Area type="monotone" dataKey="outputTokens" stackId="tokens" stroke="#61a9ff" fill="#5c9bff" fillOpacity={0.58} />
          <Area type="monotone" dataKey="cacheReadTokens" stackId="tokens" stroke="#57c3ca" fill="#67d4d7" fillOpacity={0.5} />
          <Area type="monotone" dataKey="cacheCreationTokens" stackId="tokens" stroke="#8aaeff" fill="#9fbfff" fillOpacity={0.44} />
        </AreaChart>
        )}
      </ChartViewport>
    </div>
  );
}
