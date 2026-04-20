import { useParams } from 'react-router-dom'
import { usePlayer, usePlayerMatches, usePlayerHeroes } from '../hooks/usePlayer'
import StatCard from '../components/ui/StatCard'
import { StatCardSkeleton, MatchRowSkeleton } from '../components/ui/Skeleton'
import { AlertCircle } from 'lucide-react'
import clsx from 'clsx'

export default function PlayerPage() {
  const { steamId = '' } = useParams<{ steamId: string }>()

  const player  = usePlayer(steamId)
  const matches = usePlayerMatches(steamId)
  const heroes  = usePlayerHeroes(steamId)

  if (player.isError) {
    return (
      <div className="flex flex-col items-center gap-3 py-20 text-slate-400">
        <AlertCircle size={32} className="text-red-400" />
        <p>Player not found. Check the Steam ID and try again.</p>
      </div>
    )
  }

  return (
    <div className="space-y-8">
      {/* Player header */}
      <div className="flex items-center gap-4">
        {player.isLoading ? (
          <div className="skeleton h-16 w-16 rounded-full" />
        ) : (
          <img
            src={player.data?.avatar_url}
            alt={player.data?.persona_name}
            className="h-16 w-16 rounded-full border-2 border-brand-600/50"
          />
        )}
        <div>
          {player.isLoading ? (
            <>
              <div className="skeleton h-7 w-40 mb-2" />
              <div className="skeleton h-4 w-24" />
            </>
          ) : (
            <>
              <h1 className="text-2xl font-bold text-white">{player.data?.persona_name}</h1>
              <p className="text-sm text-slate-400">{steamId}</p>
            </>
          )}
        </div>
      </div>

      {/* Stats row */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        {player.isLoading
          ? Array.from({ length: 4 }).map((_, i) => <StatCardSkeleton key={i} />)
          : (
            <>
              <StatCard label="Matches"  value={matches.data?.length ?? 0} />
              <StatCard label="KDA"      value="—" accent />
              <StatCard label="Win Rate" value="—" />
              <StatCard label="Avg Kills" value="—" />
            </>
          )
        }
      </div>

      {/* Recent matches */}
      <div>
        <h2 className="text-lg font-semibold text-white mb-3">Recent Matches</h2>
        <div className="card p-0 overflow-hidden">
          {matches.isLoading
            ? Array.from({ length: 5 }).map((_, i) => <MatchRowSkeleton key={i} />)
            : matches.data?.length === 0
              ? <p className="p-6 text-slate-500 text-sm">No matches found.</p>
              : matches.data?.map(match => (
                <div
                  key={match.match_id}
                  className="flex items-center gap-4 px-5 py-4 border-b border-surface-600 last:border-0"
                >
                  <div className={clsx(
                    'w-1.5 h-10 rounded-full shrink-0',
                    match.won ? 'bg-green-500' : 'bg-red-500'
                  )} />
                  <div className="flex-1">
                    <p className="font-medium text-white">{match.hero_name}</p>
                    <p className="text-xs text-slate-500">
                      {match.kills}/{match.deaths}/{match.assists}
                    </p>
                  </div>
                  <span className={clsx(
                    'text-xs font-semibold px-2.5 py-1 rounded-full',
                    match.won
                      ? 'bg-green-500/10 text-green-400'
                      : 'bg-red-500/10 text-red-400'
                  )}>
                    {match.won ? 'WIN' : 'LOSS'}
                  </span>
                </div>
              ))
          }
        </div>
      </div>

      {/* Hero breakdown */}
      <div>
        <h2 className="text-lg font-semibold text-white mb-3">Hero Performance</h2>
        <div className="card p-0 overflow-hidden">
          <table className="w-full text-sm">
            <thead className="border-b border-surface-600">
              <tr className="text-slate-400 text-left">
                <th className="px-5 py-3 font-medium">Hero</th>
                <th className="px-5 py-3 font-medium">Matches</th>
                <th className="px-5 py-3 font-medium">Win Rate</th>
                <th className="px-5 py-3 font-medium">Avg Kills</th>
              </tr>
            </thead>
            <tbody>
              {heroes.isLoading
                ? Array.from({ length: 3 }).map((_, i) => (
                    <tr key={i} className="border-b border-surface-600">
                      {Array.from({ length: 4 }).map((_, j) => (
                        <td key={j} className="px-5 py-4">
                          <div className="skeleton h-4 w-16" />
                        </td>
                      ))}
                    </tr>
                  ))
                : heroes.data?.map(hero => (
                    <tr key={hero.hero_id} className="border-b border-surface-600 last:border-0 hover:bg-surface-700/50">
                      <td className="px-5 py-4 font-medium text-white">{hero.hero_name}</td>
                      <td className="px-5 py-4 text-slate-400">{hero.matches}</td>
                      <td className="px-5 py-4">
                        <span className={clsx(
                          'font-semibold',
                          hero.win_rate >= 50 ? 'text-green-400' : 'text-red-400'
                        )}>
                          {hero.win_rate.toFixed(1)}%
                        </span>
                      </td>
                      <td className="px-5 py-4 text-slate-400">{hero.avg_kills.toFixed(1)}</td>
                    </tr>
                  ))
              }
            </tbody>
          </table>
        </div>
      </div>
    </div>
  )
}
