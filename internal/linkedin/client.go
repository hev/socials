package linkedin

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hev/socials/internal/config"
)

const baseURL = "https://api.linkedin.com"

type Client struct {
	httpClient  *http.Client
	accessToken string
	personURN   string
}

func NewClient(cfg *config.LinkedInConfig) *Client {
	return &Client{
		httpClient:  &http.Client{},
		accessToken: cfg.AccessToken,
		personURN:   cfg.PersonURN,
	}
}

func (c *Client) doRequest(method, url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("LinkedIn-Version", "202602")
	req.Header.Set("X-Restli-Protocol-Version", "2.0.0")
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
		return nil, fmt.Errorf("authentication failed (401): check your LinkedIn access token")
	}
	if resp.StatusCode == 403 {
		return nil, fmt.Errorf("forbidden (403): check your LinkedIn API permissions")
	}
	if resp.StatusCode == 429 {
		return nil, fmt.Errorf("rate limited (429): too many requests, try again later")
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiErr struct {
			Message string `json:"message"`
			Status  int    `json:"status"`
		}
		if json.Unmarshal(data, &apiErr) == nil && apiErr.Message != "" {
			return nil, fmt.Errorf("linkedin API error (%d): %s", resp.StatusCode, apiErr.Message)
		}
		return nil, fmt.Errorf("linkedin API error (%d): %s", resp.StatusCode, string(data))
	}

	return data, nil
}
