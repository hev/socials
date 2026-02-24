package cmd

import (
	"fmt"

	"github.com/hev/socials/internal/linkedin"
	"github.com/hev/socials/internal/output"
	"github.com/hev/socials/internal/twitter"
	"github.com/spf13/cobra"
)

var messagesCount int

var messagesCmd = &cobra.Command{
	Use:   "messages [twitter|linkedin]",
	Short: "View your direct messages",
	Long:  "View your recent direct messages on Twitter or LinkedIn.",
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
			messages, err := client.GetDirectMessages(messagesCount)
			if err != nil {
				return err
			}
			return output.Print(messages, jsonOutput)

		case "linkedin":
			if !cfg.HasLinkedIn() {
				return fmt.Errorf("linkedin not configured, run 'socials config init'")
			}
			client := linkedin.NewClient(&cfg.LinkedIn)
			messages, err := client.GetMessages(messagesCount)
			if err != nil {
				return err
			}
			return output.Print(messages, jsonOutput)

		default:
			return fmt.Errorf("unknown network: %s (use 'twitter' or 'linkedin')", network)
		}
	},
}

func init() {
	messagesCmd.Flags().IntVarP(&messagesCount, "count", "n", 10, "Number of messages to show")
}
