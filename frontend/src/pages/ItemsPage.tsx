import { useEffect, useMemo, useState } from 'react'
import type { Image, Item, ItemStats } from '@/types'
import { getImage, getItems, getItemStats } from '@/api/items'
import './ItemsPage.css'

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
type CostTierLabel = 'I' | 'II' | 'III' | 'IV'

const CATEGORY_TO_SLOT: Record<Exclude<ItemCategoryLabel, 'All'>, Item['item_slot_type']> = {
  Weapon: 'weapon',
  Vitality: 'vitality',
  Spirit: 'spirit',
}

const COST_TIER_RANGES: Record<CostTierLabel, {min: number; max: number}> = {
  I: { min:800, max:1600 },
  II: { min:1600, max:3200 },
  III: { min:3200, max:6400 },
  IV: { min:6400, max: Infinity },
}

export default function ItemsPage() {
  const [items, setItems] = useState<Item[]>([])
  const [stats, setStats] = useState<ItemStats[]>([])
  const [goldSvg, setGoldSvg] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const [sortBy, setSortBy] = useState<'win_rate' | 'pick_rate'>('win_rate')
  const [category, setCategory] = useState<ItemCategoryLabel>('All')
  const [selectedTiers, setSelectedTiers] = useState<CostTierLabel[]>([])

  const rows = useMemo(() => {
    let data = toItemWithStats(items, stats)

    if (category !== 'All') {
      const selectedSlotType = CATEGORY_TO_SLOT[category]
      data = data.filter((row) => row.item.item_slot_type === selectedSlotType)
    }

    if (selectedTiers.length > 0) {
      data = data.filter((row) => {
        const cost = row.item.cost
        if (typeof cost !== 'number') {
          return false
        }

        return selectedTiers.some((tier) => {
          const range = COST_TIER_RANGES[tier]
          return cost >= range.min && cost < range.max
        })
      })
    }

    if (sortBy === 'pick_rate') {
      return [...data].sort((a, b) => b.pick_rate - a.pick_rate)
    }

    return [...data].sort((a, b) => b.win_rate - a.win_rate)
  }, [items, stats, sortBy, category, selectedTiers])

  useEffect(() => {
    let active = true

    async function load() {
      try {
        setLoading(true)
        setError(null)

        const [itemsRes, statsRes, imagesRes] = await Promise.all([
          getItems(),
          getItemStats(),
          getImage(),
        ])

        if (!active) {
          return
        }

        setItems(itemsRes)
        setStats(statsRes)

        const goldIcon = imagesRes.find((img: Image) => typeof img.gold_svg == 'string' && img.gold_svg.length > 0)?.gold_svg ?? null
        setGoldSvg(goldIcon)
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

  return (
    <main className="items-page">
      <section className="items-hero">
        <h1 className="items-title">
          <span className="items-title-top">Item</span>
          <span className="items-title-bottom">Trends</span>
        </h1>
        <p className="page-intro">
          See which items are performing best by win and pick rate across recent matches.
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

        <div className="items-tier-group" role="group" aria-label="Filter items by cost tier">
          {(['I','II','III','IV'] as const).map((tier) => (
            <button
              key={tier}
              type="button"
              className={selectedTiers.includes(tier) ? 'items-tier-button active' : 'items-tier-button'}
              onClick={() => 
                setSelectedTiers((prev) =>
                  prev.includes(tier) ? prev.filter((t) => t != tier) : [...prev, tier]
          )
        }
        aria-pressed={selectedTiers.includes(tier)}
            >
              {tier}
            </button>
          ))}
        </div>
      </section>
      <section className="items-results">
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
                </tr>
              </thead>
              <tbody>
                {rows.map((row) => (
                  <tr key={row.item.id}>
                    <td>
                      <div className="item-name-cell">
                        {row.item.shop_image ? (
                          <img
                            src={row.item.shop_image}
                            alt={`${row.item.name} icon`}
                            className="item-icon"
                            loading="lazy"
                          />
                        ) : (
                          <span className="item-icon item-icon-fallback" aria-hidden="true">
                            ?
                          </span>
                        )}
                        <strong>{row.item.name}</strong>
                      </div>
                    </td>
                    <td>
                        <span className="item-cost">
                          <img
                            src={goldSvg ?? '/Souls.png'}
                            alt="Souls"
                            className="souls-icon"
                            loading="lazy"
                          />
                          <span className="item-cost-value">
                            {typeof row.item.cost === 'number' ? row.item.cost.toLocaleString() : 'N/A'}
                          </span>
                        </span>
                      </td>
                    <td>{pct(row.win_rate)}</td>
                    <td>{pct(row.pick_rate)}</td>
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