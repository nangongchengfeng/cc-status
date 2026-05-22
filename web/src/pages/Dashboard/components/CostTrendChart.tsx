import type { DashboardInterval, DashboardTrendPoint } from '@/types/dashboard';
import { formatBucketLabel, formatMetricValue } from '@/utils/format';
import { CartesianGrid, Line, LineChart, ResponsiveContainer, Tooltip, XAxis, YAxis } from 'recharts';

interface CostTrendChartProps {
  trend: DashboardTrendPoint[];
  interval: DashboardInterval;
}

export function CostTrendChart({ trend, interval }: CostTrendChartProps) {
  if (trend.length === 0) {
    return <div className="grid h-[280px] place-items-center rounded-[24px] border border-dashed border-white/15 bg-white/[0.03] text-sm text-[#cab99d]">这段时间还没花钱。</div>;
  }

  return (
    <div className="h-[280px] rounded-[24px] border border-white/10 bg-black/10 p-4">
      <ResponsiveContainer width="100%" height="100%">
        <LineChart data={trend}>
          <CartesianGrid stroke="rgba(255,255,255,0.08)" vertical={false} />
          <XAxis dataKey="bucket" tickFormatter={(value) => formatBucketLabel(value, interval)} stroke="#d4c5a8" tickLine={false} axisLine={false} />
          <YAxis stroke="#d4c5a8" tickFormatter={(value) => `$${value}`} tickLine={false} axisLine={false} width={56} />
          <Tooltip
            contentStyle={{ background: '#161515', border: '1px solid rgba(255,255,255,0.08)', borderRadius: '16px', color: '#fff7ea' }}
            formatter={(value) => formatMetricValue(String(value ?? 0), 'currency')}
            labelFormatter={(label) => formatBucketLabel(label, interval)}
          />
          <Line type="monotone" dataKey="totalCostUsd" stroke="#e49a61" strokeWidth={3} dot={false} />
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
}
