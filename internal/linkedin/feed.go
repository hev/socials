package linkedin

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/hev/socials/internal/output"
)

type postsResponse struct {
	Elements []struct {
		ID           string `json:"id"`
		Commentary   string `json:"commentary"`
		Author       string `json:"author"`
		CreatedAt    int64  `json:"createdAt"`
		LifecycleState string `json:"lifecycleState"`
		Distribution struct {
			FeedDistribution string `json:"feedDistribution"`
		} `json:"distribution"`
		SocialDetail *struct {
			TotalSocialActivityCounts struct {
				NumLikes    int `json:"numLikes"`
				NumComments int `json:"numComments"`
			} `json:"totalSocialActivityCounts"`
		} `json:"socialDetail,omitempty"`
	} `json:"elements"`
}

func (c *Client) GetPosts(count int) ([]output.LinkedInPost, error) {
	if count <= 0 {
		count = 10
	}

	params := url.Values{}
	params.Set("author", c.personURN)
	params.Set("q", "author")
	params.Set("count", fmt.Sprintf("%d", count))
	params.Set("sortBy", "LAST_MODIFIED")

	reqURL := fmt.Sprintf("%s/rest/posts?%s", baseURL, params.Encode())

	data, err := c.doRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get posts: %w", err)
	}

	var resp postsResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse posts: %w", err)
	}

	posts := make([]output.LinkedInPost, 0, len(resp.Elements))
	for _, p := range resp.Elements {
		var likes, comments int
		if p.SocialDetail != nil {
			likes = p.SocialDetail.TotalSocialActivityCounts.NumLikes
			comments = p.SocialDetail.TotalSocialActivityCounts.NumComments
		}

		createdAt := time.UnixMilli(p.CreatedAt).Format(time.RFC3339)

		posts = append(posts, output.LinkedInPost{
			ID:         p.ID,
			Text:       p.Commentary,
			AuthorURN:  p.Author,
			AuthorName: "You",
			CreatedAt:  createdAt,
			Likes:      likes,
			Comments:   comments,
		})
	}

	return posts, nil
}
