# Build handoff

## Day 1 — complete (October 11, 2026)

Implemented: Gin server, embedded responsive web shell, environment validation, graceful shutdown, health/readiness endpoints, PostgreSQL schema, serialized transactional migrations, non-root Docker image, Compose startup ordering, and GitHub Actions checks.

Validation completed:
- `go test -race ./...`, `go vet ./...`, and `go build ./cmd/...` passed with Go 1.27.2.
- Real PostgreSQL integration test passed: repeated migrations, watchlist/notification deduplication, and user deletion cascades.
- Docker image built and Compose migrations completed successfully, including a second migration run.
- Live HTTP checks passed for health, readiness, status, HTML, CSS, and JavaScript. `/.env` returned 404.
- Browser preview showed API Online and Database Connected; layout inspected.
- GitHub Actions workflow is configured; remote CI results are not yet verified.

Local checkout: `/Users/venyo/fortnite-item-tracker`.
Preview: http://localhost:18080 (8080 was occupied). The running Compose project is `fortnite-day1`, started using the sample development environment:

```sh
PORT=18080 docker compose --env-file .env.example -p fortnite-day1 up -d
PORT=18080 docker compose --env-file .env.example -p fortnite-day1 down
```

The app and its database were left running for preview. The isolated test database container was removed. A checksum-verified temporary Go SDK is available at `/private/tmp/fortnite-toolchain/go/bin/go`; it may disappear when temporary files are cleaned.

Remaining functionality: real shop ingestion, OAuth2/JWT, watchlists, Discord alerts, and public deployment. These belong to the later milestones.

## Next session: Day 2

1. Read this file, README, and repository status before editing.
2. Select and document the Fortnite shop/cosmetics provider and its credentials/terms.
3. Implement a bounded HTTP client with timeout and fixture tests for provider failures.
4. Persist cosmetics and shop snapshots transactionally. Retain the last successful snapshot on failure.
5. Add search/shop API routes and connect the browser UI to real data.
6. Add periodic refresh with cancellation and overlap prevention.
7. Verify, update this handoff, and make 3–5 coherent commits.

## Later decisions

- Discord application, redirect URL, bot token, and testing server are needed for Days 3–4.
- Hosting target and production secrets are needed for actual deployment on Day 5.
- Do not describe unimplemented features or deployment as complete.
