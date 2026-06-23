type Color = 'teal' | 'amber' | 'green' | 'gray'

const colorMap: Record<Color, string> = {
  teal:  'bg-teal-100 text-teal-700',
  amber: 'bg-amber-100 text-amber-800',
  green: 'bg-emerald-100 text-emerald-700',
  gray:  'bg-slate-100 text-slate-600',
}

export function Badge({ children, color = 'teal' }: { children: React.ReactNode; color?: Color }) {
  return (
    <span className={`inline-block px-2.5 py-0.5 rounded-full text-xs font-semibold tracking-wide ${colorMap[color]}`}>
      {children}
    </span>
  )
}