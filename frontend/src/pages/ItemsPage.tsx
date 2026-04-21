import { useEffect, useMemo, useState } from 'react'
import type { Item, ItemStats } from '@/types'
import { getItems, getItemStats } from '@/api/items'

type ItemWithStats = {
  item: Item
  stats: ItemStats
  win_rate: number
  pick_rate: number
}

function toItemWithStats(items: Item[], stats: ItemStats[]): ItemWithStats[] {
  const bestStatsByItemId = new Map<number, ItemStats>()

  for (const row of stats) {
    const current = bestStatsByItemId.get(row.item_id)
    if (!current || row.bucket > current.bucket) {
      bestStatsByItemId.set(row.item_id, row)
    }
  }

  const merged: ItemWithStats[] = []

  for (const item of items) {
    const stat = bestStatsByItemId.get(item.id)
    if (!stat) {
      continue
    }

    const totalMatches = stat.wins + stat.losses
    const winRate = totalMatches > 0 ? stat.wins / totalMatches : 0

    merged.push({
      item,
      stats: stat,
      win_rate: winRate,
      pick_rate: 0,
    })
  }

  const totalPlayers = merged.reduce((sum, row) => sum + row.stats.players, 0)

  for (const row of merged) {
    row.pick_rate = totalPlayers > 0 ? row.stats.players / totalPlayers : 0
  }

  return merged.sort((a, b) => b.win_rate - a.win_rate)
}

function pct(value: number): string {
  return `${(value * 100).toFixed(1)}%`
}

export default function ItemsPage() {
  const [items, setItems] = useState<Item[]>([])
  const [stats, setStats] = useState<ItemStats[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    let active = true

    async function load() {
      try {
        setLoading(true)
        setError(null)

        const [itemsRes, statsRes] = await Promise.all([getItems(), getItemStats()])

        if (!active) {
          return
        }

        setItems(itemsRes)
        setStats(statsRes)
      } catch (err) {
        if (!active) {
          return
        }

        setError(err instanceof Error ? err.message : 'Failed to load item data')
      } finally {
        if (active) {
          setLoading(false)
        }
      }
    }

    load()

    return () => {
      active = false
    }
  }, [])

  const rows = useMemo(() => toItemWithStats(items, stats), [items, stats])

  return (
    <main>
      <h1>Items</h1>
      <p>Win rate and pick rate from live item stats.</p>

      {loading ? <p>Loading items...</p> : null}
      {error ? <p>{error}</p> : null}

      {!loading && !error ? (
        <table>
          <thead>
            <tr>
              <th>Item</th>
              <th>Class</th>
              <th>Win Rate</th>
              <th>Pick Rate</th>
              <th>Matches</th>
              <th>Players</th>
            </tr>
          </thead>
          <tbody>
            {rows.map((row) => (
              <tr key={row.item.id}>
                <td>{row.item.name}</td>
                <td>{row.item.class_name}</td>
                <td>{pct(row.win_rate)}</td>
                <td>{pct(row.pick_rate)}</td>
                <td>{row.stats.matches}</td>
                <td>{row.stats.players}</td>
              </tr>
            ))}
          </tbody>
        </table>
      ) : null}
    </main>
  )
}