package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		h.writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	
	r.Get("/api/v1/players/{id}", h.GetPlayer)
    r.Get("/api/v1/players/{id}/matches", h.GetPlayerMatches)
    r.Get("/api/v1/players/{id}/heroes", h.GetPlayerHeroes)
    r.Get("/api/v1/heroes", h.GetHeroes)
}