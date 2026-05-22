import type { DashboardTopClient } from '@/types/dashboard';
import { formatMetricValue, truncateLabel } from '@/utils/format';
import { Bar, BarChart, CartesianGrid, ResponsiveContainer, Tooltip, XAxis, YAxis } from 'recharts';

interface ClientRankingProps {
  items: DashboardTopClient[];
}

export function ClientRanking({ items }: ClientRankingProps) {
  if (items.length === 0) {
    return <div className="grid h-[320px] place-items-center rounded-[24px] border border-dashed border-white/15 bg-white/[0.03] text-sm text-[#cab99d]">客户端排行还没数据。</div>;
  }

  const data = items.map((item) => ({
    name: item.clientId,
    totalCostUsd: Number(item.totalCostUsd),
  }));

  return (
    <div className="h-[320px] rounded-[24px] border border-white/10 bg-black/10 p-4">
      <ResponsiveContainer width="100%" height="100%">
        <BarChart data={data} layout="vertical" margin={{ left: 8, right: 8 }}>
          <CartesianGrid stroke="rgba(255,255,255,0.06)" horizontal={false} />
          <XAxis type="number" stroke="#d4c5a8" tickLine={false} axisLine={false} />
          <YAxis dataKey="name" type="category" stroke="#d4c5a8" tickFormatter={(value) => truncateLabel(String(value), 14)} tickLine={false} axisLine={false} width={118} />
          <Tooltip
            contentStyle={{ background: '#161515', border: '1px solid rgba(255,255,255,0.08)', borderRadius: '16px', color: '#fff7ea' }}
            formatter={(value) => formatMetricValue(String(value ?? 0), 'currency')}
          />
          <Bar dataKey="totalCostUsd" radius={[0, 10, 10, 0]} fill="#63b59c" />
        </BarChart>
      </ResponsiveContainer>
    </div>
  );
}
