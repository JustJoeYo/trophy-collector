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
type ItemCategoryLabel = 'All' | 'Weapon' | 'Vitality' | 'Spirit'

const CATEGORY_TO_SLOT: Record<Exclude<ItemCategoryLabel, 'All'>, Item['item_slot_type']> = {
  Weapon: 'weapon',
  Vitality: 'vitality',
  Spirit: 'spirit',
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

  const [sortBy, setSortBy] = useState<'win_rate' | 'pick_rate'>('win_rate')
  const [category, setCategory] = useState<ItemCategoryLabel>('All')

  const rows = useMemo(() => {
    let data = toItemWithStats(items, stats)

    if (category !== 'All') {
      const selectedSlotType = CATEGORY_TO_SLOT[category]
      data = data.filter((row) => row.item.item_slot_type === selectedSlotType)
    }

    if (sortBy === 'pick_rate') {
      return [...data].sort((a, b) => b.pick_rate - a.pick_rate)
    }

    return [...data].sort((a, b) => b.win_rate - a.win_rate)
  }, [items, stats, sortBy, category])

  return (
    <main className="items-page">
      <section>
        <p className="eyebrow">Items</p>
        <h1>Item win and pick rates</h1>
        <p className="page-intro">
          See items that are winning the most games and how popular they are among players. The win and pick rates are calculated based on the most recent bucket of matches for each item, which includes matches from the last few weeks.
        </p>
      </section>
      
      <section className="items-summary" aria-label="Item summary">
        <article className="summary-card">
          <span className="summary-label">Loaded items</span>
          <strong>{rows.length}</strong>
        </article>
        <article className="summary-card">
          <span className="summary-label">Best win rate</span>
          <strong>
            {rows.length > 0 ? `${rows[0].item.name} (${pct(rows[0].win_rate)})` : 'N/A'}
          </strong>
        </article>
        <article className="summary-card">
          <span className="summary-label">Total sampled matches</span>
          <strong>{rows.reduce((sum,row) => sum + row.stats.matches, 0).toLocaleString()}</strong>
        </article>
      </section>
      <section className="items-controls" aria-label="Item category filters"> 
        <div className="items-category-group" role="group" aria-label="Filter items by category">
          {(['All', 'Weapon', 'Vitality', 'Spirit'] as const).map((label) => (
            <button
              key={label}
              type="button"
              className={category === label ? 'items-category-button active' : 'items-category-button'}
              onClick={() => setCategory(label)}
            >
              {label}
            </button>
          ))}
        </div>
      </section>
      <section className="items-controls" aria-label="Sort controls">
        <label className="items-controls-label" htmlFor="item-sort">Sort by</label>
        <select
          className="items-controls-select"
          id="item-sort"
          value={sortBy}
          onChange={(event) => setSortBy(event.target.value as 'win_rate' | 'pick_rate')}
        >
          <option value="win_rate">Win rate</option>
          <option value="pick_rate">Pick rate</option>
        </select>
      </section>
      <section>
        {loading ? <p className="status-message">Loading items...</p> : null}
        {error ? <p className="status-message status-error">{error}</p> : null}

        {!loading && !error && rows.length === 0 ? (
          <p className="status-message">No item stats available yet.</p>
        ) : null}

        {!loading && !error && rows.length > 0 ? (
          <section className="table-card">
            <table>
              <thead>
                <tr>
                  <th>Item</th>
                  <th>Cost</th>
                  <th>Win Rate</th>
                  <th>Pick Rate</th>
                  <th>Matches</th>
                  <th>Players</th>
                </tr>
              </thead>
              <tbody>
                {rows.map((row) => (
                  <tr key={row.item.id}>
                    <td>
                      <strong>{row.item.name}</strong>
                    </td>
                    <td>{typeof row.item.cost === 'number' ? row.item.cost.toLocaleString() : 'N/A'}</td>
                    <td>{pct(row.win_rate)}</td>
                    <td>{pct(row.pick_rate)}</td>
                    <td>{row.stats.matches}</td>
                    <td>{row.stats.players}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </section>
        ) : null}
      </section>
    </main>
  )
}