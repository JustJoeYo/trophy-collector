import { useEffect, useRef, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { searchPlayers, getPlayerAvatar } from '@/api/players'
import type { PlayerSearchResult } from '@/api/players'

const STEAM64_REGEX = /^7656119\d{10}$/
const STEAM2_REGEX  = /^STEAM_[0-5]:[01]:\d+$/i
const STEAM3_REGEX  = /^\[U:1:\d+\]$/
const STEAM32_REGEX = /^\d{1,10}$/
const PROFILES_URL  = /steamcommunity\.com\/profiles\/(\d+)/
const ID_URL        = /steamcommunity\.com\/id\/([\w-]+)/

function avatarColor(name: string): string {
  let hash = 0
  for (let i = 0; i < name.length; i++) hash = name.charCodeAt(i) + ((hash << 5) - hash)
  const h = Math.abs(hash) % 360
  return `hsl(${h}, 40%, 28%)`
}

function extractDirectId(input: string): string | null {
  const t = input.trim()
  const profileMatch = t.match(PROFILES_URL)
  if (profileMatch) return profileMatch[1]
  if (
    STEAM64_REGEX.test(t) ||
    STEAM2_REGEX.test(t)  ||
    STEAM3_REGEX.test(t)  ||
    STEAM32_REGEX.test(t)
  ) return t
  return null
}

function extractVanityName(input: string): string | null {
  const m = input.trim().match(ID_URL)
  return m ? m[1] : null
}

function isSteamIdLike(input: string): boolean {
  const t = input.trim()
  return (
    STEAM64_REGEX.test(t) ||
    STEAM2_REGEX.test(t)  ||
    STEAM3_REGEX.test(t)  ||
    STEAM32_REGEX.test(t) ||
    PROFILES_URL.test(t)
  )
}

export default function PlayerSearch() {
  const [query, setQuery]       = useState('')
  const [results, setResults]   = useState<PlayerSearchResult[]>([])
  const [avatars, setAvatars]   = useState<Record<number, string>>({})
  const [loading, setLoading]   = useState(false)
  const [open, setOpen]         = useState(false)
  const [error, setError]       = useState<string | null>(null)
  const navigate  = useNavigate()
  const debounce  = useRef<ReturnType<typeof setTimeout> | null>(null)
  const wrapRef   = useRef<HTMLDivElement>(null)

  useEffect(() => {
    function onClickOutside(e: MouseEvent) {
      if (wrapRef.current && !wrapRef.current.contains(e.target as Node)) {
        setOpen(false)
      }
    }
    document.addEventListener('mousedown', onClickOutside)
    return () => document.removeEventListener('mousedown', onClickOutside)
  }, [])

  function handleChange(value: string) {
    setQuery(value)
    setOpen(false)
    setResults([])
    setError(null)

    if (debounce.current) clearTimeout(debounce.current)

    const trimmed = value.trim()
    if (!trimmed || isSteamIdLike(trimmed) || extractVanityName(trimmed)) return

    if (trimmed.length < 2) return

    debounce.current = setTimeout(async () => {
      setLoading(true)
      try {
        const data = await searchPlayers(trimmed)
        setResults(data)
        setOpen(data.length > 0)
        setAvatars({})
        data.forEach(r => {
          getPlayerAvatar(r.account_id)
            .then(res => setAvatars(prev => ({ ...prev, [r.account_id]: res.avatar_url })))
            .catch(() => {})
        })
      } catch {
        setResults([])
      } finally {
        setLoading(false)
      }
    }, 300)
  }

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    const trimmed = query.trim()
    if (!trimmed) return

    if (extractVanityName(trimmed)) {
      setError('Custom profile URLs aren\'t supported. Find your Steam64 ID at steamid.io then paste it here.')
      return
    }

    const direct = extractDirectId(trimmed)
    if (direct) {
      setOpen(false)
      navigate(`/player/${encodeURIComponent(direct)}`)
      return
    }

    if (results.length === 1) {
      selectResult(results[0])
      return
    }

    if (results.length > 1) {
      setOpen(true)
      return
    }

    searchPlayers(trimmed).then(data => {
      if (data.length === 1) {
        selectResult(data[0])
      } else {
        setResults(data)
        setOpen(data.length > 0)
      }
    }).catch(() => {})
  }

  function selectResult(result: PlayerSearchResult) {
    setOpen(false)
    setQuery(result.account_name)
    navigate(`/player/${result.account_id}`)
  }

  return (
    <div className="ps-wrap" ref={wrapRef}>
      <form className="ps-form" onSubmit={handleSubmit}>
        <label htmlFor="player-search">Search by Steam ID</label>
        <div className="ps-input-row">
          <input
            id="player-search"
            type="text"
            autoComplete="off"
            value={query}
            onChange={e => handleChange(e.target.value)}
            onFocus={() => results.length > 0 && setOpen(true)}
            placeholder="Steam ID, profile link, or player name"
          />
          <button type="submit">{loading ? '…' : 'Search'}</button>
        </div>
      </form>
      {error && <p className="ps-error" role="alert">{error}</p>}
      {open && results.length > 0 && (
        <ul className="ps-dropdown" role="listbox">
          {results.map(r => (
            <li
              key={r.account_id}
              className="ps-dropdown-item"
              role="option"
              onMouseDown={() => selectResult(r)}
            >
              {avatars[r.account_id] ? (
                <img
                  className="ps-dropdown-avatar ps-dropdown-avatar--img"
                  src={avatars[r.account_id]}
                  alt=""
                  aria-hidden="true"
                />
              ) : (
                <span
                  className="ps-dropdown-avatar"
                  style={{ background: avatarColor(r.account_name) }}
                  aria-hidden="true"
                >
                  {r.account_name.charAt(0).toUpperCase()}
                </span>
              )}
              <span className="ps-dropdown-info">
                <span className="ps-dropdown-name">{r.account_name}</span>
                <span className="ps-dropdown-meta">{r.region} · Rank {r.rank}</span>
              </span>
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}
