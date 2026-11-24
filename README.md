# KepKep backend

Go backend for KepKep (ເກັບເກັບ), a Lao personal finance app.
Plan and phases: [KEPKEP_BACKEND_PLAN.md](KEPKEP_BACKEND_PLAN.md).

## Requirements

- Go 1.25+
- Docker (Postgres locally, and testcontainers in tests)
- `make tools` installs goose, sqlc and golangci-lint (macOS, Homebrew)

## Getting started

```sh
cp .env.example .env
docker compose up -d
make run
```

Run `make help` to see all commands.
