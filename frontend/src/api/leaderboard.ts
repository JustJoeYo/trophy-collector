import { fetchJson } from './client'
import type { Leaderboard, PlayerScoreboard } from '@/types'

export const getLeaderboard = (region: string) =>
  fetchJson<Leaderboard>(`/api/v1/leaderboard/${region}`)

export const getPlayerScoreboard = (sortBy = 'wins', limit = 20) =>
  fetchJson<PlayerScoreboard[]>(`/api/v1/scoreboard/players?sort_by=${sortBy}&limit=${limit}`)
