import { useHeroes } from '../hooks/useHeroes'
import { AlertCircle } from 'lucide-react'
import clsx from 'clsx'

export default function HeroesPage() {
  const { data: heroes, isLoading, isError } = useHeroes()

  if (isError) {
    return (
      <div className="flex flex-col items-center gap-3 py-20 text-slate-400">
        <AlertCircle size={32} className="text-red-400" />
        <p>Failed to load hero data.</p>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-white">Hero Tier List</h1>

      <div className="card p-0 overflow-hidden">
        <table className="w-full text-sm">
          <thead className="border-b border-surface-600">
            <tr className="text-slate-400 text-left">
              <th className="px-5 py-3 font-medium">#</th>
              <th className="px-5 py-3 font-medium">Hero</th>
              <th className="px-5 py-3 font-medium">Win Rate</th>
              <th className="px-5 py-3 font-medium">Pick Rate</th>
              <th className="px-5 py-3 font-medium">Avg Kills</th>
            </tr>
          </thead>
          <tbody>
            {isLoading
              ? Array.from({ length: 10 }).map((_, i) => (
                  <tr key={i} className="border-b border-surface-600">
                    {Array.from({ length: 5 }).map((_, j) => (
                      <td key={j} className="px-5 py-4">
                        <div className="skeleton h-4 w-16" />
                      </td>
                    ))}
                  </tr>
                ))
              : heroes?.map((hero, i) => (
                  <tr key={hero.id} className="border-b border-surface-600 last:border-0 hover:bg-surface-700/50">
                    <td className="px-5 py-4 text-slate-500">{i + 1}</td>
                    <td className="px-5 py-4 font-medium text-white">{hero.name}</td>
                    <td className="px-5 py-4">
                      <span className={clsx(
                        'font-semibold',
                        hero.win_rate >= 52 ? 'text-green-400'
                          : hero.win_rate >= 48 ? 'text-slate-300'
                          : 'text-red-400'
                      )}>
                        {hero.win_rate.toFixed(1)}%
                      </span>
                    </td>
                    <td className="px-5 py-4 text-slate-400">{hero.pick_rate.toFixed(1)}%</td>
                    <td className="px-5 py-4 text-slate-400">{hero.avg_kills.toFixed(1)}</td>
                  </tr>
                ))
            }
          </tbody>
        </table>
      </div>
    </div>
  )
}
