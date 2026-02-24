package twitter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/dghubble/oauth1"
	"github.com/hev/socials/internal/config"
)

const baseURL = "https://api.twitter.com/2"

type Client struct {
	httpClient *http.Client
	userID     string
}

func NewClient(cfg *config.TwitterConfig) *Client {
	oauthConfig := oauth1.NewConfig(cfg.APIKey, cfg.APIKeySecret)
	token := oauth1.NewToken(cfg.AccessToken, cfg.AccessTokenSecret)
	return &Client{
		httpClient: oauthConfig.Client(oauth1.NoContext, token),
		userID:     cfg.UserID,
	}
}

func (c *Client) doRequest(method, url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode == 401 {
		return nil, fmt.Errorf("authentication failed (401): check your Twitter API tokens")
	}
	if resp.StatusCode == 403 {
		return nil, fmt.Errorf("forbidden (403): check your Twitter API access level")
	}
	if resp.StatusCode == 429 {
		return nil, fmt.Errorf("rate limited (429): too many requests, try again later")
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiErr struct {
			Detail string `json:"detail"`
			Title  string `json:"title"`
		}
		if json.Unmarshal(data, &apiErr) == nil && apiErr.Detail != "" {
			return nil, fmt.Errorf("twitter API error (%d): %s", resp.StatusCode, apiErr.Detail)
		}
		return nil, fmt.Errorf("twitter API error (%d): %s", resp.StatusCode, string(data))
	}

	return data, nil
}
