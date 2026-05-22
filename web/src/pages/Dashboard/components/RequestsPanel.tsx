interface RequestRow {
  time: string;
  model: string;
  input: string;
  output: string;
  cost: string;
  client: string;
}

interface RequestsPanelProps {
  items: readonly RequestRow[];
}

const columns = ['时间', '模型', '输入', '输出', '费用（USD）', '客户端'] as const;

export function RequestsPanel({ items }: RequestsPanelProps) {
  return (
    <section className="flex h-full min-h-0 flex-col rounded-[18px] border border-[#25314d] bg-[linear-gradient(180deg,rgba(7,17,40,0.96),rgba(6,14,34,0.95))] px-[16px] py-[14px] shadow-[0_16px_42px_rgba(0,0,0,0.34)]">
      <div className="mb-4 flex items-center justify-between gap-3">
        <div className="flex items-center gap-2 text-[14px] font-semibold text-white">
          <span className="inline-block h-4 w-1 rounded-full bg-[linear-gradient(180deg,#c27cff,#7c8fff)]" />
          <h2>最近请求</h2>
        </div>
        <button
          type="button"
          className="cursor-pointer rounded-[8px] border border-white/8 bg-white/4 px-[10px] py-[5px] text-[10px] text-[#a6b3d1] transition duration-200 ease-out hover:bg-white/8"
        >
          查看全部
        </button>
      </div>
      <div className="min-h-0 flex-1 overflow-hidden rounded-[13px] border border-white/6 bg-[linear-gradient(180deg,rgba(17,28,58,0.92),rgba(10,19,43,0.96))]">
        <table className="w-full border-collapse text-left text-[10px] text-[#bfd0ee]">
          <thead>
            <tr className="border-b border-white/6 text-[#8fa4cc]">
              {columns.map((column) => (
                <th key={column} className="px-4 py-[11px] font-medium tracking-[0.01em]">
                  {column}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {items.map((item) => (
              <tr key={`${item.time}-${item.client}`} className="border-b border-white/4 transition-colors duration-200 ease-out hover:bg-white/[0.02] last:border-b-0">
                <td className="px-4 py-[10px] text-[#90a4cc]">{item.time}</td>
                <td className="max-w-[210px] truncate px-4 py-[10px] text-[#d9e5ff]">{item.model}</td>
                <td className="font-mono-display px-4 py-[10px] text-[#c7d4f2]">{item.input}</td>
                <td className="font-mono-display px-4 py-[10px] text-[#c7d4f2]">{item.output}</td>
                <td className="font-mono-display px-4 py-[10px] text-[#36dca1]">{item.cost}</td>
                <td className="max-w-[200px] truncate px-4 py-[10px] text-[#8ea2ca]">{item.client}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </section>
  );
}
