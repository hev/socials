package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/hev/socials/internal/config"
	"github.com/hev/socials/internal/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long:  "Initialize, view, or update your socials configuration.",
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration",
	Long:  "Set up your API tokens for Twitter and LinkedIn.",
	RunE: func(cmd *cobra.Command, args []string) error {
		reader := bufio.NewReader(os.Stdin)

		fmt.Println("Socials Configuration Setup")
		fmt.Println("===========================")
		fmt.Println()
		fmt.Println("Twitter (leave blank to skip):")

		cfg := &config.Config{}

		cfg.Twitter.APIKey = prompt(reader, "  API Key: ")
		cfg.Twitter.APIKeySecret = prompt(reader, "  API Key Secret: ")
		cfg.Twitter.AccessToken = prompt(reader, "  Access Token: ")
		cfg.Twitter.AccessTokenSecret = prompt(reader, "  Access Token Secret: ")
		cfg.Twitter.UserID = prompt(reader, "  User ID: ")

		fmt.Println()
		fmt.Println("LinkedIn (leave blank to skip):")

		cfg.LinkedIn.AccessToken = prompt(reader, "  Access Token: ")
		cfg.LinkedIn.PersonURN = prompt(reader, "  Person URN (e.g. urn:li:person:abc123): ")

		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		configPath, _ := config.ConfigPath()
		fmt.Printf("\nConfig saved to %s\n", configPath)
		return nil
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long:  "Display the current configuration with secrets redacted.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfg == nil {
			return fmt.Errorf("config not found, run 'socials config init' first")
		}

		display := output.ConfigDisplay{
			Twitter: output.ConfigTwitterDisplay{
				APIKey:            output.Redact(cfg.Twitter.APIKey),
				APIKeySecret:      output.Redact(cfg.Twitter.APIKeySecret),
				AccessToken:       output.Redact(cfg.Twitter.AccessToken),
				AccessTokenSecret: output.Redact(cfg.Twitter.AccessTokenSecret),
				UserID:            cfg.Twitter.UserID,
			},
			LinkedIn: output.ConfigLinkedInDisplay{
				AccessToken: output.Redact(cfg.LinkedIn.AccessToken),
				PersonURN:   cfg.LinkedIn.PersonURN,
			},
		}

		return output.Print(display, jsonOutput)
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Long: `Set a configuration value. Keys use dot notation.

Examples:
  socials config set twitter.api_key YOUR_KEY
  socials config set linkedin.access_token YOUR_TOKEN
  socials config set twitter.user_id 12345`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		value := args[1]

		// Load existing config or create new
		existingCfg, err := config.Load()
		if err != nil {
			existingCfg = &config.Config{}
		}

		viper.Set(key, value)

		// Re-unmarshal to get updated config
		if err := viper.Unmarshal(existingCfg); err != nil {
			return fmt.Errorf("failed to update config: %w", err)
		}

		if err := config.Save(existingCfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("Set %s\n", key)
		return nil
	},
}

func init() {
	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetCmd)
}

func prompt(reader *bufio.Reader, label string) string {
	fmt.Print(label)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}
