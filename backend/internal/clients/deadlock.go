package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/JustJoeYo/trophy-collector/internal/models"
)

type DeadlockClient interface {
	GetPlayerMatches(ctx context.Context, accountID uint32, limit int) ([]models.Match, error)
	GetPlayerMatchesPage(ctx context.Context, accountID uint32, minMatchID *uint64, limit int, since *time.Time) ([]models.Match, error)
	GetPlayerMetrics(ctx context.Context, accountID uint32) (*models.PlayerMetrics, error)
	GetActiveMatches(ctx context.Context, accountIDs []uint32) ([]models.Match, error)

	GetHeroes(ctx context.Context) ([]models.Hero, error)
	GetHeroStats(ctx context.Context) ([]models.HeroStats, error)
	GetHeroBanStats(ctx context.Context) ([]models.HeroBanStats, error)
	GetHeroBuildStats(ctx context.Context, heroID uint32) ([]models.HeroBuildStats, error)
	GetHeroCounterStats(ctx context.Context) ([]models.HeroCounterStats, error)
	GetHeroSynergyStats(ctx context.Context) ([]models.HeroSynergyStats, error)
	GetAbilityOrderStats(ctx context.Context, heroID uint32) ([]models.AbilityOrderStats, error)
	GetHeroScoreboard(ctx context.Context, sortBy string, limit int) ([]models.HeroScoreboard, error)

	GetItems(ctx context.Context) ([]models.Item, error)
	GetItemStats(ctx context.Context) ([]models.ItemStats, error)

	GetLeaderboard(ctx context.Context, region string) (*models.Leaderboard, error)
	GetHeroLeaderboard(ctx context.Context, region string, heroID uint32) (*models.Leaderboard, error)
	GetPlayerScoreboard(ctx context.Context, sortBy string, limit int) ([]models.PlayerScoreboard, error)

	GetGameStats(ctx context.Context) ([]models.GameStats, error)
	GetKillDeathStats(ctx context.Context) ([]models.KillDeathStats, error)
	GetBadgeDistribution(ctx context.Context) ([]models.BadgeDistribution, error)

	GetBuilds(ctx context.Context, heroID uint32, limit int) ([]models.Build, error)
	GetRanks(ctx context.Context) ([]models.Rank, error)
}

type deadlockClient struct {
	baseURL    string
	assetsURL  string
	httpClient *http.Client
}

