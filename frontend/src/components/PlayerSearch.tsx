import { useState } from 'react'
import { useNavigate } from 'react-router-dom'

const STEAM64_REGEX = /^7656119\d{10}$/
const STEAM2_REGEX = /^STEAM_[0-5]:[01]:\d+$/i
const STEAM3_REGEX = /^\[U:1:\d+\]$/
const STEAM32_REGEX = /^\d{1,10}$/
const PROFILE_URL_REGEX = /steamcommunity\.com\/profiles\/(\d+)/
const CUSTOM_URL_REGEX = /steamcommunity\.com\/id\/[\w-]+/

function extractSteamId(input: string): string | null {
  const trimmed = input.trim()

  const profileMatch = trimmed.match(PROFILE_URL_REGEX)
  if (profileMatch) return profileMatch[1]

  if (CUSTOM_URL_REGEX.test(trimmed)) return null

  if (
    STEAM64_REGEX.test(trimmed) ||
    STEAM2_REGEX.test(trimmed) ||
    STEAM3_REGEX.test(trimmed) ||
    STEAM32_REGEX.test(trimmed)
  ) {
    return trimmed
  }

  return null
}

export default function PlayerSearch() {
  const [query, setQuery] = useState('')
  const [error, setError] = useState<string | null>(null)
  const navigate = useNavigate()

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setError(null)
    const id = extractSteamId(query)
    if (!id) {
      setError('Custom URL profiles are not supported — paste your Steam64 ID, Steam32 ID, STEAM_0:x:y, [U:1:x], or a steamcommunity.com/profiles/ link.')
      return
    }
    navigate(`/player/${encodeURIComponent(id)}`)
  }

  return (
    <form onSubmit={handleSubmit}>
      <label htmlFor="player-search">Search by Steam ID</label>
      <input
        id="player-search"
        type="text"
        value={query}
        onChange={e => { setQuery(e.target.value); setError(null) }}
        placeholder="Steam64, Steam32, STEAM_0:x:y, [U:1:x], or profile URL"
      />
      <button type="submit">Search</button>
      {error && <p role="alert">{error}</p>}
    </form>
  )
}
