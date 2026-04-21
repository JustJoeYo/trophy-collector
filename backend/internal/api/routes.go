package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		h.writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	r.Get("/api/v1/players/{id}/profile", h.GetPlayerProfile)
	r.Get("/api/v1/players/{id}/sync-status", h.GetSyncStatus)
	r.Get("/api/v1/players/{id}/matches", h.GetPlayerMatches)
	r.Get("/api/v1/players/{id}/stats", h.GetPlayerStats)
	r.Get("/api/v1/players/{id}/metrics", h.GetPlayerMetrics)
	r.Get("/api/v1/players/{id}/active", h.GetActiveMatches)

	r.Get("/api/v1/heroes", h.GetHeroes)
	r.Get("/api/v1/heroes/stats", h.GetHeroStats)
	r.Get("/api/v1/heroes/ban-stats", h.GetHeroBanStats)
	r.Get("/api/v1/heroes/counter-stats", h.GetHeroCounterStats)
	r.Get("/api/v1/heroes/synergy-stats", h.GetHeroSynergyStats)
	r.Get("/api/v1/heroes/{heroId}/build-stats", h.GetHeroBuildStats)
	r.Get("/api/v1/heroes/{heroId}/ability-order-stats", h.GetAbilityOrderStats)
	r.Get("/api/v1/heroes/{heroId}/builds", h.GetBuilds)

	r.Get("/api/v1/items", h.GetItems)
	r.Get("/api/v1/items/stats", h.GetItemStats)

	r.Get("/api/v1/leaderboard/{region}", h.GetLeaderboard)
	r.Get("/api/v1/leaderboard/{region}/{heroId}", h.GetHeroLeaderboard)

	r.Get("/api/v1/scoreboard/players", h.GetPlayerScoreboard)
	r.Get("/api/v1/scoreboard/heroes", h.GetHeroScoreboard)

	r.Get("/api/v1/analytics/game-stats", h.GetGameStats)
	r.Get("/api/v1/analytics/kill-death-stats", h.GetKillDeathStats)
	r.Get("/api/v1/analytics/badge-distribution", h.GetBadgeDistribution)

	r.Get("/api/v1/ranks", h.GetRanks)
}
