import type { Player, PlayerStats, Match, HeroStats, Hero } from '../types/api'

const BASE_URL = '/api/v1'

// Central fetch wrapper — handles errors consistently
async function request<T>(path: string): Promise<T> {
  const res = await fetch(`${BASE_URL}${path}`)

  if (!res.ok) {
    const err = await res.json().catch(() => ({ message: 'Unknown error' }))
    throw new Error(err.message ?? `Request failed: ${res.status}`)
  }

  return res.json() as Promise<T>
}

// Player
export const getPlayer = (steamId: string): Promise<Player> =>
  request<Player>(`/player/${steamId}`)

export const getPlayerStats = (steamId: string): Promise<PlayerStats> =>
  request<PlayerStats>(`/player/${steamId}`)

export const getPlayerMatches = (steamId: string): Promise<Match[]> =>
  request<Match[]>(`/player/${steamId}/matches`)

export const getPlayerHeroes = (steamId: string): Promise<HeroStats[]> =>
  request<HeroStats[]>(`/player/${steamId}/heroes`)

// Game data
export const getHeroes = (): Promise<Hero[]> =>
  request<Hero[]>('/heroes')

export const getLeaderboard = (): Promise<Player[]> =>
  request<Player[]>('/leaderboard')
