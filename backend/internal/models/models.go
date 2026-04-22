package models

type Hero struct {
	HeroID    uint32 `json:"id"`
	ClassName string `json:"class_name"`
	Name      string `json:"name"`
}

type MatchPlayer struct {
	AccountID         uint32 `json:"account_id"`
	HeroID            uint32 `json:"hero_id"`
	Team              string `json:"team"`
	Kills             int    `json:"kills"`
	Deaths            int    `json:"deaths"`
	Assists           int    `json:"assists"`
	NetWorth          int    `json:"net_worth"`
	PlayerLevel       int    `json:"player_level"`
	AssignedLane      int    `json:"assigned_lane"`
	LastHits          int    `json:"last_hits"`
	Denies            int    `json:"denies"`
	AbandonMatchTimeS int    `json:"abandon_match_time_s"`
}

type Match struct {
	MatchID      uint64        `json:"match_id"`
	MatchOutcome string        `json:"match_outcome"`
	WinningTeam  string        `json:"winning_team"`
	GameMode     string        `json:"game_mode"`
	MatchMode    string        `json:"match_mode"`
	DurationS    int           `json:"duration_s"`
	StartTime    string        `json:"start_time"`
	Players      []MatchPlayer `json:"players"`
}

type LeaderboardEntry struct {
	AccountName        string   `json:"account_name"`
	PossibleAccountIDs []uint32 `json:"possible_account_ids"`
	Rank               int      `json:"rank"`
	TopHeroIDs         []uint32 `json:"top_hero_ids"`
	BadgeLevel         int      `json:"badge_level"`
	RankedRank         int      `json:"ranked_rank"`
	RankedSubrank      int      `json:"ranked_subrank"`
}

type Leaderboard struct {
	Entries []LeaderboardEntry `json:"entries"`
}

type HeroStats struct {
	HeroID       uint32 `json:"hero_id"`
	Wins         int    `json:"wins"`
	Losses       int    `json:"losses"`
	Matches      int    `json:"matches"`
	TotalKills   int    `json:"total_kills"`
	TotalDeaths  int    `json:"total_deaths"`
	TotalAssists int    `json:"total_assists"`
}

type HeroBanStats struct {
	HeroID uint32 `json:"hero_id"`
	Bucket int    `json:"bucket"`
	Bans   int    `json:"bans"`
}

type HeroBuildStats struct {
	HeroID      uint32 `json:"hero_id"`
	HeroBuildID uint32 `json:"hero_build_id"`
	Wins        int    `json:"wins"`
	Losses      int    `json:"losses"`
	Matches     int    `json:"matches"`
	Players     int    `json:"players"`
}

type HeroCounterStats struct {
	HeroID        uint32  `json:"hero_id"`
	EnemyHeroID   uint32  `json:"enemy_hero_id"`
	Wins          int     `json:"wins"`
	MatchesPlayed int     `json:"matches_played"`
	Kills         int     `json:"kills"`
	EnemyKills    int     `json:"enemy_kills"`
	Deaths        int     `json:"deaths"`
	EnemyDeaths   int     `json:"enemy_deaths"`
	Assists       int     `json:"assists"`
	EnemyAssists  int     `json:"enemy_assists"`
	Denies        int     `json:"denies"`
	EnemyDenies   int     `json:"enemy_denies"`
	LastHits      int     `json:"last_hits"`
	EnemyLastHits int     `json:"enemy_last_hits"`
	NetWorth      float64 `json:"networth"`
	EnemyNetWorth float64 `json:"enemy_networth"`
}

type HeroSynergyStats struct {
	HeroID1       uint32  `json:"hero_id1"`
	HeroID2       uint32  `json:"hero_id2"`
	Wins          int     `json:"wins"`
	MatchesPlayed int     `json:"matches_played"`
	Kills1        int     `json:"kills1"`
	Kills2        int     `json:"kills2"`
	Deaths1       int     `json:"deaths1"`
	Deaths2       int     `json:"deaths2"`
	Assists1      int     `json:"assists1"`
	Assists2      int     `json:"assists2"`
	NetWorth1     float64 `json:"networth1"`
	NetWorth2     float64 `json:"networth2"`
}

type AbilityOrderStats struct {
	Abilities    []uint32 `json:"abilities"`
	Wins         int      `json:"wins"`
	Losses       int      `json:"losses"`
	Matches      int      `json:"matches"`
	Players      int      `json:"players"`
	TotalKills   int      `json:"total_kills"`
	TotalDeaths  int      `json:"total_deaths"`
	TotalAssists int      `json:"total_assists"`
}

