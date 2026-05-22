import { Area, AreaChart, CartesianGrid, Line, LineChart, ResponsiveContainer, Tooltip, XAxis, YAxis } from 'recharts';

interface CostPoint {
  label: string;
  cost: number;
}

interface TokenPoint {
  label: string;
  input: number;
  output: number;
}

interface TrendPanelProps {
  costTrend: readonly CostPoint[];
  tokenTrend: readonly TokenPoint[];
}

function PanelFrame(props: { title: string; children: React.ReactNode }) {
  return (
    <section className="rounded-[13px] border border-white/6 bg-[linear-gradient(180deg,rgba(18,29,59,0.9),rgba(10,20,43,0.94))] px-[14px] py-[12px]">
      <div className="mb-3 flex items-center justify-between">
        <h3 className="text-[12px] font-semibold text-white">{props.title}</h3>
        <div className="rounded-[8px] border border-white/8 bg-white/4 px-[10px] py-[5px] text-[10px] text-[#a6b3d1]">近7天</div>
      </div>
      <div className="h-[188px]">{props.children}</div>
    </section>
  );
}

export function TrendPanel({ costTrend, tokenTrend }: TrendPanelProps) {
  return (
    <div className="grid gap-3 md:grid-cols-2">
      <PanelFrame title="费用趋势（USD）">
        <ResponsiveContainer width="100%" height="100%">
          <AreaChart data={costTrend}>
            <defs>
              <linearGradient id="costFill" x1="0" y1="0" x2="0" y2="1">
                <stop offset="0%" stopColor="#f6b24b" stopOpacity={0.56} />
                <stop offset="52%" stopColor="#f6b24b" stopOpacity={0.18} />
                <stop offset="100%" stopColor="#f6b24b" stopOpacity={0} />
              </linearGradient>
            </defs>
            <CartesianGrid vertical={false} stroke="rgba(148,163,184,0.12)" />
            <XAxis dataKey="label" stroke="#7283ab" tickLine={false} axisLine={false} fontSize={11} />
            <YAxis stroke="#7283ab" tickLine={false} axisLine={false} fontSize={11} width={30} ticks={[0, 3, 6, 9, 12]} domain={[0, 12]} />
            <Tooltip
              contentStyle={{ background: '#0b1430', border: '1px solid rgba(255,255,255,0.08)', borderRadius: '14px', color: '#fff' }}
            />
            <Area type="monotone" dataKey="cost" stroke="none" fill="url(#costFill)" />
            <Line
              type="monotone"
              dataKey="cost"
              stroke="#f4a72f"
              strokeWidth={2.5}
              dot={{ r: 2.2, fill: '#f7c15f', strokeWidth: 0 }}
              activeDot={{ r: 5, fill: '#ffcf70', stroke: '#f6b24b', strokeWidth: 2 }}
            />
          </AreaChart>
        </ResponsiveContainer>
      </PanelFrame>

      <PanelFrame title="Token 趋势（万）">
        <ResponsiveContainer width="100%" height="100%">
          <AreaChart data={tokenTrend}>
            <defs>
              <linearGradient id="tokenInput" x1="0" y1="0" x2="0" y2="1">
                <stop offset="0%" stopColor="#f0b343" stopOpacity={0.8} />
                <stop offset="100%" stopColor="#f0b343" stopOpacity={0.04} />
              </linearGradient>
              <linearGradient id="tokenOutput" x1="0" y1="0" x2="0" y2="1">
                <stop offset="0%" stopColor="#d95b70" stopOpacity={0.86} />
                <stop offset="100%" stopColor="#d95b70" stopOpacity={0.08} />
              </linearGradient>
            </defs>
            <CartesianGrid vertical={false} stroke="rgba(148,163,184,0.12)" />
            <XAxis dataKey="label" stroke="#7283ab" tickLine={false} axisLine={false} fontSize={11} />
            <YAxis stroke="#7283ab" tickLine={false} axisLine={false} fontSize={11} width={34} ticks={[0, 100, 200, 300, 400]} domain={[0, 400]} />
            <Tooltip
              contentStyle={{ background: '#0b1430', border: '1px solid rgba(255,255,255,0.08)', borderRadius: '14px', color: '#fff' }}
            />
            <Area type="monotone" dataKey="output" stroke="#d95b70" strokeWidth={1.6} fill="url(#tokenOutput)" fillOpacity={1} />
            <Area type="monotone" dataKey="input" stroke="#f0b343" strokeWidth={1.6} fill="url(#tokenInput)" fillOpacity={1} />
          </AreaChart>
        </ResponsiveContainer>
      </PanelFrame>
    </div>
  );
}
