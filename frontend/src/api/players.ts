import { fetchJson } from './client'
import type { PlayerProfile, PlayerStats, PlayerMatchSummary, PlayerMetrics, Match } from '@/types'
import { API_BASE } from './client'

export type SyncStatus = { synced: boolean; syncing: boolean; total_matches: number }
export type ProfileResult =
  | { state: 'ok'; profile: PlayerProfile }
  | { state: 'syncing'; message: string }

export async function getPlayerProfile(id: string): Promise<ProfileResult> {
  const res = await fetch(`${API_BASE}/api/v1/players/${id}/profile`)
  if (res.status === 202) {
    const body = await res.json() as { message: string }
    return { state: 'syncing', message: body.message }
  }
  if (!res.ok) throw new Error(`${res.status} ${res.statusText}`)
  const profile = await res.json() as PlayerProfile
  return { state: 'ok', profile }
}

export const getSyncStatus = (id: string) =>
  fetchJson<SyncStatus>(`/api/v1/players/${id}/sync-status`)

export interface PlayerSearchResult {
  account_name: string
  account_id: number
  rank: number
  region: string
  badge_level: number
}

export const searchPlayers = (q: string) =>
  fetchJson<PlayerSearchResult[]>(`/api/v1/players/search?q=${encodeURIComponent(q)}`)

export const getPlayerAvatar = (id: number) =>
  fetchJson<{ avatar_url: string }>(`/api/v1/players/${id}/avatar`)

export const getPlayerStats = (id: string, matches = 20) =>
  fetchJson<PlayerStats>(`/api/v1/players/${id}/stats?matches=${matches}`)

export const getPlayerMatches = (id: string, limit = 20) =>
  fetchJson<PlayerMatchSummary[]>(`/api/v1/players/${id}/matches?limit=${limit}`)

export const getPlayerMetrics = (id: string) =>
  fetchJson<PlayerMetrics>(`/api/v1/players/${id}/metrics`)

export const getActiveMatches = (id: string) =>
  fetchJson<Match[]>(`/api/v1/players/${id}/active`)
