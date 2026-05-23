import type { DashboardInterval, DashboardTrendPoint } from '@/types/dashboard';
import { ChartViewport } from '@/pages/Dashboard/components/ChartViewport';
import { formatBucketLabel, formatMetricValue, getModelDisplayName } from '@/utils/format';
import { CartesianGrid, Line, LineChart, Tooltip, XAxis, YAxis, TooltipProps } from 'recharts';

interface CostTrendChartProps {
  trend: DashboardTrendPoint[];
  interval: DashboardInterval;
}

function CustomTooltip({ active, payload, label, interval }: any & { interval: DashboardInterval }) {
  if (!active || !payload || payload.length === 0) {
    return null;
  }

  const point = payload[0].payload as DashboardTrendPoint;

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
      <p style={{ fontWeight: 600, marginBottom: 8 }}>{formatBucketLabel(label, interval)}</p>
      <p style={{ marginBottom: 8 }}>
        <span style={{ fontWeight: 500 }}>总费用：</span>
        {formatMetricValue(point.totalCostUsd, 'currency')}
      </p>
      {point.modelCosts && point.modelCosts.length > 0 && (
        <div>
          <p style={{ fontWeight: 500, marginBottom: 4, fontSize: 12, opacity: 0.8 }}>各模型费用：</p>
          {point.modelCosts.map((modelCost) => (
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

export function CostTrendChart({ trend, interval }: CostTrendChartProps) {
  if (trend.length === 0) {
    return (
      <div className="grid h-[300px] place-items-center rounded-[28px] border border-dashed border-[#cfe0f1] bg-[linear-gradient(145deg,rgba(255,255,255,0.72),rgba(232,243,252,0.72))] text-sm text-[#6a86a3]">
        当前时间范围暂无费用数据。
      </div>
    );
  }

  return (
    <div className="h-[300px] min-w-0 rounded-[28px] border border-white/80 bg-[linear-gradient(145deg,rgba(255,255,255,0.8),rgba(237,246,252,0.96))] p-4 shadow-[inset_0_1px_0_rgba(255,255,255,0.72)]">
      <ChartViewport>
        {({ width, height }) => (
        <LineChart width={width} height={height} data={trend}>
          <CartesianGrid stroke="rgba(120,155,193,0.18)" vertical={false} />
          <XAxis
            dataKey="bucket"
            tickFormatter={(value) => formatBucketLabel(value, interval)}
            stroke="#7191b0"
            tickLine={false}
            axisLine={false}
          />
          <YAxis
            stroke="#7191b0"
            tickFormatter={(value) => `$${value}`}
            tickLine={false}
            axisLine={false}
            width={56}
          />
          <Tooltip content={<CustomTooltip interval={interval} />} />
          <Line
            type="monotone"
            dataKey="totalCostUsd"
            stroke="#3f8cff"
            strokeWidth={3}
            dot={false}
            activeDot={{ r: 4, strokeWidth: 0, fill: '#1f74ff' }}
          />
        </LineChart>
        )}
      </ChartViewport>
    </div>
  );
}
