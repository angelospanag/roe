# ROE

**Roe** is a lightweight, self-hosted RSS feed aggregator inspired by the Greek word *ροή* (flow) — Go + Huma + SQLC
+ PostgreSQL backend, Next.js frontend.

## Features

- Full CRUD operations for feeds and posts
- Read/unread tracking & bulk operations
- Feed refresh (all or individual), with feed-URL validation on add
- Post filtering (by feed, read status) & pagination
- Unread counts (global or per-feed)
- Type-safe queries (SQLC) & auto-generated OpenAPI docs, with a TypeScript client generated from the spec via
  [hey-api](https://heyapi.dev)

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

See [`CLAUDE.md`](CLAUDE.md), [`api/CLAUDE.md`](api/CLAUDE.md), and [`ui/CLAUDE.md`](ui/CLAUDE.md) for architecture
notes and the full command reference (or run `mise tasks`).

## Getting Started

[mise](https://mise.jdx.dev/) manages the pinned toolchain. [Podman](https://podman.io/) is a separate prerequisite.

```bash
curl https://mise.run | sh          # macOS / Linux
eval "$(mise activate zsh)"         # add to ~/.zshrc

mise trust    # one-time, confirms you trust this repo's mise.toml
mise install  # downloads Go, Bun, golangci-lint, and sqlc
```

API: **http://localhost:8080** • UI: **http://localhost:3000** • Docs: **http://localhost:8080/docs**

## Development

```bash
# Terminal 1
mise -C api run migrate  # first time only, against an empty database
mise -C api run dev

# Terminal 2
mise -C ui run install   # first time only
mise -C ui run dev
```

The Next.js dev server proxies `/api/*` to `:8080` — no CORS config needed.

Or run the whole stack in containers: `mise run compose:up` (see [`CLAUDE.md`](CLAUDE.md)).

## API

Interactive OpenAPI docs (Swagger UI) at `http://localhost:8080/docs`.

| Resource | Endpoints |
|---|---|
| Feeds | `POST /feeds` · `GET /feeds` · `GET\|PUT\|DELETE /feeds/{id}` · `POST /feeds/refresh` · `POST /feeds/{id}/mark-all-read` · `GET /feeds/{id}/unread/count` |
| Posts | `GET /posts` · `GET /posts/{id}` · `PATCH /posts/{id}/read` · `GET /posts/unread/count` |

## Configuration

```bash
# api/.env
DATABASE_URL="postgres://postgres:postgres@localhost:5432/roe_backend?sslmode=disable"
PORT="8080"
```

```bash
# ui — server-side only (no NEXT_PUBLIC_ prefix), set when running outside Podman Compose
API_URL="http://localhost:8080"
```
