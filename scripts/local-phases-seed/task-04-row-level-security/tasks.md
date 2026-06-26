# Tasks — Row-level security

## Preparation

- [x] Read [README.md](README.md) and [official_source.md](official_source.md)
- [x] Run `./agent-harness/resolve-rules.sh owasp authz tenant migration`
- [x] Docker available for Postgres testcontainer

## Migrations

- [x] Add `services/attendance/migrations/` with sqlx
- [x] Migration: create minimal `employees` table with `tenant_id`
- [x] Migration: `ALTER TABLE employees ENABLE ROW LEVEL SECURITY`
- [x] Migration: `CREATE POLICY tenant_isolation ON employees USING (...)`
- [x] Repeat for `punch_records` and `face_embeddings` stubs per `docs/DATA-MODEL.md`

## Go infrastructure

- [x] `internal/infrastructure/postgres/tenant_tx.go` — `WithTenant(ctx, tenantID, fn)`
- [x] Executes `SET LOCAL app.tenant_id = $1` inside transaction
- [x] Repository example: `GetEmployee` uses `WithTenant`

## Integration tests

- [x] testcontainers-go Postgres 16
- [x] Seed tenant A and tenant B employees
- [x] Assert tenant A session cannot SELECT tenant B row
- [x] Assert same query with correct tenant returns row

## Validation

- [x] `go test ./...` including integration tag if used: `-tags=integration`
- [x] Migrations apply cleanly on empty DB
- [x] Document `DATABASE_URL` in `.env.example` if new vars needed

## Completion

- [x] All steps above marked `[x]`
- [x] Update `.local/phases/README.md` active task
