interface DashboardHeroProps {
  updatedAtLabel: string;
}

const wavePoints = [
  { x: 220, y: 66, radius: 2.8, color: '#8b5cf6', opacity: 0.95 },
  { x: 260, y: 62, radius: 2.3, color: '#8b5cf6', opacity: 0.82 },
  { x: 308, y: 60, radius: 2.2, color: '#7c8fff', opacity: 0.7 },
  { x: 358, y: 57, radius: 1.8, color: '#6f7dff', opacity: 0.56 },
  { x: 410, y: 55, radius: 1.8, color: '#6c73ff', opacity: 0.42 },
  { x: 462, y: 54, radius: 1.9, color: '#ffbd74', opacity: 0.72 },
  { x: 514, y: 55, radius: 2.2, color: '#ffa94d', opacity: 0.75 },
  { x: 568, y: 59, radius: 2, color: '#ff9e4e', opacity: 0.7 },
  { x: 620, y: 63, radius: 2.3, color: '#ff8f52', opacity: 0.8 },
  { x: 672, y: 68, radius: 2.8, color: '#ff8448', opacity: 0.92 },
];

const farWavePoints = [
  { x: 180, y: 72, radius: 1.2, color: '#8f74ff', opacity: 0.56 },
  { x: 232, y: 67, radius: 1.1, color: '#866eff', opacity: 0.48 },
  { x: 286, y: 63, radius: 1, color: '#7d82ff', opacity: 0.42 },
  { x: 344, y: 61, radius: 0.9, color: '#7a7dff', opacity: 0.36 },
  { x: 404, y: 60, radius: 0.8, color: '#7773ff', opacity: 0.3 },
  { x: 468, y: 61, radius: 0.9, color: '#ffb86f', opacity: 0.46 },
  { x: 530, y: 64, radius: 1, color: '#ffab60', opacity: 0.5 },
  { x: 596, y: 69, radius: 1.1, color: '#ff9e56', opacity: 0.54 },
  { x: 662, y: 75, radius: 1.2, color: '#ff9150', opacity: 0.58 },
];

export function DashboardHero({ updatedAtLabel }: DashboardHeroProps) {
  return (
    <header className="relative min-h-[166px] overflow-hidden rounded-[20px] border border-[#23314e] bg-[linear-gradient(180deg,rgba(6,15,36,0.98),rgba(5,12,30,0.96))] px-[22px] pb-[16px] pt-[18px] shadow-[0_14px_36px_rgba(0,0,0,0.28)]">
      <div className="pointer-events-none absolute inset-x-0 top-0 h-full bg-[radial-gradient(circle_at_61%_4%,rgba(93,118,255,0.16),transparent_15%),radial-gradient(circle_at_67%_6%,rgba(255,149,77,0.22),transparent_9%)]" />
      <div className="pointer-events-none absolute left-[56.9%] top-[-8px] h-[98px] w-[98px] rounded-full bg-[radial-gradient(circle_at_36%_32%,rgba(255,164,104,0.98),rgba(88,88,255,0.92)_40%,rgba(14,17,43,0.98)_74%)] shadow-[0_0_50px_rgba(102,112,255,0.52)]" />
      <div className="pointer-events-none absolute left-[53.5%] top-[2px] h-[106px] w-[190px] rounded-full bg-[radial-gradient(circle_at_40%_35%,rgba(255,162,74,0.16),transparent_52%)] blur-2xl" />
      <svg className="pointer-events-none absolute inset-x-[24%] top-[18px] h-[96px] w-[68%] opacity-95" viewBox="0 0 720 110" preserveAspectRatio="none" aria-hidden="true">
        <defs>
          <linearGradient id="hero-wave-gradient" x1="0%" y1="0%" x2="100%" y2="0%">
            <stop offset="0%" stopColor="#8b5cf6" stopOpacity="0.8" />
            <stop offset="42%" stopColor="#6f7dff" stopOpacity="0.22" />
            <stop offset="68%" stopColor="#ffb35c" stopOpacity="0.62" />
            <stop offset="100%" stopColor="#ff8d50" stopOpacity="0.92" />
          </linearGradient>
          <linearGradient id="hero-wave-far-gradient" x1="0%" y1="0%" x2="100%" y2="0%">
            <stop offset="0%" stopColor="#8b5cf6" stopOpacity="0.28" />
            <stop offset="50%" stopColor="#6f7dff" stopOpacity="0.16" />
            <stop offset="100%" stopColor="#ff9d60" stopOpacity="0.34" />
          </linearGradient>
          <filter id="hero-glow">
            <feGaussianBlur stdDeviation="3" result="coloredBlur" />
            <feMerge>
              <feMergeNode in="coloredBlur" />
              <feMergeNode in="SourceGraphic" />
            </feMerge>
          </filter>
        </defs>
        <path d="M18,67 C92,45 170,44 246,57 C314,69 396,67 482,52 C556,39 627,44 700,66" fill="none" stroke="url(#hero-wave-far-gradient)" strokeWidth="1.2" opacity="0.8" />
        <path d="M24,74 C95,54 165,53 238,64 C307,75 390,73 478,57 C556,43 628,48 704,72" fill="none" stroke="url(#hero-wave-gradient)" strokeWidth="1.6" opacity="0.94" filter="url(#hero-glow)" />
        {farWavePoints.map((point) => (
          <circle
            key={`far-${point.x}-${point.y}`}
            cx={point.x}
            cy={point.y}
            r={point.radius}
            fill={point.color}
            fillOpacity={point.opacity}
          />
        ))}
        {wavePoints.map((point) => (
          <circle
            key={`near-${point.x}-${point.y}`}
            cx={point.x}
            cy={point.y}
            r={point.radius}
            fill={point.color}
            fillOpacity={point.opacity}
          />
        ))}
      </svg>

      <div className="relative flex items-start justify-between gap-6">
        <div>
          <h1 className="text-[36px] font-semibold tracking-[-0.05em] text-white">
            <span className="bg-[linear-gradient(90deg,#c576ff,#ff8eb1_42%,#ffffff_66%)] bg-clip-text text-transparent">Claude</span>{' '}
            用量看板
          </h1>
          <p className="mt-[8px] text-[12px] tracking-[0.01em] text-[#9fb0d4]">洞察使用趋势，优化资源配置，驱动智能增长</p>
        </div>
        <div className="mt-1 flex items-center gap-2 rounded-full border border-white/8 bg-[#101a34]/90 px-4 py-[7px] text-[10px] text-[#aab7d5] shadow-[0_8px_24px_rgba(0,0,0,0.3)]">
          <span className="h-2 w-2 rounded-full bg-emerald-400 shadow-[0_0_10px_rgba(74,222,128,0.8)]" />
          <span>数据更新于： {updatedAtLabel}</span>
        </div>
      </div>
    </header>
  );
}
