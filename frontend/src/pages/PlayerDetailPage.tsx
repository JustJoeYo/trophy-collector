import { useEffect, useRef, useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { getPlayerProfile, getSyncStatus } from '@/api'
import type { SyncStatus } from '@/api/players'
import type { PlayerProfile } from '@/types'

function formatDuration(seconds: number): string {
  const m = Math.floor(seconds / 60)
  const s = seconds % 60
  return `${m}:${s.toString().padStart(2, '0')}`
}

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString()
}

export default function PlayerDetailPage() {
  const { id } = useParams<{ id: string }>()
  const [profile, setProfile] = useState<PlayerProfile | null>(null)
  const [syncStatus, setSyncStatus] = useState<SyncStatus | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const pollRef = useRef<ReturnType<typeof setInterval> | null>(null)

  function stopPolling() {
    if (pollRef.current) {
      clearInterval(pollRef.current)
      pollRef.current = null
    }
  }

  useEffect(() => {
    if (!id) return

    setLoading(true)
    setError(null)
    setProfile(null)
    setSyncStatus(null)

    async function fetchProfile() {
      if (!id) return
      try {
        const result = await getPlayerProfile(id)
        if (result.state === 'ok') {
          setProfile(result.profile)
          setLoading(false)
          const status = await getSyncStatus(id)
          setSyncStatus(status)
          if (status.syncing) {
            pollRef.current = setInterval(async () => {
              try {
                const s = await getSyncStatus(id)
                setSyncStatus(s)
                if (!s.syncing) {
                  const fresh = await getPlayerProfile(id)
                  if (fresh.state === 'ok') setProfile(fresh.profile)
                  stopPolling()
                }
              } catch { /* keep polling */ }
            }, 3000)
          }
        } else {
          setLoading(false)
          pollRef.current = setInterval(async () => {
            try {
              const status = await getSyncStatus(id)
              setSyncStatus(status)
              const retry = await getPlayerProfile(id)
              if (retry.state === 'ok') {
                setProfile(retry.profile)
                if (!status.syncing) stopPolling()
              }
            } catch { /* keep polling */ }
          }, 3000)
        }
      } catch {
        setError('Player not found or has no recorded matches.')
        setLoading(false)
      }
    }

    fetchProfile()
    return () => stopPolling()
  }, [id])

  if (loading) return <main><p>Loading...</p></main>

  if (error) return (
    <main>
      <p>{error}</p>
      <Link to="/">← Back</Link>
    </main>
  )

  if (!profile) return (
    <main>
      <h1>Indexing match history...</h1>
      <p>We're pulling your full match history from the database. This takes about 30–60 seconds for first-time lookups.</p>
      {syncStatus && <p>{syncStatus.total_matches} matches indexed so far...</p>}
      <p>This page will update automatically.</p>
      <Link to="/">← Back</Link>
    </main>
  )

  const { overview, heroes, lanes, awards, recent_matches } = profile

  return (
    <main>
      <p><Link to="/">← Back</Link></p>
      <h1>Account {profile.account_id}</h1>
      <p>
        {overview.matches} matches sampled
        {syncStatus?.syncing && ` — indexing in progress (${syncStatus.total_matches} indexed so far)`}
      </p>

      <section>
        <h2>Overview</h2>
        <table>
          <tbody>
            <tr><td>Win Rate</td><td>{overview.win_rate.toFixed(1)}%</td></tr>
            <tr><td>Record</td><td>{overview.wins}W / {overview.losses}L</td></tr>
            <tr><td>KDA</td><td>{overview.kda.toFixed(2)}</td></tr>
            <tr><td>Avg Kills</td><td>{overview.avg_kills.toFixed(1)}</td></tr>
            <tr><td>Avg Deaths</td><td>{overview.avg_deaths.toFixed(1)}</td></tr>
            <tr><td>Avg Assists</td><td>{overview.avg_assists.toFixed(1)}</td></tr>
            <tr><td>Avg Net Worth</td><td>{Math.round(overview.avg_net_worth).toLocaleString()}</td></tr>
            <tr><td>Avg Last Hits</td><td>{overview.avg_last_hits.toFixed(1)}</td></tr>
            <tr><td>Avg Duration</td><td>{formatDuration(Math.round(overview.avg_duration_s))}</td></tr>
            <tr><td>Abandons</td><td>{overview.abandons}</td></tr>
          </tbody>
        </table>
      </section>

      <section>
        <h2>Heroes</h2>
        <table>
          <thead>
            <tr>
              <th>Hero ID</th><th>Matches</th><th>Win Rate</th>
              <th>KDA</th><th>Avg K</th><th>Avg D</th><th>Avg A</th>
            </tr>
          </thead>
          <tbody>
            {[...heroes].sort((a, b) => b.matches - a.matches).map(h => (
              <tr key={h.hero_id}>
                <td>{h.hero_id}</td>
                <td>{h.matches}</td>
                <td>{h.win_rate.toFixed(1)}%</td>
                <td>{h.kda.toFixed(2)}</td>
                <td>{h.avg_kills.toFixed(1)}</td>
                <td>{h.avg_deaths.toFixed(1)}</td>
                <td>{h.avg_assists.toFixed(1)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </section>

      <section>
        <h2>Lanes</h2>
        <table>
          <thead>
            <tr><th>Lane</th><th>Matches</th><th>Win Rate</th><th>KDA</th></tr>
          </thead>
          <tbody>
            {[...lanes].sort((a, b) => b.matches - a.matches).map(l => (
              <tr key={l.lane}>
                <td>Lane {l.lane}</td>
                <td>{l.matches}</td>
                <td>{l.win_rate.toFixed(1)}%</td>
                <td>{l.kda.toFixed(2)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </section>

      <section>
        <h2>Awards</h2>
        <table>
          <tbody>
            <tr><td>Most Kills</td><td>{awards.most_kills.value} (Match {awards.most_kills.match_id})</td></tr>
            <tr><td>Most Assists</td><td>{awards.most_assists.value} (Match {awards.most_assists.match_id})</td></tr>
            <tr><td>Most Last Hits</td><td>{awards.most_last_hits.value} (Match {awards.most_last_hits.match_id})</td></tr>
            <tr><td>Best KDA</td><td>{awards.best_kda.value.toFixed(2)} (Match {awards.best_kda.match_id})</td></tr>
            <tr><td>Highest Net Worth</td><td>{awards.highest_net_worth.value.toLocaleString()} (Match {awards.highest_net_worth.match_id})</td></tr>
            <tr><td>Longest Game</td><td>{formatDuration(awards.longest_game.value)} (Match {awards.longest_game.match_id})</td></tr>
          </tbody>
        </table>
      </section>

      <section>
        <h2>Recent Matches</h2>
        <table>
          <thead>
            <tr>
              <th>Date</th><th>Hero ID</th><th>Result</th>
              <th>K</th><th>D</th><th>A</th>
              <th>Net Worth</th><th>Duration</th>
            </tr>
          </thead>
          <tbody>
            {recent_matches.map(m => (
              <tr key={m.match_id}>
                <td>{formatDate(m.start_time)}</td>
                <td>{m.hero_id}</td>
                <td>{m.won ? 'Win' : 'Loss'}</td>
                <td>{m.kills}</td>
                <td>{m.deaths}</td>
                <td>{m.assists}</td>
                <td>{m.net_worth.toLocaleString()}</td>
                <td>{formatDuration(m.duration_s)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </section>
    </main>
  )
}
