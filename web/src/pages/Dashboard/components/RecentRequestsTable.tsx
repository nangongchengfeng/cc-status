import type { RecentLogItem } from '@/types/logs';
import { formatMetricValue, formatRecentRequestTime, truncateLabel } from '@/utils/format';

interface RecentRequestsTableProps {
  items: RecentLogItem[];
}

export function RecentRequestsTable({ items }: RecentRequestsTableProps) {
  if (items.length === 0) {
    return (
      <div className="grid h-[320px] place-items-center rounded-[28px] border border-dashed border-[#cfe0f1] bg-[linear-gradient(145deg,rgba(255,255,255,0.72),rgba(232,243,252,0.72))] text-sm text-[#6a86a3]">
        当前时间范围还没有最近请求。
      </div>
    );
  }

  return (
    <div className="overflow-x-auto rounded-[28px] border border-white/80 bg-[linear-gradient(145deg,rgba(255,255,255,0.8),rgba(237,246,252,0.96))] shadow-[inset_0_1px_0_rgba(255,255,255,0.72)]">
      <table className="min-w-full text-left text-sm text-[#17324b]">
        <thead className="bg-white/70 text-xs uppercase tracking-[0.18em] text-[#6c92b4]">
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
            <tr
              key={item.id}
              className="border-t border-[#d8e7f4] transition-all duration-300 ease-[cubic-bezier(0.22,1,0.36,1)] hover:bg-[rgba(121,186,255,0.08)]"
            >
              <td className="whitespace-nowrap px-4 py-3 text-[#62819e]">{formatRecentRequestTime(item.createdAt)}</td>
              <td className="max-w-[240px] px-4 py-3">
                <span className="block truncate" title={item.model}>
                  {truncateLabel(item.model, 24)}
                </span>
              </td>
              <td className="whitespace-nowrap px-4 py-3">{formatMetricValue(item.inputTokens, 'number')}</td>
              <td className="whitespace-nowrap px-4 py-3">{formatMetricValue(item.outputTokens, 'number')}</td>
              <td className="whitespace-nowrap px-4 py-3 font-medium text-[#0e73c8]">{formatMetricValue(item.totalCostUsd, 'currency')}</td>
              <td className="max-w-[220px] px-4 py-3">
                <span className="block truncate text-[#62819e]" title={item.clientId}>
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
