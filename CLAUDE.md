# Roe

Lightweight, self-hosted RSS feed aggregator — inspired by the Greek word *ροή* (flow). Add feeds, refresh them, and
read posts with read/unread tracking. Target audience: individuals who want a small, self-hosted feed reader, not a
multi-tenant SaaS.

## Repo layout

```
roe/
├── api/              Go backend (Huma + SQLC + PostgreSQL)
│   └── openapi.yaml  Generated OpenAPI spec (`mise -C api run schema`)
├── ui/               Next.js 16 frontend
│   └── client/       Generated TypeScript client (`mise run generate`)
├── compose.yml       Full-stack Podman Compose (api + ui + postgres)
└── .github/          CI (lint · vuln · test · typecheck · build)
```

See `api/CLAUDE.md` and `ui/CLAUDE.md` for per-service details.

## Running the full stack

```bash
# Requires: Podman
podman compose up
# API → http://localhost:8080
# UI  → http://localhost:3000
```

The `postgres` service has no migrations baked in — apply `api/db/migrations/*.up.sql` once against it
(`scripts/setup-db.sh`, or `psql` directly) before the API can serve requests.

For local development, run each service separately — see their CLAUDE.md files.

## CI (GitHub Actions)

Six jobs on push/PR to `main`:
- **lint** — golangci-lint on the Go code
- **vuln** — govulncheck on Go dependencies
- **lint-frontend** — Biome check on the TypeScript code
- **test** — `go test ./...` (routes are tested against a mocked `db.Querier`, no live database needed)
- **typecheck-frontend** — `tsc --noEmit`
- **build** — Podman image builds (depends on the five above)

Task runner is `mise` — see `mise.toml` for task definitions.

## Key constraints

- `api/.env` is gitignored. Never commit it.
- `ui/client/` is generated from `api/openapi.yaml` by `@hey-api/openapi-ts` — do not hand-edit it; run
  `mise run generate` (root) after changing any Huma route/type in `api/`.
- The Go API routes are unprefixed (`/feeds`, `/posts`, ...). The UI's `app/api/[...path]/route.ts` owns the `/api`
  prefix and strips it before forwarding upstream — the browser only ever calls relative `/api/*` paths.
