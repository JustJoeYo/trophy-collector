package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/JustJoeYo/trophy-collector/internal/cache"
	"github.com/JustJoeYo/trophy-collector/internal/models"
	"github.com/go-chi/chi/v5"
)

type mockDeadlockClient struct {
    heroes []models.Hero
    err    error
}

func (m *mockDeadlockClient) GetPlayer(ctx context.Context, accountID uint32) (*models.Player, error) {
    return nil, m.err
}
func (m *mockDeadlockClient) GetPlayerMatches(ctx context.Context, accountID uint32) ([]models.Match, error) {
    return nil, m.err
}
func (m *mockDeadlockClient) GetPlayerHeroes(ctx context.Context, accountID uint32) ([]models.Hero, error) {
    return nil, m.err
}
func (m *mockDeadlockClient) GetHeroes(ctx context.Context) ([]models.Hero, error) {
    return m.heroes, m.err
}

type mockCache struct{}

func (m *mockCache) Get(ctx context.Context, key string) (string, error) {
    return "", cache.ErrCacheMiss
}

func (m *mockCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
    return nil
}
func (m *mockCache) Delete(ctx context.Context, key string) error {
    return nil
}

func TestGetHeroes_Success(t *testing.T) {
    mockClient := &mockDeadlockClient{
        heroes: []models.Hero{
            {HeroID: 1, Name: "Infernus", ClassName: "hero_inferno"},
            {HeroID: 2, Name: "Seven", ClassName: "hero_gigawatt"},
        },
    }

    handler := NewHandler(mockClient, &mockCache{})

    r := chi.NewRouter()
    handler.RegisterRoutes(r)

    req := httptest.NewRequest(http.MethodGet, "/api/v1/heroes", nil)
    rec := httptest.NewRecorder()

    r.ServeHTTP(rec, req)

    if rec.Code != http.StatusOK {
        t.Errorf("expected status 200, got %d", rec.Code)
    }

    var heroes []models.Hero
    if err := json.NewDecoder(rec.Body).Decode(&heroes); err != nil {
        t.Fatalf("failed to decode response: %v", err)
    }

    if len(heroes) != 2 {
        t.Errorf("expected 2 heroes, got %d", len(heroes))
    }
}

func TestGetHeroes_ClientError(t *testing.T) {
    mockClient := &mockDeadlockClient{
        err: fmt.Errorf("api unavailable"),
    }

    handler := NewHandler(mockClient, &mockCache{})

    r := chi.NewRouter()
    handler.RegisterRoutes(r)

    req := httptest.NewRequest(http.MethodGet, "/api/v1/heroes", nil)
    rec := httptest.NewRecorder()

    r.ServeHTTP(rec, req)

    if rec.Code != http.StatusInternalServerError {
        t.Errorf("expected status 500, got %d", rec.Code)
    }
}
