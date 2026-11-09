# ai-analysis

Personal system. The repository is currently at **Phase 1: Local Foundation**.

Phase 1 validates the application boundary without calling real job websites. It loads and validates configuration and a candidate profile, connects to PostgreSQL, applies embedded SQL migrations, persists an independent record for every invocation, and checks the optional TypeScript browser-worker boundary.

## Architecture

```text
Go CLI
  -> YAML configuration + candidate profile
  -> PostgreSQL migrations
  -> persisted run lifecycle
  -> optional TypeScript browser-worker health protocol
```

- **Go** owns the CLI, configuration, profile, database, locking, logging, and later the collection/matching pipeline.
- **TypeScript + Playwright** owns only browser-dependent collection. Phase 1 implements its typed NDJSON protocol and health action; login and collection explicitly return `not_implemented` until Phase 2.
- **PostgreSQL** owns durable state.
- **pgweb** provides local database inspection at <http://localhost:8081>.

## Prerequisites

- Go 1.23+
- Node.js 20+
- Docker Desktop / Docker Engine with Compose

## Run Phase 1

From this repository:

```bash
make setup
make up
make dry-run
```

Expected output includes:

```text
Phase 1 dry run completed
Status:      completed
PostgreSQL:  connected
Run record:  persisted
```

Open <http://localhost:8081>, select the `runs` table, and inspect the saved dry-run record.

Stop local services with:

```bash
make down
```

## Validate everything

```bash
make verify
```

The browser worker can be checked independently:

```bash
make browser-health
```

It responds with one JSON object. `executableInstalled` can be false in Phase 1 because no browser is launched yet; Phase 2 will install Chromium before authenticated collection is enabled.

## Configuration

- `config/config.example.yaml` is the committed application example.
- `config/config.local.yaml` is the ignored local application configuration.
- `config/profile.yaml` is the editable, non-secret candidate profile.
- `.env.example` documents the database connection variable.

The Makefile defaults to the database created by `compose.yaml`. To use another PostgreSQL database:

```bash
make dry-run DATABASE_URL='postgres://user:password@host/database?sslmode=require'
```

## CLI

```bash
go run ./cmd/jobfinder help
go run ./cmd/jobfinder version
DATABASE_URL='postgres://...' go run ./cmd/jobfinder run --dry-run
```

Running without `--dry-run` fails deliberately. Real collection is Phase 2.

## Phase boundaries

### Included now — Phase 1

- Go service and CLI
- Strict YAML decoding and validation
- Editable candidate profile
- PostgreSQL and embedded migration runner
- Persisted `runs` table
- Independent concurrent run records
- Structured JSON logs via `log/slog`
- Docker Compose with PostgreSQL and pgweb
- TypeScript browser-worker protocol and tests
- Unit tests

### Deferred to Phase 2

- Public ATS adapters
- Workday discovery
- Safe sitemap/JSON-LD crawler
- Real Playwright authentication and collection
- Job normalization, hashing, and persistence
- Source health

Matching and HTML/Markdown job reports arrive in Phase 3; unattended scheduling arrives in Phase 4.

## Browser-worker protocol

Go invokes the worker with newline-delimited JSON over standard input/output. This avoids exposing a second local HTTP port.

Request:

```json
{"id":"phase1-health","action":"health"}
```

Response:

```json
{"id":"phase1-health","status":"ok","result":{"worker":"ai-analysis-browser-worker","version":"0.1.0","phase":1,"engine":"chromium","executableInstalled":false}}
```

Authentication state will later live under `data/auth/<source-id>/`. That directory is excluded from Git. Passwords, OTPs, cookies, and browser storage must never be printed or committed.

## Database inspection

Useful pgweb query:

```sql
SELECT id, mode, status, started_at, finished_at, metadata, error_message
FROM runs
ORDER BY started_at DESC;
```
