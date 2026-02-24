package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

func Print(data any, jsonMode bool) error {
	if jsonMode {
		return PrintJSON(data)
	}
	return PrintHuman(data)
}

func PrintJSON(data any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

func PrintHuman(data any) error {
	switch v := data.(type) {
	case []Tweet:
		for _, t := range v {
			fmt.Printf("@%s · %s\n", t.AuthorUsername, formatTime(t.CreatedAt))
			fmt.Println(t.Text)
			fmt.Printf("♥ %d  🔁 %d  💬 %d\n\n", t.Likes, t.Retweets, t.Replies)
		}
	case []LinkedInPost:
		for _, p := range v {
			fmt.Printf("%s · %s\n", p.AuthorName, formatTime(p.CreatedAt))
			fmt.Println(p.Text)
			fmt.Printf("👍 %d  💬 %d\n\n", p.Likes, p.Comments)
		}
	case []DirectMessage:
		for _, m := range v {
			fmt.Printf("[%s] %s: %s\n", formatTime(m.CreatedAt), m.SenderName, m.Text)
		}
	case []LinkedInMessage:
		for _, m := range v {
			fmt.Printf("[%s] %s: %s\n", formatTime(m.CreatedAt), m.SenderName, m.Text)
		}
	case PostResult:
		fmt.Printf("Posted to %s\n", v.Network)
		if v.ID != "" {
			fmt.Printf("ID: %s\n", v.ID)
		}
		if v.URL != "" {
			fmt.Printf("URL: %s\n", v.URL)
		}
	case []PostResult:
		for _, r := range v {
			fmt.Printf("Posted to %s\n", r.Network)
			if r.ID != "" {
				fmt.Printf("  ID: %s\n", r.ID)
			}
			if r.URL != "" {
				fmt.Printf("  URL: %s\n", r.URL)
			}
		}
	case DryRunResult:
		fmt.Printf("=== Dry Run: %s ===\n", v.Network)
		for i, chunk := range v.Chunks {
			if len(v.Chunks) > 1 {
				fmt.Printf("--- Part %d/%d (%d chars) ---\n", i+1, len(v.Chunks), len(chunk))
			}
			fmt.Println(chunk)
		}
		fmt.Println()
	case []DryRunResult:
		for _, r := range v {
			PrintHuman(r)
		}
	case ConfigDisplay:
		fmt.Println("Twitter:")
		fmt.Printf("  API Key:       %s\n", v.Twitter.APIKey)
		fmt.Printf("  API Secret:    %s\n", v.Twitter.APIKeySecret)
		fmt.Printf("  Access Token:  %s\n", v.Twitter.AccessToken)
		fmt.Printf("  Access Secret: %s\n", v.Twitter.AccessTokenSecret)
		fmt.Printf("  User ID:       %s\n", v.Twitter.UserID)
		fmt.Println("LinkedIn:")
		fmt.Printf("  Access Token:  %s\n", v.LinkedIn.AccessToken)
		fmt.Printf("  Person URN:    %s\n", v.LinkedIn.PersonURN)
	default:
		return PrintJSON(data)
	}
	return nil
}

func formatTime(t string) string {
	parsed, err := time.Parse(time.RFC3339, t)
	if err != nil {
		return t
	}
	return parsed.Local().Format("Jan 2 15:04")
}

func Redact(s string) string {
	if len(s) == 0 {
		return "(not set)"
	}
	if len(s) <= 8 {
		return strings.Repeat("*", len(s))
	}
	return s[:4] + strings.Repeat("*", len(s)-8) + s[len(s)-4:]
}

// Output types

type Tweet struct {
	ID             string `json:"id"`
	Text           string `json:"text"`
	AuthorID       string `json:"author_id"`
	AuthorUsername string `json:"author_username"`
	CreatedAt      string `json:"created_at"`
	Likes          int    `json:"likes"`
	Retweets       int    `json:"retweets"`
	Replies        int    `json:"replies"`
}

type LinkedInPost struct {
	ID        string `json:"id"`
	Text      string `json:"text"`
	AuthorURN string `json:"author_urn"`
	AuthorName string `json:"author_name"`
	CreatedAt string `json:"created_at"`
	Likes     int    `json:"likes"`
	Comments  int    `json:"comments"`
}

type DirectMessage struct {
	ID         string `json:"id"`
	Text       string `json:"text"`
	SenderID   string `json:"sender_id"`
	SenderName string `json:"sender_name"`
	CreatedAt  string `json:"created_at"`
}

type LinkedInMessage struct {
	ID         string `json:"id"`
	Text       string `json:"text"`
	SenderURN  string `json:"sender_urn"`
	SenderName string `json:"sender_name"`
	CreatedAt  string `json:"created_at"`
}

type PostResult struct {
	Network string `json:"network"`
	ID      string `json:"id"`
	URL     string `json:"url,omitempty"`
	Text    string `json:"text"`
}

type DryRunResult struct {
	Network string   `json:"network"`
	Chunks  []string `json:"chunks"`
}

type ConfigDisplay struct {
	Twitter  ConfigTwitterDisplay  `json:"twitter"`
	LinkedIn ConfigLinkedInDisplay `json:"linkedin"`
}

type ConfigTwitterDisplay struct {
	APIKey            string `json:"api_key"`
	APIKeySecret      string `json:"api_key_secret"`
	AccessToken       string `json:"access_token"`
	AccessTokenSecret string `json:"access_token_secret"`
	UserID            string `json:"user_id"`
}

type ConfigLinkedInDisplay struct {
	AccessToken string `json:"access_token"`
	PersonURN   string `json:"person_urn"`
}
