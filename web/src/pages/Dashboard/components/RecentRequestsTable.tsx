import type { RecentLogItem } from '@/types/logs';
import { formatMetricValue, formatRecentRequestTime, truncateLabel } from '@/utils/format';

interface RecentRequestsTableProps {
  items: RecentLogItem[];
}

export function RecentRequestsTable({ items }: RecentRequestsTableProps) {
  if (items.length === 0) {
    return <div className="grid h-[320px] place-items-center rounded-[24px] border border-dashed border-white/15 bg-white/[0.03] text-sm text-[#cab99d]">当前时间范围还没有最近请求。</div>;
  }

  return (
    <div className="overflow-x-auto rounded-[24px] border border-white/10 bg-black/10">
      <table className="min-w-full text-left text-sm text-[#f3e9d7]">
        <thead className="bg-white/[0.04] text-xs uppercase tracking-[0.18em] text-[#d4c5a8]/80">
          <tr>
            <th className="px-4 py-3 font-medium">时间</th>
            <th className="px-4 py-3 font-medium">模型</th>
            <th className="px-4 py-3 font-medium">输入</th>
            <th className="px-4 py-3 font-medium">输出</th>
            <th className="px-4 py-3 font-medium">总费用</th>
            <th className="px-4 py-3 font-medium">客户端</th>
          </tr>
        </thead>
        <tbody>
          {items.map((item) => (
            <tr key={item.id} className="border-t border-white/8">
              <td className="whitespace-nowrap px-4 py-3 text-[#d8ccb8]">{formatRecentRequestTime(item.createdAt)}</td>
              <td className="max-w-[240px] px-4 py-3">
                <span className="block truncate" title={item.model}>
                  {truncateLabel(item.model, 24)}
                </span>
              </td>
              <td className="whitespace-nowrap px-4 py-3">{formatMetricValue(item.inputTokens, 'number')}</td>
              <td className="whitespace-nowrap px-4 py-3">{formatMetricValue(item.outputTokens, 'number')}</td>
              <td className="whitespace-nowrap px-4 py-3 text-[#63b59c]">{formatMetricValue(item.totalCostUsd, 'currency')}</td>
              <td className="max-w-[220px] px-4 py-3">
                <span className="block truncate text-[#d8ccb8]" title={item.clientId}>
                  {truncateLabel(item.clientId, 20)}
                </span>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
