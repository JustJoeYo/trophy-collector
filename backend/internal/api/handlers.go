package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/JustJoeYo/trophy-collector/internal/cache"
	"github.com/JustJoeYo/trophy-collector/internal/clients"
	"github.com/JustJoeYo/trophy-collector/internal/db"
	"github.com/JustJoeYo/trophy-collector/internal/models"
)

type Handler struct {
	deadlock clients.DeadlockClient
	cache    cache.Cache
	db       *db.DB
}

func NewHandler(deadlock clients.DeadlockClient, cache cache.Cache, database *db.DB) *Handler {
	return &Handler{
		deadlock: deadlock,
		cache:    cache,
		db:       database,
	}
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode response", "error", err)
	}
}

func (h *Handler) cacheGet(r *http.Request, w http.ResponseWriter, key string) bool {
	if cached, err := h.cache.Get(r.Context(), key); err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(cached))
		return true
	} else {
		slog.Debug("cache miss", "key", key, "error", err)
		return false
	}
}

func (h *Handler) cacheSet(r *http.Request, key string, data interface{}, ttl time.Duration) {
	if encoded, err := json.Marshal(data); err == nil {
		h.cache.Set(r.Context(), key, string(encoded), ttl)
	}
}

func (h *Handler) parseAccountID(w http.ResponseWriter, r *http.Request) (uint32, bool) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid account id"})
		return 0, false
	}
	return uint32(id), true
}

func (h *Handler) parseHeroID(w http.ResponseWriter, r *http.Request) (uint32, bool) {
	idStr := chi.URLParam(r, "heroId")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid hero id"})
		return 0, false
	}
	return uint32(id), true
}

func (h *Handler) parseLimit(r *http.Request, defaultVal, maxVal int) int {
	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		return defaultVal
	}
	if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 && parsed <= maxVal {
		return parsed
	}
	return defaultVal
}

func (h *Handler) GetPlayerMatches(w http.ResponseWriter, r *http.Request) {
	accountID, ok := h.parseAccountID(w, r)
	if !ok {
		return
	}

	cacheKey := fmt.Sprintf("player-matches:%d", accountID)
	if h.cacheGet(r, w, cacheKey) {
		return
	}

	if h.db != nil {
		needsSync, _ := h.db.NeedsSync(r.Context(), accountID)
		if needsSync {
			go h.db.SyncPlayer(context.Background(), h.deadlock, accountID)
		}
		summaries, err := h.db.GetPlayerMatchHistory(r.Context(), accountID)
		if err == nil && len(summaries) > 0 {
			h.cacheSet(r, cacheKey, summaries, 5*time.Minute)
			h.writeJSON(w, http.StatusOK, summaries)
			return
		}
	}

	limit := h.parseLimit(r, 50, 50)
	matches, err := h.deadlock.GetPlayerMatches(r.Context(), accountID, limit)
	if err != nil {
		slog.Error("failed to fetch matches", "account_id", accountID, "error", err)
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to fetch matches"})
		return
	}

	summaries := make([]models.PlayerMatchSummary, 0, len(matches))
	for _, match := range matches {
		for _, player := range match.Players {
			if player.AccountID == accountID {
				summaries = append(summaries, models.PlayerMatchSummary{
					MatchID:     match.MatchID,
					HeroID:      player.HeroID,
					Won:         player.Team == match.WinningTeam,
					Kills:       player.Kills,
					Deaths:      player.Deaths,
					Assists:     player.Assists,
					NetWorth:    player.NetWorth,
					LastHits:    player.LastHits,
					Denies:      player.Denies,
					PlayerLevel: player.PlayerLevel,
					DurationS:   match.DurationS,
					StartTime:   match.StartTime,
					GameMode:    match.GameMode,
				})
				break
			}
		}
	}

	h.cacheSet(r, cacheKey, summaries, 5*time.Minute)
	h.writeJSON(w, http.StatusOK, summaries)
}

