import { useEffect, useState } from 'react'
import { getHeroes, getHeroStats } from '@/api'
import type { Hero, HeroStats } from '@/types'

export default function HeroPage() {
  const [heroes, setHeroes] = useState<Hero[]>([])
  const [stats, setStats] = useState<Map<number, HeroStats>>(new Map())
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    Promise.all([getHeroes(), getHeroStats()])
      .then(([heroList, statList]) => {
        setHeroes(heroList)
        setStats(new Map(statList.map(s => [s.hero_id, s])))
      })
      .finally(() => setLoading(false))
  }, [])

  if (loading) return <main><p>Loading...</p></main>

  return (
    <main>
      <h1>Heroes</h1>
      <table>
        <thead>
          <tr>
            <th>Name</th>
            <th>Matches</th>
            <th>Win Rate</th>
            <th>Avg KDA</th>
          </tr>
        </thead>
        <tbody>
          {heroes.map(hero => {
            const s = stats.get(hero.id)
            const winRate = s ? (s.wins / s.matches * 100).toFixed(1) : '--'
            const kda = s
              ? ((s.total_kills + s.total_assists) / Math.max(s.total_deaths, 1) / s.matches).toFixed(2)
              : '--'
            return (
              <tr key={hero.id}>
                <td>{hero.name}</td>
                <td>{s?.matches.toLocaleString() ?? '--'}</td>
                <td>{winRate}{s ? '%' : ''}</td>
                <td>{kda}</td>
              </tr>
            )
          })}
        </tbody>
      </table>
    </main>
  )
}
