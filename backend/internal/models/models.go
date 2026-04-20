package models

// Player represents a Deadlock player's profile
type Player struct {
	SteamID     string `json:"steam_id"`
	PersonaName string `json:"persona_name"`
	AvatarURL   string `json:"avatar_url"`
	ProfileURL  string `json:"profile_url"`
}

// PlayerStats represents aggregated stats for a player
type PlayerStats struct {
	SteamID      string      `json:"steam_id"`
	Wins         int         `json:"wins"`
	Losses       int         `json:"losses"`
	WinRate      float64     `json:"win_rate"`
	KDA          float64     `json:"kda"`
	AvgKills     float64     `json:"avg_kills"`
	AvgDeaths    float64     `json:"avg_deaths"`
	AvgAssists   float64     `json:"avg_assists"`
	HeroStats    []HeroStats `json:"hero_stats"`
}

// HeroStats represents a player's performance on a specific hero
type HeroStats struct {
	HeroID   int     `json:"hero_id"`
	HeroName string  `json:"hero_name"`
	Matches  int     `json:"matches"`
	Wins     int     `json:"wins"`
	WinRate  float64 `json:"win_rate"`
	AvgKills float64 `json:"avg_kills"`
}

// Match represents a single completed match
type Match struct {
	MatchID   int64   `json:"match_id"`
	HeroID    int     `json:"hero_id"`
	HeroName  string  `json:"hero_name"`
	Won       bool    `json:"won"`
	Kills     int     `json:"kills"`
	Deaths    int     `json:"deaths"`
	Assists   int     `json:"assists"`
	Duration  int     `json:"duration_seconds"`
	StartedAt int64   `json:"started_at"`
}

// Hero represents a Deadlock hero with aggregate stats
type Hero struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	WinRate  float64 `json:"win_rate"`
	PickRate float64 `json:"pick_rate"`
	AvgKills float64 `json:"avg_kills"`
}

// APIError is a consistent error response shape
type APIError struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}