type ItemStats struct {
	ItemID              uint32  `json:"item_id"`
	Bucket              int     `json:"bucket"`
	Wins                int     `json:"wins"`
	Losses              int     `json:"losses"`
	Matches             int     `json:"matches"`
	Players             int     `json:"players"`
	AvgBuyTimeS         float64 `json:"avg_buy_time_s"`
	AvgSellTimeS        float64 `json:"avg_sell_time_s"`
	AvgBuyTimeRelative  float64 `json:"avg_buy_time_relative"`
	AvgSellTimeRelative float64 `json:"avg_sell_time_relative"`
}

type GameStats struct {
	Bucket       int     `json:"bucket"`
	TotalMatches int     `json:"total_matches"`
	AvgDurationS float64 `json:"avg_duration_s"`
	AvgKills     float64 `json:"avg_kills"`
	AvgDeaths    float64 `json:"avg_deaths"`
	AvgAssists   float64 `json:"avg_assists"`
	AvgKDRatio   float64 `json:"avg_kd_ratio"`
	AvgNetWorth  float64 `json:"avg_net_worth"`
	AvgLastHits  float64 `json:"avg_last_hits"`
	AvgDenies    float64 `json:"avg_denies"`
}

type KillDeathStats struct {
	PositionX  int `json:"position_x"`
	PositionY  int `json:"position_y"`
	KillerTeam int `json:"killer_team"`
	Deaths     int `json:"deaths"`
	Kills      int `json:"kills"`
}

type BadgeDistribution struct {
	BadgeLevel   int `json:"badge_level"`
	TotalMatches int `json:"total_matches"`
}

type HeroScoreboard struct {
	Rank    int     `json:"rank"`
	HeroID  uint32  `json:"hero_id"`
	Value   float64 `json:"value"`
	Matches int     `json:"matches"`
}

type PlayerScoreboard struct {
	Rank      int     `json:"rank"`
	AccountID uint32  `json:"account_id"`
	Value     float64 `json:"value"`
	Matches   int     `json:"matches"`
}

type BuildDetails struct {
	HeroID               uint32 `json:"hero_id"`
	HeroBuildID          uint32 `json:"hero_build_id"`
	AuthorAccountID      uint32 `json:"author_account_id"`
	LastUpdatedTimestamp int64  `json:"last_updated_timestamp"`
	Name                 string `json:"name"`
	Description          string `json:"description"`
	Version              int    `json:"version"`
	NumFavorites         int    `json:"num_favorites"`
	NumIgnores           int    `json:"num_ignores"`
}

type Build struct {
	HeroBuild BuildDetails `json:"hero_build"`
}

type RankImages struct {
	Large     string `json:"large"`
	LargeWebp string `json:"large_webp"`
	Small     string `json:"small"`
	SmallWebp string `json:"small_webp"`
}

type Rank struct {
	Tier   int        `json:"tier"`
	Name   string     `json:"name"`
	Images RankImages `json:"images"`
	Color  string     `json:"color"`
}

type MetricStat struct {
	Avg          float64 `json:"avg"`
	Std          float64 `json:"std"`
	Percentile1  float64 `json:"percentile1"`
	Percentile5  float64 `json:"percentile5"`
	Percentile25 float64 `json:"percentile25"`
	Percentile50 float64 `json:"percentile50"`
	Percentile75 float64 `json:"percentile75"`
	Percentile90 float64 `json:"percentile90"`
	Percentile95 float64 `json:"percentile95"`
	Percentile99 float64 `json:"percentile99"`
}

type PlayerMetrics struct {
	TeammateHealing     MetricStat `json:"teammate_healing"`
	SelfHealingPerMin   MetricStat `json:"self_healing_per_min"`
	CritShotRate        MetricStat `json:"crit_shot_rate"`
	PlayerDamage        MetricStat `json:"player_damage"`
	Healing             MetricStat `json:"healing"`
	KillsPlusAssists    MetricStat `json:"kills_plus_assists"`
	Denies              MetricStat `json:"denies"`
	NeutralDamagePerMin MetricStat `json:"neutral_damage_per_min"`
}

