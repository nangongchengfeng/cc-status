import type { DashboardTopClient } from '@/types/dashboard';
import { ChartViewport } from '@/pages/Dashboard/components/ChartViewport';
import { formatMetricValue, truncateLabel, getModelDisplayName } from '@/utils/format';
import { Bar, BarChart, CartesianGrid, Tooltip, XAxis, YAxis } from 'recharts';

interface ClientRankingProps {
  items: DashboardTopClient[];
}

function CustomTooltip({ active, payload }: any) {
  if (!active || !payload || payload.length === 0) {
    return null;
  }

  const clientData = payload[0].payload as { originalItem: DashboardTopClient };
  const item = clientData.originalItem;

  return (
    <div
      style={{
        background: 'rgba(245, 250, 255, 0.96)',
        border: '1px solid rgba(138, 176, 214, 0.36)',
        borderRadius: '18px',
        color: '#17324b',
        boxShadow: '0 18px 42px rgba(104, 153, 204, 0.18)',
        padding: '12px 16px',
      }}
    >
      <p style={{ fontWeight: 600, marginBottom: 8 }}>{item.clientId}</p>
      <p style={{ marginBottom: 8 }}>
        <span style={{ fontWeight: 500 }}>总费用：</span>
        {formatMetricValue(item.totalCostUsd, 'currency')}
      </p>
      {item.modelCosts && item.modelCosts.length > 0 && (
        <div>
          <p style={{ fontWeight: 500, marginBottom: 4, fontSize: 12, opacity: 0.8 }}>各模型费用：</p>
          {item.modelCosts.map((modelCost) => (
            <p key={modelCost.model} style={{ fontSize: 12, margin: 2 }}>
              {getModelDisplayName({ displayName: modelCost.displayName, model: modelCost.model })}：
              {formatMetricValue(modelCost.costUsd, 'currency')}
            </p>
          ))}
        </div>
      )}
    </div>
  );
}

export function ClientRanking({ items }: ClientRankingProps) {
  if (items.length === 0) {
    return (
      <div className="grid h-[320px] place-items-center rounded-[28px] border border-dashed border-[#cfe0f1] bg-[linear-gradient(145deg,rgba(255,255,255,0.72),rgba(232,243,252,0.72))] text-sm text-[#6a86a3]">
        当前时间范围暂无客户端排行数据。
      </div>
    );
  }

  const data = items.map((item) => ({
    name: item.clientId,
    totalCostUsd: Number(item.totalCostUsd),
    originalItem: item,
  }));

  return (
    <div className="h-[320px] min-w-0 rounded-[28px] border border-white/80 bg-[linear-gradient(145deg,rgba(255,255,255,0.8),rgba(237,246,252,0.96))] p-4 shadow-[inset_0_1px_0_rgba(255,255,255,0.72)]">
      <ChartViewport>
        {({ width, height }) => (
          <BarChart width={width} height={height} data={data} layout="vertical" margin={{ left: 8, right: 8 }}>
            <CartesianGrid stroke="rgba(120,155,193,0.18)" horizontal={false} />
            <XAxis type="number" stroke="#7191b0" tickLine={false} axisLine={false} />
            <YAxis dataKey="name" type="category" stroke="#7191b0" tickFormatter={(value) => truncateLabel(String(value), 14)} tickLine={false} axisLine={false} width={118} />
            <Tooltip content={<CustomTooltip />} />
            <Bar dataKey="totalCostUsd" radius={[0, 12, 12, 0]} fill="#69c5d0" />
          </BarChart>
        )}
      </ChartViewport>
    </div>
  );
}
