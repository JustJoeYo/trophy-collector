import clsx from 'clsx'

interface StatCardProps {
  label: string
  value: string | number
  sub?: string
  accent?: boolean
}

export default function StatCard({ label, value, sub, accent }: StatCardProps) {
  return (
    <div className={clsx('card flex flex-col gap-1', accent && 'border-brand-600/50')}>
      <span className="stat-label">{label}</span>
      <span className={clsx('stat-value', accent && 'text-brand-400')}>{value}</span>
      {sub && <span className="text-xs text-slate-500">{sub}</span>}
    </div>
  )
}
