package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/JustJoeYo/trophy-collector/internal/cache"
	"github.com/JustJoeYo/trophy-collector/internal/models"
)

type mockDeadlockClient struct {
	heroes      []models.Hero
	items       []models.Item
	matches     []models.Match
	leaderboard *models.Leaderboard
	heroStats   []models.HeroStats
	scoreboard  []models.PlayerScoreboard
	ranks       []models.Rank
	err         error
}

func (m *mockDeadlockClient) GetPlayerMatches(ctx context.Context, accountID uint32, limit int) ([]models.Match, error) {
	return m.matches, m.err
}
func (m *mockDeadlockClient) GetPlayerMetrics(ctx context.Context, accountID uint32) (*models.PlayerMetrics, error) {
	return nil, m.err
}
func (m *mockDeadlockClient) GetActiveMatches(ctx context.Context, accountIDs []uint32) ([]models.Match, error) {
	return m.matches, m.err
}
func (m *mockDeadlockClient) GetHeroes(ctx context.Context) ([]models.Hero, error) {
	return m.heroes, m.err
}
func (m *mockDeadlockClient) GetHeroStats(ctx context.Context) ([]models.HeroStats, error) {
	return m.heroStats, m.err
}
func (m *mockDeadlockClient) GetHeroBanStats(ctx context.Context) ([]models.HeroBanStats, error) {
	return nil, m.err
}
func (m *mockDeadlockClient) GetHeroBuildStats(ctx context.Context, heroID uint32) ([]models.HeroBuildStats, error) {
	return nil, m.err
}
func (m *mockDeadlockClient) GetHeroCounterStats(ctx context.Context) ([]models.HeroCounterStats, error) {
	return nil, m.err
}
func (m *mockDeadlockClient) GetHeroSynergyStats(ctx context.Context) ([]models.HeroSynergyStats, error) {
	return nil, m.err
}
func (m *mockDeadlockClient) GetAbilityOrderStats(ctx context.Context, heroID uint32) ([]models.AbilityOrderStats, error) {
	return nil, m.err
}
func (m *mockDeadlockClient) GetHeroScoreboard(ctx context.Context, sortBy string, limit int) ([]models.HeroScoreboard, error) {
	return nil, m.err
}
func (m *mockDeadlockClient) GetItems(ctx context.Context) ([]models.Item, error) {
	return m.items, m.err
}
func (m *mockDeadlockClient) GetItemStats(ctx context.Context) ([]models.ItemStats, error) {
	return nil, m.err
}
func (m *mockDeadlockClient) GetLeaderboard(ctx context.Context, region string) (*models.Leaderboard, error) {
	return m.leaderboard, m.err
}
func (m *mockDeadlockClient) GetHeroLeaderboard(ctx context.Context, region string, heroID uint32) (*models.Leaderboard, error) {
	return m.leaderboard, m.err
}
func (m *mockDeadlockClient) GetPlayerScoreboard(ctx context.Context, sortBy string, limit int) ([]models.PlayerScoreboard, error) {
	return m.scoreboard, m.err
}
func (m *mockDeadlockClient) GetGameStats(ctx context.Context) ([]models.GameStats, error) {
	return nil, m.err
}
func (m *mockDeadlockClient) GetKillDeathStats(ctx context.Context) ([]models.KillDeathStats, error) {
	return nil, m.err
}
func (m *mockDeadlockClient) GetBadgeDistribution(ctx context.Context) ([]models.BadgeDistribution, error) {
	return nil, m.err
}
func (m *mockDeadlockClient) GetBuilds(ctx context.Context, heroID uint32, limit int) ([]models.Build, error) {
	return nil, m.err
}
func (m *mockDeadlockClient) GetRanks(ctx context.Context) ([]models.Rank, error) {
	return m.ranks, m.err
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

func TestGetPlayerMatches_InvalidID(t *testing.T) {
	handler := NewHandler(&mockDeadlockClient{}, &mockCache{})
	r := chi.NewRouter()
	handler.RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/players/notanid/matches", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestGetLeaderboard_Success(t *testing.T) {
	mockClient := &mockDeadlockClient{
		leaderboard: &models.Leaderboard{
			Entries: []models.LeaderboardEntry{
				{AccountName: "TestPlayer", Rank: 1, BadgeLevel: 10},
			},
		},
	}

	handler := NewHandler(mockClient, &mockCache{})
	r := chi.NewRouter()
	handler.RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/leaderboard/Europe", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}
