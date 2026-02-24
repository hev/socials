package linkedin

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hev/socials/internal/output"
)

type conversationsResponse struct {
	Elements []struct {
		ID    string `json:"id"`
		Events []struct {
			EventContent struct {
				MessageEvent struct {
					Body string `json:"body"`
				} `json:"messageEvent"`
			} `json:"eventContent"`
			From struct {
				Member string `json:"com.linkedin.voyager.messaging.MessagingMember"`
			} `json:"from"`
			CreatedAt int64 `json:"createdAt"`
		} `json:"events"`
	} `json:"elements"`
}

func (c *Client) GetMessages(count int) ([]output.LinkedInMessage, error) {
	if count <= 0 {
		count = 10
	}

	reqURL := fmt.Sprintf("%s/rest/conversations?q=participant&count=%d", baseURL, count)

	data, err := c.doRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	var resp conversationsResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse messages: %w", err)
	}

	var messages []output.LinkedInMessage
	for _, conv := range resp.Elements {
		for _, event := range conv.Events {
			createdAt := time.UnixMilli(event.CreatedAt).Format(time.RFC3339)
			messages = append(messages, output.LinkedInMessage{
				ID:         conv.ID,
				Text:       event.EventContent.MessageEvent.Body,
				SenderURN:  event.From.Member,
				SenderName: event.From.Member,
				CreatedAt:  createdAt,
			})
		}
	}

	if len(messages) > count {
		messages = messages[:count]
	}

	return messages, nil
}
