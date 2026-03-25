# Go Telegram Bot

Telegram bot built with pure Go standard library — no third-party Telegram SDK. Supports both long-polling and webhook modes.

## Stack

- **Go** — standard library `net/http`, no external Telegram SDK
- **Storage** — pluggable interface (in-memory for dev, Redis-ready)
- **Docker** — multi-stage build
- **Two modes** — long-polling (dev) and webhook (production)

## Features

- `/start` — Welcome message
- `/help` — List commands
- `/note <text>` — Save a note
- `/list` — Show saved notes
- `/stats` — Usage statistics

## Quick Start

```bash
# Set your bot token
export TELEGRAM_BOT_TOKEN=your-token-here

# Run directly
go run ./cmd/bot

# Or with Docker
docker compose up -d
```

## Architecture

```
cmd/bot/              # Entry point, mode selection
internal/
├── bot/              # Core bot: polling, webhook, Telegram API client
│   ├── bot.go        # Bot struct, polling loop, SendMessage
│   └── types.go      # Telegram API types (Update, Message, User, Chat)
├── handlers/         # Command router and handlers
│   └── router.go     # /start, /help, /note, /list, /stats
└── storage/          # Pluggable storage interface
    ├── store.go      # Interface definition
    ├── memory.go     # In-memory implementation
    └── redis.go      # Redis implementation (TODO)
```

## Webhook Mode

```bash
BOT_MODE=webhook PORT=8080 go run ./cmd/bot
```

Set webhook URL via Telegram API:
```bash
curl "https://api.telegram.org/bot$TOKEN/setWebhook?url=https://your-domain.com/webhook"
```

## TODO

- [ ] Redis storage implementation
- [ ] Inline keyboard support
- [ ] Note deletion and editing
- [ ] Reminder scheduling
- [ ] Rate limiting per user