func (h *Handler) GetPlayerStats(w http.ResponseWriter, r *http.Request) {
	accountID, ok := h.parseAccountID(w, r)
	if !ok {
		return
	}

	limit := h.parseLimit(r, 20, 50)
	cacheKey := fmt.Sprintf("player-stats:%d:%d", accountID, limit)
	if h.cacheGet(r, w, cacheKey) {
		return
	}

	matches, err := h.deadlock.GetPlayerMatches(r.Context(), accountID, limit)
	if err != nil {
		slog.Error("failed to fetch player stats", "account_id", accountID, "error", err)
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to fetch player stats"})
		return
	}

	stats := models.PlayerStats{AccountID: accountID}
	totalDuration := 0
	totalNetWorth := 0

	for _, match := range matches {
		for _, player := range match.Players {
			if player.AccountID != accountID {
				continue
			}
			stats.MatchesSampled++
			stats.TotalKills += player.Kills
			stats.TotalDeaths += player.Deaths
			stats.TotalAssists += player.Assists
			totalDuration += match.DurationS
			totalNetWorth += player.NetWorth
			if player.Team == match.WinningTeam {
				stats.Wins++
			} else {
				stats.Losses++
			}
			break
		}
	}

	if stats.MatchesSampled == 0 {
		h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "no matches found for account"})
		return
	}

	n := float64(stats.MatchesSampled)
	deaths := float64(stats.TotalDeaths)
	if deaths == 0 {
		deaths = 1
	}

	stats.KDA = float64(stats.TotalKills+stats.TotalAssists) / deaths
	stats.WinRate = float64(stats.Wins) / n * 100
	stats.AvgKills = float64(stats.TotalKills) / n
	stats.AvgDeaths = float64(stats.TotalDeaths) / n
	stats.AvgAssists = float64(stats.TotalAssists) / n
	stats.AvgDurationS = float64(totalDuration) / n
	stats.AvgNetWorth = float64(totalNetWorth) / n

	h.cacheSet(r, cacheKey, stats, 5*time.Minute)
	h.writeJSON(w, http.StatusOK, stats)
}

