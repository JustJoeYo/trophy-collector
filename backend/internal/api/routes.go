package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/JustJoeYo/trophy-collector/internal/config"
)

// NewRouter builds and returns the v1 API router.
// Each handler will be wired up here as we build them out.
func NewRouter(cfg *config.Config) http.Handler {
	r := chi.NewRouter()

	// Log every request with method, path, status, and duration
	r.Use(middleware.Logger)

	h := NewHandler(cfg)

	// Player routes
	r.Get("/player/{steamid}", h.GetPlayer)
	r.Get("/player/{steamid}/matches", h.GetPlayerMatches)
	r.Get("/player/{steamid}/heroes", h.GetPlayerHeroes)

	// Game data routes
	r.Get("/heroes", h.GetHeroes)
	r.Get("/leaderboard", h.GetLeaderboard)

	return r
}
