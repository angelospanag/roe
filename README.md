# ROE

**Roe** is a lightweight RSS feed aggregator inspired by the Greek word *ροή* (flow) — built with Go 1.26, Huma, SQLC,
and PostgreSQL to keep your feeds moving.

## Features

- Full CRUD operations for feeds and posts
- Read/unread tracking & bulk operations
- Feed refresh (all or individual)
- Post filtering (by feed, read status) & pagination
- Unread counts (global or per-feed)
- Type-safe queries (SQLC) & auto-generated OpenAPI docs

## Tech Stack

Go 1.26 • Huma v2 • SQLC • PostgreSQL • pgx/v5 • slog

## Project Structure

```
internal/
├── feed/          # Feed models, service, routes
├── post/          # Post models, routes  
└── db/            # SQLC generated code
db/
├── migrations/    # SQL migrations
└── queries/       # SQLC queries
cmd/api/main.go    # Entry point
```

## Getting Started

[mise](https://mise.jdx.dev/) manages the pinned toolchain
([Go 1.26.2](https://go.dev), [golangci-lint](https://golangci-lint.run/)), [sqlc](https://sqlc.dev/).
[Docker](https://www.docker.com/) is a separate prerequisite (install Docker Desktop or
equivalent).

```bash
# macOS / Linux
curl https://mise.run | sh

# Windows
winget install jdx.mise
```

Activate mise in your shell so the pinned versions take precedence over any system
installs (Homebrew, etc.). In `~/.zshrc`:

```zsh
eval "$(mise activate zsh)"
```

Then, in the repo:

```bash
mise trust    # one-time, confirms you trust this repo's mise.toml
mise install  # downloads and pins Go and golangci-lint
```

API: **http://localhost:8080** • Docs: **http://localhost:8080/docs**

## Development

| Command             | Description                                 |
|---------------------|---------------------------------------------|
| `mise run dev`      | Run without building a binary               |
| `mise run build`    | Build for current platform                  |
| `mise run test`     | Run tests (requires Docker)            ß    |
| `mise run fmt`      | Format code via `golangci-lint fmt`         |
| `mise run lint`     | Run linters via `golangci-lint run`         |
| `mise run vuln`     | Scan dependencies for known vulnerabilities |
| `mise run deps`     | Update and tidy dependencies                |
| `mise run generate` | Generate sqlc bindings from SQL schema      |
| `mise run clean`    | Remove build artifacts                      |

## API Endpoints

### Feeds

| Method | Endpoint                    | Description      |
|--------|-----------------------------|------------------|
| POST   | `/feeds`                    | Create feed      |
| GET    | `/feeds`                    | List feeds       |
| GET    | `/feeds/{id}`               | Get feed         |
| PUT    | `/feeds/{id}`               | Update feed      |
| DELETE | `/feeds/{id}`               | Delete feed      |
| POST   | `/feeds/refresh`            | Refresh feeds    |
| POST   | `/feeds/{id}/mark-all-read` | Mark all as read |
| GET    | `/feeds/{id}/unread/count`  | Count unread     |

### Posts

| Method | Endpoint              | Description      |
|--------|-----------------------|------------------|
| GET    | `/posts`              | List posts       |
| GET    | `/posts/{id}`         | Get post         |
| PATCH  | `/posts/{id}/read`    | Mark read/unread |
| GET    | `/posts/unread/count` | Count unread     |

## Examples

```bash
# Add feed
curl -X POST http://localhost:8080/feeds \
  -H "Content-Type: application/json" \
  -d '{"title":"Hacker News","url":"https://news.ycombinator.com/rss","description":"News","link":"https://news.ycombinator.com"}'

# Refresh feeds
curl -X POST http://localhost:8080/feeds/refresh -H "Content-Type: application/json" -d '{}'

# List posts
curl "http://localhost:8080/posts?limit=20"

# Filter by feed and status
curl "http://localhost:8080/posts?feed_id=1&unread_only=true"

# Mark post as read
curl -X PATCH http://localhost:8080/posts/1/read -H "Content-Type: application/json" -d '{"is_read": true}'
```

## Configuration

Environment variables:

```bash
DATABASE_URL="postgres://postgres:postgres@localhost:5432/roe_backend?sslmode=disable"
PORT="8080"
```

## Database

**Feeds**: id, title, url (unique), description, link, timestamps

**Posts**: id, feed_id (FK), title, description, content, link, author, published_at, guid, is_read, timestamps

- Unique: (feed_id, guid)
- Indexes: feed_id, is_read, published_at
