# socials

CLI for Twitter and LinkedIn. Built for humans and AI agents.

## Install

```bash
go install github.com/hev/socials@latest
```

Or build from source:

```bash
git clone https://github.com/hev/socials.git
cd socials
make build
```

Requires Go 1.24+.

## Setup

```bash
socials config init
```

This creates `~/.config/socials/config.yaml` and walks you through adding API credentials for Twitter and/or LinkedIn.

## Usage

```bash
# Read your feeds
socials feed twitter --count 20
socials feed linkedin

# Post (supports markdown)
socials post --file post.md --network twitter,linkedin
socials post --file post.md --dry-run  # preview without posting

# Direct messages
socials messages twitter --count 10
socials messages linkedin

# Config
socials config show
socials config set twitter.api_key <key>
```

All commands support `--json` for structured output.

## License

MIT
