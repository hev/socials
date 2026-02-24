package linkedin

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hev/socials/internal/output"
)

type createPostRequest struct {
	Author       string       `json:"author"`
	Commentary   string       `json:"commentary"`
	Visibility   string       `json:"visibility"`
	Distribution distribution `json:"distribution"`
	LifecycleState string    `json:"lifecycleState"`
}

type distribution struct {
	FeedDistribution string `json:"feedDistribution"`
}

type createPostResponse struct {
	ID string `json:"id"`
}

func (c *Client) CreatePost(text string) (*output.PostResult, error) {
	reqBody := createPostRequest{
		Author:       c.personURN,
		Commentary:   text,
		Visibility:   "PUBLIC",
		Distribution: distribution{FeedDistribution: "MAIN_FEED"},
		LifecycleState: "PUBLISHED",
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal post: %w", err)
	}

	data, err := c.doRequest("POST", baseURL+"/rest/posts", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	var resp createPostResponse
	// LinkedIn returns 201 with the ID in the response header or body
	if len(data) > 0 {
		json.Unmarshal(data, &resp)
	}

	return &output.PostResult{
		Network: "linkedin",
		ID:      resp.ID,
		Text:    text,
	}, nil
}
