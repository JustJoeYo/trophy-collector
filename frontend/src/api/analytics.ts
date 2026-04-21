import { fetchJson } from './client'
import type { GameStats, KillDeathStats, BadgeDistribution, Rank } from '@/types'

export const getGameStats = () =>
  fetchJson<GameStats[]>('/api/v1/analytics/game-stats')

export const getKillDeathStats = () =>
  fetchJson<KillDeathStats[]>('/api/v1/analytics/kill-death-stats')

export const getBadgeDistribution = () =>
  fetchJson<BadgeDistribution[]>('/api/v1/analytics/badge-distribution')

export const getRanks = () =>
  fetchJson<Rank[]>('/api/v1/ranks')
