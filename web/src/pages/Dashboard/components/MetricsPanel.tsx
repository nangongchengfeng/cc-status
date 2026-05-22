type MetricAccent = 'purple' | 'amber' | 'sky' | 'emerald';

interface MetricItem {
  title: string;
  value: string;
  note: string;
  delta: string;
  accent: MetricAccent;
}

interface MetricsPanelProps {
  items: readonly MetricItem[];
}

const accentStyles: Record<MetricAccent, { glow: string; icon: string; delta: string; symbol: string }> = {
  purple: {
    glow: 'shadow-[0_0_22px_rgba(154,123,255,0.28)]',
    icon: 'bg-[radial-gradient(circle_at_35%_35%,rgba(196,157,255,0.95),rgba(106,79,223,0.9)_55%,rgba(58,44,117,0.95))]',
    delta: 'text-[#ff6b7a]',
    symbol: '◎',
  },
  amber: {
    glow: 'shadow-[0_0_22px_rgba(255,183,77,0.28)]',
    icon: 'bg-[radial-gradient(circle_at_35%_35%,rgba(255,214,120,0.98),rgba(255,166,48,0.92)_55%,rgba(115,64,21,0.96))]',
    delta: 'text-[#ff6b7a]',
    symbol: '$',
  },
  sky: {
    glow: 'shadow-[0_0_22px_rgba(69,202,255,0.28)]',
    icon: 'bg-[radial-gradient(circle_at_35%_35%,rgba(122,232,255,0.98),rgba(34,169,255,0.9)_55%,rgba(23,74,118,0.96))]',
    delta: 'text-[#57d8ff]',
    symbol: '∿',
  },
  emerald: {
    glow: 'shadow-[0_0_22px_rgba(65,223,164,0.28)]',
    icon: 'bg-[radial-gradient(circle_at_35%_35%,rgba(148,255,209,0.98),rgba(55,211,153,0.88)_55%,rgba(22,100,82,0.96))]',
    delta: 'text-[#4ce6a6]',
    symbol: '◉',
  },
};

export function MetricsPanel({ items }: MetricsPanelProps) {
  return (
    <section className="rounded-[18px] border border-[#25314d] bg-[linear-gradient(180deg,rgba(7,17,40,0.96),rgba(6,14,34,0.95))] px-[16px] py-[14px] shadow-[0_16px_42px_rgba(0,0,0,0.34)]">
      <div className="mb-4 flex items-center gap-2 text-[14px] font-semibold text-white">
        <span className="inline-block h-4 w-1 rounded-full bg-[linear-gradient(180deg,#c27cff,#7c8fff)]" />
        <h2>核心指标</h2>
      </div>
      <div className="grid gap-[12px] md:grid-cols-2 xl:grid-cols-4">
        {items.map((item) => {
          const style = accentStyles[item.accent];

          return (
            <article
              key={item.title}
              className="min-h-[92px] rounded-[13px] border border-white/6 bg-[linear-gradient(180deg,rgba(18,29,59,0.88),rgba(11,20,43,0.95))] px-4 py-[18px] shadow-[inset_0_1px_0_rgba(255,255,255,0.04)]"
            >
              <div className="flex items-center gap-3">
                <div
                  className={`grid h-10 w-10 place-items-center rounded-full text-[13px] font-semibold text-white ring-1 ring-white/10 ${style.icon} ${style.glow}`}
                >
                  {style.symbol}
                </div>
                <div>
                  <p className="text-[10px] leading-4 text-[#7e90b8]">{item.title}</p>
                  <p className="font-mono-display mt-[2px] text-[22px] font-semibold tracking-[-0.03em] text-white">{item.value}</p>
                </div>
              </div>
              <div className="mt-[10px] flex items-center gap-2 text-[10px]">
                <span className="text-[#6e7e9e]">{item.note}</span>
                <span className={style.delta}>{item.delta}</span>
              </div>
            </article>
          );
        })}
      </div>
    </section>
  );
}
