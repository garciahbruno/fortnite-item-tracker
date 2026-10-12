# Fortnite Item Tracker

A Go web app and Discord bot for tracking Fortnite cosmetics and getting notified when they return to the shop.

**Status: Day 1 foundation.** Shop ingestion, accounts, watchlists, and notifications are planned; they are not implemented yet. This project is not deployed yet and is not affiliated with Epic Games.

## Five-session build plan

| Day | Deliverable | Acceptance criteria |
| --- | --- | --- |
| 1 | Foundation | Gin API, embedded web shell, PostgreSQL migrations, Docker Compose, health checks |
| 2 | Shop and cosmetics | Provider client with timeouts; persistent shop snapshots; cosmetic search; scheduled refresh; explicit stale/error states |
| 3 | Authentication and tracking | Discord OAuth2 with state validation; secure JWT cookies; authenticated watchlist CRUD; ownership checks |
| 4 | Discord notifications | Bot setup; explicit opt-in; notification worker; retries and deduplication; alert preferences |
| 5 | Release | End-to-end verification; responsive UI polish; deployment instructions; operational checks |

Each session targets 3–5 meaningful commits. Resume with “start day 2” (or the next day). See `docs/PROGRESS.md` for the handoff. No unattended schedule is configured.

## Local startup

Install Docker with Compose and start its engine, then:

```sh
cp .env.example .env
docker compose up --build
```

Open http://localhost:8080. PostgreSQL is private to the Compose network; only the app is published, on localhost. The migration service runs before the app starts and safely skips applied migrations. Data persists in the `postgres_data` volume.

```sh
curl http://localhost:8080/healthz  # process is alive
curl http://localhost:8080/readyz   # database responds
curl http://localhost:8080/api/v1/status
docker compose down              # keeps database data
```

Never commit `.env`. The sample password is for local development only.

## Go development

Use Go 1.25 or later. Set `DATABASE_URL` to a reachable PostgreSQL database. `PORT` defaults to 8080 and `GIN_MODE` to release. Environment files are loaded by Compose, not the Go executable.

```sh
go mod download
go run ./cmd/migrate
go run ./cmd/server
go test ./...
go vet ./...
```

The API fails startup if the database is unavailable. Runtime readiness checks return 503 during a database outage; liveness remains 200. Shutdown drains HTTP requests before closing the connection pool.

## Layout

- `cmd/server`: HTTP process and graceful shutdown
- `cmd/migrate`: transactional migration command
- `internal/config`: validated environment configuration
- `internal/database`: PostgreSQL pool and embedded migration runner
- `internal/httpapi`: routes and health checks
- `internal/web`: embedded browser assets

Discord and Fortnite credentials are not needed until their integrations are implemented.