type Item struct {
	ItemID       uint32 `json:"id"`
	ClassName    string `json:"class_name"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	ItemSlotType string `json:"item_slot_type"`
	Cost         uint32 `json:"cost"`
}

type HeroPerformance struct {
	HeroID         uint32  `json:"hero_id"`
	Matches        int     `json:"matches"`
	Wins           int     `json:"wins"`
	Losses         int     `json:"losses"`
	WinRate        float64 `json:"win_rate"`
	AvgKills       float64 `json:"avg_kills"`
	AvgDeaths      float64 `json:"avg_deaths"`
	AvgAssists     float64 `json:"avg_assists"`
	KDA            float64 `json:"kda"`
	AvgNetWorth    float64 `json:"avg_net_worth"`
	AvgLastHits    float64 `json:"avg_last_hits"`
	AvgDenies      float64 `json:"avg_denies"`
	AvgPlayerLevel float64 `json:"avg_player_level"`
}

type LanePerformance struct {
	Lane       int     `json:"lane"`
	Matches    int     `json:"matches"`
	Wins       int     `json:"wins"`
	Losses     int     `json:"losses"`
	WinRate    float64 `json:"win_rate"`
	AvgKills   float64 `json:"avg_kills"`
	AvgDeaths  float64 `json:"avg_deaths"`
	AvgAssists float64 `json:"avg_assists"`
	KDA        float64 `json:"kda"`
}

type BestGame struct {
	MatchID uint64  `json:"match_id"`
	HeroID  uint32  `json:"hero_id"`
	Value   float64 `json:"value"`
}

type Awards struct {
	MostKills       BestGame `json:"most_kills"`
	BestKDA         BestGame `json:"best_kda"`
	HighestNetWorth BestGame `json:"highest_net_worth"`
	MostAssists     BestGame `json:"most_assists"`
	MostLastHits    BestGame `json:"most_last_hits"`
	LongestGame     BestGame `json:"longest_game"`
}

type FrequentPlayer struct {
	AccountID  uint32  `json:"account_id"`
	Matches    int     `json:"matches"`
	Wins       int     `json:"wins"`
	WinRate    float64 `json:"win_rate"`
	AsTeammate bool    `json:"as_teammate"`
}

type PlayerOverview struct {
	Matches        int     `json:"matches"`
	Wins           int     `json:"wins"`
	Losses         int     `json:"losses"`
	WinRate        float64 `json:"win_rate"`
	TotalKills     int     `json:"total_kills"`
	TotalDeaths    int     `json:"total_deaths"`
	TotalAssists   int     `json:"total_assists"`
	KDA            float64 `json:"kda"`
	AvgKills       float64 `json:"avg_kills"`
	AvgDeaths      float64 `json:"avg_deaths"`
	AvgAssists     float64 `json:"avg_assists"`
	AvgNetWorth    float64 `json:"avg_net_worth"`
	AvgLastHits    float64 `json:"avg_last_hits"`
	AvgDenies      float64 `json:"avg_denies"`
	AvgPlayerLevel float64 `json:"avg_player_level"`
	AvgDurationS   float64 `json:"avg_duration_s"`
	Abandons       int     `json:"abandons"`
}

type PlayerProfile struct {
	AccountID       uint32               `json:"account_id"`
	MatchesSampled  int                  `json:"matches_sampled"`
	Overview        PlayerOverview       `json:"overview"`
	Heroes          []HeroPerformance    `json:"heroes"`
	Lanes           []LanePerformance    `json:"lanes"`
	Awards          Awards               `json:"awards"`
	FrequentPlayers []FrequentPlayer     `json:"frequent_players"`
	RecentMatches   []PlayerMatchSummary `json:"recent_matches"`
}

type PlayerMatchSummary struct {
	MatchID     uint64 `json:"match_id"`
	HeroID      uint32 `json:"hero_id"`
	Won         bool   `json:"won"`
	Kills       int    `json:"kills"`
	Deaths      int    `json:"deaths"`
	Assists     int    `json:"assists"`
	NetWorth    int    `json:"net_worth"`
	LastHits    int    `json:"last_hits"`
	Denies      int    `json:"denies"`
	PlayerLevel int    `json:"player_level"`
	DurationS   int    `json:"duration_s"`
	StartTime   string `json:"start_time"`
	GameMode    string `json:"game_mode"`
}

type PlayerStats struct {
	AccountID      uint32  `json:"account_id"`
	MatchesSampled int     `json:"matches_sampled"`
	Wins           int     `json:"wins"`
	Losses         int     `json:"losses"`
	WinRate        float64 `json:"win_rate"`
	TotalKills     int     `json:"total_kills"`
	TotalDeaths    int     `json:"total_deaths"`
	TotalAssists   int     `json:"total_assists"`
	KDA            float64 `json:"kda"`
	AvgKills       float64 `json:"avg_kills"`
	AvgDeaths      float64 `json:"avg_deaths"`
	AvgAssists     float64 `json:"avg_assists"`
	AvgNetWorth    float64 `json:"avg_net_worth"`
	AvgDurationS   float64 `json:"avg_duration_s"`
}