func (h *Handler) GetPlayerProfile(w http.ResponseWriter, r *http.Request) {
	accountID, ok := h.parseAccountID(w, r)
	if !ok {
		return
	}

	cacheKey := fmt.Sprintf("player-profile:%d", accountID)
	if h.cacheGet(r, w, cacheKey) {
		return
	}

	if h.db != nil {
		needsSync, _ := h.db.NeedsSync(r.Context(), accountID)
		if needsSync {
			go h.db.SyncPlayer(context.Background(), h.deadlock, accountID)
		}
		profile, err := h.db.GetPlayerProfile(r.Context(), accountID)
		if err == nil && profile != nil {
			h.cacheSet(r, cacheKey, profile, 5*time.Minute)
			h.writeJSON(w, http.StatusOK, profile)
			return
		}
	}

	limit := 50
	matches, err := h.deadlock.GetPlayerMatches(r.Context(), accountID, limit)
	if err != nil {
		slog.Error("failed to fetch player profile", "account_id", accountID, "error", err)
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to fetch player profile"})
		return
	}

	type heroAcc struct {
		wins, losses, kills, deaths, assists, netWorth, lastHits, denies, level int
	}
	type laneAcc struct {
		wins, losses, kills, deaths, assists int
	}
	type playerAcc struct {
		wins, matches int
		asTeammate    bool
	}

	heroMap := map[uint32]*heroAcc{}
	laneMap := map[int]*laneAcc{}
	playerMap := map[uint32]*playerAcc{}

	overview := models.PlayerOverview{}
	awards := models.Awards{}
	recentMatches := make([]models.PlayerMatchSummary, 0, len(matches))
	totalDuration, totalNetWorth, totalLastHits, totalDenies, totalLevel := 0, 0, 0, 0, 0

	for _, match := range matches {
		var me *models.MatchPlayer
		for i := range match.Players {
			if match.Players[i].AccountID == accountID {
				me = &match.Players[i]
				break
			}
		}
		if me == nil {
			continue
		}

		won := me.Team == match.WinningTeam
		overview.Matches++
		overview.TotalKills += me.Kills
		overview.TotalDeaths += me.Deaths
		overview.TotalAssists += me.Assists
		totalDuration += match.DurationS
		totalNetWorth += me.NetWorth
		totalLastHits += me.LastHits
		totalDenies += me.Denies
		totalLevel += me.PlayerLevel
		if me.AbandonMatchTimeS > 0 {
			overview.Abandons++
		}
		if won {
			overview.Wins++
		} else {
			overview.Losses++
		}

		if _, ok := heroMap[me.HeroID]; !ok {
			heroMap[me.HeroID] = &heroAcc{}
		}
		h := heroMap[me.HeroID]
		h.kills += me.Kills
		h.deaths += me.Deaths
		h.assists += me.Assists
		h.netWorth += me.NetWorth
		h.lastHits += me.LastHits
		h.denies += me.Denies
		h.level += me.PlayerLevel
		if won {
			h.wins++
		} else {
			h.losses++
		}

		if _, ok := laneMap[me.AssignedLane]; !ok {
			laneMap[me.AssignedLane] = &laneAcc{}
		}
		l := laneMap[me.AssignedLane]
		l.kills += me.Kills
		l.deaths += me.Deaths
		l.assists += me.Assists
		if won {
			l.wins++
		} else {
			l.losses++
		}

		for _, p := range match.Players {
			if p.AccountID == accountID {
				continue
			}
			teammate := p.Team == me.Team
			if _, ok := playerMap[p.AccountID]; !ok {
				playerMap[p.AccountID] = &playerAcc{asTeammate: teammate}
			}
			acc := playerMap[p.AccountID]
			acc.matches++
			if won && teammate {
				acc.wins++
			}
		}

		kda := float64(me.Kills+me.Assists) / math.Max(float64(me.Deaths), 1)
		if me.Kills > int(awards.MostKills.Value) {
			awards.MostKills = models.BestGame{MatchID: match.MatchID, HeroID: me.HeroID, Value: float64(me.Kills)}
		}
		if me.Assists > int(awards.MostAssists.Value) {
			awards.MostAssists = models.BestGame{MatchID: match.MatchID, HeroID: me.HeroID, Value: float64(me.Assists)}
		}
		if me.LastHits > int(awards.MostLastHits.Value) {
			awards.MostLastHits = models.BestGame{MatchID: match.MatchID, HeroID: me.HeroID, Value: float64(me.LastHits)}
		}
		if me.NetWorth > int(awards.HighestNetWorth.Value) {
			awards.HighestNetWorth = models.BestGame{MatchID: match.MatchID, HeroID: me.HeroID, Value: float64(me.NetWorth)}
		}
		if kda > awards.BestKDA.Value {
			awards.BestKDA = models.BestGame{MatchID: match.MatchID, HeroID: me.HeroID, Value: kda}
		}
		if match.DurationS > int(awards.LongestGame.Value) {
			awards.LongestGame = models.BestGame{MatchID: match.MatchID, HeroID: me.HeroID, Value: float64(match.DurationS)}
		}

		recentMatches = append(recentMatches, models.PlayerMatchSummary{
			MatchID:     match.MatchID,
			HeroID:      me.HeroID,
			Won:         won,
			Kills:       me.Kills,
			Deaths:      me.Deaths,
			Assists:     me.Assists,
			NetWorth:    me.NetWorth,
			LastHits:    me.LastHits,
			Denies:      me.Denies,
			PlayerLevel: me.PlayerLevel,
			DurationS:   match.DurationS,
			StartTime:   match.StartTime,
			GameMode:    match.GameMode,
		})
	}

	if overview.Matches == 0 {
		h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "no matches found for account"})
		return
	}

	n := float64(overview.Matches)
	deaths := float64(overview.TotalDeaths)
	if deaths == 0 {
		deaths = 1
	}
	overview.KDA = float64(overview.TotalKills+overview.TotalAssists) / deaths
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
		hDeaths := float64(acc.deaths)
		if hDeaths == 0 {
			hDeaths = 1
		}
		heroes = append(heroes, models.HeroPerformance{
			HeroID:         heroID,
			Matches:        total,
			Wins:           acc.wins,
			Losses:         acc.losses,
			WinRate:        float64(acc.wins) / fn * 100,
			AvgKills:       float64(acc.kills) / fn,
			AvgDeaths:      float64(acc.deaths) / fn,
			AvgAssists:     float64(acc.assists) / fn,
			KDA:            float64(acc.kills+acc.assists) / hDeaths,
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
		lDeaths := float64(acc.deaths)
		if lDeaths == 0 {
			lDeaths = 1
		}
		lanes = append(lanes, models.LanePerformance{
			Lane:       lane,
			Matches:    total,
			Wins:       acc.wins,
			Losses:     acc.losses,
			WinRate:    float64(acc.wins) / fn * 100,
			AvgKills:   float64(acc.kills) / fn,
			AvgDeaths:  float64(acc.deaths) / fn,
			AvgAssists: float64(acc.assists) / fn,
			KDA:        float64(acc.kills+acc.assists) / lDeaths,
		})
	}

	frequentPlayers := make([]models.FrequentPlayer, 0)
	for pid, acc := range playerMap {
		if acc.matches < 2 {
			continue
		}
		wr := 0.0
		if acc.asTeammate {
			wr = float64(acc.wins) / float64(acc.matches) * 100
		}
		frequentPlayers = append(frequentPlayers, models.FrequentPlayer{
			AccountID:  pid,
			Matches:    acc.matches,
			Wins:       acc.wins,
			WinRate:    wr,
			AsTeammate: acc.asTeammate,
		})
	}

	profile := models.PlayerProfile{
		AccountID:       accountID,
		MatchesSampled:  overview.Matches,
		Overview:        overview,
		Heroes:          heroes,
		Lanes:           lanes,
		Awards:          awards,
		FrequentPlayers: frequentPlayers,
		RecentMatches:   recentMatches,
	}

	h.cacheSet(r, cacheKey, profile, 5*time.Minute)
	h.writeJSON(w, http.StatusOK, profile)
}

