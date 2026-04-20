import { useLeaderboard } from '../hooks/useHeroes'
import { useNavigate } from 'react-router-dom'
import { Crown } from 'lucide-react'

export default function LeaderboardPage() {
  const { data: players, isLoading } = useLeaderboard()
  const navigate = useNavigate()

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-white">Leaderboard</h1>

      <div className="card p-0 overflow-hidden">
        <table className="w-full text-sm">
          <thead className="border-b border-surface-600">
            <tr className="text-slate-400 text-left">
              <th className="px-5 py-3 font-medium">Rank</th>
              <th className="px-5 py-3 font-medium">Player</th>
            </tr>
          </thead>
          <tbody>
            {isLoading
              ? Array.from({ length: 10 }).map((_, i) => (
                  <tr key={i} className="border-b border-surface-600">
                    <td className="px-5 py-4"><div className="skeleton h-4 w-8" /></td>
                    <td className="px-5 py-4"><div className="skeleton h-4 w-32" /></td>
                  </tr>
                ))
              : players?.map((player, i) => (
                  <tr
                    key={player.steam_id}
                    className="border-b border-surface-600 last:border-0 hover:bg-surface-700/50 cursor-pointer"
                    onClick={() => navigate(`/player/${player.steam_id}`)}
                  >
                    <td className="px-5 py-4">
                      {i < 3 ? (
                        <Crown size={16} className={
                          i === 0 ? 'text-yellow-400'
                            : i === 1 ? 'text-slate-300'
                            : 'text-amber-600'
                        } />
                      ) : (
                        <span className="text-slate-500">{i + 1}</span>
                      )}
                    </td>
                    <td className="px-5 py-4">
                      <div className="flex items-center gap-3">
                        <img src={player.avatar_url} alt="" className="h-8 w-8 rounded-full" />
                        <span className="font-medium text-white">{player.persona_name}</span>
                      </div>
                    </td>
                  </tr>
                ))
            }
          </tbody>
        </table>
      </div>
    </div>
  )
}
