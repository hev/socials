package twitter

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hev/socials/internal/output"
)

type dmEventsResponse struct {
	Data []struct {
		ID        string `json:"id"`
		Text      string `json:"text"`
		EventType string `json:"event_type"`
		SenderID  string `json:"sender_id"`
		CreatedAt string `json:"created_at"`
	} `json:"data"`
	Includes struct {
		Users []struct {
			ID       string `json:"id"`
			Username string `json:"username"`
			Name     string `json:"name"`
		} `json:"users"`
	} `json:"includes"`
}

func (c *Client) GetDirectMessages(count int) ([]output.DirectMessage, error) {
	if count <= 0 {
		count = 10
	}

	url := fmt.Sprintf("%s/dm_events?max_results=%d&dm_event.fields=created_at,sender_id&event_types=MessageCreate&expansions=sender_id&user.fields=username,name",
		baseURL, count)

	data, err := c.doRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get DMs: %w", err)
	}

	var resp dmEventsResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse DMs: %w", err)
	}

	userMap := make(map[string]string)
	for _, u := range resp.Includes.Users {
		userMap[u.ID] = u.Username
	}

	messages := make([]output.DirectMessage, 0, len(resp.Data))
	for _, m := range resp.Data {
		createdAt := m.CreatedAt
		if createdAt == "" {
			createdAt = time.Now().Format(time.RFC3339)
		}
		messages = append(messages, output.DirectMessage{
			ID:         m.ID,
			Text:       m.Text,
			SenderID:   m.SenderID,
			SenderName: userMap[m.SenderID],
			CreatedAt:  createdAt,
		})
	}

	return messages, nil
}
