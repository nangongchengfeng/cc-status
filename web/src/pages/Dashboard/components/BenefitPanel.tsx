type BenefitAccent = 'amber' | 'emerald' | 'sky';

interface BenefitItem {
  title: string;
  value: string;
  note: string;
  delta: string;
  accent: BenefitAccent;
}

interface BenefitPanelProps {
  description: string;
  items: readonly BenefitItem[];
}

const valueClassByAccent: Record<BenefitAccent, string> = {
  amber: 'text-[#ffb642]',
  emerald: 'text-[#33e2a4]',
  sky: 'text-[#48a8ff]',
};

export function BenefitPanel({ description, items }: BenefitPanelProps) {
  return (
    <section className="rounded-[18px] border border-[#25314d] bg-[linear-gradient(180deg,rgba(7,17,40,0.96),rgba(6,14,34,0.95))] px-[16px] py-[14px] shadow-[0_16px_42px_rgba(0,0,0,0.34)]">
      <div className="mb-4 flex items-start gap-2">
        <span className="mt-1 inline-block h-4 w-1 rounded-full bg-[linear-gradient(180deg,#c27cff,#7c8fff)]" />
        <div>
          <h2 className="text-[14px] font-semibold text-white">综合效益</h2>
          <p className="mt-1 text-[10px] leading-5 text-[#91a2c6]">{description}</p>
        </div>
      </div>
      <div className="grid gap-[12px] md:grid-cols-3">
        {items.map((item) => (
          <article
            key={item.title}
            className="min-h-[92px] rounded-[13px] border border-white/6 bg-[linear-gradient(180deg,rgba(18,29,59,0.88),rgba(11,20,43,0.95))] px-4 py-[18px]"
          >
            <p className="text-[10px] leading-4 text-[#7b8eb8]">{item.title}</p>
            <p className={`font-mono-display mt-[6px] text-[22px] font-semibold tracking-[-0.03em] ${valueClassByAccent[item.accent]}`}>{item.value}</p>
            <div className="mt-[8px] flex items-center gap-2 text-[10px]">
              <span className="text-[#6e7e9e]">{item.note}</span>
              <span className="text-[#91a2c6]">{item.delta}</span>
            </div>
          </article>
        ))}
      </div>
    </section>
  );
}
