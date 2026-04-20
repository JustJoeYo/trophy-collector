package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/JustJoeYo/trophy-collector/internal/models"
)

// DeadlockClient defines the interface for all Deadlock API calls.
// Using an interface means we can swap in a mock for testing.
type DeadlockClient interface {
	GetPlayerMatches(ctx context.Context, steamID string) ([]models.Match, error)
	GetPlayerHeroStats(ctx context.Context, steamID string) ([]models.HeroStats, error)
	GetHeroes(ctx context.Context) ([]models.Hero, error)
}

type deadlockClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewDeadlockClient creates a new client for deadlock-api.com
func NewDeadlockClient(baseURL string) DeadlockClient {
	return &deadlockClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *deadlockClient) GetPlayerMatches(ctx context.Context, steamID string) ([]models.Match, error) {
	url := fmt.Sprintf("%s/v1/players/%s/match-history", c.baseURL, steamID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("deadlock API returned status %d", resp.StatusCode)
	}

	var matches []models.Match
	if err := json.NewDecoder(resp.Body).Decode(&matches); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return matches, nil
}

func (c *deadlockClient) GetPlayerHeroStats(ctx context.Context, steamID string) ([]models.HeroStats, error) {
	url := fmt.Sprintf("%s/v1/players/%s/hero-stats", c.baseURL, steamID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("deadlock API returned status %d", resp.StatusCode)
	}

	var stats []models.HeroStats
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return stats, nil
}

func (c *deadlockClient) GetHeroes(ctx context.Context) ([]models.Hero, error) {
	url := fmt.Sprintf("%s/v1/heroes", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("deadlock API returned status %d", resp.StatusCode)
	}

	var heroes []models.Hero
	if err := json.NewDecoder(resp.Body).Decode(&heroes); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return heroes, nil
}
