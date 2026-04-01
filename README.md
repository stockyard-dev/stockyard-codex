# Stockyard Codex

**Internal wiki — Markdown, full-text search, version history, nested pages**

Part of the [Stockyard](https://stockyard.dev) family of self-hosted developer tools.

## Quick Start

```bash
docker run -p 9120:9120 -v codex_data:/data ghcr.io/stockyard-dev/stockyard-codex
```

Or with docker-compose:

```bash
docker-compose up -d
```

Open `http://localhost:9120` in your browser.

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `9120` | HTTP port |
| `DATA_DIR` | `./data` | SQLite database directory |
| `CODEX_LICENSE_KEY` | *(empty)* | Pro license key |

## Free vs Pro

| | Free | Pro |
|-|------|-----|
| Limits | 20 pages | Unlimited pages and history |
| Price | Free | $4.99/mo |

Get a Pro license at [stockyard.dev/tools/](https://stockyard.dev/tools/).

## Category

Developer Tools

## License

Apache 2.0
