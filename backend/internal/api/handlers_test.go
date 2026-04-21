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

func newRouter(client *mockDeadlockClient) *chi.Mux {
	r := chi.NewRouter()
	NewHandler(client, &mockCache{}, nil).RegisterRoutes(r)
	return r
}

type mockDeadlockClient struct {
	heroes        []models.Hero
	items         []models.Item
	matches       []models.Match
	leaderboard   *models.Leaderboard
	heroStats     []models.HeroStats
	heroBanStats  []models.HeroBanStats
	heroBuildStats []models.HeroBuildStats
	heroCounters  []models.HeroCounterStats
	heroSynergies []models.HeroSynergyStats
	abilityStats  []models.AbilityOrderStats
	heroScoreboard []models.HeroScoreboard
	itemStats     []models.ItemStats
	scoreboard    []models.PlayerScoreboard
	gameStats     []models.GameStats
	kdStats       []models.KillDeathStats
	badgeDist     []models.BadgeDistribution
	builds        []models.Build
	ranks         []models.Rank
	metrics       *models.PlayerMetrics
	err           error
}

func (m *mockDeadlockClient) GetPlayerMatches(ctx context.Context, accountID uint32, limit int) ([]models.Match, error) {
	return m.matches, m.err
}
func (m *mockDeadlockClient) GetPlayerMetrics(ctx context.Context, accountID uint32) (*models.PlayerMetrics, error) {
	return m.metrics, m.err
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
	return m.heroBanStats, m.err
}
func (m *mockDeadlockClient) GetHeroBuildStats(ctx context.Context, heroID uint32) ([]models.HeroBuildStats, error) {
	return m.heroBuildStats, m.err
}
func (m *mockDeadlockClient) GetHeroCounterStats(ctx context.Context) ([]models.HeroCounterStats, error) {
	return m.heroCounters, m.err
}
func (m *mockDeadlockClient) GetHeroSynergyStats(ctx context.Context) ([]models.HeroSynergyStats, error) {
	return m.heroSynergies, m.err
}
func (m *mockDeadlockClient) GetAbilityOrderStats(ctx context.Context, heroID uint32) ([]models.AbilityOrderStats, error) {
	return m.abilityStats, m.err
}
func (m *mockDeadlockClient) GetHeroScoreboard(ctx context.Context, sortBy string, limit int) ([]models.HeroScoreboard, error) {
	return m.heroScoreboard, m.err
}
func (m *mockDeadlockClient) GetItems(ctx context.Context) ([]models.Item, error) {
	return m.items, m.err
}
func (m *mockDeadlockClient) GetItemStats(ctx context.Context) ([]models.ItemStats, error) {
	return m.itemStats, m.err
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
	return m.gameStats, m.err
}
func (m *mockDeadlockClient) GetKillDeathStats(ctx context.Context) ([]models.KillDeathStats, error) {
	return m.kdStats, m.err
}
func (m *mockDeadlockClient) GetBadgeDistribution(ctx context.Context) ([]models.BadgeDistribution, error) {
	return m.badgeDist, m.err
}
func (m *mockDeadlockClient) GetBuilds(ctx context.Context, heroID uint32, limit int) ([]models.Build, error) {
	return m.builds, m.err
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

func (m *mockDeadlockClient) GetPlayerMatchesPage(ctx context.Context, accountID uint32, minMatchID *uint64, limit int, since *time.Time) ([]models.Match, error) {
	return m.matches, m.err
}

func assertStatus(t *testing.T, rec *httptest.ResponseRecorder, expected int) {
	t.Helper()
	if rec.Code != expected {
		t.Errorf("expected status %d, got %d", expected, rec.Code)
	}
}

func assertJSONArrayLen(t *testing.T, rec *httptest.ResponseRecorder, expected int) {
	t.Helper()
	var result []interface{}
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(result) != expected {
		t.Errorf("expected %d items, got %d", expected, len(result))
	}
}

func assertErrorBody(t *testing.T, rec *httptest.ResponseRecorder) {
	t.Helper()
	var result map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}
	if _, ok := result["error"]; !ok {
		t.Error("expected error field in response body")
	}
}

