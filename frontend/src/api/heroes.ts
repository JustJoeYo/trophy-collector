import { fetchJson } from './client'
import type {
  Hero, HeroStats, HeroBanStats, HeroBuildStats,
  HeroCounterStats, HeroSynergyStats, AbilityOrderStats,
  HeroScoreboard, Build, Leaderboard,
} from '@/types'

export const getHeroes = () =>
  fetchJson<Hero[]>('/api/v1/heroes')

export const getHeroStats = () =>
  fetchJson<HeroStats[]>('/api/v1/heroes/stats')

export const getHeroBanStats = () =>
  fetchJson<HeroBanStats[]>('/api/v1/heroes/ban-stats')

export const getHeroCounterStats = () =>
  fetchJson<HeroCounterStats[]>('/api/v1/heroes/counter-stats')

export const getHeroSynergyStats = () =>
  fetchJson<HeroSynergyStats[]>('/api/v1/heroes/synergy-stats')

export const getHeroBuildStats = (heroId: number) =>
  fetchJson<HeroBuildStats[]>(`/api/v1/heroes/${heroId}/build-stats`)

export const getAbilityOrderStats = (heroId: number) =>
  fetchJson<AbilityOrderStats[]>(`/api/v1/heroes/${heroId}/ability-order-stats`)

export const getBuilds = (heroId: number, limit = 10) =>
  fetchJson<Build[]>(`/api/v1/heroes/${heroId}/builds?limit=${limit}`)

export const getHeroScoreboard = (sortBy = 'wins', limit = 20) =>
  fetchJson<HeroScoreboard[]>(`/api/v1/scoreboard/heroes?sort_by=${sortBy}&limit=${limit}`)

export const getHeroLeaderboard = (region: string, heroId: number) =>
  fetchJson<Leaderboard>(`/api/v1/leaderboard/${region}/${heroId}`)
