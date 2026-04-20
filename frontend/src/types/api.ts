// All API response types — these mirror the Go models exactly.
// Zero `any` types. If it comes from the API, it lives here.

export interface Player {
  steam_id: string
  persona_name: string
  avatar_url: string
  profile_url: string
}

export interface PlayerStats {
  steam_id: string
  wins: number
  losses: number
  win_rate: number
  kda: number
  avg_kills: number
  avg_deaths: number
  avg_assists: number
  hero_stats: HeroStats[]
}

export interface HeroStats {
  hero_id: number
  hero_name: string
  matches: number
  wins: number
  win_rate: number
  avg_kills: number
}

export interface Match {
  match_id: number
  hero_id: number
  hero_name: string
  won: boolean
  kills: number
  deaths: number
  assists: number
  duration_seconds: number
  started_at: number
}

export interface Hero {
  id: number
  name: string
  win_rate: number
  pick_rate: number
  avg_kills: number
}

export interface APIError {
  error: string
  message: string
  code: number
}
