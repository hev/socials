package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Twitter  TwitterConfig  `mapstructure:"twitter"`
	LinkedIn LinkedInConfig `mapstructure:"linkedin"`
}

type TwitterConfig struct {
	APIKey            string `mapstructure:"api_key"`
	APIKeySecret      string `mapstructure:"api_key_secret"`
	AccessToken       string `mapstructure:"access_token"`
	AccessTokenSecret string `mapstructure:"access_token_secret"`
	UserID            string `mapstructure:"user_id"`
}

type LinkedInConfig struct {
	AccessToken string `mapstructure:"access_token"`
	PersonURN   string `mapstructure:"person_urn"`
}

func ConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, ".config", "socials"), nil
}

func ConfigPath() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.yaml"), nil
}

func Load() (*Config, error) {
	dir, err := ConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get config dir: %w", err)
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(dir)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("config not found, run 'socials config init' first")
		}
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}

func LoadFrom(path string) (*Config, error) {
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config from %s: %w", path, err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}

func (c *Config) HasTwitter() bool {
	return c.Twitter.APIKey != "" &&
		c.Twitter.APIKeySecret != "" &&
		c.Twitter.AccessToken != "" &&
		c.Twitter.AccessTokenSecret != ""
}

func (c *Config) HasLinkedIn() bool {
	return c.LinkedIn.AccessToken != "" && c.LinkedIn.PersonURN != ""
}

func Save(cfg *Config) error {
	dir, err := ConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get config dir: %w", err)
	}

	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create config dir: %w", err)
	}

	viper.Set("twitter.api_key", cfg.Twitter.APIKey)
	viper.Set("twitter.api_key_secret", cfg.Twitter.APIKeySecret)
	viper.Set("twitter.access_token", cfg.Twitter.AccessToken)
	viper.Set("twitter.access_token_secret", cfg.Twitter.AccessTokenSecret)
	viper.Set("twitter.user_id", cfg.Twitter.UserID)
	viper.Set("linkedin.access_token", cfg.LinkedIn.AccessToken)
	viper.Set("linkedin.person_urn", cfg.LinkedIn.PersonURN)

	configPath := filepath.Join(dir, "config.yaml")
	if err := viper.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	if err := os.Chmod(configPath, 0600); err != nil {
		return fmt.Errorf("failed to set config permissions: %w", err)
	}

	return nil
}
