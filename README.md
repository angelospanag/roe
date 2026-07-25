# ROE

**Roe** is a lightweight RSS feed aggregator inspired by the Greek word *ροή* (flow) — built with Go 1.26, Huma, SQLC,
and PostgreSQL on the backend, and Next.js on the frontend, to keep your feeds moving.

## Features

- Full CRUD operations for feeds and posts
- Read/unread tracking & bulk operations
- Feed refresh (all or individual)
- Post filtering (by feed, read status) & pagination
- Unread counts (global or per-feed)
- Type-safe queries (SQLC) & auto-generated OpenAPI docs
- TypeScript client generated from the OpenAPI spec via [hey-api](https://heyapi.dev)

## Repo layout

```
roe/
├── api/              Go backend (Huma + SQLC + PostgreSQL)
│   └── openapi.yaml  Generated OpenAPI spec (`mise -C api run schema`)
├── ui/               Next.js 16 frontend
│   └── client/       Generated TypeScript client (`mise run generate`)
├── compose.yml       Full-stack Podman Compose (api + ui + postgres)
└── .github/          CI (backend: lint · vuln · test — frontend: lint · vuln · typecheck — build)
```

## Stack

| Directory     | Purpose                | Technologies                                                            |
| ------------- | ----------------------- | ------------------------------------------------------------------------ |
| `api/`        | Feed/post CRUD, refresh | Go 1.26, Huma v2, chi, SQLC, pgx/v5                                      |
| `ui/`         | Feed reader             | Next.js 16, React 19, TypeScript 5, Tailwind CSS v4, TanStack Query v5, hey-api |
| `compose.yml` | Full-stack orchestration| Podman Compose, PostgreSQL 18                                            |

## Getting Started

[mise](https://mise.jdx.dev/) manages the pinned toolchain.
[Podman](https://podman.io/) is a separate prerequisite (install Podman Desktop or equivalent; version 4+ for the built-in `podman compose`).

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
mise install  # downloads Go, Bun, golangci-lint, and sqlc
```

API: **http://localhost:8080** • UI: **http://localhost:3000** • Docs: **http://localhost:8080/docs**

## Development

Run the API and UI in separate terminals:

```bash
# Terminal 1 — start PostgreSQL and run migrations (see api/README or scripts/setup-db.sh)
mise -C api run dev

# Terminal 2
mise -C ui run install   # first time only
mise -C ui run dev
```

The Next.js dev server proxies `/api/*` to `:8080` — no CORS config needed.

Root `mise.toml` exposes cross-project aggregates and Podman tasks. Service-specific tasks run from inside `api/` or `ui/` (or via `mise -C <dir>`).

Root-level commands run against **both** the backend (`api/`) and frontend (`ui/`):

| Command                 | Description                                      |
| ------------------------ | ------------------------------------------------ |
| `mise run fmt`          | Format all code (Go + TypeScript)                |
| `mise run lint`         | Lint all code (Go + TypeScript)                  |
| `mise run vuln`         | Scan all dependencies for known vulnerabilities (Go + TypeScript) |
| `mise run test`         | Run Go tests                                     |
| `mise run deps`         | Update all dependencies (Go + Bun)               |
| `mise run generate`     | Regenerate OpenAPI schema then TypeScript client |
| `mise run compose:up`   | Build + start full stack with Podman Compose     |
| `mise run compose:down` | Stop all services                                |
| `mise run compose:logs` | Follow compose logs                              |

| API / backend command (`cd api/`) | Description                           |
| ----------------------------------- | -------------------------------------- |
| `mise run dev`           | Start the Go API on :8080             |
| `mise run build`         | Build for current platform            |
| `mise run test`          | Run tests                             |
| `mise run fmt`           | Format code via `golangci-lint fmt`   |
| `mise run lint`          | Run linters via `golangci-lint run`   |
| `mise run vuln`          | Scan Go dependencies for known vulnerabilities (govulncheck) |
| `mise run deps`          | Update and tidy dependencies          |
| `mise run schema`        | Generate OpenAPI spec to openapi.yaml |
| `mise run generate`      | Generate sqlc bindings from SQL schema|
| `mise run clean`         | Remove build artifacts                |

| UI / frontend command (`cd ui/`) | Description                           |
| ----------------------------------- | -------------------------------------- |
| `mise run dev`         | Start the Next.js dev server on :3000 |
| `mise run build`       | Production Next.js build              |
| `mise run typecheck`   | TypeScript type check                 |
| `mise run fmt`         | Format code with Biome                |
| `mise run lint`        | Biome check                           |
| `mise run vuln`        | Scan TypeScript dependencies for known vulnerabilities (`bun audit`) |
| `mise run generate`    | Generate TypeScript client from ../api/openapi.yaml |

## API Endpoints

OpenAPI docs available at `http://localhost:8080/docs`.

### Feeds

| Method | Endpoint                    | Description      |
|--------|-----------------------------|-------------------|
| POST   | `/feeds`                    | Create feed       |
| GET    | `/feeds`                    | List feeds        |
| GET    | `/feeds/{id}`               | Get feed          |
| PUT    | `/feeds/{id}`               | Update feed       |
| DELETE | `/feeds/{id}`               | Delete feed       |
| POST   | `/feeds/refresh`            | Refresh feeds     |
| POST   | `/feeds/{id}/mark-all-read` | Mark all as read  |
| GET    | `/feeds/{id}/unread/count`  | Count unread      |

### Posts

| Method | Endpoint              | Description       |
|--------|------------------------|-------------------|
| GET    | `/posts`              | List posts         |
| GET    | `/posts/{id}`         | Get post           |
| PATCH  | `/posts/{id}/read`    | Mark read/unread   |
| GET    | `/posts/unread/count` | Count unread       |

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

```bash
# api/.env
DATABASE_URL="postgres://postgres:postgres@localhost:5432/roe_backend?sslmode=disable"
PORT="8080"
```

```bash
# ui — set when running outside Podman Compose (server-side only, no NEXT_PUBLIC_ prefix)
API_URL="http://localhost:8080"
```

## Database

**Feeds**: id, title, url (unique), description, link, timestamps

**Posts**: id, feed_id (FK), title, description, content, link, author, published_at, guid, is_read, timestamps

- Unique: (feed_id, guid)
- Indexes: feed_id, is_read, published_at
