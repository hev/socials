# Socials

CLI tool for managing Twitter and LinkedIn from the terminal. Built for AI agent consumption.

## Stack

- Go 1.24+
- Cobra for CLI
- Viper for config
- Direct HTTP to Twitter v2 and LinkedIn REST APIs
- dghubble/oauth1 for Twitter OAuth 1.0a
- goldmark for markdown parsing

## Key Paths

- Config: `~/.config/socials/config.yaml`
- Binary: `socials`

## Build

```bash
go build -o socials .
```

## Commands

```
socials feed twitter|linkedin [--count N] [--json]
socials post --file post.md [--network twitter,linkedin] [--dry-run] [--json]
socials messages twitter|linkedin [--count N] [--json]
socials config init|show|set
```

## Code Style

- Descriptive names, no abbreviations
- Early returns over deep nesting
- Errors are values, handle them explicitly
- All commands support --json for agent consumption
