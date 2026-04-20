export interface Player {
  account_id: string
  player_name: string
  avatar_url: string
  wins: number
  losses: number
}

export interface Match {
  match_id: number
  hero_id: number
  outcome: string
  kills: number
  deaths: number
  assists: number
  duration_secs: number
}

export interface Hero {
  hero_id: number
  name: string
  image_url: string
}