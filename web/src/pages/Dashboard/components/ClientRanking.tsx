import type { DashboardTopClient } from '@/types/dashboard';
import { ChartViewport } from '@/pages/Dashboard/components/ChartViewport';
import { formatMetricValue, truncateLabel } from '@/utils/format';
import { Bar, BarChart, CartesianGrid, Tooltip, XAxis, YAxis } from 'recharts';

interface ClientRankingProps {
  items: DashboardTopClient[];
}

export function ClientRanking({ items }: ClientRankingProps) {
  if (items.length === 0) {
    return (
      <div className="grid h-[320px] place-items-center rounded-[28px] border border-dashed border-[#cfe0f1] bg-[linear-gradient(145deg,rgba(255,255,255,0.72),rgba(232,243,252,0.72))] text-sm text-[#6a86a3]">
        客户端排行还没数据。
      </div>
    );
  }

  const data = items.map((item) => ({
    name: item.clientId,
    totalCostUsd: Number(item.totalCostUsd),
  }));

  return (
    <div className="h-[320px] min-w-0 rounded-[28px] border border-white/80 bg-[linear-gradient(145deg,rgba(255,255,255,0.8),rgba(237,246,252,0.96))] p-4 shadow-[inset_0_1px_0_rgba(255,255,255,0.72)]">
      <ChartViewport>
        {({ width, height }) => (
        <BarChart width={width} height={height} data={data} layout="vertical" margin={{ left: 8, right: 8 }}>
          <CartesianGrid stroke="rgba(120,155,193,0.18)" horizontal={false} />
          <XAxis type="number" stroke="#7191b0" tickLine={false} axisLine={false} />
          <YAxis dataKey="name" type="category" stroke="#7191b0" tickFormatter={(value) => truncateLabel(String(value), 14)} tickLine={false} axisLine={false} width={118} />
          <Tooltip
            contentStyle={{
              background: 'rgba(245, 250, 255, 0.96)',
              border: '1px solid rgba(138, 176, 214, 0.36)',
              borderRadius: '18px',
              color: '#17324b',
              boxShadow: '0 18px 42px rgba(104, 153, 204, 0.18)',
            }}
            formatter={(value) => formatMetricValue(String(value ?? 0), 'currency')}
          />
          <Bar dataKey="totalCostUsd" radius={[0, 12, 12, 0]} fill="#69c5d0" />
        </BarChart>
        )}
      </ChartViewport>
    </div>
  );
}
