package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/JustJoeYo/trophy-collector/internal/cache"
	"github.com/JustJoeYo/trophy-collector/internal/clients"
)

type Handler struct {
	deadlock clients.DeadlockClient
	cache    cache.Cache
}

func NewHandler(deadlock clients.DeadlockClient, cache cache.Cache) *Handler {
	return &Handler{
		deadlock: deadlock,
		cache:    cache,
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

	limit := h.parseLimit(r, 10, 50)
	cacheKey := fmt.Sprintf("matches:%d:%d", accountID, limit)
	if h.cacheGet(r, w, cacheKey) {
		return
	}

	matches, err := h.deadlock.GetPlayerMatches(r.Context(), accountID, limit)
	if err != nil {
		slog.Error("failed to fetch matches", "account_id", accountID, "error", err)
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to fetch matches"})
		return
	}

	h.cacheSet(r, cacheKey, matches, 5*time.Minute)
	h.writeJSON(w, http.StatusOK, matches)
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
