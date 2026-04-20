package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/JustJoeYo/trophy-collector/internal/config"
	"github.com/JustJoeYo/trophy-collector/internal/models"
)

// Handler holds dependencies for all HTTP handlers.
// We pass config in here — later we'll also pass the
// cache client and API clients through this struct.
type Handler struct {
	cfg *config.Config
}

func NewHandler(cfg *config.Config) *Handler {
	return &Handler{cfg: cfg}
}

// GetPlayer returns a player's profile by Steam ID
func (h *Handler) GetPlayer(w http.ResponseWriter, r *http.Request) {
	steamID := chi.URLParam(r, "steamid")

	// TODO: call Steam client + Deadlock API client
	// For now return a placeholder so the route works end to end
	slog.Info("GetPlayer called", "steam_id", steamID)

	player := models.Player{
		SteamID:     steamID,
		PersonaName: "placeholder",
		AvatarURL:   "",
		ProfileURL:  "",
	}

	writeJSON(w, http.StatusOK, player)
}

// GetPlayerMatches returns recent match history for a player
func (h *Handler) GetPlayerMatches(w http.ResponseWriter, r *http.Request) {
	steamID := chi.URLParam(r, "steamid")
	slog.Info("GetPlayerMatches called", "steam_id", steamID)

	// TODO: call Deadlock API client
	writeJSON(w, http.StatusOK, []models.Match{})
}

// GetPlayerHeroes returns per-hero stats for a player
func (h *Handler) GetPlayerHeroes(w http.ResponseWriter, r *http.Request) {
	steamID := chi.URLParam(r, "steamid")
	slog.Info("GetPlayerHeroes called", "steam_id", steamID)

	// TODO: call Deadlock API client
	writeJSON(w, http.StatusOK, []models.HeroStats{})
}

// GetHeroes returns the full hero tier list
func (h *Handler) GetHeroes(w http.ResponseWriter, r *http.Request) {
	slog.Info("GetHeroes called")

	// TODO: call Deadlock API client
	writeJSON(w, http.StatusOK, []models.Hero{})
}

// GetLeaderboard returns top ranked players
func (h *Handler) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	slog.Info("GetLeaderboard called")

	// TODO: call Deadlock API client
	writeJSON(w, http.StatusOK, []models.Player{})
}

// writeJSON is a helper that writes a JSON response with the correct headers.
// Centralizing this means every handler responds consistently.
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode response", "error", err)
	}
}

// writeError writes a consistent JSON error response
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, models.APIError{
		Error:   http.StatusText(status),
		Message: message,
		Code:    status,
	})
}