func (h *Handler) GetPlayerMetrics(w http.ResponseWriter, r *http.Request) {
	accountID, ok := h.parseAccountID(w, r)
	if !ok {
		return
	}

	cacheKey := fmt.Sprintf("metrics:%d", accountID)
	if h.cacheGet(r, w, cacheKey) {
		return
	}

	metrics, err := h.deadlock.GetPlayerMetrics(r.Context(), accountID)
	if err != nil {
		slog.Error("failed to fetch player metrics", "account_id", accountID, "error", err)
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to fetch player metrics"})
		return
	}

	h.cacheSet(r, cacheKey, metrics, 5*time.Minute)
	h.writeJSON(w, http.StatusOK, metrics)
}

func (h *Handler) GetActiveMatches(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid account id"})
		return
	}

	cacheKey := fmt.Sprintf("active:%d", id)
	if h.cacheGet(r, w, cacheKey) {
		return
	}

	matches, err := h.deadlock.GetActiveMatches(r.Context(), []uint32{uint32(id)})
	if err != nil {
		slog.Error("failed to fetch active matches", "account_id", id, "error", err)
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to fetch active matches"})
		return
	}

	h.cacheSet(r, cacheKey, matches, 30*time.Second)
	h.writeJSON(w, http.StatusOK, matches)
}

func (h *Handler) GetHeroes(w http.ResponseWriter, r *http.Request) {
	cacheKey := "heroes:all"
	if h.cacheGet(r, w, cacheKey) {
		return
	}

	heroes, err := h.deadlock.GetHeroes(r.Context())
	if err != nil {
		slog.Error("failed to fetch heroes", "error", err)
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to fetch heroes"})
		return
	}

	h.cacheSet(r, cacheKey, heroes, 15*time.Minute)
	h.writeJSON(w, http.StatusOK, heroes)
}

func (h *Handler) GetHeroStats(w http.ResponseWriter, r *http.Request) {
	cacheKey := "hero-stats:all"
	if h.cacheGet(r, w, cacheKey) {
		return
	}

	stats, err := h.deadlock.GetHeroStats(r.Context())
	if err != nil {
		slog.Error("failed to fetch hero stats", "error", err)
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to fetch hero stats"})
		return
	}

	h.cacheSet(r, cacheKey, stats, 10*time.Minute)
	h.writeJSON(w, http.StatusOK, stats)
}

func (h *Handler) GetHeroBanStats(w http.ResponseWriter, r *http.Request) {
	cacheKey := "hero-ban-stats:all"
	if h.cacheGet(r, w, cacheKey) {
		return
	}

	stats, err := h.deadlock.GetHeroBanStats(r.Context())
	if err != nil {
		slog.Error("failed to fetch hero ban stats", "error", err)
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to fetch hero ban stats"})
		return
	}

	h.cacheSet(r, cacheKey, stats, 10*time.Minute)
	h.writeJSON(w, http.StatusOK, stats)
}

func (h *Handler) GetHeroBuildStats(w http.ResponseWriter, r *http.Request) {
	heroID, ok := h.parseHeroID(w, r)
	if !ok {
		return
	}

	cacheKey := fmt.Sprintf("hero-build-stats:%d", heroID)
	if h.cacheGet(r, w, cacheKey) {
		return
	}

	stats, err := h.deadlock.GetHeroBuildStats(r.Context(), heroID)
	if err != nil {
		slog.Error("failed to fetch hero build stats", "hero_id", heroID, "error", err)
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to fetch hero build stats"})
		return
	}

	h.cacheSet(r, cacheKey, stats, 10*time.Minute)
	h.writeJSON(w, http.StatusOK, stats)
}

