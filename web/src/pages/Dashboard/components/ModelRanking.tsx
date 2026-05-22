import type { DashboardTopModel } from '@/types/dashboard';
import { formatMetricValue, getModelDisplayName } from '@/utils/format';
import { Bar, BarChart, CartesianGrid, ResponsiveContainer, Tooltip, XAxis, YAxis } from 'recharts';

interface ModelRankingProps {
  items: DashboardTopModel[];
}

export function ModelRanking({ items }: ModelRankingProps) {
  if (items.length === 0) {
    return <div className="grid h-[320px] place-items-center rounded-[24px] border border-dashed border-white/15 bg-white/[0.03] text-sm text-[#cab99d]">模型排行还没数据。</div>;
  }

  const data = items.map((item) => ({
    name: getModelDisplayName({ displayName: item.displayName, model: item.model }),
    totalTokens: item.totalTokens,
  }));

  return (
    <div className="h-[320px] rounded-[24px] border border-white/10 bg-black/10 p-4">
      <ResponsiveContainer width="100%" height="100%">
        <BarChart data={data} layout="vertical" margin={{ left: 8, right: 8 }}>
          <CartesianGrid stroke="rgba(255,255,255,0.06)" horizontal={false} />
          <XAxis type="number" stroke="#d4c5a8" tickLine={false} axisLine={false} />
          <YAxis dataKey="name" type="category" stroke="#d4c5a8" tickLine={false} axisLine={false} width={118} />
          <Tooltip
            contentStyle={{ background: '#161515', border: '1px solid rgba(255,255,255,0.08)', borderRadius: '16px', color: '#fff7ea' }}
            formatter={(value) => formatMetricValue(Number(value ?? 0), 'number')}
          />
          <Bar dataKey="totalTokens" radius={[0, 10, 10, 0]} fill="#d49a4e" />
        </BarChart>
      </ResponsiveContainer>
    </div>
  );
}
