package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/JustJoeYo/trophy-collector/internal/models"
)

const steamAPIBase = "https://api.steampowered.com"

// SteamClient defines the interface for Steam Web API calls
type SteamClient interface {
	GetPlayerSummary(ctx context.Context, steamID string) (*models.Player, error)
	ResolveVanityURL(ctx context.Context, vanityURL string) (string, error)
}

type steamClient struct {
	apiKey     string
	httpClient *http.Client
}

func NewSteamClient(apiKey string) SteamClient {
	return &steamClient{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetPlayerSummary fetches a player's Steam profile by SteamID64
func (c *steamClient) GetPlayerSummary(ctx context.Context, steamID string) (*models.Player, error) {
	url := fmt.Sprintf(
		"%s/ISteamUser/GetPlayerSummaries/v2/?key=%s&steamids=%s",
		steamAPIBase, c.apiKey, steamID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	// Steam API response shape
	var result struct {
		Response struct {
			Players []struct {
				SteamID     string `json:"steamid"`
				PersonaName string `json:"personaname"`
				Avatar      string `json:"avatarfull"`
				ProfileURL  string `json:"profileurl"`
			} `json:"players"`
		} `json:"response"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	if len(result.Response.Players) == 0 {
		return nil, fmt.Errorf("player not found: %s", steamID)
	}

	p := result.Response.Players[0]
	return &models.Player{
		SteamID:     p.SteamID,
		PersonaName: p.PersonaName,
		AvatarURL:   p.Avatar,
		ProfileURL:  p.ProfileURL,
	}, nil
}

// ResolveVanityURL converts a Steam username to a SteamID64
func (c *steamClient) ResolveVanityURL(ctx context.Context, vanityURL string) (string, error) {
	url := fmt.Sprintf(
		"%s/ISteamUser/ResolveVanityURL/v1/?key=%s&vanityurl=%s",
		steamAPIBase, c.apiKey, vanityURL,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Response struct {
			SteamID string `json:"steamid"`
			Success int    `json:"success"`
			Message string `json:"message"`
		} `json:"response"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decoding response: %w", err)
	}

	if result.Response.Success != 1 {
		return "", fmt.Errorf("could not resolve vanity URL %q: %s", vanityURL, result.Response.Message)
	}

	return result.Response.SteamID, nil
}
