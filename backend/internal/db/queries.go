package db

import (
	"context"
	"math"

	"github.com/JustJoeYo/trophy-collector/internal/models"
)

func (db *DB) GetPlayerProfile(ctx context.Context, accountID uint32) (*models.PlayerProfile, error) {
	rows, err := db.pool.Query(ctx, `
		SELECT
			m.match_id, m.game_mode, m.duration_s, m.start_time, m.winning_team,
			pm.hero_id, pm.kills, pm.deaths, pm.assists,
			pm.net_worth, pm.last_hits, pm.denies, pm.player_level,
			pm.assigned_lane, pm.abandon_match_time_s, pm.won
		FROM player_matches pm
		JOIN matches m ON m.match_id = pm.match_id
		WHERE pm.account_id = $1
		ORDER BY m.start_time DESC
	`, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type heroAcc struct {
		wins, losses, kills, deaths, assists, netWorth, lastHits, denies, level int
	}
	type laneAcc struct {
		wins, losses, kills, deaths, assists int
	}

	heroMap := map[uint32]*heroAcc{}
	laneMap := map[int]*laneAcc{}

	overview := models.PlayerOverview{}
	awards := models.Awards{}
	recentMatches := make([]models.PlayerMatchSummary, 0)
	totalNetWorth, totalLastHits, totalDenies, totalLevel, totalDuration := 0, 0, 0, 0, 0

	for rows.Next() {
		var (
			matchID           uint64
			gameMode          string
			durationS         int
			startTime         interface{}
			winningTeam       string
			heroID            uint64
			kills             int
			deaths            int
			assists           int
			netWorth          int
			lastHits          int
			denies            int
			playerLevel       int
			assignedLane      int
			abandonMatchTimeS int
			won               bool
			startTimeStr      string
		)

		err := rows.Scan(
			&matchID, &gameMode, &durationS, &startTime, &winningTeam,
			&heroID, &kills, &deaths, &assists,
			&netWorth, &lastHits, &denies, &playerLevel,
			&assignedLane, &abandonMatchTimeS, &won,
		)
		if err != nil {
			return nil, err
		}

		if t, ok := startTime.(interface{ Format(string) string }); ok {
			startTimeStr = t.Format("2006-01-02 15:04:05")
		}

		overview.Matches++
		overview.TotalKills += kills
		overview.TotalDeaths += deaths
		overview.TotalAssists += assists
		totalNetWorth += netWorth
		totalLastHits += lastHits
		totalDenies += denies
		totalLevel += playerLevel
		totalDuration += durationS
		if abandonMatchTimeS > 0 {
			overview.Abandons++
		}
		if won {
			overview.Wins++
		} else {
			overview.Losses++
		}

		if _, ok := heroMap[uint32(heroID)]; !ok {
			heroMap[uint32(heroID)] = &heroAcc{}
		}
		h := heroMap[uint32(heroID)]
		h.kills += kills
		h.deaths += deaths
		h.assists += assists
		h.netWorth += netWorth
		h.lastHits += lastHits
		h.denies += denies
		h.level += playerLevel
		if won {
			h.wins++
		} else {
			h.losses++
		}

		if _, ok := laneMap[assignedLane]; !ok {
			laneMap[assignedLane] = &laneAcc{}
		}
		l := laneMap[assignedLane]
		l.kills += kills
		l.deaths += deaths
		l.assists += assists
		if won {
			l.wins++
		} else {
			l.losses++
		}

		kda := float64(kills+assists) / math.Max(float64(deaths), 1)
		if kills > int(awards.MostKills.Value) {
			awards.MostKills = models.BestGame{MatchID: matchID, HeroID: uint32(heroID), Value: float64(kills)}
		}
		if assists > int(awards.MostAssists.Value) {
			awards.MostAssists = models.BestGame{MatchID: matchID, HeroID: uint32(heroID), Value: float64(assists)}
		}
		if lastHits > int(awards.MostLastHits.Value) {
			awards.MostLastHits = models.BestGame{MatchID: matchID, HeroID: uint32(heroID), Value: float64(lastHits)}
		}
		if netWorth > int(awards.HighestNetWorth.Value) {
			awards.HighestNetWorth = models.BestGame{MatchID: matchID, HeroID: uint32(heroID), Value: float64(netWorth)}
		}
		if kda > awards.BestKDA.Value {
			awards.BestKDA = models.BestGame{MatchID: matchID, HeroID: uint32(heroID), Value: kda}
		}
		if durationS > int(awards.LongestGame.Value) {
			awards.LongestGame = models.BestGame{MatchID: matchID, HeroID: uint32(heroID), Value: float64(durationS)}
		}

		recentMatches = append(recentMatches, models.PlayerMatchSummary{
			MatchID:     matchID,
			HeroID:      uint32(heroID),
			Won:         won,
			Kills:       kills,
			Deaths:      deaths,
			Assists:     assists,
			NetWorth:    netWorth,
			LastHits:    lastHits,
			Denies:      denies,
			PlayerLevel: playerLevel,
			DurationS:   durationS,
			StartTime:   startTimeStr,
			GameMode:    gameMode,
		})
	}

	if overview.Matches == 0 {
		return nil, nil
	}

	n := float64(overview.Matches)
	d := math.Max(float64(overview.TotalDeaths), 1)
	overview.KDA = float64(overview.TotalKills+overview.TotalAssists) / d
	overview.WinRate = float64(overview.Wins) / n * 100
	overview.AvgKills = float64(overview.TotalKills) / n
	overview.AvgDeaths = float64(overview.TotalDeaths) / n
	overview.AvgAssists = float64(overview.TotalAssists) / n
	overview.AvgNetWorth = float64(totalNetWorth) / n
	overview.AvgLastHits = float64(totalLastHits) / n
	overview.AvgDenies = float64(totalDenies) / n
	overview.AvgPlayerLevel = float64(totalLevel) / n
	overview.AvgDurationS = float64(totalDuration) / n

	heroes := make([]models.HeroPerformance, 0, len(heroMap))
	for heroID, acc := range heroMap {
		total := acc.wins + acc.losses
		fn := float64(total)
		hd := math.Max(float64(acc.deaths), 1)
		heroes = append(heroes, models.HeroPerformance{
			HeroID:         heroID,
			Matches:        total,
			Wins:           acc.wins,
			Losses:         acc.losses,
			WinRate:        float64(acc.wins) / fn * 100,
			AvgKills:       float64(acc.kills) / fn,
			AvgDeaths:      float64(acc.deaths) / fn,
			AvgAssists:     float64(acc.assists) / fn,
			KDA:            float64(acc.kills+acc.assists) / hd,
			AvgNetWorth:    float64(acc.netWorth) / fn,
			AvgLastHits:    float64(acc.lastHits) / fn,
			AvgDenies:      float64(acc.denies) / fn,
			AvgPlayerLevel: float64(acc.level) / fn,
		})
	}

	lanes := make([]models.LanePerformance, 0, len(laneMap))
	for lane, acc := range laneMap {
		total := acc.wins + acc.losses
		fn := float64(total)
		ld := math.Max(float64(acc.deaths), 1)
		lanes = append(lanes, models.LanePerformance{
			Lane:       lane,
			Matches:    total,
			Wins:       acc.wins,
			Losses:     acc.losses,
			WinRate:    float64(acc.wins) / fn * 100,
			AvgKills:   float64(acc.kills) / fn,
			AvgDeaths:  float64(acc.deaths) / fn,
			AvgAssists: float64(acc.assists) / fn,
			KDA:        float64(acc.kills+acc.assists) / ld,
		})
	}

	return &models.PlayerProfile{
		AccountID:      accountID,
		MatchesSampled: overview.Matches,
		Overview:       overview,
		Heroes:         heroes,
		Lanes:          lanes,
		Awards:         awards,
		RecentMatches:  recentMatches,
	}, nil
}

func (db *DB) GetPlayerMatchHistory(ctx context.Context, accountID uint32) ([]models.PlayerMatchSummary, error) {
	rows, err := db.pool.Query(ctx, `
		SELECT
			m.match_id, m.game_mode, m.duration_s, m.start_time,
			pm.hero_id, pm.kills, pm.deaths, pm.assists,
			pm.net_worth, pm.last_hits, pm.denies, pm.player_level, pm.won
		FROM player_matches pm
		JOIN matches m ON m.match_id = pm.match_id
		WHERE pm.account_id = $1
		ORDER BY m.start_time DESC
	`, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	matches := make([]models.PlayerMatchSummary, 0)
	for rows.Next() {
		var (
			matchID                                      uint64
			gameMode                                     string
			durationS, kills, deaths, assists            int
			netWorth, lastHits, denies, playerLevel      int
			heroID                                       uint64
			won                                          bool
			startTime                                    interface{}
			startTimeStr                                 string
		)

		err := rows.Scan(
			&matchID, &gameMode, &durationS, &startTime,
			&heroID, &kills, &deaths, &assists,
			&netWorth, &lastHits, &denies, &playerLevel, &won,
		)
		if err != nil {
			return nil, err
		}

		if t, ok := startTime.(interface{ Format(string) string }); ok {
			startTimeStr = t.Format("2006-01-02 15:04:05")
		}

		matches = append(matches, models.PlayerMatchSummary{
			MatchID:     matchID,
			HeroID:      uint32(heroID),
			Won:         won,
			Kills:       kills,
			Deaths:      deaths,
			Assists:     assists,
			NetWorth:    netWorth,
			LastHits:    lastHits,
			Denies:      denies,
			PlayerLevel: playerLevel,
			DurationS:   durationS,
			StartTime:   startTimeStr,
			GameMode:    gameMode,
		})
	}

	return matches, nil
}

func (db *DB) GetKnownPlayers(ctx context.Context) ([]uint32, error) {
	rows, err := db.pool.Query(ctx, "SELECT account_id FROM players")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var players []uint32
	for rows.Next() {
		var id uint32
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		players = append(players, id)
	}
	return players, nil
}
