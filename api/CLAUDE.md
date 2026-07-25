# API — Roe Go backend

Go service that stores feeds and posts in PostgreSQL, fetches/parses RSS via `gofeed`, and exposes a Huma-generated
REST + OpenAPI API.

## Stack

| Concern | Package |
|---|---|
| HTTP framework | Huma v2 + chi v5 |
| CLI / lifecycle | Huma v2 `humacli` + `spf13/cobra` (adds the `openapi` subcommand) |
| Config | `caarlos0/env/v11` (struct tags on `internal/config.Config`) |
| Database | PostgreSQL via `pgx/v5` (`pgxpool`) |
| Query layer | SQLC — hand-written SQL in `db/queries`, generated code in `internal/db` |
| Feed parsing | `mmcdole/gofeed` |
| Module path | `github.com/angelospanag/roe` |

## Directory structure

```
api/
├── cmd/api/main.go             Entry point: config → pgxpool → router → server; `openapi` subcommand prints the spec
├── internal/
│   ├── config/
│   │   └── config.go            Env-var config struct (PORT, DATABASE_URL)
│   ├── middleware/
│   │   └── request_id.go        Request ID + access-log middleware (chi)
│   ├── feed/
│   │   ├── routes.go            Huma route registration for /feeds*
│   │   ├── service.go           Feed refresh logic (fetch, parse, upsert posts)
│   │   └── models.go            Request/response types
│   ├── post/
│   │   ├── routes.go            Huma route registration for /posts*
│   │   ├── service.go           Read/unread + listing logic
│   │   └── models.go            Request/response types
│   └── db/                      SQLC-generated code — do not hand-edit
│       ├── db.go, models.go, querier.go, feeds.sql.go, posts.sql.go
├── db/
│   ├── migrations/               Plain SQL, applied manually (golang-migrate or psql)
│   └── queries/                  Hand-written SQL — source for SQLC codegen
└── scripts/setup-db.sh           Runs migrations against $DB_* env vars via golang-migrate
```

## Route registration pattern

`feed.RegisterRoutes` and `post.RegisterRoutes` both take `(api huma.API, querier db.Querier, logger *slog.Logger)` —
routes are registered against the `db.Querier` **interface**, not a concrete `*db.Queries`, so tests register routes
against a hand-rolled `mockQuerier` (see `internal/feed/routes_test.go`) with no live database required.

`cmd/api/main.go`'s `newRouter(queries, logger)` builds the chi router + Huma API once; both the real server and the
`openapi` subcommand call it (the subcommand passes `nil` for `queries`, which is safe because registration never
calls query methods — only request handlers do).

## Logging

All logging is structured JSON via `slog.NewJSONHandler` (configured once in `cmd/api/main.go`,
`slog.SetDefault`'d so anything not request-scoped still logs JSON).

`internal/middleware/request_id.go`'s `RequestID` middleware generates a **UUIDv7** (`uuid.NewV7()` — time-ordered,
so request IDs sort chronologically and stay grep-friendly across log aggregation) per request, derives a logger
via `base.With("request_id", id.String())`, and stores it in the request context. It also logs one structured
`"request"` line per completed request (`method`, `path`, `status`, `duration_ms`, `remote_addr`).

Route handlers must fetch the request-scoped logger with `apimiddleware.LoggerFromContext(ctx)` — never
`slog.Default()` or a package-level logger — so every log line a handler emits carries the same `request_id` as the
access log for that request. See any handler in `internal/feed/routes.go` or `internal/post/routes.go` for the
pattern:

```go
func(ctx context.Context, input *X) (*Y, error) {
    logger := apimiddleware.LoggerFromContext(ctx)
    logger.Info("doing the thing", "some_field", value)
    ...
}
```

## humacli execution model

`humacli.New`'s setup closure runs for **every** cobra (sub)command, not just the default serve path — only
`hooks.OnStart`/`hooks.OnStop` are gated to the actual server run. That's why `pool.Ping` (the only part of startup
that touches the network) lives inside `hooks.OnStart` rather than the closure body: it keeps `go run ./cmd/api
openapi` working without a live Postgres connection.

## Environment variables

| Variable | Required | Default | Notes |
|---|---|---|---|
| `PORT` | no | `8080` | |
| `DATABASE_URL` | no | `postgres://postgres:postgres@localhost:5432/roe_backend?sslmode=disable` | |

## Dev

```bash
cd api
cp .env.example .env   # defaults work against a local Postgres on :5432
go run ./cmd/api        # starts on :8080
```

## Testing

```bash
cd api
go test ./...
```

Tests mock `db.Querier` directly — no database or service container needed in CI.

## Adding a new endpoint

1. Add input/output structs in the relevant `internal/{feed,post}/routes.go`.
2. Register with `huma.Register(api, huma.Operation{...}, handler)`. Set `DefaultStatus` explicitly when the
   convention isn't Huma's default (e.g. `http.StatusCreated` for a POST that creates a resource) — Huma already
   defaults to `204 No Content` for handlers whose output struct has no `Body` field.
3. Call through the `db.Querier` passed into `RegisterRoutes`.
4. Run `mise run schema` (or `mise -C api run schema` from the root) to regenerate `openapi.yaml`, then
   `mise -C ui run generate` to regenerate the TypeScript client.

## Adding a query

1. Write the SQL in `db/queries/{feeds,posts}.sql` (named query, sqlc annotation comment).
2. Run `mise run generate` (from `api/`) to regenerate `internal/db/*.sql.go`.
3. Add the new method to any mock implementations of `db.Querier` used in tests.

Huma generates OpenAPI docs automatically at `/docs` (Swagger UI).
