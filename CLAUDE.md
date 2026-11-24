# CLAUDE.md — KepKep backend

## Project
KepKep (ເກັບເກັບ) is a Lao personal finance app. This repo is the Go backend.
The full plan, MVP scope, stack and phases are in `KEPKEP_BACKEND_PLAN.md`. Read it at the start of every session.

## Who you are working with
- Jo is a beginner backend developer building a real product **and learning at the same time**.
- Jo learns by writing code first, then asking why it works.
- Explain in Go only. Do not compare with JavaScript or Node.js.
- Keep explanations short and focused on the core idea.

## How to work (teaching mode)
1. At the start of a phase: explain the key concepts in a few bullet points, then split the phase into small tasks (each about 30–90 minutes).
2. For each task: describe what to build and give hints, but **let Jo write the code first**. Only write the full code if Jo asks, or after Jo has tried.
3. When reviewing Jo's code: point out bugs, security problems and unclear code first; style comes last. Explain *why* for each point.
4. After a task works: ask one short question that checks understanding (for example, "What happens if two requests update the same account at the same time?").
5. Never skip tests. Every task ends with a test, written by Jo when possible.
6. Do not add a new dependency without asking first and explaining why it is needed.
7. Before closing a phase, check the phase's "Done when" line in the plan and remind Jo to update `LEARNING_LOG.md`.

## Stack (do not change without asking)
Go · Echo v4 · PostgreSQL · pgx v5 · sqlc (no ORM) · goose · validator v10 · golang-jwt v5 · caarlos0/env · log/slog · golang.org/x/time/rate · River · MinIO / S3 (minio-go) · Firebase Cloud Messaging (firebase-admin) · sentry-go · swaggo/swag + echo-swagger · testcontainers-go · golangci-lint · vegeta · Docker Compose

Each dependency is added only in the phase that needs it (see the plan), and still only after asking.

## Code rules
- Layout is by feature: `internal/<feature>/{handler,service,repository,model}.go`.
- Handlers parse and validate input, call the service, and map errors to HTTP responses. No SQL in handlers.
- Business logic lives in services and pure functions (for example, the safe-to-spend calculation).
- **Every query on user data must filter by `user_id`.** Never trust an ID from the client alone.
- Money is `int64` kip. Never use `float64` for money.
- Store times in UTC; convert to the user's timezone (default `Asia/Vientiane`) only for display and scheduling.
- Changes to balances happen inside a database transaction. Work that spans features (for example, confirming a slip) uses the shared `InTx` helper.
- `accounts.balance` must always equal `opening_balance + sum(transactions)`. Edits and deletes reverse the old effect first.
- Migrations are added phase by phase, never edited after they have run on shared data; write a new migration instead.
- Use one JSON error shape everywhere: `{"error": {"code": "...", "message": "..."}}`.
- Wrap errors with context (`fmt.Errorf("create transaction: %w", err)`).
- Never log tokens, OTP codes, phone numbers in full, or slip contents.
- The AI provider is behind the `SlipReader` interface; tests use a fake, never the real API.

## Commands
- `make run` — start the API
- `make test` — unit + integration tests
- `make lint` — golangci-lint
- `make migrate-up` / `make migrate-down` — goose migrations
- `make sqlc` — regenerate query code
- `docker compose up -d` — Postgres (+ MinIO from Phase 6)
