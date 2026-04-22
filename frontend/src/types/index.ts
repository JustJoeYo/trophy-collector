export interface Hero {
  id: number
  class_name: string
  name: string
}

export interface MatchPlayer {
  account_id: number
  hero_id: number
  team: string
  kills: number
  deaths: number
  assists: number
  net_worth: number
  player_level: number
  assigned_lane: number
  last_hits: number
  denies: number
  abandon_match_time_s: number
}

export interface Match {
  match_id: number
  match_outcome: string
  winning_team: string
  game_mode: string
  match_mode: string
  duration_s: number
  start_time: string
  players: MatchPlayer[]
}

export interface PlayerMatchSummary {
  match_id: number
  hero_id: number
  won: boolean
  kills: number
  deaths: number
  assists: number
  net_worth: number
  last_hits: number
  denies: number
  player_level: number
  duration_s: number
  start_time: string
  game_mode: string
}

export interface PlayerOverview {
  matches: number
  wins: number
  losses: number
  win_rate: number
  total_kills: number
  total_deaths: number
  total_assists: number
  kda: number
  avg_kills: number
  avg_deaths: number
  avg_assists: number
  avg_net_worth: number
  avg_last_hits: number
  avg_denies: number
  avg_player_level: number
  avg_duration_s: number
  abandons: number
}

export interface HeroPerformance {
  hero_id: number
  matches: number
  wins: number
  losses: number
  win_rate: number
  avg_kills: number
  avg_deaths: number
  avg_assists: number
  kda: number
  avg_net_worth: number
  avg_last_hits: number
  avg_denies: number
  avg_player_level: number
}

export interface LanePerformance {
  lane: number
  matches: number
  wins: number
  losses: number
  win_rate: number
  avg_kills: number
  avg_deaths: number
  avg_assists: number
  kda: number
}

export interface BestGame {
  match_id: number
  hero_id: number
  value: number
}

export interface Awards {
  most_kills: BestGame
  best_kda: BestGame
  highest_net_worth: BestGame
  most_assists: BestGame
  most_last_hits: BestGame
  longest_game: BestGame
}

export interface FrequentPlayer {
  account_id: number
  matches: number
  wins: number
  win_rate: number
  as_teammate: boolean
}

export interface PlayerProfile {
  account_id: number
  matches_sampled: number
  overview: PlayerOverview
  heroes: HeroPerformance[]
  lanes: LanePerformance[]
  awards: Awards
  frequent_players: FrequentPlayer[]
  recent_matches: PlayerMatchSummary[]
}

export interface PlayerStats {
  account_id: number
  matches_sampled: number
  wins: number
  losses: number
  win_rate: number
  total_kills: number
  total_deaths: number
  total_assists: number
  kda: number
  avg_kills: number
  avg_deaths: number
  avg_assists: number
  avg_net_worth: number
  avg_duration_s: number
}

export interface LeaderboardEntry {
  account_name: string
  possible_account_ids: number[]
  rank: number
  top_hero_ids: number[]
  badge_level: number
  ranked_rank: number
  ranked_subrank: number
}

export interface Leaderboard {
  entries: LeaderboardEntry[]
}

export interface HeroStats {
  hero_id: number
  wins: number
  losses: number
  matches: number
  total_kills: number
  total_deaths: number
  total_assists: number
}

export interface HeroBanStats {
  hero_id: number
  bucket: number
  bans: number
}

export interface HeroBuildStats {
  hero_id: number
  hero_build_id: number
  wins: number
  losses: number
  matches: number
  players: number
}

export interface HeroCounterStats {
  hero_id: number
  enemy_hero_id: number
  wins: number
  matches_played: number
  kills: number
  enemy_kills: number
  deaths: number
  enemy_deaths: number
  assists: number
  enemy_assists: number
  denies: number
  enemy_denies: number
  last_hits: number
  enemy_last_hits: number
  networth: number
  enemy_networth: number
}

export interface HeroSynergyStats {
  hero_id1: number
  hero_id2: number
  wins: number
  matches_played: number
  kills1: number
  kills2: number
  deaths1: number
  deaths2: number
  assists1: number
  assists2: number
  networth1: number
  networth2: number
}

export interface AbilityOrderStats {
  abilities: number[]
  wins: number
  losses: number
  matches: number
  players: number
  total_kills: number
  total_deaths: number
  total_assists: number
}

export interface Item {
  id: number
  class_name: string
  name: string
  item_slot_type: 'weapon' | 'vitality' | 'spirit'
  cost: number
}

export interface ItemStats {
  item_id: number
  bucket: number
  wins: number
  losses: number
  matches: number
  players: number
  avg_buy_time_s: number
  avg_sell_time_s: number
  avg_buy_time_relative: number
  avg_sell_time_relative: number
}

export interface GameStats {
  bucket: number
  total_matches: number
  avg_duration_s: number
  avg_kills: number
  avg_deaths: number
  avg_assists: number
  avg_kd_ratio: number
  avg_net_worth: number
  avg_last_hits: number
  avg_denies: number
}

export interface KillDeathStats {
  position_x: number
  position_y: number
  killer_team: number
  deaths: number
  kills: number
}

export interface BadgeDistribution {
  badge_level: number
  total_matches: number
}

export interface HeroScoreboard {
  rank: number
  hero_id: number
  value: number
  matches: number
}

export interface PlayerScoreboard {
  rank: number
  account_id: number
  value: number
  matches: number
}

export interface BuildDetails {
  hero_id: number
  hero_build_id: number
  author_account_id: number
  last_updated_timestamp: number
  name: string
  description: string
  version: number
  num_favorites: number
  num_ignores: number
}

export interface Build {
  hero_build: BuildDetails
}

export interface RankImages {
  large: string
  large_webp: string
  small: string
  small_webp: string
}

export interface Rank {
  tier: number
  name: string
  images: RankImages
  color: string
}

export interface MetricStat {
  avg: number
  std: number
  percentile1: number
  percentile5: number
  percentile10: number
  percentile25: number
  percentile50: number
  percentile75: number
  percentile90: number
  percentile95: number
  percentile99: number
}

export interface PlayerMetrics {
  teammate_healing: MetricStat
  self_healing_per_min: MetricStat
  crit_shot_rate: MetricStat
  player_damage: MetricStat
  healing: MetricStat
  kills_plus_assists: MetricStat
  denies: MetricStat
  neutral_damage_per_min: MetricStat
}
