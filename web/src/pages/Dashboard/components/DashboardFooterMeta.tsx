interface DashboardFooterMetaProps {
  left: string;
  right: string;
}

export function DashboardFooterMeta({ left, right }: DashboardFooterMetaProps) {
  return (
    <footer className="flex items-center justify-between gap-6 rounded-[14px] border border-white/6 bg-[linear-gradient(180deg,rgba(6,14,33,0.84),rgba(6,14,33,0.78))] px-4 py-[9px] text-[10px] text-[#7385ab]">
      <div className="flex items-center gap-2">
        <span className="grid h-4 w-4 place-items-center rounded-full border border-white/10 text-[10px] text-[#9ab0dc]">i</span>
        <span>{left}</span>
      </div>
      <span>{right}</span>
    </footer>
  );
}