func (h *Handler) GetHeroCounterStats(w http.ResponseWriter, r *http.Request) {
	cacheKey := "hero-counter-stats:all"
	if h.cacheGet(r, w, cacheKey) {
		return
	}

	stats, err := h.deadlock.GetHeroCounterStats(r.Context())
	if err != nil {
		slog.Error("failed to fetch hero counter stats", "error", err)
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to fetch hero counter stats"})
		return
	}

	h.cacheSet(r, cacheKey, stats, 10*time.Minute)
	h.writeJSON(w, http.StatusOK, stats)
}

func (h *Handler) GetHeroSynergyStats(w http.ResponseWriter, r *http.Request) {
	cacheKey := "hero-synergy-stats:all"
	if h.cacheGet(r, w, cacheKey) {
		return
	}

	stats, err := h.deadlock.GetHeroSynergyStats(r.Context())
	if err != nil {
		slog.Error("failed to fetch hero synergy stats", "error", err)
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to fetch hero synergy stats"})
		return
	}

	h.cacheSet(r, cacheKey, stats, 10*time.Minute)
	h.writeJSON(w, http.StatusOK, stats)
}

func (h *Handler) GetAbilityOrderStats(w http.ResponseWriter, r *http.Request) {
	heroID, ok := h.parseHeroID(w, r)
	if !ok {
		return
	}

	cacheKey := fmt.Sprintf("ability-order-stats:%d", heroID)
	if h.cacheGet(r, w, cacheKey) {
		return
	}

	stats, err := h.deadlock.GetAbilityOrderStats(r.Context(), heroID)
	if err != nil {
		slog.Error("failed to fetch ability order stats", "hero_id", heroID, "error", err)
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to fetch ability order stats"})
		return
	}

	h.cacheSet(r, cacheKey, stats, 10*time.Minute)
	h.writeJSON(w, http.StatusOK, stats)
}

func (h *Handler) GetHeroScoreboard(w http.ResponseWriter, r *http.Request) {
	sortBy := r.URL.Query().Get("sort_by")
	if sortBy == "" {
		sortBy = "wins"
	}
	limit := h.parseLimit(r, 20, 100)

	cacheKey := fmt.Sprintf("scoreboard:heroes:%s:%d", sortBy, limit)
	if h.cacheGet(r, w, cacheKey) {
		return
	}

	scoreboard, err := h.deadlock.GetHeroScoreboard(r.Context(), sortBy, limit)
	if err != nil {
		slog.Error("failed to fetch hero scoreboard", "sort_by", sortBy, "error", err)
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to fetch hero scoreboard"})
		return
	}

	h.cacheSet(r, cacheKey, scoreboard, 10*time.Minute)
	h.writeJSON(w, http.StatusOK, scoreboard)
}

func (h *Handler) GetItems(w http.ResponseWriter, r *http.Request) {
	cacheKey := "items:all"
	if h.cacheGet(r, w, cacheKey) {
		return
	}

	items, err := h.deadlock.GetItems(r.Context())
	if err != nil {
		slog.Error("failed to fetch items", "error", err)
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to fetch items"})
		return
	}

	h.cacheSet(r, cacheKey, items, 15*time.Minute)
	h.writeJSON(w, http.StatusOK, items)
}

func (h *Handler) GetItemStats(w http.ResponseWriter, r *http.Request) {
	cacheKey := "item-stats:all"
	if h.cacheGet(r, w, cacheKey) {
		return
	}

	stats, err := h.deadlock.GetItemStats(r.Context())
	if err != nil {
		slog.Error("failed to fetch item stats", "error", err)
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to fetch item stats"})
		return
	}

	h.cacheSet(r, cacheKey, stats, 10*time.Minute)
	h.writeJSON(w, http.StatusOK, stats)
}

func (h *Handler) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	region := chi.URLParam(r, "region")

	cacheKey := "leaderboard:" + region
	if h.cacheGet(r, w, cacheKey) {
		return
	}

	leaderboard, err := h.deadlock.GetLeaderboard(r.Context(), region)
	if err != nil {
		slog.Error("failed to fetch leaderboard", "region", region, "error", err)
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to fetch leaderboard"})
		return
	}

	h.cacheSet(r, cacheKey, leaderboard, 10*time.Minute)
	h.writeJSON(w, http.StatusOK, leaderboard)
}

