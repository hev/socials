package cmd

import (
	"fmt"
	"os"

	"github.com/hev/socials/internal/config"
	"github.com/spf13/cobra"
)

var (
	cfg        *config.Config
	verbose    bool
	jsonOutput bool
	configPath string
)

var rootCmd = &cobra.Command{
	Use:   "socials",
	Short: "Manage Twitter and LinkedIn from the terminal",
	Long: `Socials is a CLI tool for managing Twitter and LinkedIn.
Built for AI agents and humans alike, with structured JSON output
for programmatic consumption.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		switch cmd.Name() {
		case "init", "set", "help", "completion":
			return nil
		}

		var err error
		if configPath != "" {
			cfg, err = config.LoadFrom(configPath)
		} else {
			cfg, err = config.Load()
		}
		if err != nil {
			// Allow commands that can work without config (e.g. post --dry-run)
			if verbose {
				fmt.Fprintf(os.Stderr, "Warning: %s\n", err)
			}
		}
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Show verbose output")
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output as JSON")
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "Path to config file")

	rootCmd.AddCommand(feedCmd)
	rootCmd.AddCommand(postCmd)
	rootCmd.AddCommand(messagesCmd)
	rootCmd.AddCommand(configCmd)
}