func NewDeadlockClient(baseURL string, assetsURL string) DeadlockClient {
	return &deadlockClient{
		baseURL:   baseURL,
		assetsURL: assetsURL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *deadlockClient) fetch(ctx context.Context, url string, target interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

func (c *deadlockClient) GetPlayerMatchesPage(ctx context.Context, accountID uint32, minMatchID *uint64, limit int, since *time.Time) ([]models.Match, error) {
	url := fmt.Sprintf("%s/v1/matches/metadata?account_ids=%d&limit=%d&include_player_info=true", c.baseURL, accountID, limit)
	if minMatchID != nil {
		url += fmt.Sprintf("&min_match_id=%d", *minMatchID)
	}
	if since != nil {
		url += fmt.Sprintf("&min_unix_timestamp=%d", since.Unix())
	}
	var matches []models.Match
	if err := c.fetch(ctx, url, &matches); err != nil {
		return nil, fmt.Errorf("GetPlayerMatchesPage %d: %w", accountID, err)
	}
	return matches, nil
}

func (c *deadlockClient) GetPlayerMatches(ctx context.Context, accountID uint32, limit int) ([]models.Match, error) {
	minTimestamp := time.Now().AddDate(0, 0, -90).Unix()
	url := fmt.Sprintf("%s/v1/matches/metadata?account_ids=%d&limit=%d&include_player_info=true&min_unix_timestamp=%d", c.baseURL, accountID, limit, minTimestamp)
	var matches []models.Match
	if err := c.fetch(ctx, url, &matches); err != nil {
		return nil, fmt.Errorf("GetPlayerMatches %d: %w", accountID, err)
	}
	return matches, nil
}

func (c *deadlockClient) GetPlayerMetrics(ctx context.Context, accountID uint32) (*models.PlayerMetrics, error) {
	url := fmt.Sprintf("%s/v1/analytics/player-stats/metrics?account_ids=%d", c.baseURL, accountID)
	var metrics models.PlayerMetrics
	if err := c.fetch(ctx, url, &metrics); err != nil {
		return nil, fmt.Errorf("GetPlayerMetrics %d: %w", accountID, err)
	}
	return &metrics, nil
}

func (c *deadlockClient) GetActiveMatches(ctx context.Context, accountIDs []uint32) ([]models.Match, error) {
	ids := make([]string, len(accountIDs))
	for i, id := range accountIDs {
		ids[i] = fmt.Sprintf("%d", id)
	}
	url := fmt.Sprintf("%s/v1/matches/active?account_ids=%s", c.baseURL, strings.Join(ids, ","))
	var matches []models.Match
	if err := c.fetch(ctx, url, &matches); err != nil {
		return nil, fmt.Errorf("GetActiveMatches: %w", err)
	}
	return matches, nil
}

func (c *deadlockClient) GetHeroes(ctx context.Context) ([]models.Hero, error) {
	url := fmt.Sprintf("%s/v2/heroes", c.assetsURL)
	var heroes []models.Hero
	if err := c.fetch(ctx, url, &heroes); err != nil {
		return nil, fmt.Errorf("GetHeroes: %w", err)
	}
	return heroes, nil
}

func (c *deadlockClient) GetHeroStats(ctx context.Context) ([]models.HeroStats, error) {
	url := fmt.Sprintf("%s/v1/analytics/hero-stats", c.baseURL)
	var stats []models.HeroStats
	if err := c.fetch(ctx, url, &stats); err != nil {
		return nil, fmt.Errorf("GetHeroStats: %w", err)
	}
	return stats, nil
}

func (c *deadlockClient) GetHeroBanStats(ctx context.Context) ([]models.HeroBanStats, error) {
	url := fmt.Sprintf("%s/v1/analytics/hero-ban-stats", c.baseURL)
	var stats []models.HeroBanStats
	if err := c.fetch(ctx, url, &stats); err != nil {
		return nil, fmt.Errorf("GetHeroBanStats: %w", err)
	}
	return stats, nil
}

func (c *deadlockClient) GetHeroBuildStats(ctx context.Context, heroID uint32) ([]models.HeroBuildStats, error) {
	url := fmt.Sprintf("%s/v1/analytics/hero-build-stats/%d", c.baseURL, heroID)
	var stats []models.HeroBuildStats
	if err := c.fetch(ctx, url, &stats); err != nil {
		return nil, fmt.Errorf("GetHeroBuildStats %d: %w", heroID, err)
	}
	return stats, nil
}

func (c *deadlockClient) GetHeroCounterStats(ctx context.Context) ([]models.HeroCounterStats, error) {
	url := fmt.Sprintf("%s/v1/analytics/hero-counter-stats", c.baseURL)
	var stats []models.HeroCounterStats
	if err := c.fetch(ctx, url, &stats); err != nil {
		return nil, fmt.Errorf("GetHeroCounterStats: %w", err)
	}
	return stats, nil
}

func (c *deadlockClient) GetHeroSynergyStats(ctx context.Context) ([]models.HeroSynergyStats, error) {
	url := fmt.Sprintf("%s/v1/analytics/hero-synergy-stats", c.baseURL)
	var stats []models.HeroSynergyStats
	if err := c.fetch(ctx, url, &stats); err != nil {
		return nil, fmt.Errorf("GetHeroSynergyStats: %w", err)
	}
	return stats, nil
}

func (c *deadlockClient) GetAbilityOrderStats(ctx context.Context, heroID uint32) ([]models.AbilityOrderStats, error) {
	url := fmt.Sprintf("%s/v1/analytics/ability-order-stats?hero_id=%d", c.baseURL, heroID)
	var stats []models.AbilityOrderStats
	if err := c.fetch(ctx, url, &stats); err != nil {
		return nil, fmt.Errorf("GetAbilityOrderStats %d: %w", heroID, err)
	}
	return stats, nil
}

func (c *deadlockClient) GetHeroScoreboard(ctx context.Context, sortBy string, limit int) ([]models.HeroScoreboard, error) {
	url := fmt.Sprintf("%s/v1/analytics/scoreboards/heroes?sort_by=%s&limit=%d", c.baseURL, sortBy, limit)
	var scoreboard []models.HeroScoreboard
	if err := c.fetch(ctx, url, &scoreboard); err != nil {
		return nil, fmt.Errorf("GetHeroScoreboard: %w", err)
	}
	return scoreboard, nil
}

func (c *deadlockClient) GetItems(ctx context.Context) ([]models.Item, error) {
	url := fmt.Sprintf("%s/v2/items", c.assetsURL)
	var items []models.Item
	if err := c.fetch(ctx, url, &items); err != nil {
		return nil, fmt.Errorf("GetItems: %w", err)
	}
	return items, nil
}

func (c *deadlockClient) GetItemStats(ctx context.Context) ([]models.ItemStats, error) {
	url := fmt.Sprintf("%s/v1/analytics/item-stats", c.baseURL)
	var stats []models.ItemStats
	if err := c.fetch(ctx, url, &stats); err != nil {
		return nil, fmt.Errorf("GetItemStats: %w", err)
	}
	return stats, nil
}

func (c *deadlockClient) GetLeaderboard(ctx context.Context, region string) (*models.Leaderboard, error) {
	url := fmt.Sprintf("%s/v1/leaderboard/%s", c.baseURL, region)
	var leaderboard models.Leaderboard
	if err := c.fetch(ctx, url, &leaderboard); err != nil {
		return nil, fmt.Errorf("GetLeaderboard %s: %w", region, err)
	}
	return &leaderboard, nil
}

func (c *deadlockClient) GetHeroLeaderboard(ctx context.Context, region string, heroID uint32) (*models.Leaderboard, error) {
	url := fmt.Sprintf("%s/v1/leaderboard/%s/%d", c.baseURL, region, heroID)
	var leaderboard models.Leaderboard
	if err := c.fetch(ctx, url, &leaderboard); err != nil {
		return nil, fmt.Errorf("GetHeroLeaderboard %s/%d: %w", region, heroID, err)
	}
	return &leaderboard, nil
}

func (c *deadlockClient) GetPlayerScoreboard(ctx context.Context, sortBy string, limit int) ([]models.PlayerScoreboard, error) {
	url := fmt.Sprintf("%s/v1/analytics/scoreboards/players?sort_by=%s&limit=%d", c.baseURL, sortBy, limit)
	var scoreboard []models.PlayerScoreboard
	if err := c.fetch(ctx, url, &scoreboard); err != nil {
		return nil, fmt.Errorf("GetPlayerScoreboard: %w", err)
	}
	return scoreboard, nil
}

func (c *deadlockClient) GetGameStats(ctx context.Context) ([]models.GameStats, error) {
	url := fmt.Sprintf("%s/v1/analytics/game-stats", c.baseURL)
	var stats []models.GameStats
	if err := c.fetch(ctx, url, &stats); err != nil {
		return nil, fmt.Errorf("GetGameStats: %w", err)
	}
	return stats, nil
}

func (c *deadlockClient) GetKillDeathStats(ctx context.Context) ([]models.KillDeathStats, error) {
	url := fmt.Sprintf("%s/v1/analytics/kill-death-stats", c.baseURL)
	var stats []models.KillDeathStats
	if err := c.fetch(ctx, url, &stats); err != nil {
		return nil, fmt.Errorf("GetKillDeathStats: %w", err)
	}
	return stats, nil
}

func (c *deadlockClient) GetBadgeDistribution(ctx context.Context) ([]models.BadgeDistribution, error) {
	url := fmt.Sprintf("%s/v1/analytics/badge-distribution", c.baseURL)
	var dist []models.BadgeDistribution
	if err := c.fetch(ctx, url, &dist); err != nil {
		return nil, fmt.Errorf("GetBadgeDistribution: %w", err)
	}
	return dist, nil
}

func (c *deadlockClient) GetBuilds(ctx context.Context, heroID uint32, limit int) ([]models.Build, error) {
	url := fmt.Sprintf("%s/v1/builds?hero_id=%d&limit=%d", c.baseURL, heroID, limit)
	var builds []models.Build
	if err := c.fetch(ctx, url, &builds); err != nil {
		return nil, fmt.Errorf("GetBuilds: %w", err)
	}
	return builds, nil
}

func (c *deadlockClient) GetRanks(ctx context.Context) ([]models.Rank, error) {
	url := fmt.Sprintf("%s/v2/ranks", c.assetsURL)
	var ranks []models.Rank
	if err := c.fetch(ctx, url, &ranks); err != nil {
		return nil, fmt.Errorf("GetRanks: %w", err)
	}
	return ranks, nil
}
