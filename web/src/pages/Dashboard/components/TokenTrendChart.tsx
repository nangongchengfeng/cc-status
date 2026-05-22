import type { DashboardInterval, DashboardTrendPoint } from '@/types/dashboard';
import { formatBucketLabel } from '@/utils/format';
import { Area, AreaChart, CartesianGrid, ResponsiveContainer, Tooltip, XAxis, YAxis } from 'recharts';

interface TokenTrendChartProps {
  trend: DashboardTrendPoint[];
  interval: DashboardInterval;
}

export function TokenTrendChart({ trend, interval }: TokenTrendChartProps) {
  if (trend.length === 0) {
    return <div className="grid h-[280px] place-items-center rounded-[24px] border border-dashed border-white/15 bg-white/[0.03] text-sm text-[#cab99d]">还没有 token 趋势。</div>;
  }

  return (
    <div className="h-[280px] rounded-[24px] border border-white/10 bg-black/10 p-4">
      <ResponsiveContainer width="100%" height="100%">
        <AreaChart data={trend}>
          <CartesianGrid stroke="rgba(255,255,255,0.08)" vertical={false} />
          <XAxis dataKey="bucket" tickFormatter={(value) => formatBucketLabel(value, interval)} stroke="#d4c5a8" tickLine={false} axisLine={false} />
          <YAxis stroke="#d4c5a8" tickLine={false} axisLine={false} width={56} />
          <Tooltip
            contentStyle={{ background: '#161515', border: '1px solid rgba(255,255,255,0.08)', borderRadius: '16px', color: '#fff7ea' }}
            labelFormatter={(label) => formatBucketLabel(label, interval)}
          />
          <Area type="monotone" dataKey="inputTokens" stackId="tokens" stroke="#d49a4e" fill="#d49a4e" fillOpacity={0.7} />
          <Area type="monotone" dataKey="outputTokens" stackId="tokens" stroke="#63b59c" fill="#63b59c" fillOpacity={0.65} />
          <Area type="monotone" dataKey="cacheReadTokens" stackId="tokens" stroke="#d85f4d" fill="#d85f4d" fillOpacity={0.55} />
          <Area type="monotone" dataKey="cacheCreationTokens" stackId="tokens" stroke="#a78a6f" fill="#a78a6f" fillOpacity={0.45} />
        </AreaChart>
      </ResponsiveContainer>
    </div>
  );
}
