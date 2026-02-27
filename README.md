# ROE

**Roe** is a lightweight RSS feed aggregator inspired by the Greek word *ροή* (flow) — built with Go 1.26, Huma, SQLC, and PostgreSQL to keep your feeds moving.

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

## Prerequisites

- Go 1.26+
- Task: `brew install go-task` (or [install from releases](https://taskfile.dev/installation/))
- Docker (for PostgreSQL)

## Quick Start

```bash
# Setup prerequisites (one time)
task setup

# Start PostgreSQL
docker-compose up -d

# Run migrations
./scripts/setup-db.sh

# Run API
task run
```

Or build and run the binary:

```bash
task build
./bin/app
```

API: **http://localhost:8080** • Docs: **http://localhost:8080/docs**

## Common Commands

```bash
task run            # Run without building
task build          # Build for current platform
task build:all      # Build for all platforms
task test           # Run tests
task lint           # Run linters
task fmt            # Format code
task clean          # Remove build artifacts
```

## API Endpoints

### Feeds
| Method | Endpoint                    | Description      |
| ------ | --------------------------- | ---------------- |
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
| ------ | --------------------- | ---------------- |
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

## Development

```bash
# Regenerate SQLC code after changing db/queries/*.sql
sqlc generate

# Build
go build -o bin/api cmd/api/main.go

# Test
go test ./...
```
