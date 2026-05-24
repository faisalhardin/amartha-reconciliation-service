# amartha-reconciliation-service

Go HTTP API for bank reconciliation: upload system transactions and bank statement CSVs, match them asynchronously, and retrieve reconciliation summaries.

## Prerequisites

- Go 1.25+
- Docker (for local PostgreSQL)
- `psql` (PostgreSQL client, for migrations)
- `make`

## Quick start

### 1. Environment

```bash
cp .env-dist .env
```

Edit `.env` if your database credentials differ. Defaults match the Docker Compose Postgres service.

### 2. Database

Start Postgres:

```bash
make db-up
```

Apply migrations (creates/updates tables under `schema/reconciliation/`):

```bash
make migrate
```

Or run both in one step:

```bash
make setup
```

### 3. Run the API

From the repository root (required so `/demo` and migrations resolve paths correctly):

```bash
make run
```

Server listens on `http://localhost:8080` (override with `PORT` in `.env`).

### 4. Interactive API demo

Open in a browser:

**http://localhost:8080/demo**

The demo page calls all main endpoints: bulk upload (transactions + bank statement), list files/statements/transactions, and reconciliation summary. After a bank statement upload, `fileId` is filled automatically for list/summary fields.

## Makefile targets

| Target | Description |
|--------|-------------|
| `make setup` | Start Postgres (`db-up`), wait, then `migrate` |
| `make db-up` | Start Postgres via Docker Compose |
| `make db-down` | Stop Postgres container |
| `make migrate` | Apply all SQL files in `schema/reconciliation/` (sorted by name) |
| `make env-setup` | Copy `.env-dist` → `.env` if missing |
| `make run` | Run the API (`go run cmd/api/main.go`) |
| `make build` | Build binary to `./main` |
| `make test` | Run unit tests |
| `make clean` | Remove `./main` binary |

Legacy aliases: `docker-run` = `db-up`, `docker-down` = `db-down`.

## Testing workflow (manual)

Use sample CSVs in `files/samples/` (e.g. `transactions_bca_t3.csv`, `bank_statement_bca_t3.csv` with timeframe `1779620400`–`1779638400`).

Recommended order:

1. **Upload transactions** — `POST /v1/transaction/bca/upload` (CSV)
2. **Upload bank statement** — `POST /v1/bank-statement/upload` with `bank_code`, `start_date`, `end_date`, and CSV file
3. Wait until the file **status** is `COMPLETED` (async worker; check via list files)
4. **Reconciliation summary** — `GET /v1/reconciliation/summary/{fileId}` (runs automatically after processing; no manual run required)

### curl examples

```bash
# Health
curl -sS http://localhost:8080/health

# Upload transactions
curl -sS -X POST "http://localhost:8080/v1/transaction/bca/upload" \
  -F "file=@files/samples/transactions_bca_t3.csv"

# Upload bank statement
curl -sS -X POST "http://localhost:8080/v1/bank-statement/upload" \
  -F "bank_code=bca" \
  -F "start_date=1779620400" \
  -F "end_date=1779638400" \
  -F "file=@files/samples/bank_statement_bca_t3.csv"

# List bank statement files
curl -sS "http://localhost:8080/v1/bank-statement/bca/files?limit=10"

# List bank statements for a file
curl -sS "http://localhost:8080/v1/bank-statement/bca?fileId=<FILE_ID>&limit=30"

# List transactions
curl -sS "http://localhost:8080/v1/transaction/bca?limit=30"

# Reconciliation summary (after file is COMPLETED)
curl -sS "http://localhost:8080/v1/reconciliation/summary/<FILE_ID>"
```

## API overview

| Method | Path | Description |
|--------|------|-------------|
| GET | `/health` | Health check |
| GET | `/demo` | Browser API showcase |
| POST | `/v1/transaction/{bankCode}/upload` | Bulk insert transactions (CSV) |
| GET | `/v1/transaction/{bankCode}` | List transactions (`limit`, `offset`) |
| POST | `/v1/bank-statement/upload` | Upload bank statement file (multipart) |
| GET | `/v1/bank-statement/{bankCode}/files` | List uploaded statement files |
| GET | `/v1/bank-statement/{bankCode}` | List statement rows (`fileId`, `limit`) |
| GET | `/v1/reconciliation/summary/{fileID}` | Reconciliation summary for a file |
| POST | `/v1/reconciliation/{bankCode}/run` | Run reconciliation manually (`fileId` query) |
| GET | `/v1/reconciliation/{bankCode}` | List reconciliation rows (`fileId`, optional `status`) |

## Database migrations

SQL migrations live in [`schema/reconciliation/`](schema/reconciliation/). They are applied in **filename sort order** by `make migrate`:

1. `20260522_amt_mst_transaction.sql` — transactions + bank statements
2. `20260524_amt_mst_bank_statement_file.sql` — statement file metadata
3. `20260525_amt_mst_bank_statement_file_status.sql` — file status, statement file link
4. `20260526_amt_mst_bank_statement_file_update_time.sql` — file `update_time`
5. `20260527_amt_trx_reconciliation.sql` — reconciliation results

Migrations are idempotent where possible (`IF NOT EXISTS`, etc.). Re-running `make migrate` on an already-migrated database is safe for normal development.

Connection settings for `make migrate` are read from `.env` (`DB_HOST`, `DB_PORT`, `DB_USERNAME`, `DB_PASSWORD`, `DB_NAME`).

## Troubleshooting

- **`make migrate` fails: `psql: command not found`** — Install the PostgreSQL client (`brew install libpq` on macOS, then add `psql` to PATH).
- **`connection refused` on migrate** — Run `make db-up` first and wait a few seconds.
- **`column does not exist` on upload** — Run `make migrate` against the same database as the app (check `.env`).
- **`404 page not found` on `/demo`** — Start the server from the repo root with `make run`.
- **Summary returns 404** — File must be `COMPLETED` and reconciliation must have finished; confirm rows with list files / list reconciliation.

## Project layout

See [CLAUDE.md](CLAUDE.md) for architecture, layering, and conventions.
