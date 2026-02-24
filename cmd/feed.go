package cmd

import (
	"fmt"

	"github.com/hev/socials/internal/linkedin"
	"github.com/hev/socials/internal/output"
	"github.com/hev/socials/internal/twitter"
	"github.com/spf13/cobra"
)

var feedCount int

var feedCmd = &cobra.Command{
	Use:   "feed [twitter|linkedin]",
	Short: "View your feed",
	Long:  "View your home timeline (Twitter) or your posts (LinkedIn).",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		network := args[0]

		if cfg == nil {
			return fmt.Errorf("config not found, run 'socials config init' first")
		}

		switch network {
		case "twitter":
			if !cfg.HasTwitter() {
				return fmt.Errorf("twitter not configured, run 'socials config init'")
			}
			client := twitter.NewClient(&cfg.Twitter)
			tweets, err := client.GetTimeline(feedCount)
			if err != nil {
				return err
			}
			return output.Print(tweets, jsonOutput)

		case "linkedin":
			if !cfg.HasLinkedIn() {
				return fmt.Errorf("linkedin not configured, run 'socials config init'")
			}
			client := linkedin.NewClient(&cfg.LinkedIn)
			posts, err := client.GetPosts(feedCount)
			if err != nil {
				return err
			}
			return output.Print(posts, jsonOutput)

		default:
			return fmt.Errorf("unknown network: %s (use 'twitter' or 'linkedin')", network)
		}
	},
}

func init() {
	feedCmd.Flags().IntVarP(&feedCount, "count", "n", 10, "Number of items to show")
}
