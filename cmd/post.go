package cmd

import (
	"fmt"
	"strings"

	"github.com/hev/socials/internal/linkedin"
	"github.com/hev/socials/internal/markdown"
	"github.com/hev/socials/internal/output"
	"github.com/hev/socials/internal/twitter"
	"github.com/spf13/cobra"
)

var (
	postFile    string
	postNetwork string
	postDryRun  bool
)

var postCmd = &cobra.Command{
	Use:   "post",
	Short: "Post content from a markdown file",
	Long: `Post content from a markdown file to Twitter and/or LinkedIn.
Markdown is converted to platform-appropriate formatting.
Long posts are automatically split into Twitter threads.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if postFile == "" {
			return fmt.Errorf("--file is required")
		}

		content, err := markdown.ParseFile(postFile)
		if err != nil {
			return fmt.Errorf("failed to read post file: %w", err)
		}

		networks := strings.Split(postNetwork, ",")
		for i := range networks {
			networks[i] = strings.TrimSpace(networks[i])
		}

		if postDryRun {
			return doDryRun(content, networks)
		}

		return doPost(content, networks)
	},
}

func doDryRun(content string, networks []string) error {
	var results []output.DryRunResult

	for _, network := range networks {
		switch network {
		case "twitter":
			chunks := markdown.ToTwitter(content)
			results = append(results, output.DryRunResult{
				Network: "twitter",
				Chunks:  chunks,
			})
		case "linkedin":
			text := markdown.ToLinkedIn(content)
			results = append(results, output.DryRunResult{
				Network: "linkedin",
				Chunks:  []string{text},
			})
		default:
			return fmt.Errorf("unknown network: %s", network)
		}
	}

	return output.Print(results, jsonOutput)
}

func doPost(content string, networks []string) error {
	if cfg == nil {
		return fmt.Errorf("config not found, run 'socials config init' first")
	}

	var results []output.PostResult

	for _, network := range networks {
		switch network {
		case "twitter":
			if !cfg.HasTwitter() {
				return fmt.Errorf("twitter not configured, run 'socials config init'")
			}
			client := twitter.NewClient(&cfg.Twitter)
			chunks := markdown.ToTwitter(content)

			if len(chunks) == 1 {
				result, err := client.PostTweet(chunks[0])
				if err != nil {
					return fmt.Errorf("failed to post to twitter: %w", err)
				}
				results = append(results, *result)
			} else {
				threadResults, err := client.PostThread(chunks)
				if err != nil {
					return fmt.Errorf("failed to post twitter thread: %w", err)
				}
				results = append(results, threadResults...)
			}

		case "linkedin":
			if !cfg.HasLinkedIn() {
				return fmt.Errorf("linkedin not configured, run 'socials config init'")
			}
			client := linkedin.NewClient(&cfg.LinkedIn)
			text := markdown.ToLinkedIn(content)

			result, err := client.CreatePost(text)
			if err != nil {
				return fmt.Errorf("failed to post to linkedin: %w", err)
			}
			results = append(results, *result)

		default:
			return fmt.Errorf("unknown network: %s", network)
		}
	}

	return output.Print(results, jsonOutput)
}

func init() {
	postCmd.Flags().StringVarP(&postFile, "file", "f", "", "Path to markdown file to post")
	postCmd.Flags().StringVarP(&postNetwork, "network", "n", "twitter", "Networks to post to (comma-separated: twitter,linkedin)")
	postCmd.Flags().BoolVar(&postDryRun, "dry-run", false, "Preview the post without publishing")
}