func TestHealth(t *testing.T) {
	r := newRouter(&mockDeadlockClient{})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusOK)
}

func TestGetHeroes_Success(t *testing.T) {
	r := newRouter(&mockDeadlockClient{
		heroes: []models.Hero{
			{HeroID: 1, Name: "Infernus", ClassName: "hero_inferno"},
			{HeroID: 2, Name: "Seven", ClassName: "hero_gigawatt"},
		},
	})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/heroes", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusOK)
	assertJSONArrayLen(t, rec, 2)
}

func TestGetHeroes_ClientError(t *testing.T) {
	r := newRouter(&mockDeadlockClient{err: fmt.Errorf("api down")})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/heroes", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusInternalServerError)
	assertErrorBody(t, rec)
}

func TestGetHeroStats_Success(t *testing.T) {
	r := newRouter(&mockDeadlockClient{
		heroStats: []models.HeroStats{{HeroID: 1, Wins: 100, Losses: 50}},
	})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/heroes/stats", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusOK)
	assertJSONArrayLen(t, rec, 1)
}

func TestGetHeroStats_ClientError(t *testing.T) {
	r := newRouter(&mockDeadlockClient{err: fmt.Errorf("api down")})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/heroes/stats", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusInternalServerError)
	assertErrorBody(t, rec)
}

func TestGetHeroBanStats_Success(t *testing.T) {
	r := newRouter(&mockDeadlockClient{
		heroBanStats: []models.HeroBanStats{{HeroID: 1, Bans: 500}},
	})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/heroes/ban-stats", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusOK)
	assertJSONArrayLen(t, rec, 1)
}

func TestGetHeroBanStats_ClientError(t *testing.T) {
	r := newRouter(&mockDeadlockClient{err: fmt.Errorf("api down")})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/heroes/ban-stats", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusInternalServerError)
	assertErrorBody(t, rec)
}

func TestGetHeroBuildStats_Success(t *testing.T) {
	r := newRouter(&mockDeadlockClient{
		heroBuildStats: []models.HeroBuildStats{{HeroID: 1, HeroBuildID: 100, Wins: 50}},
	})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/heroes/1/build-stats", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusOK)
	assertJSONArrayLen(t, rec, 1)
}

func TestGetHeroBuildStats_InvalidHeroID(t *testing.T) {
	r := newRouter(&mockDeadlockClient{})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/heroes/notanid/build-stats", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusBadRequest)
	assertErrorBody(t, rec)
}

func TestGetHeroBuildStats_ClientError(t *testing.T) {
	r := newRouter(&mockDeadlockClient{err: fmt.Errorf("api down")})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/heroes/1/build-stats", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusInternalServerError)
	assertErrorBody(t, rec)
}

func TestGetHeroCounterStats_Success(t *testing.T) {
	r := newRouter(&mockDeadlockClient{
		heroCounters: []models.HeroCounterStats{{HeroID: 1, EnemyHeroID: 2, Wins: 100}},
	})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/heroes/counter-stats", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusOK)
	assertJSONArrayLen(t, rec, 1)
}

func TestGetHeroCounterStats_ClientError(t *testing.T) {
	r := newRouter(&mockDeadlockClient{err: fmt.Errorf("api down")})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/heroes/counter-stats", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusInternalServerError)
	assertErrorBody(t, rec)
}

func TestGetHeroSynergyStats_Success(t *testing.T) {
	r := newRouter(&mockDeadlockClient{
		heroSynergies: []models.HeroSynergyStats{{HeroID1: 1, HeroID2: 2, Wins: 200}},
	})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/heroes/synergy-stats", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusOK)
	assertJSONArrayLen(t, rec, 1)
}

func TestGetHeroSynergyStats_ClientError(t *testing.T) {
	r := newRouter(&mockDeadlockClient{err: fmt.Errorf("api down")})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/heroes/synergy-stats", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusInternalServerError)
	assertErrorBody(t, rec)
}

func TestGetAbilityOrderStats_Success(t *testing.T) {
	r := newRouter(&mockDeadlockClient{
		abilityStats: []models.AbilityOrderStats{{Wins: 50, Losses: 30}},
	})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/heroes/1/ability-order-stats", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusOK)
	assertJSONArrayLen(t, rec, 1)
}

