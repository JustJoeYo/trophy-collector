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

export interface Item {
  item_id: number;
  class_name: string;
  name: string;
}

export interface ItemStats {
  item_id: number;
  bucket: number;
  wins: number;
  losses: number;
  matches: number;
  players: number;
  avg_buy_time_s: number;
  avg_sell_time_s: number;
  avg_buy_time_relative: number;
  avg_sell_time_relative: number;
}

export interface ItemWithStats extends Item {
  stats: ItemStats;
  win_rate: number;
  pick_rate: number;
}