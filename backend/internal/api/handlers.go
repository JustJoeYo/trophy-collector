package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/JustJoeYo/trophy-collector/internal/cache"
	"github.com/JustJoeYo/trophy-collector/internal/clients"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	deadlock clients.DeadlockClient
	cache cache.Cache
}

func NewHandler(deadlock clients.DeadlockClient, cache cache.Cache) *Handler {
	return &Handler{
		deadlock: deadlock,
		cache: cache,
	}
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode response", "error", err)
	}
}

func (h *Handler) GetPlayer (w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid account id"})
		return
	}
	accountID := uint32(id)

	cacheKey := "player:" + idStr
	if cached, err := h.cache.Get(r.Context(), cacheKey); err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(cached))
		return
	}

	player, err := h.deadlock.GetPlayer(r.Context(), accountID)
	if err != nil {
		slog.Error("failed to fetch player", "account_id", accountID, "error", err)
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to fetch player"})
		return
	}
	
	if data, err := json.Marshal(player); err == nil {
		h.cache.Set(r.Context(), cacheKey, string(data), 5*time.Minute)
	}

	h.writeJSON(w, http.StatusOK, player)
}	

func (h *Handler) GetPlayerMatches(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid account id"})
        return
    }
    accountID := uint32(id)

	cacheKey := "matches:" + idStr
    if cached, err := h.cache.Get(r.Context(), cacheKey); err == nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(cached))
        return
    }

	matches, err := h.deadlock.GetPlayerMatches(r.Context(), accountID)
    if err != nil {
        slog.Error("failed to fetch matches", "account_id", accountID, "error", err)
        h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to fetch matches"})
        return
    }

    if data, err := json.Marshal(matches); err == nil {
        h.cache.Set(r.Context(), cacheKey, string(data), 5*time.Minute)
    }

    h.writeJSON(w, http.StatusOK, matches)
}

func (h *Handler) GetPlayerHeroes(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid account id"})
        return
    }
    accountID := uint32(id)

    cacheKey := "heroes:" + idStr
    if cached, err := h.cache.Get(r.Context(), cacheKey); err == nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(cached))
        return
    }

    heroes, err := h.deadlock.GetPlayerHeroes(r.Context(), accountID)
    if err != nil {
        slog.Error("failed to fetch player heroes", "account_id", accountID, "error", err)
        h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to fetch player heroes"})
        return
    }

    if data, err := json.Marshal(heroes); err == nil {
        h.cache.Set(r.Context(), cacheKey, string(data), 5*time.Minute)
    }

    h.writeJSON(w, http.StatusOK, heroes)
}

func (h *Handler) GetHeroes(w http.ResponseWriter, r *http.Request) {
    cacheKey := "heroes:all"
    if cached, err := h.cache.Get(r.Context(), cacheKey); err == nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(cached))
        return
    }

    heroes, err := h.deadlock.GetHeroes(r.Context())
    if err != nil {
        slog.Error("failed to fetch heroes", "error", err)
        h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to fetch heroes"})
        return
    }

    if data, err := json.Marshal(heroes); err == nil {
        h.cache.Set(r.Context(), cacheKey, string(data), 15*time.Minute)
    }

    h.writeJSON(w, http.StatusOK, heroes)
}