func TestGetAbilityOrderStats_InvalidHeroID(t *testing.T) {
	r := newRouter(&mockDeadlockClient{})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/heroes/notanid/ability-order-stats", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusBadRequest)
	assertErrorBody(t, rec)
}

func TestGetAbilityOrderStats_ClientError(t *testing.T) {
	r := newRouter(&mockDeadlockClient{err: fmt.Errorf("api down")})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/heroes/1/ability-order-stats", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusInternalServerError)
	assertErrorBody(t, rec)
}

func TestGetHeroScoreboard_Success(t *testing.T) {
	r := newRouter(&mockDeadlockClient{
		heroScoreboard: []models.HeroScoreboard{{Rank: 1, HeroID: 1, Value: 500}},
	})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/scoreboard/heroes?sort_by=wins", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusOK)
	assertJSONArrayLen(t, rec, 1)
}

func TestGetHeroScoreboard_DefaultSortBy(t *testing.T) {
	r := newRouter(&mockDeadlockClient{
		heroScoreboard: []models.HeroScoreboard{{Rank: 1, HeroID: 1, Value: 500}},
	})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/scoreboard/heroes", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusOK)
}

func TestGetHeroScoreboard_ClientError(t *testing.T) {
	r := newRouter(&mockDeadlockClient{err: fmt.Errorf("api down")})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/scoreboard/heroes", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusInternalServerError)
	assertErrorBody(t, rec)
}

func TestGetItems_Success(t *testing.T) {
	r := newRouter(&mockDeadlockClient{
		items: []models.Item{{ItemID: 1, Name: "Basic Magazine"}},
	})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/items", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusOK)
	assertJSONArrayLen(t, rec, 1)
}

func TestGetItems_ClientError(t *testing.T) {
	r := newRouter(&mockDeadlockClient{err: fmt.Errorf("api down")})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/items", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusInternalServerError)
	assertErrorBody(t, rec)
}

func TestGetItemStats_Success(t *testing.T) {
	r := newRouter(&mockDeadlockClient{
		itemStats: []models.ItemStats{{ItemID: 1, Wins: 1000, Losses: 800}},
	})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/items/stats", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusOK)
	assertJSONArrayLen(t, rec, 1)
}

func TestGetItemStats_ClientError(t *testing.T) {
	r := newRouter(&mockDeadlockClient{err: fmt.Errorf("api down")})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/items/stats", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusInternalServerError)
	assertErrorBody(t, rec)
}

func TestGetPlayerProfile_Success(t *testing.T) {
	r := newRouter(&mockDeadlockClient{
		matches: []models.Match{
			{
				MatchID:     1,
				WinningTeam: "Team0",
				DurationS:   2000,
				GameMode:    "Normal",
				StartTime:   "2024-11-02 03:40:37",
				Players: []models.MatchPlayer{
					{AccountID: 12345, HeroID: 1, Team: "Team0", Kills: 8, Deaths: 3, Assists: 5, NetWorth: 40000, LastHits: 200, AssignedLane: 1},
					{AccountID: 99999, HeroID: 2, Team: "Team0", Kills: 5, Deaths: 4, Assists: 8},
				},
			},
			{
				MatchID:     2,
				WinningTeam: "Team1",
				DurationS:   1800,
				GameMode:    "Normal",
				StartTime:   "2024-11-03 05:00:00",
				Players: []models.MatchPlayer{
					{AccountID: 12345, HeroID: 1, Team: "Team0", Kills: 4, Deaths: 6, Assists: 3, NetWorth: 30000, LastHits: 150, AssignedLane: 1},
					{AccountID: 99999, HeroID: 2, Team: "Team1", Kills: 9, Deaths: 2, Assists: 6},
				},
			},
		},
	})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/players/12345/profile", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusOK)

	var profile models.PlayerProfile
	if err := json.NewDecoder(rec.Body).Decode(&profile); err != nil {
		t.Fatalf("failed to decode profile: %v", err)
	}
	if profile.Overview.Matches != 2 {
		t.Errorf("expected 2 matches, got %d", profile.Overview.Matches)
	}
	if profile.Overview.Wins != 1 {
		t.Errorf("expected 1 win, got %d", profile.Overview.Wins)
	}
	if len(profile.Heroes) != 1 {
		t.Errorf("expected 1 hero entry, got %d", len(profile.Heroes))
	}
	if len(profile.Lanes) != 1 {
		t.Errorf("expected 1 lane entry, got %d", len(profile.Lanes))
	}
	if profile.Awards.MostKills.Value != 8 {
		t.Errorf("expected most kills = 8, got %f", profile.Awards.MostKills.Value)
	}
	if len(profile.RecentMatches) != 2 {
		t.Errorf("expected 2 recent matches, got %d", len(profile.RecentMatches))
	}
}

