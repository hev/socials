package twitter

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hev/socials/internal/output"
)

type timelineResponse struct {
	Data []struct {
		ID        string `json:"id"`
		Text      string `json:"text"`
		AuthorID  string `json:"author_id"`
		CreatedAt string `json:"created_at"`
		PublicMetrics struct {
			LikeCount    int `json:"like_count"`
			RetweetCount int `json:"retweet_count"`
			ReplyCount   int `json:"reply_count"`
		} `json:"public_metrics"`
	} `json:"data"`
	Includes struct {
		Users []struct {
			ID       string `json:"id"`
			Username string `json:"username"`
		} `json:"users"`
	} `json:"includes"`
}

func (c *Client) GetTimeline(count int) ([]output.Tweet, error) {
	if count <= 0 {
		count = 10
	}

	url := fmt.Sprintf("%s/users/%s/timelines/reverse_chronological?max_results=%d&tweet.fields=created_at,public_metrics,author_id&expansions=author_id&user.fields=username",
		baseURL, c.userID, count)

	data, err := c.doRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get timeline: %w", err)
	}

	var resp timelineResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse timeline: %w", err)
	}

	userMap := make(map[string]string)
	for _, u := range resp.Includes.Users {
		userMap[u.ID] = u.Username
	}

	tweets := make([]output.Tweet, 0, len(resp.Data))
	for _, t := range resp.Data {
		createdAt := t.CreatedAt
		if createdAt == "" {
			createdAt = time.Now().Format(time.RFC3339)
		}
		tweets = append(tweets, output.Tweet{
			ID:             t.ID,
			Text:           t.Text,
			AuthorID:       t.AuthorID,
			AuthorUsername: userMap[t.AuthorID],
			CreatedAt:      createdAt,
			Likes:          t.PublicMetrics.LikeCount,
			Retweets:       t.PublicMetrics.RetweetCount,
			Replies:        t.PublicMetrics.ReplyCount,
		})
	}

	return tweets, nil
}
