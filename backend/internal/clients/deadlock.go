package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/JustJoeYo/trophy-collector/internal/models"
)


type DeadlockClient interface {
    GetPlayer(ctx context.Context, accountID uint32) (*models.Player, error)
    GetPlayerMatches(ctx context.Context, accountID uint32) ([]models.Match, error)
    GetPlayerHeroes(ctx context.Context, accountID uint32) ([]models.Hero, error)
    GetHeroes(ctx context.Context) ([]models.Hero, error)
}

type deadlockClient struct {
    baseURL    string
    assetsURL  string
    httpClient *http.Client
}

func NewDeadlockClient(baseURL string, assetsURL string) DeadlockClient {
    return &deadlockClient{
        baseURL:    baseURL,
        assetsURL:  assetsURL,
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

func (c *deadlockClient) GetPlayer(ctx context.Context, accountID uint32) (*models.Player, error) {
	url := fmt.Sprintf("%s/v1/players/%d", c.baseURL, accountID)
	var player models.Player
	if err := c.fetch(ctx, url, &player); err != nil {
		return nil, fmt.Errorf("GetPlayer %d: %w", accountID, err)
	}
	return &player, nil
}

func (c *deadlockClient) GetPlayerMatches(ctx context.Context, accountID uint32) ([]models.Match, error) {
    url := fmt.Sprintf("%s/v1/players/%d/matches", c.baseURL, accountID)
    var matches []models.Match
    if err := c.fetch(ctx, url, &matches); err != nil {
        return nil, fmt.Errorf("GetPlayerMatches %d: %w", accountID, err)
    }
    return matches, nil
}

func (c *deadlockClient) GetPlayerHeroes(ctx context.Context, accountID uint32) ([]models.Hero, error) {
    url := fmt.Sprintf("%s/v1/players/%d/heroes", c.baseURL, accountID)
    var heroes []models.Hero
    if err := c.fetch(ctx, url, &heroes); err != nil {
        return nil, fmt.Errorf("GetPlayerHeroes %d: %w", accountID, err)
    }
    return heroes, nil
}

func (c *deadlockClient) GetHeroes(ctx context.Context) ([]models.Hero, error) {
    url := fmt.Sprintf("%s/v2/heroes", c.assetsURL)
    var heroes []models.Hero
    if err := c.fetch(ctx, url, &heroes); err != nil {
        return nil, fmt.Errorf("GetHeroes: %w", err)
    }
    return heroes, nil
}