func (h *Handler) GetHeroLeaderboard(w http.ResponseWriter, r *http.Request) {
	region := chi.URLParam(r, "region")
	heroID, ok := h.parseHeroID(w, r)
	if !ok {
		return
	}

	cacheKey := fmt.Sprintf("leaderboard:%s:%d", region, heroID)
	if h.cacheGet(r, w, cacheKey) {
		return
	}

	leaderboard, err := h.deadlock.GetHeroLeaderboard(r.Context(), region, heroID)
	if err != nil {
		slog.Error("failed to fetch hero leaderboard", "region", region, "hero_id", heroID, "error", err)
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to fetch hero leaderboard"})
		return
	}

	h.cacheSet(r, cacheKey, leaderboard, 10*time.Minute)
	h.writeJSON(w, http.StatusOK, leaderboard)
}

func (h *Handler) GetPlayerScoreboard(w http.ResponseWriter, r *http.Request) {
	sortBy := r.URL.Query().Get("sort_by")
	if sortBy == "" {
		sortBy = "wins"
	}
	limit := h.parseLimit(r, 20, 100)

	cacheKey := fmt.Sprintf("scoreboard:players:%s:%d", sortBy, limit)
	if h.cacheGet(r, w, cacheKey) {
		return
	}

	scoreboard, err := h.deadlock.GetPlayerScoreboard(r.Context(), sortBy, limit)
	if err != nil {
		slog.Error("failed to fetch player scoreboard", "sort_by", sortBy, "error", err)
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to fetch player scoreboard"})
		return
	}

	h.cacheSet(r, cacheKey, scoreboard, 10*time.Minute)
	h.writeJSON(w, http.StatusOK, scoreboard)
}

func (h *Handler) GetGameStats(w http.ResponseWriter, r *http.Request) {
	cacheKey := "game-stats:all"
	if h.cacheGet(r, w, cacheKey) {
		return
	}

	stats, err := h.deadlock.GetGameStats(r.Context())
	if err != nil {
		slog.Error("failed to fetch game stats", "error", err)
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to fetch game stats"})
		return
	}

	h.cacheSet(r, cacheKey, stats, 10*time.Minute)
	h.writeJSON(w, http.StatusOK, stats)
}

func (h *Handler) GetKillDeathStats(w http.ResponseWriter, r *http.Request) {
	cacheKey := "kill-death-stats:all"
	if h.cacheGet(r, w, cacheKey) {
		return
	}

	stats, err := h.deadlock.GetKillDeathStats(r.Context())
	if err != nil {
		slog.Error("failed to fetch kill death stats", "error", err)
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to fetch kill death stats"})
		return
	}

	h.cacheSet(r, cacheKey, stats, 10*time.Minute)
	h.writeJSON(w, http.StatusOK, stats)
}

func (h *Handler) GetBadgeDistribution(w http.ResponseWriter, r *http.Request) {
	cacheKey := "badge-distribution:all"
	if h.cacheGet(r, w, cacheKey) {
		return
	}

	dist, err := h.deadlock.GetBadgeDistribution(r.Context())
	if err != nil {
		slog.Error("failed to fetch badge distribution", "error", err)
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to fetch badge distribution"})
		return
	}

	h.cacheSet(r, cacheKey, dist, 10*time.Minute)
	h.writeJSON(w, http.StatusOK, dist)
}

func (h *Handler) GetBuilds(w http.ResponseWriter, r *http.Request) {
	heroID, ok := h.parseHeroID(w, r)
	if !ok {
		return
	}
	limit := h.parseLimit(r, 10, 50)

	cacheKey := fmt.Sprintf("builds:%d:%d", heroID, limit)
	if h.cacheGet(r, w, cacheKey) {
		return
	}

	builds, err := h.deadlock.GetBuilds(r.Context(), heroID, limit)
	if err != nil {
		slog.Error("failed to fetch builds", "hero_id", heroID, "error", err)
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to fetch builds"})
		return
	}

	h.cacheSet(r, cacheKey, builds, 5*time.Minute)
	h.writeJSON(w, http.StatusOK, builds)
}

func (h *Handler) GetRanks(w http.ResponseWriter, r *http.Request) {
	cacheKey := "ranks:all"
	if h.cacheGet(r, w, cacheKey) {
		return
	}

	ranks, err := h.deadlock.GetRanks(r.Context())
	if err != nil {
		slog.Error("failed to fetch ranks", "error", err)
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to fetch ranks"})
		return
	}

	h.cacheSet(r, cacheKey, ranks, 60*time.Minute)
	h.writeJSON(w, http.StatusOK, ranks)
}
