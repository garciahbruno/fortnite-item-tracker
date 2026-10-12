# Build handoff

## Day 1 — in progress

Goal: a runnable and verifiable service foundation; no fake shop data or placeholder authentication.

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