func TestGetPlayerProfile_NoMatches(t *testing.T) {
	r := newRouter(&mockDeadlockClient{matches: []models.Match{}})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/players/12345/profile", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusNotFound)
	assertErrorBody(t, rec)
}

func TestGetPlayerProfile_InvalidID(t *testing.T) {
	r := newRouter(&mockDeadlockClient{})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/players/notanid/profile", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusBadRequest)
	assertErrorBody(t, rec)
}

func TestGetPlayerProfile_ClientError(t *testing.T) {
	r := newRouter(&mockDeadlockClient{err: fmt.Errorf("api down")})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/players/12345/profile", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusInternalServerError)
	assertErrorBody(t, rec)
}

func TestGetPlayerMatches_Success(t *testing.T) {
	r := newRouter(&mockDeadlockClient{
		matches: []models.Match{
			{
				MatchID:     123,
				GameMode:    "Normal",
				WinningTeam: "Team0",
				DurationS:   2000,
				StartTime:   "2024-11-02 03:40:37",
				Players: []models.MatchPlayer{
					{AccountID: 12345, HeroID: 1, Team: "Team0", Kills: 5, Deaths: 3, Assists: 7},
				},
			},
		},
	})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/players/12345/matches", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusOK)
	assertJSONArrayLen(t, rec, 1)
}

func TestGetPlayerMatches_PlayerNotInMatch(t *testing.T) {
	r := newRouter(&mockDeadlockClient{
		matches: []models.Match{
			{
				MatchID: 123,
				Players: []models.MatchPlayer{
					{AccountID: 99999, HeroID: 1, Team: "Team0"},
				},
			},
		},
	})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/players/12345/matches", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusOK)
	assertJSONArrayLen(t, rec, 0)
}

func TestGetPlayerStats_Success(t *testing.T) {
	r := newRouter(&mockDeadlockClient{
		matches: []models.Match{
			{
				MatchID:     1,
				WinningTeam: "Team0",
				DurationS:   2000,
				Players: []models.MatchPlayer{
					{AccountID: 12345, Team: "Team0", Kills: 10, Deaths: 2, Assists: 5},
				},
			},
			{
				MatchID:     2,
				WinningTeam: "Team1",
				DurationS:   1800,
				Players: []models.MatchPlayer{
					{AccountID: 12345, Team: "Team0", Kills: 4, Deaths: 6, Assists: 3},
				},
			},
		},
	})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/players/12345/stats", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusOK)

	var stats models.PlayerStats
	if err := json.NewDecoder(rec.Body).Decode(&stats); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if stats.MatchesSampled != 2 {
		t.Errorf("expected 2 matches sampled, got %d", stats.MatchesSampled)
	}
	if stats.Wins != 1 {
		t.Errorf("expected 1 win, got %d", stats.Wins)
	}
	if stats.Losses != 1 {
		t.Errorf("expected 1 loss, got %d", stats.Losses)
	}
	if stats.TotalKills != 14 {
		t.Errorf("expected 14 total kills, got %d", stats.TotalKills)
	}
}

func TestGetPlayerStats_NoMatches(t *testing.T) {
	r := newRouter(&mockDeadlockClient{matches: []models.Match{}})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/players/12345/stats", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusNotFound)
	assertErrorBody(t, rec)
}

func TestGetPlayerStats_InvalidID(t *testing.T) {
	r := newRouter(&mockDeadlockClient{})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/players/notanid/stats", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusBadRequest)
	assertErrorBody(t, rec)
}

