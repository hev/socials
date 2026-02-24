package twitter

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hev/socials/internal/output"
)

type createTweetRequest struct {
	Text  string       `json:"text"`
	Reply *replyConfig `json:"reply,omitempty"`
}

type replyConfig struct {
	InReplyToTweetID string `json:"in_reply_to_tweet_id"`
}

type createTweetResponse struct {
	Data struct {
		ID   string `json:"id"`
		Text string `json:"text"`
	} `json:"data"`
}

func (c *Client) PostTweet(text string) (*output.PostResult, error) {
	body, err := json.Marshal(createTweetRequest{Text: text})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal tweet: %w", err)
	}

	data, err := c.doRequest("POST", baseURL+"/tweets", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to post tweet: %w", err)
	}

	var resp createTweetResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &output.PostResult{
		Network: "twitter",
		ID:      resp.Data.ID,
		URL:     fmt.Sprintf("https://twitter.com/i/status/%s", resp.Data.ID),
		Text:    resp.Data.Text,
	}, nil
}

func (c *Client) PostThread(chunks []string) ([]output.PostResult, error) {
	var results []output.PostResult
	var lastID string

	for i, chunk := range chunks {
		var body []byte
		var err error

		if i == 0 {
			body, err = json.Marshal(createTweetRequest{Text: chunk})
		} else {
			body, err = json.Marshal(createTweetRequest{
				Text:  chunk,
				Reply: &replyConfig{InReplyToTweetID: lastID},
			})
		}
		if err != nil {
			return results, fmt.Errorf("failed to marshal tweet %d: %w", i+1, err)
		}

		data, err := c.doRequest("POST", baseURL+"/tweets", bytes.NewReader(body))
		if err != nil {
			return results, fmt.Errorf("failed to post tweet %d: %w", i+1, err)
		}

		var resp createTweetResponse
		if err := json.Unmarshal(data, &resp); err != nil {
			return results, fmt.Errorf("failed to parse response for tweet %d: %w", i+1, err)
		}

		lastID = resp.Data.ID
		results = append(results, output.PostResult{
			Network: "twitter",
			ID:      resp.Data.ID,
			URL:     fmt.Sprintf("https://twitter.com/i/status/%s", resp.Data.ID),
			Text:    resp.Data.Text,
		})
	}

	return results, nil
}
