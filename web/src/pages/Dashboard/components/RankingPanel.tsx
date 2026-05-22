interface RankingItem {
  rank: number;
  name: string;
  value: string;
  width: string;
}

interface RankingPanelProps {
  modelItems: readonly RankingItem[];
  clientItems: readonly RankingItem[];
}

function RankingBlock(props: {
  title: string;
  topLabel: string;
  items: readonly RankingItem[];
  barClassName: string;
  valueClassName: string;
}) {
  return (
    <section className="rounded-[13px] border border-white/6 bg-[linear-gradient(180deg,rgba(18,29,59,0.9),rgba(10,20,43,0.94))] px-[14px] py-[12px]">
      <div className="mb-4 flex items-center gap-3 text-[12px] font-semibold">
        <h3 className={props.valueClassName}>{props.title}</h3>
        <span className="text-[10px] text-[#8ea1c7]">{props.topLabel}</span>
      </div>
      <div className="space-y-[13px]">
        {props.items.map((item) => (
          <div key={`${props.title}-${item.rank}`} className="grid grid-cols-[14px_1fr] items-center gap-3 text-[10px] text-[#c4d0e9]">
            <span className="text-[#9ca9c7]">{item.rank}</span>
            <div>
              <div className="mb-[6px] flex items-center justify-between gap-3">
                <span className="truncate text-[#cfd9f3]">{item.name}</span>
                <span className="font-mono-display text-[#8ea1c7]">{item.value}</span>
              </div>
              <div className="h-[12px] rounded-full bg-[#0b142f]">
                <div className={`h-full rounded-full ${props.barClassName}`} style={{ width: item.width }} />
              </div>
            </div>
          </div>
        ))}
      </div>
    </section>
  );
}

export function RankingPanel({ modelItems, clientItems }: RankingPanelProps) {
  return (
    <div className="grid gap-3 md:grid-cols-2">
      <RankingBlock
        title="模型排行"
        topLabel="TOP 5"
        items={modelItems}
        valueClassName="text-[#be8cff]"
        barClassName="bg-[linear-gradient(90deg,#6a58ff,#b96efb)]"
      />
      <RankingBlock
        title="客户端排行"
        topLabel="TOP 5"
        items={clientItems}
        valueClassName="text-[#52b7ff]"
        barClassName="bg-[linear-gradient(90deg,#2a89d9,#53b7ff)]"
      />
    </div>
  );
}