func TestGetPlayerStats_ClientError(t *testing.T) {
	r := newRouter(&mockDeadlockClient{err: fmt.Errorf("api down")})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/players/12345/stats", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusInternalServerError)
	assertErrorBody(t, rec)
}

func TestGetPlayerStats_ZeroDeaths(t *testing.T) {
	r := newRouter(&mockDeadlockClient{
		matches: []models.Match{
			{
				MatchID:     1,
				WinningTeam: "Team0",
				Players: []models.MatchPlayer{
					{AccountID: 12345, Team: "Team0", Kills: 10, Deaths: 0, Assists: 5},
				},
			},
		},
	})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/players/12345/stats", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusOK)

	var stats models.PlayerStats
	if err := json.NewDecoder(rec.Body).Decode(&stats); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if stats.KDA != 15.0 {
		t.Errorf("expected KDA 15.0 (deaths clamped to 1), got %f", stats.KDA)
	}
}

func TestGetPlayerMatches_InvalidID(t *testing.T) {
	r := newRouter(&mockDeadlockClient{})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/players/notanid/matches", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusBadRequest)
	assertErrorBody(t, rec)
}

func TestGetPlayerMatches_ClientError(t *testing.T) {
	r := newRouter(&mockDeadlockClient{err: fmt.Errorf("api down")})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/players/12345/matches", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusInternalServerError)
	assertErrorBody(t, rec)
}

func TestGetPlayerMatches_LimitClamped(t *testing.T) {
	r := newRouter(&mockDeadlockClient{matches: []models.Match{}})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/players/12345/matches?limit=9999", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusOK)
}

func TestGetPlayerMetrics_Success(t *testing.T) {
	r := newRouter(&mockDeadlockClient{
		metrics: &models.PlayerMetrics{},
	})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/players/12345/metrics", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusOK)
}

func TestGetPlayerMetrics_InvalidID(t *testing.T) {
	r := newRouter(&mockDeadlockClient{})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/players/notanid/metrics", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusBadRequest)
	assertErrorBody(t, rec)
}

func TestGetPlayerMetrics_ClientError(t *testing.T) {
	r := newRouter(&mockDeadlockClient{err: fmt.Errorf("api down")})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/players/12345/metrics", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusInternalServerError)
	assertErrorBody(t, rec)
}

func TestGetActiveMatches_Success(t *testing.T) {
	r := newRouter(&mockDeadlockClient{matches: []models.Match{}})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/players/12345/active", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusOK)
}

func TestGetActiveMatches_InvalidID(t *testing.T) {
	r := newRouter(&mockDeadlockClient{})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/players/notanid/active", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusBadRequest)
	assertErrorBody(t, rec)
}

func TestGetActiveMatches_ClientError(t *testing.T) {
	r := newRouter(&mockDeadlockClient{err: fmt.Errorf("api down")})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/players/12345/active", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusInternalServerError)
	assertErrorBody(t, rec)
}

func TestGetLeaderboard_Success(t *testing.T) {
	r := newRouter(&mockDeadlockClient{
		leaderboard: &models.Leaderboard{
			Entries: []models.LeaderboardEntry{{AccountName: "Player1", Rank: 1}},
		},
	})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/leaderboard/Europe", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusOK)
}

func TestGetLeaderboard_ClientError(t *testing.T) {
	r := newRouter(&mockDeadlockClient{err: fmt.Errorf("api down")})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/leaderboard/Europe", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusInternalServerError)
	assertErrorBody(t, rec)
}

func TestGetHeroLeaderboard_Success(t *testing.T) {
	r := newRouter(&mockDeadlockClient{
		leaderboard: &models.Leaderboard{
			Entries: []models.LeaderboardEntry{{AccountName: "Player1", Rank: 1}},
		},
	})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/leaderboard/Europe/1", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusOK)
}

func TestGetHeroLeaderboard_InvalidHeroID(t *testing.T) {
	r := newRouter(&mockDeadlockClient{})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/leaderboard/Europe/notanid", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusBadRequest)
	assertErrorBody(t, rec)
}

func TestGetHeroLeaderboard_ClientError(t *testing.T) {
	r := newRouter(&mockDeadlockClient{err: fmt.Errorf("api down")})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/leaderboard/Europe/1", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusInternalServerError)
	assertErrorBody(t, rec)
}

func TestGetPlayerScoreboard_Success(t *testing.T) {
	r := newRouter(&mockDeadlockClient{
		scoreboard: []models.PlayerScoreboard{{Rank: 1, AccountID: 123, Value: 500}},
	})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/scoreboard/players?sort_by=wins", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusOK)
	assertJSONArrayLen(t, rec, 1)
}

func TestGetPlayerScoreboard_DefaultSortBy(t *testing.T) {
	r := newRouter(&mockDeadlockClient{
		scoreboard: []models.PlayerScoreboard{{Rank: 1, AccountID: 123}},
	})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/scoreboard/players", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusOK)
}

func TestGetPlayerScoreboard_ClientError(t *testing.T) {
	r := newRouter(&mockDeadlockClient{err: fmt.Errorf("api down")})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/scoreboard/players", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusInternalServerError)
	assertErrorBody(t, rec)
}

func TestGetGameStats_Success(t *testing.T) {
	r := newRouter(&mockDeadlockClient{
		gameStats: []models.GameStats{{TotalMatches: 1000, AvgDurationS: 2200}},
	})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/analytics/game-stats", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusOK)
	assertJSONArrayLen(t, rec, 1)
}

func TestGetGameStats_ClientError(t *testing.T) {
	r := newRouter(&mockDeadlockClient{err: fmt.Errorf("api down")})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/analytics/game-stats", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusInternalServerError)
	assertErrorBody(t, rec)
}

func TestGetKillDeathStats_Success(t *testing.T) {
	r := newRouter(&mockDeadlockClient{
		kdStats: []models.KillDeathStats{{PositionX: 100, PositionY: 200, Kills: 50}},
	})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/analytics/kill-death-stats", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusOK)
	assertJSONArrayLen(t, rec, 1)
}

func TestGetKillDeathStats_ClientError(t *testing.T) {
	r := newRouter(&mockDeadlockClient{err: fmt.Errorf("api down")})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/analytics/kill-death-stats", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusInternalServerError)
	assertErrorBody(t, rec)
}

func TestGetBadgeDistribution_Success(t *testing.T) {
	r := newRouter(&mockDeadlockClient{
		badgeDist: []models.BadgeDistribution{{BadgeLevel: 10, TotalMatches: 5000}},
	})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/analytics/badge-distribution", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusOK)
	assertJSONArrayLen(t, rec, 1)
}

func TestGetBadgeDistribution_ClientError(t *testing.T) {
	r := newRouter(&mockDeadlockClient{err: fmt.Errorf("api down")})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/analytics/badge-distribution", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusInternalServerError)
	assertErrorBody(t, rec)
}

func TestGetBuilds_Success(t *testing.T) {
	r := newRouter(&mockDeadlockClient{
		builds: []models.Build{{HeroBuild: models.BuildDetails{HeroID: 1, Name: "Test Build"}}},
	})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/heroes/1/builds", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusOK)
	assertJSONArrayLen(t, rec, 1)
}

func TestGetBuilds_InvalidHeroID(t *testing.T) {
	r := newRouter(&mockDeadlockClient{})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/heroes/notanid/builds", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusBadRequest)
	assertErrorBody(t, rec)
}

func TestGetBuilds_ClientError(t *testing.T) {
	r := newRouter(&mockDeadlockClient{err: fmt.Errorf("api down")})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/heroes/1/builds", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusInternalServerError)
	assertErrorBody(t, rec)
}

func TestGetRanks_Success(t *testing.T) {
	r := newRouter(&mockDeadlockClient{
		ranks: []models.Rank{{Tier: 1, Name: "Initiate"}},
	})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/ranks", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusOK)
	assertJSONArrayLen(t, rec, 1)
}

func TestGetRanks_ClientError(t *testing.T) {
	r := newRouter(&mockDeadlockClient{err: fmt.Errorf("api down")})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/ranks", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assertStatus(t, rec, http.StatusInternalServerError)
	assertErrorBody(t, rec)
}
